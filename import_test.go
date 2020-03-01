package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	file := strings.NewReader(
		`"email","name","password","class"
"s10512345@st.fcjh.tc.edu.tw","Justin","j_password","207"
"s10642236@st.fcjh.tc.edu.tw","Kevin","k_pwd","109"
"s10443256@st.fcjs.tc.edu.tw","Anna","elsa","303"`,
	)
	_, err := importAccounts(testH.db, file, Student)
	assert.Equal(t, nil, err)
	var students []Account
	testH.db.Preload("Class").Find(&students)
	t.Log(students)

	file = strings.NewReader(
		`"email","name","password","class"
		"test@email","someone","password"
		`,
	)
	_, err = importAccounts(testH.db, file, Student)
	assert.NotNil(t, err, "err should not be nil")
}
