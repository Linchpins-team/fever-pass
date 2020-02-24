package main

import (
	"database/sql"
	"net/http"
	"strconv"
)

func (h Handler) newRecordPage(w http.ResponseWriter, r *http.Request) {
	var records []Record
	if acct, ok := r.Context().Value(KeyAccount).(Account); ok {
		err := h.db.Where("account_id = ?", acct.ID).Order("time desc").Limit(20).Find(&records).Error
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
	h.tpl.ExecuteTemplate(w, "new.htm", page)
}

func (h Handler) listRecordsPage(w http.ResponseWriter, r *http.Request) {
	type recordT struct {
		Record
		Recorder sql.NullString
	}
	records := make([]recordT, 0, 20)
	p, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		p = 0
	}
	rows, err := h.db.Table("records").Select(
		"records.id, records.user_id, records.pass, records.time, accounts.name",
	).Joins(
		"left join accounts on records.account_id = accounts.id",
	).Order("time desc").Offset(p * 20).Limit(20).Rows()

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var record recordT
		err := rows.Scan(&record.ID, &record.UserID, &record.Pass, &record.Time, &record.Recorder)
		if err != nil {
			panic(err)
		}
		records = append(records, record)
	}

	page := struct {
		Page    int
		Records []recordT
	}{p, records}
	if err := h.tpl.ExecuteTemplate(w, "list.htm", page); err != nil {
		panic(err)
	}
}
