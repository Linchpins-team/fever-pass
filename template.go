package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/jinzhu/gorm"
)

func (h *Handler) loadTemplates() {
	h.tpls = make(map[string]*template.Template)
	mainTmpl := template.New("main")
	mainTmpl.Funcs(template.FuncMap{
		"formatTime":   formatTime,
		"formatDate":   formatDate,
		"weekdayColor": weekdayColor,
		"add":          add,
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

func (h Handler) HTML(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	log.Println(page)
	acct, ok := r.Context().Value(KeyAccount).(Account)
	pageData := struct {
		Data  interface{}
		Login bool
		Account
	}{
		data,
		ok,
		acct,
	}
	if tpl, ok := h.tpls[page]; ok {
		if err := tpl.ExecuteTemplate(w, "main", pageData); err != nil {
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
		err := h.listRecord(acct).Limit(20).Find(&records).Error
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
	h.HTML(w, r, "new.htm", page)
}

func (h Handler) listRecordsPage(w http.ResponseWriter, r *http.Request) {
	records := make([]Record, 0, 20)
	p, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		p = 1
	}
	// acct must have value
	acct := r.Context().Value(KeyAccount).(Account)
	err = h.listRecord(acct).Offset((p - 1) * 20).Limit(20).Find(&records).Error
	if err != nil {
		panic(err)
	}

	page := struct {
		Page    int
		Records []Record
	}{p, records}
	h.HTML(w, r, "list.htm", page)
}

func (h Handler) page(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.HTML(w, r, path, nil)
	}
}

func (h Handler) check(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("account_id")
	var record Record
	page := struct {
		Record
		Found bool
	}{}
	err := h.db.Where("account_id = ? and time > ?", userID, today()).Order("id desc").First(&record).Error
	if gorm.IsRecordNotFoundError(err) {
		page.Found = false
		h.HTML(w, r, "check.htm", page)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	page.Found = true
	page.Record = record

	h.HTML(w, r, "check.htm", page)
}

func (h Handler) listAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts []Account
	err := h.db.Find(&accounts).Error
	if err != nil {
		panic(err)
	}

	h.HTML(w, r, "account_list.htm", accounts)
}
