package main

import "time"

type Record struct {
	ID          uint32
	UserID      string
	Pass        bool
	Temperature float64
	Time        time.Time
}
