package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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
	r.Use(logger)

	r.HandleFunc("/api/records", h.auth(h.newRecord, Self, None)).Methods("POST")
	r.HandleFunc("/api/records/{id}", h.auth(h.deleteRecord, Self, None)).Methods("DELETE")

	r.HandleFunc("/api/accounts/{id}", h.auth(h.deleteAccount, None, All)).Methods("DELETE")
	r.HandleFunc("/api/accounts/{id}", h.auth(h.updateAccountAuthority, None, Group)).Methods("PUT")
	r.HandleFunc("/api/accounts", h.auth(h.findAccountByClassAndNumber, Self, Self)).Methods("GET")
	r.HandleFunc("/api/stats", h.auth(h.statsList, Group, None))

	r.HandleFunc("/api/login", h.login)

	r.HandleFunc("/new", h.auth(h.newRecordPage, Group, None))
	r.HandleFunc("/list", h.auth(h.listRecordsPage, Self, None))
	r.HandleFunc("/accounts", h.auth(h.listAccountsPage, Self, Self))
	r.HandleFunc("/stats", h.auth(h.stats, Group, None))
	r.HandleFunc("/import", h.auth(h.importPage, None, All)).Methods("GET")
	r.HandleFunc("/import", h.auth(h.importHandler, None, All)).Methods("POST")
	r.HandleFunc("/export", h.auth(h.exportCSV, Group, None))
	r.HandleFunc("/accounts/{id}", h.auth(h.profile, Self, Self))
	r.HandleFunc("/login", h.page("login.htm")).Methods("GET")
	r.HandleFunc("/login", h.login).Methods("POST")

	r.HandleFunc("/doc/{title}", h.doc)

	r.HandleFunc("/", h.index).Methods("GET")
	r.HandleFunc("/", h.auth(h.newSelfRecord, Self, None)).Methods("POST")
	r.Handle("/reset", h.auth(h.resetPage, None, Self)).Methods("GET")
	r.Handle("/reset", h.auth(h.resetPassword, None, Self)).Methods("POST")
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/register", h.auth(h.registerPage, None, All)).Methods("GET")
	r.HandleFunc("/register", h.auth(h.register, None, All)).Methods("POST")
	r.HandleFunc("/manifest.json", serveFile("manifest.json"))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	h.router = r
}

func serveFile(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	}
}

func (h Handler) Router() *mux.Router {
	return h.router
}

func logger(next http.Handler) http.Handler {
	return handlers.CustomLoggingHandler(os.Stdout, next, logFormat)
}

func logFormat(writer io.Writer, params handlers.LogFormatterParams) {
	acct, ok := session(params.Request)
	if !ok {
		acct.ID = "-"
	}

	_, err := fmt.Fprintf(writer, `%s - %s %s "%s %s %s" %d %d`,
		params.Request.RemoteAddr, acct.ID, params.TimeStamp.Format("2006-01-02 03:04:05"),
		params.Request.Method, params.URL.RequestURI(), params.Request.Proto,
		params.StatusCode, params.Size,
	)
	fmt.Fprintln(writer)
	if err != nil {
		fmt.Fprintln(writer, err)
	}
}
