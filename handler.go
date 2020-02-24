package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	db     *gorm.DB
	router *mux.Router
	tpl    *template.Template
}

type ContextKey uint32

const (
	KeyAccount ContextKey = iota
)

func NewHandler(db *gorm.DB) Handler {
	h := Handler{
		db: db,
	}
	h.tpl = template.New("")
	h.tpl.Funcs(template.FuncMap{
		"formatTime": formatTime,
	})
	h.tpl.ParseGlob("templates/*.htm")
	h.newRouter()
	return h
}

func (h *Handler) newRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/api/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hi")
	})

	r.HandleFunc("/api/check", h.findRecord).Methods("GET")

	r.HandleFunc("/api/records", h.auth(h.newRecord, Editor)).Methods("POST")
	r.HandleFunc("/api/records", h.auth(h.listRecord, Editor)).Methods("GET")
	r.HandleFunc("/api/records/{id}", h.auth(h.deleteRecord, Editor)).Methods("DELETE")

	r.HandleFunc("/api/login", h.login)
	r.HandleFunc("/api/logout", logout)
	r.HandleFunc("/api/register", h.auth(h.register, Admin))

	r.HandleFunc("/admin/new", h.auth(h.newRecordPage, Editor))
	r.HandleFunc("/admin/list", h.auth(h.listRecordsPage, Editor))
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	h.router = r
}

func (h Handler) Router() *mux.Router {
	return h.router
}
