package models

type Worker struct {
	ID       string  `json:"id" gorm:"primaryKey;uniqueIndex"`
	Name     string  `json:"name" gorm:"size:100"`
	LastName string  `json:"last_name" gorm:"size:100"`
	Email    string  `json:"email" gorm:"size:100"`
	Nie      string  `json:"nie" gorm:"uniqueIndex"`
	Cargo    string  `json:"cargo" gorm:"size:50"`
	Status   string  `json:"status" gorm:"size:25"`
	Prueba   string  `json:"prueba" gorm:"size:25"`
	StoreID  *string `json:"store_id" gorm:"size:50"`
	UserID   string  `json:"user_id" gorm:"not null"`
	Store    Store   `json:"store" gorm:"foreignKey:StoreID;references:ID"`
	User     User    `json:"-" gorm:"foreignKey:UserID;references:ID"`
}
