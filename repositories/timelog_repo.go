package repositories

import (
	"github.com/javimartzs/worker-hub-backend/models"
	"gorm.io/gorm"
)

type TimelogRepository struct {
	db *gorm.DB
}

func NewTimelogRepository(db *gorm.DB) *TimelogRepository {
	return &TimelogRepository{db: db}
}

// CreateTimelog - Crea un registro horario
// --------------------------------------------------------------------
func (r *TimelogRepository) CreateTimelog(timelog *models.Timelog) error {
	return r.db.Create(timelog).Error
}
