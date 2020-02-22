package main

type Role uint32

const (
	Unknown Role = iota
	Admin
	Editor
	User
)

type Account struct {
	ID       uint32
	Name     string
	Password []byte

	Role Role
}
