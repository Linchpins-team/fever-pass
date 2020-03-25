package main

import (
	"testing"
)

func TestImport(t *testing.T) {
	importTestData(testH.db)
	var students []Account
	testH.db.Preload("Class").Find(&students)
	t.Log(students)
}
