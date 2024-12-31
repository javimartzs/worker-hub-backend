package repositories

import (
	"github.com/javimartzs/worker-hub-backend/models"
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
