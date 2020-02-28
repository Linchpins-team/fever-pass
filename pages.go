package main

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

func (h Handler) index(w http.ResponseWriter, r *http.Request) {
	acct, ok := r.Context().Value(KeyAccount).(Account)
	if ok {
		record, err := h.lastRecord(acct)
		if err == RecordNotFound {
			h.HTML(w, r, "index.htm", nil)
		}
		h.HTML(w, r, "index.htm", record)
	} else {
		h.HTML(w, r, "index.htm", nil)
	}
}

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
	if acct, ok := r.Context().Value(KeyAccount).(Account); ok {
		err := h.listRecord(acct).Limit(20).Find(&records).Error
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		http.Error(w, "cannot read account from session", 500)
		return
	}
	page := struct {
		Records []Record
	}{records}
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

	err = tx.Offset((p - 1) * 20).Limit(20).Find(&records).Error
	if err != nil {
		panic(err)
	}

	page := struct {
		Page    int
		Account Account
		Records []Record
	}{
		Page:    p,
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
func (h Handler) status(w http.ResponseWriter, r *http.Request) {
	acct := r.Context().Value(KeyAccount).(Account)
	tx := h.listAccounts(acct)
	class := r.FormValue("class")
	if class != "" {
		tx = tx.Joins("JOIN classes ON class_id = classes.id").Where("classes.name = ?", class)
	}
	if acct.Role == Teacher {
		class = acct.Class.Name
	}

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
	h.HTML(w, r, "status.htm", page)
}
