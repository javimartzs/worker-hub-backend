package dtos

type HolidayWithWorkerName struct {
	ID             int    `json:"id"`
	WorkerName     string `json:"worker_name"`
	WorkerLastName string `json:"worker_last_name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	Status         string `json:"status"`
}
