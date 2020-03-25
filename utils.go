package main

import (
	"fmt"
	"strings"
	"time"
)

func formatTime(t time.Time) string {
	return fmt.Sprintf("%s%s%s", t.Format("01/02"), weekday(t), t.Format("15:04"))
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
		return "（一）"

	case time.Tuesday:
		return "（二）"

	case time.Wednesday:
		return "（三）"

	case time.Thursday:
		return "（四）"

	case time.Friday:
		return "（五）"

	case time.Saturday:
		return "（六）"

	case time.Sunday:
		return "（日）"
	}
	return ""
}

func dashToSlash(s string) string {
	return strings.ReplaceAll(s, "-", "/")
}
