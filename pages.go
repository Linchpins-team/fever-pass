package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

func (h Handler) index(w http.ResponseWriter, r *http.Request) {
	acct, ok := r.Context().Value(KeyAccount).(Account)
	if ok {
		record, err := h.lastRecord(acct)
		if err == nil {
			h.HTML(w, r, "index.htm", record)
			return
		}
	}
	h.HTML(w, r, "index.htm", nil)

}

// get the last record today of the account
func (h Handler) lastRecord(account Account) (record Record, err error) {
	err = h.db.Set("gorm:auto_preload", true).Where("created_at > ?", today()).Order("id desc").First(&record, "account_id = ?", account.ID).Error
	if gorm.IsRecordNotFoundError(err) {
		err = RecordNotFound
		return
	} else if err != nil {
		panic(err)
	}
	return
}

func (h Handler) newRecordPage(w http.ResponseWriter, r *http.Request) {
	var records []Record
	acct, ok := r.Context().Value(KeyAccount).(Account)
	if ok {
		err := h.listRecord(acct).Limit(20).Find(&records).Error
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		http.Error(w, "cannot read account from session", 500)
		return
	}

	class := ""
	if acct.Role == Teacher {
		class = acct.Class.Name
	}

	page := struct {
		Class string
		Records []Record
	}{class, records}
	h.HTML(w, r, "new.htm", page)
}

func (h Handler) listRecordsPage(w http.ResponseWriter, r *http.Request) {
	records := make([]Record, 0, 20)
	p, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		p = 1
	}
	// acct must have value
	acct := r.Context().Value(KeyAccount).(Account)
	tx := h.listRecord(acct)
	id := r.FormValue("account_id")
	var account Account
	if id != "" {
		err = h.db.First(&account, "id = ?", id).Error
		if gorm.IsRecordNotFoundError(err) {
			http.Error(w, RecordNotFound.Error(), 404)
			return
		} else if err != nil {
			panic(err)
		}
		tx = tx.Where("account_id = ?", id)
	}

	date, err := time.ParseInLocation("2006-01-02", r.FormValue("date"), time.Local)
	if err == nil {
		tx = tx.Where("created_at > ? and created_at < ?", date, date.AddDate(0, 0, 1))
	}

	err = tx.Offset((p - 1) * 20).Limit(20).Find(&records).Error
	if err != nil {
		panic(err)
	}

	page := struct {
		Page    int
		Date    time.Time
		Account Account
		Records []Record
	}{
		Page:    p,
		Date:    date,
		Account: account,
		Records: records,
	}
	h.HTML(w, r, "list.htm", page)
}

func (h Handler) listAccountsPage(w http.ResponseWriter, r *http.Request) {
	acct := r.Context().Value(KeyAccount).(Account)
	var accounts []Account
	err := h.listAccounts(acct).Find(&accounts).Error
	if err != nil {
		panic(err)
	}

	h.HTML(w, r, "account_list.htm", accounts)
}

// 今日總覽頁面
func (h Handler) stats(w http.ResponseWriter, r *http.Request) {
	acct := r.Context().Value(KeyAccount).(Account)
	tx := h.listAccounts(acct)
	class := r.FormValue("class")
	if class != "" {
		tx = tx.Joins("JOIN classes ON class_id = classes.id").Where("classes.name = ?", class)
	}
	if acct.Role == Teacher {
		class = acct.Class.Name
	}
	tx = tx.Where("role = ?", Student)

	var accounts []Account
	err := tx.Find(&accounts).Error
	if err != nil {
		panic(err)
	}

	page := struct {
		Class                  string
		All, Unrecorded, Fever []Account
	}{
		Class: class,
		All:   accounts,
	}
	for _, account := range accounts {
		record, err := h.lastRecord(account)
		if err == RecordNotFound {
			page.Unrecorded = append(page.Unrecorded, account)
			continue
		}

		if record.Fever() {
			page.Fever = append(page.Fever, account)
		}
	}
	h.HTML(w, r, "stats.htm", page)
}

func (h Handler) resetPage(w http.ResponseWriter, r *http.Request) {
	acct := r.Context().Value(KeyAccount).(Account)
	account, err := h.getAccount(r.FormValue("account_id"))
	if err == AccountNotFound {
		account = acct
	}

	msg, ok := r.Context().Value(KeyMessage).(string)
	if !ok {
		msg = ""
	}

	if !permission(acct, account) {
		msg = "您沒有權限變更" + account.Name + "的密碼"
	}

	page := struct {
		Account
		Message string
	}{
		Account: account,
		Message: msg,
	}

	h.HTML(w, r, "reset.htm", page)
}

func addMessage(r *http.Request, msg string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, KeyMessage, msg)
	return r.WithContext(ctx)
}
