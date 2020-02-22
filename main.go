package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/gorm.sqlite")
	if err != nil {
		panic(err)
	}
	migrate(db)
	return db
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Record{}, &Account{}).Error; err != nil {
		panic(err)
	}
}

func main() {
	h := NewHandler(initDB())
	defer h.db.Close()

	srv := &http.Server{
		Handler: h.Router(),
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listen on http://%s", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
