package main

// Authority define the permission value
type Authority int

const (
	None Authority = iota
	Self
	Group
	All
)

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
