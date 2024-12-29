package models

type WorkShift struct {
	ID            int    `json:"id" gorm:"primaryKey;autoIncrement"`              // Clave primaria con autoincremento
	WorkDate      string `json:"work_date" gorm:"type:date;not null"`             // Fecha del trabajo (YYYY-MM-DD)
	StartInterval string `json:"start_interval" gorm:"type:time;not null"`        // Intervalo de entrada (HH:MM:SS)
	EndInterval   string `json:"end_interval" gorm:"type:time;not null"`          // Intervalo de salida (HH:MM:SS)
	Store         string `json:"store" gorm:"size:50"`                            // Clave foránea opcional hacia la tienda
	CellColor     string `json:"cell_color" gorm:"size:7"`                        // Color de celda en formato hexadecimal (#RRGGBB)
	WorkerID      string `json:"worker_id" gorm:"not null"`                       // Clave foránea obligatoria hacia Worker
	Worker        Worker `json:"worker" gorm:"foreignKey:WorkerID;references:ID"` // Relación con Worker
}
