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

// Handle the upload file request
func (h Handler) importHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(415)
		h.HTML(w, r, "import.htm", err.Error())
		return
	}

	if !strings.HasSuffix(header.Filename, ".csv") {
		w.WriteHeader(415)
		h.HTML(w, r, "import.htm", "請上傳 CSV 檔案")
		return
	}

	n, err := importAccounts(
		h.db,
		file,
		r.FormValue("role"),
		r.FormValue("record_authority"),
		r.FormValue("account_authority"),
	)
	if err != nil {
		w.WriteHeader(500)
		h.HTML(w, r, "import.htm", err.Error())
		return
	}
	h.HTML(w, r, "import.htm", fmt.Sprintf("成功匯入%d筆資料", n))
}

func importTestData(db *gorm.DB) {
	file, err := os.Open("testdata/test-teachers.csv")
	if err != nil {
		panic(err)
	}
	_, err = importAccounts(db, file, "teacher", "group", "self")
	if err != nil {
		panic(err)
	}
	file, err = os.Open("testdata/test-students.csv")
	if err != nil {
		panic(err)
	}
	_, err = importAccounts(db, file, "student", "self", "self")
}

func importAccounts(db *gorm.DB, r io.Reader, role, recordAuthority, accountAuthority string) (n int, err error) {
	reader := csv.NewReader(r)
	index, _ := reader.Read() // ignore column name
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		columns := parseColumns(index, row)
		columns["role"] = role
		columns["record_authority"] = recordAuthority
		columns["account_authority"] = accountAuthority
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
		columns["role"],
		columns["record_authority"],
		columns["account_authority"],
	)
	if err == AccountAlreadyExist {
		return false, nil
	}
	return true, err
}
