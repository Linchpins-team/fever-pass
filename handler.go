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

	r.HandleFunc("/api/check", h.findRecord).Methods("GET")

	r.HandleFunc("/api/records", h.auth(h.newRecord, Teacher)).Methods("POST")
	r.HandleFunc("/api/records", h.auth(h.listRecord, Teacher)).Methods("GET")
	r.HandleFunc("/api/records/{id}", h.auth(h.deleteRecord, Teacher)).Methods("DELETE")

	r.HandleFunc("/api/accounts/{id}", h.auth(h.deleteAccount, Admin)).Methods("DELETE")
	r.HandleFunc("/api/accounts/{id}", h.auth(h.updateAccount, Admin)).Methods("PUT")

	r.HandleFunc("/api/login", h.login)
	r.HandleFunc("/api/register", h.register).Methods("POST")

	r.HandleFunc("/api/url", h.auth(h.newURL, Admin)).Methods("POST")

	r.HandleFunc("/admin/new", h.auth(h.newRecordPage, Teacher))
	r.HandleFunc("/admin/list", h.auth(h.listRecordsPage, Teacher))
	r.HandleFunc("/admin/invite", h.auth(h.page("generate_url.htm"), Admin))
	r.HandleFunc("/admin/accounts", h.auth(h.listAccounts, Admin))

	r.HandleFunc("/doc/{title}", h.doc)

	r.HandleFunc("/", h.page("index.htm"))
	r.Handle("/login", h.page("login.htm"))
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/register", h.registerPage)
	r.HandleFunc("/check", h.check)

	r.HandleFunc("/qrcodes/{file}", h.auth(h.qrcode, Admin))

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	h.router = r
}

func (h Handler) Router() *mux.Router {
	return h.router
}
