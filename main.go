package main

import (
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
}

func main() {
	initDB()
	defer db.Close()
}
