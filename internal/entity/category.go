package entity

import "time"

type Category struct {
	ID          uint
	UserID      uint
	Name        string
	Description string
	CreatedAt time.Time
}