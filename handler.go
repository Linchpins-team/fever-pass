package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	db     *gorm.DB
	router *mux.Router
	config Config
	tpls   map[string]*template.Template
}

type ContextKey uint32

const (
	KeyAccount ContextKey = iota
	KeyMessage
)

func NewHandler(db *gorm.DB, config Config) Handler {
	h := Handler{
		db:     db,
		config: config,
	}
	h.loadTemplates()
	h.newRouter()
	return h
}

func (h *Handler) newRouter() {
	r := mux.NewRouter()

	r.Use(h.identify)

	r.HandleFunc("/api/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hi")
	})

	r.HandleFunc("/api/records", h.auth(h.newRecord, Student)).Methods("POST")
	r.HandleFunc("/api/records/{id}", h.auth(h.deleteRecord, Student)).Methods("DELETE")

	r.HandleFunc("/api/accounts/{id}", h.auth(h.deleteAccount, Admin)).Methods("DELETE")
	r.HandleFunc("/api/accounts/{id}", h.auth(h.updateAccount, Admin)).Methods("PUT")

	r.HandleFunc("/api/login", h.login)
	r.HandleFunc("/api/register", h.register).Methods("POST")

	r.HandleFunc("/api/url", h.auth(h.newURL, Admin)).Methods("POST")

	r.HandleFunc("/new", h.auth(h.newRecordPage, Teacher))
	r.HandleFunc("/list", h.auth(h.listRecordsPage, Student))
	r.HandleFunc("/invite", h.auth(h.page("generate_url.htm"), Admin))
	r.HandleFunc("/accounts", h.auth(h.listAccountsPage, Student))
	r.HandleFunc("/stats", h.auth(h.stats, Teacher))
	r.HandleFunc("/import", h.auth(h.page("import.htm"), Admin)).Methods("GET")
	r.HandleFunc("/import", h.auth(h.importHandler, Admin)).Methods("POST")
	r.HandleFunc("/export", h.auth(h.exportCSV, Teacher))

	r.HandleFunc("/doc/{title}", h.doc)

	r.HandleFunc("/", h.index).Methods("GET")
	r.HandleFunc("/", h.auth(h.newSelfRecord, Student)).Methods("POST")
	r.Handle("/reset", h.auth(h.resetPage, Student)).Methods("GET")
	r.Handle("/reset", h.auth(h.resetPassword, Student)).Methods("POST")
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/register", h.registerPage)

	r.HandleFunc("/qrcodes/{file}", h.auth(h.qrcode, Admin))

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	h.router = r
}

func (h Handler) Router() *mux.Router {
	return h.router
}
