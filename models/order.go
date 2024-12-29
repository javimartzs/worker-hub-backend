package models

type Order struct {
	ID       string `json:"id" gorm:"primaryKey;uniqueIndex"`
	Date     string `json:"date" gorm:"type:date;not null"`
	Status   string `json:"status" gorm:"size:50;not null"`
	Product  string `json:"product" gorm:"size:250;not null"`
	Quantity int    `json:"quantity" gorm:"not null"`
	StoreID  string `json:"store_id" gorm:"not null"`
	Store    Store  `json:"store" gorm:"foreignKey:StoreID;references:ID"`
}
