package main

// Authority define the permission value
type Authority int

const (
	None Authority = iota
	Self
	Group
	All
)

func (a Authority) String() string {
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
	switch a.RecordAuthority {
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
	switch a.AccountAuthority {
	case All:
		return true

	case Group:
		return a.ClassID == b.ClassID

	case Self:
		return a.ID == b.ID
	}
	return false
}
