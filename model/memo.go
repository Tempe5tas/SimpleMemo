package model

import "time"

type Memo struct {
	ID      uint
	UID     uint   `gorm:"not null""`
	User    User   `gorm:"foreignKey:UID"`
	Title   string `gorm:"not null"`
	Status  bool   `gorm:"default:false"`
	Content string
	Time    time.Time
}
