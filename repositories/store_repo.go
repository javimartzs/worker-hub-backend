package repositories

import (
	"errors"

	"github.com/javimartzs/worker-hub-backend/models"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

// CreateStore - Crea una nueva tienda
// --------------------------------------------------------------------
func (r *StoreRepository) CreateStore(tx *gorm.DB, store *models.Store) error {
	if tx != nil {
		return tx.Create(store).Error
	}
	return r.db.Create(store).Error
}

// FindStoreByName - Busca una tienda por su nombre
// --------------------------------------------------------------------
func (r *StoreRepository) FindStoreByName(name string) (*models.Store, error) {
	var store models.Store
	err := r.db.Where("name = ?", name).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &store, nil
}

// FindStoreByID - Busca una tienda por su ID
// --------------------------------------------------------------------
func (r *StoreRepository) FindStoreByID(storeID string) (*models.Store, error) {
	var store models.Store
	err := r.db.Where("id = ?", storeID).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &store, nil
}

// DeleteStore - Elimina una tienda
// --------------------------------------------------------------------
func (r *StoreRepository) DeleteStore(tx *gorm.DB, storeID string) error {
	if tx != nil {
		return tx.Where("id = ?", storeID).Delete(&models.Store{}).Error
	}
	return r.db.Where("id = ?", storeID).Delete(&models.Store{}).Error
}

// GetAllStores - Obtiene todas las tiendas
// --------------------------------------------------------------------
func (r *StoreRepository) GetAllStores() ([]models.Store, error) {
	var stores []models.Store
	if err := r.db.Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

// UpdateStore - Actualiza una tienda
// --------------------------------------------------------------------
func (r *StoreRepository) UpdateStore(storeID string, store *models.Store) error {
	return r.db.Model(&models.Store{}).Where("id = ?", storeID).Updates(store).Error
}
