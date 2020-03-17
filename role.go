package main

import (
	"strings"
)

type Role int

const (
	_ Role = iota
	Admin
	Staff
	Student
)

const (
	KeyRole             = "role"
	KeyRecordAuthority  = "record_authority"
	KeyAccountAuthority = "account_authority"
)

func (r Role) String() string {
	switch r {
	case Admin:
		return "管理員"

	case Staff:
		return "教職員"

	case Student:
		return "學生"
	}
	return "未知"
}

func (a *Account) defaultAuthority(role Role) {
	switch role {
	case Admin:
		a.RecordAuthority = All
		a.AccountAuthority = All

	case Staff:
		a.RecordAuthority = Self
		a.AccountAuthority = Self

	case Student:
		a.RecordAuthority = Self
		a.AccountAuthority = Self

	default:
		a.RecordAuthority = None
		a.AccountAuthority = None
	}
}

func parseAuthority(s string) Authority {
	switch strings.ToLower(s) {
	case "self":
		return Self

	case "group":
		return Group

	case "all":
		return All

	default:
		return None
	}
}

func parseRole(s string) Role {
	switch strings.ToLower(s) {
	case "admin":
		return Admin

	case "staff", "teacher":
		return Staff

	case "student":
		return Student

	default:
		return 0
	}
}
