package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	ID       string
	ExpireAt time.Time
}

func init() {
	gob.Register(Session{})
}

func (h Handler) auth(next http.HandlerFunc, role Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if role == Unknown {
			next.ServeHTTP(w, r)
		} else if acct, ok := r.Context().Value(KeyAccount).(Account); ok && acct.Role <= role {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "permission denied, require "+role.String(), 401)
		}
	}
}

func (h Handler) identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := securecookie.New(hashKey, blockKey)
		if cookie, err := r.Cookie("session"); err == nil {
			var session Session
			if err := s.Decode("session", cookie.Value, &session); err == nil {
				if time.Now().After(session.ExpireAt) {
					// session expire
					logout(w, r)
				} else {
					var acct Account
					if err = h.db.Set("gorm:auto_preload", true).First(&acct, "id = ?", session.ID).Error; err != nil {
						logout(w, r)
						http.Error(w, "account not found", 401)
						return
					}
					ctx := r.Context()
					ctx = context.WithValue(ctx, KeyAccount, acct)
					r = r.WithContext(ctx)
				}
			} else {
				logout(w, r)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func newSession(id string) *http.Cookie {
	s := securecookie.New(hashKey, blockKey)
	var encoded string
	var err error
	session := Session{
		ID:       id,
		ExpireAt: expire(),
	}
	if encoded, err = s.Encode("session", session); err != nil {
		panic(err)
	}
	return &http.Cookie{
		Name:    "session",
		Value:   encoded,
		Path:    "/",
		Expires: session.ExpireAt,
	}
}

func expire() time.Time {
	return time.Now().AddDate(0, 0, 7)
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	var acct Account
	err := h.db.First(&acct, "id = ?", r.FormValue("username")).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, "user not found", 404)
		return
	} else if err != nil {
		panic(err)
	}
	password := r.FormValue("password")
	if bcrypt.CompareHashAndPassword(acct.Password, []byte(password)) != nil {
		http.Error(w, "wrong password", 401)
		return
	}
	http.SetCookie(w, newSession(acct.ID))
}

func logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/", 303)
}

func (h Handler) register(w http.ResponseWriter, r *http.Request) {
	var err error
	var acct Account
	acct.ID = r.FormValue("account_id")
	if err = h.db.First(&acct, "id = ?", acct.ID).Error; !gorm.IsRecordNotFoundError(err) {
		h.HTML(w, r, "register.htm", "使用者帳號已經存在")
		return
	}
	acct.Name = r.FormValue("name")
	acct.Role, err = parseRole(r.FormValue("role"))
	if err != nil {
		h.errorPage(w, r, 415, "無效的身份", fmt.Sprintf("身份 '%s' 無法被解析", r.FormValue("role")))
		return
	}

	acct.Password = generatePassword(r.FormValue("password"))

	if err = h.db.Create(&acct).Error; err != nil {
		h.HTML(w, r, "register.htm", "無法註冊使用者："+err.Error())
		return
	}

	h.HTML(w, r, "register.htm", "成功註冊"+acct.Name)
}
