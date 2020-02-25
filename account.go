package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Role uint32

const (
	Unknown Role = iota
	Admin
	Editor
	User
)

func (r Role) String() string {
	switch r {
	case Admin:
		return "管理員"

	case Editor:
		return "編輯者"

	case User:
		return "使用者"

	default:
		return "未知"
	}
}

type Account struct {
	ID       uint32 `gorm:"primary_key"`
	Name     string `gorm:"unique;type:varchar(32)"`
	Password []byte

	Role Role
}

func parseRole(str string) (Role, error) {
	role, err := strconv.Atoi(str)
	if err != nil || role < 0 || role > 3 {
		return Unknown, fmt.Errorf("Cannot parse '%s' as role", str)
	}
	return Role(role), nil
}

func session(id uint32) *http.Cookie {
	s := securecookie.New(hashKey, blockKey)
	var encoded string
	var err error
	if encoded, err = s.Encode("session", id); err != nil {
		panic(err)
	}
	return &http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/",
	}
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	var acct Account
	acct.Name = r.FormValue("username")
	err := h.db.Where(acct).First(&acct).Error
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
	http.SetCookie(w, session(acct.ID))
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
	acct.Role = Editor
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

	http.SetCookie(w, session(acct.ID))
}

func generatePassword(password string) []byte {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return pwd
}

func (h Handler) auth(next http.HandlerFunc, role Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := securecookie.New(hashKey, blockKey)
		if cookie, err := r.Cookie("session"); err == nil {
			var id uint32
			if err := s.Decode("session", cookie.Value, &id); err == nil {
				var acct Account
				if err = h.db.First(&acct, id).Error; err != nil {
					http.Error(w, "account not found", 401)
					return
				}
				switch {
				case acct.Role == Unknown:
					http.Error(w, "unknown role", 401)
					return
				case acct.Role > role:
					http.Error(w, "permission denied", 401)
					return
				}
				ctx := r.Context()
				ctx = context.WithValue(ctx, KeyAccount, acct)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "session cannot be decode", 401)
			logout(w, r)
		} else {
			http.Error(w, err.Error(), 401)
		}
	}
}

func (h Handler) deleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	if id == 1 {
		http.Error(w, "cannot delete admin", 403)
		return
	}
	err = h.db.Delete(&Account{}, "id = ?", id).Error
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h Handler) updateAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	var acct Account
	err = h.db.First(&acct, "id = ?", id).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, "user not found", 404)
		return
	} else if err != nil {
		panic(err)
	}

	role, _ := parseRole(r.FormValue("role"))
	if role != Unknown {
		if acct.ID == 1 {
			http.Error(w, "cannot change admin role", 403)
			return
		}
		acct.Role = role
	}

	password := r.FormValue("password")
	if password != "" {
		acct.Password = generatePassword(password)
	}

	err = h.db.Save(&acct).Error
	if err != nil {
		panic(err)
	}
}
