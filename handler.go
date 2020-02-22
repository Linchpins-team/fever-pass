package main

import (
	"github.com/jinzhu/gorm"
)

type Handler struct {
	db *gorm.DB
}

type ContextKey uint32

const (
	KeyAccount ContextKey = iota
)
