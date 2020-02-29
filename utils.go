package main

import (
	"fmt"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format("01/02 15:04 ") + weekday(t)
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

func sub(a, b int) int {
	return a - b
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("01/02")
}

func weekday(t time.Time) string {
	switch t.Weekday() {
	case time.Monday:
		return "星期一"

	case time.Tuesday:
		return "星期二"

	case time.Wednesday:
		return "星期三"

	case time.Thursday:
		return "星期四"

	case time.Friday:
		return "星期五"

	case time.Saturday:
		return "星期六"

	case time.Sunday:
		return "星期天"
	}
	return ""
}

func weekdayColor(t time.Time) string {
	switch t.Weekday() {
	case time.Monday:
		return "red"

	case time.Tuesday:
		return "orange"

	case time.Wednesday:
		return "yellow"

	case time.Thursday:
		return "green"

	case time.Friday:
		return "blue"

	case time.Saturday:
		return "indigo"

	case time.Sunday:
		return "purple"
	}
	return ""
}
