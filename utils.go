package main

import (
	"fmt"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format("01-02 15:04")
}

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

func parseBool(str string) (bool, error) {
	switch str {
	case "on", "true":
		return true, nil
	case "off", "false", "":
		return false, nil
	default:
		return false, fmt.Errorf("'%s' cannot parse to bool", str)
	}
}

func add(a, b int) int {
	return a + b
}
