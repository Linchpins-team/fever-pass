package main

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

func (h Handler) index(w http.ResponseWriter, r *http.Request) {
	acct, ok := r.Context().Value(KeyAccount).(Account)
	if ok {
		var record Record
		err := h.db.Preload("RecordedBy").Where("created_at > ?", today()).Order("id desc").First(&record, "account_id = ?", acct.ID).Error
		if !gorm.IsRecordNotFoundError(err) && err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		h.HTML(w, r, "index.htm", record)
	} else {
		h.HTML(w, r, "index.htm", nil)
	}
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
	err = h.listRecord(acct).Offset((p - 1) * 20).Limit(20).Find(&records).Error
	if err != nil {
		panic(err)
	}

	page := struct {
		Page    int
		Records []Record
	}{p, records}
	h.HTML(w, r, "list.htm", page)
}

func (h Handler) listAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts []Account
	err := h.db.Preload("Class").Find(&accounts).Error
	if err != nil {
		panic(err)
	}

	h.HTML(w, r, "account_list.htm", accounts)
}
