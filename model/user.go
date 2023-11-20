package model

import "time"

type User struct {
	ID        uint
	Password  string
	Name      string
	Email     string
	CreatedAt time.Time
	Memo      []Memo `gorm:"foreignKey:UID"`
}

type TokenBody struct {
	ID   int
	Data string
}
