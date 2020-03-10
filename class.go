package main

type Class struct {
	ID   uint32
	Name string
}

func (c Class) String() string {
	return c.Name
}
