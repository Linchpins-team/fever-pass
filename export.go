package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (h Handler) export(records []AccountRecords, w io.Writer) (err error) {
	enc := csv.NewWriter(w)

	columns := []string{"帳號", "班級", "座號", "姓名", "體溫", "類型", "發燒", "時間"}
	err = enc.Write(columns)
	if err != nil {
		panic(err)
	}
	for _, r := range records {
		row := r.CSV()
		if err = enc.Write(row); err != nil {
			panic(err)
		}
	}
	enc.Flush()
	return err
}

/*
"account.id", "class.name","number","account.name","temperature","type","fever","created_at"
*/
func (r AccountRecords) CSV() []string {
	data := []string{
		r.Account.ID,
		r.Class.Name,
		strconv.Itoa(r.Number),
		r.Name,
	}
	if r.Recorded {
		data = append(data, []string{
			fmt.Sprint(r.Temperature),
			r.Type.String(),
			r.FeverString(),
			r.Record.CreatedAt.Format("2006-01-02 03:04:05"),
		}...)
	}
	return data
}

func (r Record) FeverString() string {
	if r.Fever() {
		return "發燒"
	}
	return "正常"
}

func (h Handler) exportCSV(w http.ResponseWriter, r *http.Request) {
	filename := time.Now().Format("fever-pass-2006-01-02_03-04-export.csv")
	w.Header().Set("Content-Type", "text/x-csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	acct, _ := session(r)
	tx, err := statsBase(h.db, acct, "")
	tx = statsQuery(h.db, tx, Complete, today())
	err = h.export(selectRecords(tx), w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
