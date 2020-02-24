package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewURL(t *testing.T) {
	body := fmt.Sprintf("valid_time=%d&max_usage=%d", 2, 100)
	rr := testHandler("POST", "/api/url", body)

	assert.Equal(t, 200, rr.Code)
	data := make(map[string]string)
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &data))

	t.Log(data)
}
