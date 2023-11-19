package model

type User struct {
	ID       uint
	Password string
	Name     string
	Email    string
}

type TokenBody struct {
	ID   int
	Data string
}
