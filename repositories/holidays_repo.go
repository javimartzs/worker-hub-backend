package repositories

import (
	"github.com/javimartzs/worker-hub-backend/models"
	"github.com/javimartzs/worker-hub-backend/models/dtos"
	"gorm.io/gorm"
)

type HolidaysRepository struct {
	db *gorm.DB
}

func NewHolidaysRepository(db *gorm.DB) *HolidaysRepository {
	return &HolidaysRepository{
		db: db,
	}
}

// CreateHoliday - Crea una nueva vacacion
// --------------------------------------------------------------------
func (r *HolidaysRepository) CreateHoliday(holiday *models.Holiday) error {
	return r.db.Create(holiday).Error
}

// GetHolidayByID - Obtiene una vacacion por su ID
// --------------------------------------------------------------------
func (r *HolidaysRepository) GetHolidayByID(holidayID string) (*models.Holiday, error) {
	var holiday models.Holiday
	if err := r.db.First(&holiday, holidayID).Error; err != nil {
		return nil, err
	}
	return &holiday, nil
}

// GetAllHolidays - Obtiene todas las vacaciones
// --------------------------------------------------------------------
func (r *HolidaysRepository) GetAllHolidays() ([]models.Holiday, error) {
	var holidays []models.Holiday
	if err := r.db.Find(&holidays).Error; err != nil {
		return nil, err
	}
	return holidays, nil
}

// GetHolidaysWithWorker - Obtiene todas las vacaciones con el nombre del trabajador
// --------------------------------------------------------------------

func (r *HolidaysRepository) GetHolidaysWithWorker() ([]dtos.HolidayWithWorkerName, error) {
	var holidays []dtos.HolidayWithWorkerName

	err := r.db.Table("holidays").
		Select("holidays.*, workers.name as worker_name, workers.last_name as worker_last_name").
		Joins("left join workers on workers.id = holidays.worker_id").
		Find(&holidays).Error

	if err != nil {
		return nil, err
	}

	return holidays, nil
}

// DeleteHoliday - Elimina una vacacion
// --------------------------------------------------------------------
func (r *HolidaysRepository) DeleteHoliday(holidayID string) error {
	return r.db.Delete(&models.Holiday{}, holidayID).Error
}

// UpdateHoliday - Actualiza una vacacion
// --------------------------------------------------------------------
func (r *HolidaysRepository) UpdateHoliday(holidayID string, holiday *models.Holiday) error {
	return r.db.Model(&models.Holiday{}).Where("id = ?", holidayID).Updates(holiday).Error
}
