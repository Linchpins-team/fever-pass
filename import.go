package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"net/http"

	"github.com/jinzhu/gorm"
)

func (h Handler) importHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(415)
		h.HTML(w, r, "import.htm", err.Error())
		return
	}
	role, err := parseRole(r.FormValue("role"))
	if err != nil {
		w.WriteHeader(415)
		h.HTML(w, r, "import.htm", err.Error())
		return
	}
	n, err := importAccounts(h.db, file, role)
	if err != nil {
		w.WriteHeader(500)
		h.HTML(w, r, "import.htm", err.Error())
		return
	}
	h.HTML(w, r, "import.htm", fmt.Sprintf("成功匯入%d筆資料", n))
}

func importTestData(db *gorm.DB) {
	file, err := os.Open("testdata/teachers.csv")
	if err != nil {
		panic(err)
	}
	_, err = importAccounts(db, file, Teacher)
	if err != nil {
		panic(err)
	}
	file, err = os.Open("testdata/students.csv")
	if err != nil {
		panic(err)
	}
	_, err = importAccounts(db, file, Student)
}

func importAccounts(db *gorm.DB, r io.Reader, role Role) (n int, err error) {
	reader := csv.NewReader(r)
	index, _ := reader.Read() // ignore column name
	tx := db.Begin()
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		added, err := createAccount(tx, parseColumns(index, row), role)
		if err != nil {
			tx.Rollback()
			return n, err
		}
		if added {
			n++
		}
	}
	return n, tx.Commit().Error
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

func createAccount(db *gorm.DB, columns map[string]string, role Role) (added bool, err error) {
	acct := Account{
		ID:   columns["id"],
		Name: columns["name"],
	}
	acct.Role = role
	password := columns["password"]

	if err := db.First(&acct, "id = ?", acct.ID).Error; !gorm.IsRecordNotFoundError(err) {
		return false, err 
	}

	var class Class
	if err := db.FirstOrCreate(&class, Class{Name: columns["class"]}).Error; err != nil {
		return false, err
	}
	acct.Class = class

	acct.Number, err = strconv.Atoi(columns["number"])
	if err != nil {
		acct.Number = 0
	}

	acct.Password = generatePassword(password)

	return true, db.Create(&acct).Error
}
