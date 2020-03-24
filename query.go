package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

// join accounts table with classes table
func joinClasses(tx *gorm.DB) *gorm.DB {
	return tx.Joins("JOIN classes ON accounts.class_id = classes.id")
}

// select specific class when query accounts
func whereClass(tx *gorm.DB, className string) *gorm.DB {
	if className == "" {
		return tx
	}
	return tx.Where("classes.name = ?", className)
}

// select specific number when query accounts
func whereNumber(tx *gorm.DB, number string) *gorm.DB {
	if number == "" {
		return tx
	}
	return tx.Where("accounts.number = ?", number)
}

// select specific date when query records
func whereDate(tx *gorm.DB, date time.Time) *gorm.DB {
	return tx.Where("records.created_at > ? and records.created_at < ?", date, date.AddDate(0, 0, 1))
}
