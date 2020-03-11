package main

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

type ListType int

const (
	UnknownListType ListType = iota
	Recorded
	Unrecorded
	Fevered
)

func (h Handler) findStudents(acct Account, class string) (accounts []Account, err error) {
	tx := h.listAccounts(acct)
	if class != "" {
		err = h.db.First(&Class{}, "name = ?", class).Error
		if err != nil {
			return
		}
		tx = joinClasses(tx).Where("classes.name = ?", class)
	}
	if acct.Role == Teacher {
		class = acct.Class.Name
	}
	tx = tx.Where("role = ?", Student)

	err = tx.Find(&accounts).Error
	if err != nil {
		panic(err)
	}
	return
}

// 今日總覽頁面
func (h Handler) stats(w http.ResponseWriter, r *http.Request) {
	acct := r.Context().Value(KeyAccount).(Account)
	class := r.FormValue("class")
	accounts, err := h.findStudents(acct, class)
	if gorm.IsRecordNotFoundError(err) {
		h.errorPage(w, r, "找不到班級", "找不到班級"+r.FormValue("class"))
		return
	}
	page := struct {
		Class                  string
		All, Unrecorded, Fever int
	}{
		Class: class,
		All:   len(accounts),
	}
	for _, account := range accounts {
		record, err := h.lastRecord(account)
		if err == RecordNotFound {
			page.Unrecorded++
			continue
		}

		if record.Fever() {
			page.Fever++
		}
	}
	h.HTML(w, r, "stats.htm", page)
}

func (h Handler) statsList(w http.ResponseWriter, r *http.Request) {
	acct := r.Context().Value(KeyAccount).(Account)
	t := parseListType(r.FormValue("type"))
	if t == UnknownListType {
		http.Error(w, "Unknown list type", 415)
		return
	}

	accounts, err := h.findStudents(acct, r.FormValue("class"))
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, "class not found", 404)
		return
	}

	result := make([]Account, 0, len(accounts))
	for _, account := range accounts {
		record, err := h.lastRecord(account)
		if err == RecordNotFound && t == Unrecorded {
			result = append(result, account)
			continue
		}
		if err == nil && t == Recorded {
			result = append(result, account)
			continue
		}
		if record.Fever() && t == Fevered {
			result = append(result, account)
		}
	}
	if err = h.tpls["stats.htm"].ExecuteTemplate(w, "account_list", result); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func parseListType(s string) ListType {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 3 {
		return UnknownListType
	}
	return ListType(n)
}
