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
	ID uint32 `gorm:"primary_key"`

	Account   Account
	AccountID uint32

	Temperature float64
	Time        time.Time

	RecordedBy Account `gorm:"foreignkey:RecorderID"`
	RecorderID uint32
}

// insert a new record into database
func (h Handler) newRecord(w http.ResponseWriter, r *http.Request) {
	var err error
	var account Account

	err = h.db.First(&account, "id = ?", r.FormValue("account_id")).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, "account not found", 404)
		return
	} else if err != nil {
		panic(err)
	}

	acct, ok := r.Context().Value(KeyAccount).(Account)
	if !ok {
		http.Error(w, "cannot read account from session", 500)
		return
	}

	if !permission(acct, account) {
		http.Error(w, "permission denied", 403)
		return
	}

	record := Record{
		Account:    account,
		RecordedBy: acct,
	}

	record.Temperature, err = strconv.ParseFloat(r.FormValue("temperature"), 64)
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	record.Time = time.Now()

	err = h.db.Create(&record).Error
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(w)
	enc.Encode(record)
}

// permission return whether A can modify B
func permission(a, b Account) bool {
	switch a.Role {
	case Admin:
		return true

	case Teacher:
		return a.ClassID == b.ClassID

	case Student:
		return a.ID == b.ID
	}
	return false
}

func (h Handler) findRecord(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("account_id")
	var record Record
	err := h.db.Where("account_id = ? and time > ?", userID, today()).Order("id desc").First(&record).Error
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

func (h Handler) listRecord(acct Account) (tx *gorm.DB) {
	tx = h.db.Order("id desc").Set("gorm:auto_preload", true)
	switch acct.Role {
	case Admin:
		return tx

	case Teacher:
		return tx.Joins(
			"JOIN accounts on records.account_id = accounts.id",
		).Where(
			"accounts.class_id = ?", acct.ClassID,
		)

	case Student:
		return tx.Where("account_id = ?", acct.ID)
	}

	return nil
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
