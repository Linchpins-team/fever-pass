package main

import (
	"context"
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

	r.Handle("/", http.FileServer(http.Dir("static")))
	h.router = r
}

func (h Handler) Router() *mux.Router {
	return h.router
}

func (h Handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := securecookie.New(hashKey, blockKey)
		if cookie, err := r.Cookie("session"); err == nil {
			var id uint32
			if err := s.Decode("session", cookie.Value, &id); err == nil {
				var acct Account
				if err = h.db.First(&acct, id).Error; err != nil {
					http.Error(w, "account not found", 401)
					return
				}
				ctx := r.Context()
				ctx = context.WithValue(ctx, KeyAccount, acct)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "session cannot be decode", 401)
			return
		}
		http.Error(w, "session not found", 401)
	})
}
