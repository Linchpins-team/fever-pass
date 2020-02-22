package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
	r.HandleFunc("/api/records", h.newRecord).Methods("POST")
	r.HandleFunc("/api/records", h.listRecord).Methods("GET")
	r.HandleFunc("/api/records/{id}", h.deleteRecord).Methods("DELETE")

	r.Handle("/", http.FileServer(http.Dir("static")))
	h.router = r
}

func (h Handler) Router() *mux.Router {
	return h.router
}
