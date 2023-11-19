package model

import "time"

type Memo struct {
	ID     uint
	Title  string
	Body   string
	Time   time.Time
	Status bool
}
