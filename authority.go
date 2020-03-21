package main

import (
	"strings"
)

type Authority struct {
	Role    Role
	Record  AuthorityLevel
	Account AuthorityLevel
}

type Role int

// AuthorityLevel define the permission value
type AuthorityLevel int

var (
	Admin = Authority{
		Role:    RoleAdmin,
		Record:  All,
		Account: All,
	}

	Nurse = Authority{
		Role:    RoleStaff,
		Record:  All,
		Account: Self,
	}

	Staff = Authority{
		Role:    RoleStaff,
		Record:  Self,
		Account: Self,
	}

	Tutor = Authority{
		Role:    RoleStaff,
		Record:  Group,
		Account: Group,
	}

	Hygiene = Authority{
		Role:    RoleStudent,
		Record:  Group,
		Account: Group,
	}

	Student = Authority{
		Role:    RoleStudent,
		Record:  Self,
		Account: Self,
	}

	Unknown = Authority{}

	Authorities = [...]Authority{
		Admin,
		Nurse,
		Staff,
		Tutor,
		Hygiene,
		Student,
	}
)

const (
	_ Role = iota
	RoleAdmin
	RoleStaff
	RoleStudent
)

var (
	Roles = [...]Role{RoleAdmin, RoleStaff, RoleStudent}
)

const (
	None AuthorityLevel = iota
	Self
	Group
	All
)

var (
	Levels = [...]AuthorityLevel{None, Self, Group, All}
)

const (
	KeyRole             = "role"
	KeyRecordAuthority  = "record_authority"
	KeyAccountAuthority = "account_authority"
)

func (a Authority) String() string {
	switch a {
	case Admin:
		return "管理員"

	case Nurse:
		return "護理師"

	case Staff:
		return "教職員"

	case Tutor:
		return "導師"

	case Hygiene:
		return "衛生股長"

	case Student:
		return "學生"

	case Unknown:
		return ""

	default:
		return a.Role.String()
	}
}

func (a Authority) Key() string {
	switch a {
	case Admin:
		return "admin"

	case Nurse:
		return "nurse"

	case Staff:
		return "staff"

	case Tutor:
		return "tutor"

	case Hygiene:
		return "hygiene"

	case Student:
		return "student"

	default:
		return ""
	}
}

func (a AuthorityLevel) String() string {
	switch a {
	case Self:
		return "個人"

	case Group:
		return "班級"

	case All:
		return "全校"

	default:
		return "無"
	}
}

func (a AuthorityLevel) Key() string {
	switch a {
	case Self:
		return "self"

	case Group:
		return "group"

	case All:
		return "all"

	default:
		return ""
	}
}

func parseAuthority(s string) Authority {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	for _, authority := range Authorities {
		if authority.Key() == s {
			return authority
		}
	}
	return Unknown
}

// recordPermission return whether A can access B's records
func recordPermission(a, b Account) bool {
	switch a.Authority.Record {
	case All:
		return true

	case Group:
		return a.ClassID == b.ClassID

	case Self:
		return a.ID == b.ID
	}
	return false
}

func accountPermission(a, b Account) bool {
	switch a.Authority.Account {
	case All:
		return true

	case Group:
		return a.ClassID == b.ClassID

	case Self:
		return a.ID == b.ID
	}
	return false
}

func (a Authority) bigger() AuthorityLevel {
	level := a.Record
	if a.Account > level {
		level = a.Account
	}
	return level
}

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

func (r Role) Key() string {
	switch r {
	case RoleAdmin:
		return "admin"

	case RoleStaff:
		return "staff"

	case RoleStudent:
		return "student"
	}
	return ""
}

func parseAuthorityLevel(s string) AuthorityLevel {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	for _, level := range Levels {
		if level.Key() == s {
			return level
		}
	}
	return None
}

func parseRole(s string) Role {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	for _, role := range Roles {
		if role.Key() == s {
			return role
		}
	}
	return 0
}
