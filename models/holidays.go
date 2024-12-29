package models

type Holiday struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	WorkerID  string `json:"worker_id" gorm:"not null"`
	StartDate string `json:"start_date" gorm:"type:date;not null"` // Formato YYYY-MM-DD
	EndDate   string `json:"end_date" gorm:"type:date;not null"`
	Status    string `json:"status" gorm:"size:50"`
	Worker    Worker `json:"worker" gorm:"foreignKey:WorkerID;references:ID"`
}
