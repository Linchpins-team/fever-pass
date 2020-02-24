package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Record struct {
	ID     uint32 `gorm:"primary_key"`
	UserID string
	Pass   bool
	Time   time.Time

	RecordedBy Account `json:"-"`
	AccountID  uint32  `json:"-"`
}

// insert a new record into database
func (h Handler) newRecord(w http.ResponseWriter, r *http.Request) {
	var record Record
	var err error
	record.UserID = r.FormValue("user_id")

	record.Pass, err = parseBool(r.FormValue("pass")) // default false
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	record.Time = time.Now()

	if acct, ok := r.Context().Value(KeyAccount).(Account); ok {
		record.AccountID = acct.ID
	} else {
		http.Error(w, "cannot read account from session", 500)
		return
	}

	err = h.db.Create(&record).Error
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(w)
	enc.Encode(record)
}

func (h Handler) findRecord(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	var record Record
	err := h.db.Where("user_id = ? and time > ?", userID, today()).Order("time desc").First(&record).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, "record not found", 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	enc := json.NewEncoder(w)
	if err = enc.Encode(&record); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h Handler) listRecord(w http.ResponseWriter, r *http.Request) {
	var records []Record
	err := h.db.Where("time > ?", today()).Order("time desc").Find(&records).Error

	enc := json.NewEncoder(w)
	if err = enc.Encode(&records); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h Handler) deleteRecord(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	err = h.db.Delete(&Record{}, "id = ?", id).Error
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
