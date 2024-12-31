package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;uniqueIndex;autoIncrement"`
	Username  string    `json:"username" gorm:"size:100;not null;uniqueIndex"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	Role      string    `json:"role" gorm:"size:50;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}
