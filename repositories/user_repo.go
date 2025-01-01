package repositories

import (
	"github.com/javimartzs/worker-hub-backend/models"
	"github.com/javimartzs/worker-hub-backend/models/dtos"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser - Crea un nuevo usuario
// --------------------------------------------------------------------
func (r *UserRepository) CreateUser(tx *gorm.DB, user *models.User) error {
	if tx != nil {
		return tx.Create(user).Error
	}
	return r.db.Create(user).Error
}

// FindUserByUsername - Busca un usuario por su nombre de usuario
// --------------------------------------------------------------------
func (r *UserRepository) FindUserByUsername(tx *gorm.DB, username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Usuario no encontrado
		}
		return nil, err // Otro error ocurrió
	}
	return &user, nil
}

// GetAllUsers - Obtiene todos los usuarios
// --------------------------------------------------------------------
func (r *UserRepository) GetAllUsers() ([]dtos.UserList, error) {
	var users []dtos.UserList
	err := r.db.Table("users").
		Select("users.id, users.username, users.role").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// DeleteUser - Elimina un usuario
// --------------------------------------------------------------------
func (r *UserRepository) DeleteUser(tx *gorm.DB, userID string) error {
	if tx != nil {
		return tx.Where("id = ?", userID).Delete(&models.User{}).Error
	}
	return r.db.Where("id = ?", userID).Delete(&models.User{}).Error
}
