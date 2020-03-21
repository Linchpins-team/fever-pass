package main

type Authority struct {
	Role    Role
	Record  AuthorityLevel
	Account AuthorityLevel
}

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
)

// AuthorityLevel define the permission value
type AuthorityLevel int

const (
	None AuthorityLevel = iota
	Self
	Group
	All
)

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
