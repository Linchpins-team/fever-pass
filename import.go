package main

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/jinzhu/gorm"
)

func (h Handler) importAccounts(r io.Reader, role Role) (err error) {
	reader := csv.NewReader(r)
	reader.Read() // ignore column name
	tx := h.db.Begin()
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		err = newAccount(tx, row, role)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func newAccount(db *gorm.DB, row []string, role Role) (err error) {
	if len(row) < 4 {
		return fmt.Errorf("row doesn't 4 column")
	}
	acct := Account{
		Name:  row[1],
		Email: row[0],
	}
	acct.Role = role
	password := row[2]

	var class Class
	class.Name = row[3]
	if err := db.FirstOrCreate(&class).Error; err != nil {
		return err
	}
	acct.Class = class

	acct.Password = generatePassword(password)

	if err := db.Create(&acct).Error; err != nil {
		return err
	}
	return nil
}
