package main

type Role int

const (
	_ Role = iota
	Admin
	Staff
	Student
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
