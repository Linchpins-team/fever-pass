package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	AccountNotFound = errors.New("找不到此帳號")
	InvalidValue    = errors.New("無效的資料")
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

func NewAccount(db *gorm.DB, id, name, password, class, number, role, recordAuthority, accountAuthority string) (account Account, err error) {
	account = Account{
		ID:       id,
		Name:     name,
		Password: generatePassword(password),
	}

	if err = db.First(&account, "id = ?", id).Error; !gorm.IsRecordNotFoundError(err) {
		return account, AccountAlreadyExist
	}

	account.Number, err = strconv.Atoi(number)
	if err != nil {
		account.Number = 0
	}

	if err = db.FirstOrCreate(&account.Class, Class{Name: class}).Error; err != nil {
		return
	}

	account.Role = parseRole(role)
	if account.Role == 0 {
		return account, fmt.Errorf("%w: 無效的身份%s", InvalidValue, role)
	}
	account.RecordAuthority = parseAuthority(recordAuthority)
	if account.RecordAuthority == None {
		return account, fmt.Errorf("%w: 無效的體溫權限%s", InvalidValue, recordAuthority)
	}
	account.AccountAuthority = parseAuthority(accountAuthority)
	if account.AccountAuthority == None {
		return account, fmt.Errorf("%w: 無效的帳號權限%s", InvalidValue, accountAuthority)
	}

	err = db.Create(&account).Error
	return
}

func (a Account) permission(recordAuthority, acctAuthority Authority) bool {
	return a.AccountAuthority >= acctAuthority && a.RecordAuthority >= recordAuthority
}

func (account *Account) updateAuthority(role, recordAuthority, accountAuthority string) {

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

	setAuthority := func(authority, max Authority, out *Authority) {
		if authority == 0 {
			http.Error(w, InvalidValue.Error(), 415)
			return
		}
		if authority > max {
			http.Error(w, PermissionDenied.Error(), 403)
			return
		}
		*out = authority
		return
	}

	if role := parseRole(r.FormValue(KeyRole)); role != 0 && acct.AccountAuthority == All {
		account.Role = role
	} else if authority := parseAuthority(r.FormValue(KeyRecordAuthority)); authority != None {
		setAuthority(authority, acct.RecordAuthority, &account.RecordAuthority)
	} else if authority := parseAuthority(r.FormValue(KeyAccountAuthority)); authority != None {
		setAuthority(authority, acct.AccountAuthority, &account.AccountAuthority)
	}

	err = h.db.Save(&account).Error
	if err != nil {
		panic(err)
	}
}

func (h Handler) listAccounts(acct Account) *gorm.DB {
	tx := joinClasses(h.db).Preload("Class").Order("role, classes.name, number asc")
	authority := acct.RecordAuthority
	if acct.AccountAuthority > authority {
		authority = acct.AccountAuthority
	}
	switch authority {
	case All:
		return tx

	case Group:
		return tx.Where("class_id = ?", acct.ClassID)

	case Self:
		return tx.Where("accounts.id = ?", acct.ID)

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
	var account Account
	tx := whereClass(h.db, r.FormValue("class"))
	tx = whereNumber(tx, r.FormValue("number"))
	err = tx.First(&account).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, AccountNotFound.Error(), 404)
		return
	} else if err != nil {
		panic(err)
	}

	acct, _ := session(r)
	if !recordPermission(acct, account) {
		http.Error(w, PermissionDenied.Error(), 403)
		return
	}

	if _, err = fmt.Fprint(w, account.ID); err != nil {
		panic(err)
	}
}
