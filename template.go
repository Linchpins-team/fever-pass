package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
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

func (h Handler) page(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.HTML(w, r, path, nil)
	}
}
