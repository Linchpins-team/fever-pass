package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	db     *gorm.DB
	router *mux.Router
}

type ContextKey uint32

const (
	KeyAccount ContextKey = iota
)

var (
	hashKey  = securecookie.GenerateRandomKey(32)
	blockKey = securecookie.GenerateRandomKey(32)
)

func NewHandler(db *gorm.DB) Handler {
	h := Handler{
		db: db,
	}
	h.newRouter()
	return h
}

func (h *Handler) newRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/api/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hi")
	})

	r.HandleFunc("/api/check", h.findRecord).Methods("GET")

	{
		s := r.PathPrefix("/api/records").Subrouter()
		s.Use(h.auth)
		s.HandleFunc("", h.newRecord).Methods("POST")
		s.HandleFunc("", h.listRecord).Methods("GET")
		s.HandleFunc("/{id}", h.deleteRecord).Methods("DELETE")
	}

	r.HandleFunc("/api/login", h.login)
	r.HandleFunc("/api/logout", logout)
	r.HandleFunc("/api/register", h.register)
	r.Handle("/", http.FileServer(http.Dir("static")))
	h.router = r
}

func (h Handler) Router() *mux.Router {
	return h.router
}
