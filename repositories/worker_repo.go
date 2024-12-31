package repositories

import (
	"errors"

	"github.com/javimartzs/worker-hub-backend/models"
	"gorm.io/gorm"
)

type WorkerRepository struct {
	db *gorm.DB
}

func NewWorkerRepository(db *gorm.DB) *WorkerRepository {
	return &WorkerRepository{db: db}
}

// CreateWorker - Crea un nuevo trabajador
// --------------------------------------------------------------------
func (r *WorkerRepository) CreateWorker(tx *gorm.DB, worker *models.Worker) error {
	if tx != nil {
		return tx.Create(worker).Error
	}
	return r.db.Create(worker).Error
}

// FindWorkerByNie - Busca un trabajador por su NIE
// --------------------------------------------------------------------
func (r *WorkerRepository) FindWorkerByNie(nie string) (*models.Worker, error) {
	var worker models.Worker
	err := r.db.Where("nie = ?", nie).First(&worker).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // No se encontró el trabajador
		}
		return nil, err // Otro error ocurrió
	}
	return &worker, nil
}

// FindWorkerByID - Busca un trabajador por su ID
// --------------------------------------------------------------------
func (r *WorkerRepository) FindWorkerByID(workerID string) (*models.Worker, error) {
	var worker models.Worker
	err := r.db.Where("id = ?", workerID).First(&worker).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &worker, nil
}

// DeleteWorker - Elimina un trabajador
// --------------------------------------------------------------------
func (r *WorkerRepository) DeleteWorker(tx *gorm.DB, workerID string) error {
	if tx != nil {
		return tx.Where("id = ?", workerID).Delete(&models.Worker{}).Error
	}
	return r.db.Where("id = ?", workerID).Delete(&models.Worker{}).Error
}
