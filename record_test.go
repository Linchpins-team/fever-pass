package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	testDB *gorm.DB
	admin  = Account{
		Role:     Admin,
		Name:     "admin",
		Password: []byte{},
	}
	mockData = []Record{
		Record{
			UserID:      "109123456",
			Pass:        true,
			Temperature: 0,
			Time:        time.Now().Add(-10 * time.Minute),
		},
		Record{
			UserID:      "108234567",
			Pass:        false,
			Temperature: 37.8,
			Time:        time.Now().Add(-5 * time.Minute),
		},
		Record{
			UserID:      "108114256",
			Pass:        true,
			Temperature: 0,
			Time:        today().Add(-1 * time.Hour),
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
	testH.db = testDB
	migrate(testDB)
	insertMockData(testDB)
}

func insertMockData(db *gorm.DB) {
	err := db.Create(&admin).Error
	if err != nil {
		panic(err)
	}
	for _, record := range mockData {
		record.RecordedBy = admin
		err := db.Create(&record).Error
		if err != nil {
			panic(err)
		}
	}
}

func TestNewRecord(t *testing.T) {
	record := Record{
		UserID:      "108222333",
		Pass:        true,
		Temperature: 0,
	}

	body := fmt.Sprintf("user_id=%s&pass=%t&temperature=%f", record.UserID, record.Pass, record.Temperature)
	rr := testHandler(testH.newRecord, "POST", "/api/records", body)
	if rr.Code != 200 {
		t.Errorf("status code is not 200, got %d\n%s\n", rr.Code, rr.Body.String())
	}

	var r Record
	err := testH.db.Where("user_id = ?", record.UserID).First(&r).Error
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, record.Pass, r.Pass)
	assert.Equal(t, record.Temperature, r.Temperature)
}

func adminSession(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, KeyAccount, admin)
	return r.WithContext(ctx)
}

func testHandler(handler func(w http.ResponseWriter, r *http.Request), method, url, body string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, strings.NewReader(body))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}

	req = adminSession(req)
	rr := httptest.NewRecorder()
	http.HandlerFunc(handler).ServeHTTP(rr, req)
	return rr
}

func TestFindRecord(t *testing.T) {
	body := "user_id=" + mockData[0].UserID
	rr := testHandler(testH.findRecord, "POST", "/api/check", body)
	if rr.Code != 200 {
		t.Errorf("status code is not 200, got %d\n%s\n", rr.Code, rr.Body.String())
	}
	t.Log(rr.Body.String())

	body = "user_id=" + mockData[1].UserID
	rr = testHandler(testH.findRecord, "POST", "/api/check", body)
	if rr.Code != 200 {
		t.Errorf("status code is not 200, got %d\n%s\n", rr.Code, rr.Body.String())
	}
	t.Log(rr.Body.String())
}

func TestListRecord(t *testing.T) {
	rr := testHandler(testH.listRecord, "GET", "/api/records", "")
	if rr.Code != 200 {
		t.Errorf("status code is not 200, got %d\n%s\n", rr.Code, rr.Body.String())
	}

	var records []Record
	dec := json.NewDecoder(rr.Body)
	if err := dec.Decode(&records); err != nil {
		t.Error(err)
	}

	for _, record := range records {
		switch record.ID {
		case mockData[2].ID:
			t.Error("expired record should not exist")
		}
	}
	t.Log(records)
}
