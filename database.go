package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

func initDB(c Config) (db *gorm.DB, err error) {
	switch c.Mode {
	case Debug:
		db, err = gorm.Open("sqlite3", "/tmp/gorm.sqlite")

	case Release:
		connection := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", c.Database.User, c.Database.Password, c.Database.Name)
		db, err = gorm.Open("mysql", connection)
	}
	if err != nil {
		return
	}
	migrate(db)
	return
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Account{}, &Record{}, &Class{}).Error; err != nil {
		panic(err)
	}
}

func setupDB(c Config, db *gorm.DB) {
	db.Unscoped().Where("id = 'admin'").Delete(&Account{})
	_, err := NewAccount(
		db,
		"admin",
		"admin",
		c.Password,
		"T",
		"0",
		Admin.Key(),
	)
	if err != nil {
		panic(err)
	}
	if UseTestData {
		importTestData(db)
	}
}
