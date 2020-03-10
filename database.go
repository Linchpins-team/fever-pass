package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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
	var admin Account
	err := db.First(&admin, "id = 'admin'").Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		panic(err)
	}
	admin.ID = "admin"
	admin.Name = "admin"
	admin.Role = Admin
	password := c.Password
	admin.Password, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = db.Save(&admin).Error
	if err != nil {
		panic(err)
	}
	if UseTestData {
		importTestData(db)
	}
}
