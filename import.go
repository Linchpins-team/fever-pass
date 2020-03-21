package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
)

func (h Handler) importPage(w http.ResponseWriter, r *http.Request) {
	page := make(map[string]interface{})
	page["authorities"] = Authorities
	page["message"] = r.Context().Value(KeyMessage)
	h.HTML(w, r, "import.htm", page)
}

// Handle the upload file request
func (h Handler) importHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		r = addMessage(r, err.Error())
		h.importPage(w, r)
		return
	}

	if !strings.HasSuffix(header.Filename, ".csv") {
		r = addMessage(r, "請上傳 CSV 檔案")
		h.importPage(w, r)
		return
	}

	n, err := importAccounts(
		h.db,
		file,
		r.FormValue("authority"),
	)
	if err != nil {
		r = addMessage(r, err.Error())
		h.importPage(w, r)
		return
	}
	r = addMessage(r, fmt.Sprintf("成功匯入%d筆資料", n))
	h.importPage(w, r)
}

func importTestData(db *gorm.DB) {
	file, err := os.Open("testdata/test-teachers.csv")
	if err != nil {
		panic(err)
	}
	_, err = importAccounts(db, file, Tutor.Key())
	if err != nil {
		panic(err)
	}
	file, err = os.Open("testdata/test-students.csv")
	if err != nil {
		panic(err)
	}
	_, err = importAccounts(db, file, Student.Key())
}

func importAccounts(db *gorm.DB, r io.Reader, authroity string) (n int, err error) {
	reader := csv.NewReader(r)
	index, _ := reader.Read() // ignore column name
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		columns := parseColumns(index, row)
		columns["authority"] = authroity
		added, err := createAccount(db, columns)
		if err != nil {
			return n, err
		}
		if added {
			n++
		}
	}
	return n, nil
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

func createAccount(db *gorm.DB, columns map[string]string) (_ bool, err error) {
	_, err = NewAccount(
		db,
		columns["id"],
		columns["name"],
		columns["password"],
		columns["class"],
		columns["number"],
		columns["authority"],
	)
	if err == AccountAlreadyExist {
		return false, nil
	}
	return true, err
}
