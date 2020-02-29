package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (h Handler) export(w io.Writer) (err error) {
	enc := csv.NewWriter(w)
	var accounts []Account
	tx := h.db.Joins("JOIN classes ON classes.id = class_id")
	tx = tx.Order("classes.name asc").Order("number asc")
	tx = tx.Preload("Class").Where("role = ?", Student)
	if err = tx.Find(&accounts).Error; err != nil {
		return
	}

	columns := []string{"班級", "座號", "姓名", "體溫", "類型", "發燒", "時間", "紀錄者"}
	err = enc.Write(columns)
	if err != nil {
		panic(err)
	}
	for _, account := range accounts {
		row := []string{
			account.Class.String(),
			strconv.Itoa(account.Number),
			account.Name,
		}
		record, err := h.lastRecord(account)
		if err == nil {
			row = append(row, record.CSV()...)
		}
		if err = enc.Write(row); err != nil {
			panic(err)
		}
	}
	enc.Flush()
	return err
}

/*
"class.name","number","account.name","temperature","type","fever","created_at","recorded_by"
*/
func (r Record) CSV() []string {
	return []string{
		fmt.Sprint(r.Temperature),
		r.Type.String(),
		r.FeverString(),
		r.CreatedAt.Format("2006-01-02 03:04:05"),
		r.RecordedBy.Name,
	}
}

func (r Record) FeverString() string {
	if r.Fever() {
		return "發燒"
	} else {
		return "正常"
	}
}

func (h Handler) exportCSV(w http.ResponseWriter, r *http.Request) {
	filename := time.Now().Format("fever-pass-2006-01-02_03-04-export.csv")
	w.Header().Set("Content-Type", "text/x-csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	err := h.export(w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
