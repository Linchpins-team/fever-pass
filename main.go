package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db *gorm.DB
)

func initDB() {
	var err error
	db, err = gorm.Open("sqlite3", "/tmp/gorm.sqlite")
	if err != nil {
		panic(err)
	}
	migrate(db)
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Record{}, &Account{}).Error; err != nil {
		panic(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/api/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hi")
	})

	/*
		r.HandleFunc("/api/records").Methods("GET").Queries("user_id")
		r.HandleFunc("/api/records").Methods("GET")
		r.HandleFunc("/api/records").Methods("POST")
		r.HandleFunc("/api/records").Methods("PUT")
	*/

	r.Handle("/", http.FileServer(http.Dir("static")))

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listen on http://%s", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
