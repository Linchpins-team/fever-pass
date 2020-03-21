package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	testDB *gorm.DB
	admin  = Account{
		Role:     RoleAdmin,
		Name:     "admin",
		Password: []byte{},
	}
	mockAccounts = []Account{
		Account{
			Name: "小明",
			Class: Class{
				Name: "103",
			},
			Role: RoleStudent,
		},
		Account{
			Name: "小美",
			Class: Class{
				Name: "103",
			},
			Role: RoleStudent,
		},
		Account{
			Name: "陳老師",
			Class: Class{
				Name: "103",
			},
		},
	}
	mockRecords = []Record{
		Record{
			AccountID:   "1",
			Temperature: 36.8,
		},
		Record{
			AccountID:   "2",
			Temperature: 37.9,
		},
		Record{
			AccountID:   "3",
			Temperature: 38.1,
		},
	}
	testH Handler
)

func TestMain(m *testing.M) {
	setupTestDB()
	os.Exit(m.Run())
}

func setupTestDB() {
	var err error
	testDB, err = gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	c := loadConfig("config.toml")
	testH = NewHandler(testDB, c)
	migrate(testDB)
	insertMockData(testDB)
}

func insertMockData(db *gorm.DB) {
	err := db.Create(&admin).Error
	if err != nil {
		panic(err)
	}
	for _, record := range mockRecords {
		record.RecordedBy = admin
		err := db.Create(&record).Error
		if err != nil {
			panic(err)
		}
	}
}

func TestNewRecord(t *testing.T) {
	record := Record{
		AccountID:   "2",
		Temperature: 37.2,
	}

	body := fmt.Sprintf("account_id=%s&temperature=%f&type=%d", record.AccountID, record.Temperature, Forehead)
	rr := testHandler("POST", "/api/records", body)
	if rr.Code != 200 {
		t.Errorf("status code is not 200, got %d\n%s\n", rr.Code, rr.Body.String())
	}
	err := json.Unmarshal(rr.Body.Bytes(), &record)
	assert.Equal(t, nil, err)

	var r Record
	err = testH.db.Where("id = ?", record.ID).First(&r).Error
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, record.Temperature, r.Temperature)
}

func adminSession(r *http.Request) *http.Request {
	r.AddCookie(newSession(admin.ID))
	return r
}

func testHandler(method, url, body string) *httptest.ResponseRecorder {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}

	req = adminSession(req)
	rr := httptest.NewRecorder()
	testH.Router().ServeHTTP(rr, req)
	return rr
}

func TestDeleteRecord(t *testing.T) {
	id := mockRecords[0].ID
	url := fmt.Sprintf("/api/records/%d", id)
	rr := testHandler("DELETE", url, "")
	if rr.Code != 200 {
		t.Errorf("status code is not 200, got %d\n%s\n", rr.Code, rr.Body.String())
	}

	var record Record
	err := testH.db.First(&record, id).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
