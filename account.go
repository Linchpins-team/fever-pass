package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	AccountNotFound = errors.New("找不到此帳號")
)

type Account struct {
	ID       string `gorm:"primary_key;type:varchar(32)"`
	Name     string `gorm:"type:varchar(32)"`
	Password []byte `json:"-"`

	Class   Class
	ClassID uint32
	Number  int

	CreatedAt time.Time
	DeletedAt *time.Time

	Role Role

	RecordAuthority  Authority
	AccountAuthority Authority
}

func (a Account) String() string {
	return a.Name
}

func (a Account) permission(tempAuthority, acctAuthority Authority) bool {
	return a.AccountAuthority >= acctAuthority && a.RecordAuthority >= tempAuthority
}

func (a *Account) defaultAuthority(role Role) {
	switch role {
	case Admin:
		a.RecordAuthority = All
		a.AccountAuthority = All

	case Staff:
		a.RecordAuthority = Group
		a.AccountAuthority = Group

	case Student:
		a.RecordAuthority = Self
		a.AccountAuthority = Self

	default:
		a.RecordAuthority = None
		a.AccountAuthority = None
	}
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

func (h Handler) updateAccountAuthority(w http.ResponseWriter, r *http.Request) {
	account, err := h.getAccount(mux.Vars(r)["id"])

	if err == AccountNotFound {
		http.Error(w, err.Error(), 404)
		return
	}

	acct, _ := session(r)
	if !accountPermission(acct, account) {
		http.Error(w, PermissionDenied.Error(), 403)
		return
	}

	if role != Unknown {
		if account.ID == "admin" {
			http.Error(w, "cannot change admin role", 403)
			return
		}
		account.Role = role
	}

	err = h.db.Save(&account).Error
	if err != nil {
		panic(err)
	}
}

func joinClasses(tx *gorm.DB) *gorm.DB {
	return tx.Table("accounts").Joins("JOIN classes ON class_id = classes.id")
}

func (h Handler) listAccounts(acct Account) *gorm.DB {
	tx := joinClasses(h.db).Preload("Class").Order("classes.name, number asc")
	switch acct.AccountAuthority {
	case All:
		return tx

	case Group:
		return tx.Where("class_id = ?", acct.ClassID)

	case Self:
		return tx.Where("id = ?", acct.ID)

	default:
		return nil
	}
}

func (h Handler) getAccount(id string) (acct Account, err error) {
	err = h.db.First(&acct, "id = ?", id).Error
	if gorm.IsRecordNotFoundError(err) {
		err = AccountNotFound
		return
	} else if err != nil {
		panic(err)
	}
	return
}

// 重設密碼
func (h Handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	acct := r.Context().Value(KeyAccount).(Account)
	account, err := h.getAccount(r.FormValue("account_id"))
	if err == AccountNotFound {
		account = acct
	}

	if !recordPermission(acct, account) {
		w.WriteHeader(403)
		h.resetPage(w, addMessage(r, "您沒有權限變更 "+account.Name+" 的密碼"))
		return
	}

	current := r.FormValue("current_password")
	if bcrypt.CompareHashAndPassword(acct.Password, []byte(current)) != nil {
		w.WriteHeader(403)
		h.resetPage(w, addMessage(r, "密碼錯誤"))
		return
	}

	account.Password = generatePassword(r.FormValue("new_password"))

	if err := h.db.Save(&account).Error; err != nil {
		w.WriteHeader(500)
		h.resetPage(w, addMessage(r, err.Error()))
		return
	}

	h.resetPage(w, addMessage(r, "已重設 "+account.Name+" 的密碼"))
}

func (h Handler) findAccountByClassAndNumber(w http.ResponseWriter, r *http.Request) {
	var err error
	acct := r.Context().Value(KeyAccount).(Account)

	var account Account
	err = joinClasses(h.db).Where(
		"classes.name = ? and number = ?", r.FormValue("class"), r.FormValue("number"),
	).Set("gorm:auto_preload", true).First(&account).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, AccountNotFound.Error(), 404)
		return
	} else if err != nil {
		panic(err)
	}

	if !recordPermission(acct, account) {
		http.Error(w, PermissionDenied.Error(), 403)
		return
	}

	if _, err = fmt.Fprint(w, account.ID); err != nil {
		panic(err)
	}
}
