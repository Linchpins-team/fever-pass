package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (h *Handler) loadTemplates() {
	h.tpls = make(map[string]*template.Template)
	mainTmpl := template.New("main")
	mainTmpl.Funcs(template.FuncMap{
		"formatTime":  formatTime,
		"formatDate":  formatDate,
		"add":         add,
		"sub":         sub,
		"dashToSlash": dashToSlash,
	})
	layoutFiles, err := filepath.Glob("templates/layouts/*.htm")
	if err != nil {
		panic(err)
	}

	customFiles, err := filepath.Glob("templates/custom/*.htm")
	if err != nil {
		panic(err)
	}

	includeFiles, err := filepath.Glob("templates/*.htm")
	if err != nil {
		panic(err)
	}

	files := make([]string, 1, len(layoutFiles)+len(customFiles)+1)
	files = append(files, layoutFiles...)
	files = append(files, customFiles...)
	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files[0] = file
		tpl := template.Must(mainTmpl.Clone())
		h.tpls[fileName] = template.Must(tpl.ParseFiles(files...))
	}
}

func (h Handler) HTML(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	acct, ok := r.Context().Value(KeyAccount).(Account)
	pageData := struct {
		Data  interface{}
		Login bool
		Account
		Config  Config
		Message string
	}{
		data,
		ok,
		acct,
		h.config,
		"",
	}
	if msg, ok := r.Context().Value(KeyMessage).(string); ok {
		pageData.Message = msg
	}
	if tpl, ok := h.tpls[page]; ok {
		if err := tpl.ExecuteTemplate(w, "main", pageData); err != nil {
			http.Error(w, err.Error(), 500)
		}
	} else {
		http.Error(w, "cannot find templates", 500)
	}
}

func (h Handler) page(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.HTML(w, r, path, nil)
	}
}

func (h Handler) message(w http.ResponseWriter, r *http.Request, title, msg string) {
	h.HTML(w, r, "message.htm", struct {
		Title, Message string
	}{title, msg})
}
