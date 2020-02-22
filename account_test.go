package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	acct = Account{
		Name: "editor",
		Role: Editor,
	}
	password = "my_password"
)

func TestRegister(t *testing.T) {
	body := fmt.Sprintf("username=%s&role=%d&password=%s", acct.Name, acct.Role, password)
	rr := testHandler("POST", "/api/register", body)
	if rr.Code != 200 {
		t.Error("response not success:", rr.Code, rr.Body.String())
		t.FailNow()
	}

	var a Account
	err := testH.db.Where("name = ?", acct.Name).First(&a).Error
	assert.Equal(t, nil, err)
	assert.Equal(t, acct.Role, a.Role)
}

func TestLogin(t *testing.T) {
	body := fmt.Sprintf("username=%s&password=%s", acct.Name, password)
	rr := testHandler("POST", "/api/login", body)
	if rr.Code != 200 {
		t.Error("response not success:", rr.Code, rr.Body.String())
		t.FailNow()
	}

	body = fmt.Sprintf("username=%s&password=%s", acct.Name, "wrong_password")
	rr = testHandler("POST", "/api/login", body)
	assert.Equal(t, 401, rr.Code)

	body = fmt.Sprintf("username=%s&password=%s", "unknown", "wrong_password")
	rr = testHandler("POST", "/api/login", body)
	assert.Equal(t, 404, rr.Code)
}
