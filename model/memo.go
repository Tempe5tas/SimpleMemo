package model

import "time"

type Memo struct {
	// Memo ID and foreign user key
	ID   uint
	UID  uint `json:"-"`
	User User `gorm:"foreignKey:UID" json:"-"`
	// Memo body
	Title   string `gorm:"not null"`
	Time    time.Time
	Status  bool `gorm:"default:false"`
	Content string
}

type MemoCreate struct {
	Title   string `binding:"required"`
	Time    string `binding:"required"`
	Status  bool   `gorm:"default:false"`
	Content string
}
