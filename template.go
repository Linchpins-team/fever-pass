package main

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
)

func (h *Handler) loadTemplates() {
	h.tpls = make(map[string]*template.Template)
	mainTmpl := template.New("main")
	mainTmpl.Funcs(template.FuncMap{
		"formatTime": formatTime,
	})
	layoutFiles, err := filepath.Glob("templates/layouts/*.htm")
	if err != nil {
		panic(err)
	}

	includeFiles, err := filepath.Glob("templates/*.htm")
	if err != nil {
		panic(err)
	}

	log.Println(includeFiles)
	log.Println(layoutFiles)
	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		tpl := template.Must(mainTmpl.Clone())
		h.tpls[fileName] = template.Must(tpl.ParseFiles(files...))
	}
	log.Println(h.tpls)
}

func (h Handler) HTML(w http.ResponseWriter, page string, data interface{}) {
	log.Println(page)
	if tpl, ok := h.tpls[page]; ok {
		if err := tpl.ExecuteTemplate(w, "main", data); err != nil {
			http.Error(w, err.Error(), 500)
		}
	} else {
		log.Println(tpl)
		http.Error(w, "cannot find templates", 500)
	}
}

func (h Handler) newRecordPage(w http.ResponseWriter, r *http.Request) {
	var records []Record
	if acct, ok := r.Context().Value(KeyAccount).(Account); ok {
		err := h.db.Where("account_id = ?", acct.ID).Order("id desc").Limit(20).Find(&records).Error
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
	h.HTML(w, "new.htm", page)
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
	).Order("id desc").Offset(p * 20).Limit(20).Rows()

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var record recordT
		err := rows.Scan(&record.ID, &record.UserID, &record.Fever, &record.Time, &record.Recorder)
		if err != nil {
			panic(err)
		}
		records = append(records, record)
	}

	page := struct {
		Page    int
		Records []recordT
	}{p, records}
	h.HTML(w, "list.htm", page)
}

func (h Handler) page(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.HTML(w, path, nil)
	}
}
