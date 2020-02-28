package main

import (
	"context"
	"encoding/gob"
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
					if err = h.db.First(&acct, "id = ?", session.ID).Error; err != nil {
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
	http.Redirect(w, r, "/login", 303)
}

func (h Handler) register(w http.ResponseWriter, r *http.Request) {
	var err error
	var url URL
	key := r.FormValue("key")
	err = h.db.Where("id = ? and expire_at > ? and last > 0", key, time.Now()).First(&url).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, "invalid key", 404)
		return
	} else if err != nil {
		panic(err)
	}

	var acct Account
	acct.Name = r.FormValue("username")
	acct.Role = Teacher
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	acct.Password = generatePassword(r.FormValue("password"))

	tx := h.db.Begin()
	url.Last--
	if err = tx.Save(&url).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	if err = tx.Create(&acct).Error; err != nil {
		tx.Rollback()
		http.Error(w, "cannot register user "+err.Error(), 500)
		return
	}

	if err = tx.Commit().Error; err != nil {
		panic(err)
	}

	http.SetCookie(w, newSession(acct.ID))
}
