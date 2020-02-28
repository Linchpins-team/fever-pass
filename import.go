package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
)

func importTestData(db *gorm.DB) {
	file, err := os.Open("testdata/teachers.csv")
	if err != nil {
		panic(err)
	}
	importAccounts(db, file, Teacher)
	file, err = os.Open("testdata/students.csv")
	if err != nil {
		panic(err)
	}
	importAccounts(db, file, Student)
}

func importAccounts(db *gorm.DB, r io.Reader, role Role) (err error) {
	reader := csv.NewReader(r)
	index, _ := reader.Read() // ignore column name
	tx := db.Begin()
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		err = newAccount(tx, parseColumns(index, row), role)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func parseColumns(index, row []string) (columns map[string]string) {
	columns = make(map[string]string)
	for i, v := range index {
		if len(row) > i {
			columns[v] = row[i]
		}
	}
	return
}

func newAccount(db *gorm.DB, columns map[string]string, role Role) (err error) {
	if len(columns) < 4 {
		return fmt.Errorf("row doesn't 4 column")
	}
	acct := Account{
		ID:   columns["id"],
		Name: columns["name"],
	}
	acct.Role = role
	password := columns["password"]

	var class Class
	if err := db.FirstOrCreate(&class, Class{Name: columns["class"]}).Error; err != nil {
		return err
	}
	acct.Class = class

	acct.Number, err = strconv.Atoi(columns["number"])
	if err != nil {
		acct.Number = 0
	}

	acct.Password = generatePassword(password)

	if err := db.Create(&acct).Error; err != nil {
		return err
	}
	return nil
}
