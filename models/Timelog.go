package models

type Timelog struct {
	ID       string `json:"id" gorm:"primaryKey;uniqueIndex"`
	StoreID  string `json:"store_id" gorm:"not null"`
	WorkerID string `json:"worker_id" gorm:"not null"`
	InOut    string `json:"in_out" gorm:"not null"`
	Timelog  string `json:"timelog" gorm:"type:datetime;autoCreateTime"`
	Store    Store  `json:"-" gorm:"foreignKey:StoreID;references:ID"`
	Worker   Worker `json:"-" gorm:"foreignKey:WorkerID;references:ID"`
}
