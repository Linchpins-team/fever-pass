package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Record struct {
	ID          uint32
	UserID      string
	Pass        bool
	Temperature float64
	Time        time.Time

	RecordedBy Account
	AccountID  uint32
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
		log.Println(acct)
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
