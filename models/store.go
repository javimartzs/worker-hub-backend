package models

type Store struct {
	ID     string `form:"id" json:"id" gorm:"uniqueIndex"`
	Name   string `form:"name" json:"name" gorm:"not null size:100"`
	City   string `form:"city" json:"city" gorm:"not null size:100"`
	Phone  int    `form:"phone" json:"phone" gorm:"size:25"`
	Status string `form:"status" json:"status" gorm:"not null size:100"`
	UserID string `json:"user_id" gorm:"not null"`
	User   User   `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
