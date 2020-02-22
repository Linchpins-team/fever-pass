package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type Record struct {
	ID          uint32
	UserID      string
	Pass        bool
	Temperature float64
	Time        time.Time

	RecordedBy Account `json:"-"`
	AccountID  uint32  `json:"-"`
}

// insert a new record into database
func (h Handler) newRecord(w http.ResponseWriter, r *http.Request) {
	var record Record
	var err error
	record.UserID = r.PostFormValue("user_id")

	record.Pass, err = strconv.ParseBool(r.PostFormValue("pass")) // default false
	if err != nil {
		http.Error(w, "Pass cannot be parsed", 415)
		return
	}

	record.Temperature, err = strconv.ParseFloat(r.PostFormValue("temperature"), 64)
	if err != nil {
		http.Error(w, "Temperature cannot be parsed", 415)
		return
	}

	record.Time = time.Now()
	if acct, ok := r.Context().Value(KeyAccount).(Account); ok {
		record.RecordedBy = acct
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
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	err := h.db.Where("user_id = ? and time > ?", userID, today).Order("time desc").First(&record).Error
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
