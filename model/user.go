package model

import "time"

type User struct {
	ID        uint
	Name      string `binding:"min=3,max=12"`
	Password  string `binding:"min=8,max=24"`
	Email     string `binding:"email"`
	CreatedAt time.Time
	Memo      []Memo `gorm:"foreignKey:UID" json:"-"`
}

type UserLogin struct {
	Name     string `binding:"min=3,max=12"`
	Password string `binding:"min=8,max=24"`
}
