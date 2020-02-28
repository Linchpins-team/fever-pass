package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Role uint32

const (
	Unknown Role = iota
	Admin
	Teacher
	Student
)

func (r Role) String() string {
	switch r {
	case Admin:
		return "管理員"

	case Teacher:
		return "老師"

	case Student:
		return "學生"

	default:
		return "未知"
	}
}

type Account struct {
	ID       string `gorm:"primary_key;type:varchar(32)"`
	Name     string `gorm:"type:varchar(32)"`
	Password []byte `json:"-"`

	Class   Class
	ClassID uint32
	Number  int

	Role Role
}

func (a Account) String() string {
	return a.Name
}

func parseRole(str string) (Role, error) {
	role, err := strconv.Atoi(str)
	if err != nil || role < 0 || role > 3 {
		return Unknown, fmt.Errorf("Cannot parse '%s' as role", str)
	}
	return Role(role), nil
}

func generatePassword(password string) []byte {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return pwd
}

func (h Handler) deleteAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "admin" {
		http.Error(w, "cannot delete admin", 403)
		return
	}

	err := h.db.Delete(&Account{}, "id = ?", id).Error
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h Handler) updateAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var acct Account
	err := h.db.First(&acct, "id = ?", id).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, "user not found", 404)
		return
	} else if err != nil {
		panic(err)
	}

	role, _ := parseRole(r.FormValue("role"))
	if role != Unknown {
		if acct.ID == "admin" {
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

func (h Handler) listAccounts(acct Account) *gorm.DB {
	tx := h.db.Preload("Class").Table("accounts")
	switch acct.Role {
	case Admin:
		return tx

	case Teacher:
		return tx.Where("class_id = ?", acct.ClassID)

	case Student:
		return tx.Where("id = ?", acct.ID)

	default:
		return nil
	}
}
