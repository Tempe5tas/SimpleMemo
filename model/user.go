package model

import "time"

type User struct {
	ID        uint
	Name      string
	Password  string
	Email     string
	CreatedAt time.Time
	Memo      []Memo `gorm:"foreignKey:UID"`
}

type TokenBody struct {
	ID   int
	Data string
}

type RegForm struct {
	Name     string `binding:"required,min=3,max=20"`
	Password string `binding:"required,min=8,max=24"`
	Email    string `binding:"required,email,max=30"`
}

type LoginForm struct {
	Name     string `binding:"required,min=3,max=20"`
	Password string `binding:"required,min=8,max=24"`
}
