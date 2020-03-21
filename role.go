package main

import (
	"strings"
)

type Role int

const (
	_ Role = iota
	RoleAdmin
	RoleStaff
	RoleStudent
)

const (
	KeyRole             = "role"
	KeyRecordAuthority  = "record_authority"
	KeyAccountAuthority = "account_authority"
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "管理員"

	case RoleStaff:
		return "教職員"

	case RoleStudent:
		return "學生"
	}
	return "未知"
}

func parseAuthority(s string) AuthorityLevel {
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
		return RoleAdmin

	case "staff", "teacher":
		return RoleStaff

	case "student":
		return RoleStudent

	default:
		return 0
	}
}
