package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	PermissionDenied = errors.New("permission denied")
	RecordNotFound   = errors.New("record not found")
)

type Record struct {
	gorm.Model

	Temperature float64

	Account    Account
	AccountID  uint32
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

	var record Record
	err = h.db.Preload("Account").First(&record, id).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, RecordNotFound.Error(), 404)
	} else if err != nil {
		panic(err)
	}

	acct := r.Context().Value(KeyAccount).(Account)
	if !permission(acct, record.Account) {
		http.Error(w, PermissionDenied.Error(), 403)
		return
	}

	err = h.db.Delete(&record).Error
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
