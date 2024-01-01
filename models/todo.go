package models

import "time"

type Todo struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Contents  string    `gorm:"not null" json:"contents"`
}
