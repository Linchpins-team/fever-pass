package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/jinzhu/gorm"
)

var (
	testfile = strings.NewReader(
		`"email","name","password","class"
"s10512345@st.fcjh.tc.edu.tw","Justin","j_password","207"
"s10642236@st.fcjh.tc.edu.tw","Kevin","k_pwd","109"
"s10412556@st.fcjh.tc.edu.tw","Elsa","anna","303"
"s10443256@st.fcjs.tc.edu.tw","Anna","elsa","303"`,
	)
)

func importAccounts(db *gorm.DB, r io.Reader, role Role) (err error) {
	reader := csv.NewReader(r)
	reader.Read() // ignore column name
	tx := db.Begin()
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
	if err := db.FirstOrCreate(&class, Class{Name: row[3]}).Error; err != nil {
		return err
	}
	acct.Class = class

	acct.Password = generatePassword(password)

	if err := db.Create(&acct).Error; err != nil {
		return err
	}
	return nil
}
