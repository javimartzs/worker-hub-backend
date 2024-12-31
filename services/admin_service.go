package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/javimartzs/worker-hub-backend/models"
	"github.com/javimartzs/worker-hub-backend/repositories"
	"github.com/javimartzs/worker-hub-backend/utils"
	"gorm.io/gorm"
)

type AdminService struct {
	userRepo   *repositories.UserRepository
	workerRepo *repositories.WorkerRepository
	db         *gorm.DB
}

func NewAdminService(
	userRepo *repositories.UserRepository,
	workerRepo *repositories.WorkerRepository,
	db *gorm.DB) *AdminService {
	return &AdminService{
		userRepo:   userRepo,
		workerRepo: workerRepo,
		db:         db,
	}
}

// generateUniqueUsername - Genera un nombre de usuario unico para un trabajador nuevo
// -------------------------------------------------------------------
func generateUniqueUsername(tx *gorm.DB, userRepo *repositories.UserRepository, name string) (string, error) {
	baseUsername := strings.ToLower(strings.ReplaceAll(name, " ", ""))
	username := baseUsername
	counter := 1

	for {
		// Verificar la existencia del usuario por su username
		existingUser, err := userRepo.FindUserByUsername(tx, username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("error al verificar el nombre de usuario: %w", err)
		}
		if existingUser == nil {
			break
		}
		username = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++
	}
	return username, nil
}

// CreateWorker - Crea un trabajador nuevo y su usuario asociado
// -------------------------------------------------------------------
func (s *AdminService) CreateWorker(worker *models.Worker) error {

	// Validaciones de los campos del trabajador
	if err := utils.ValidateWorkerFields(worker); err != nil {
		return err
	}

	// Comprobamos que el trabajador no exista ya en la base de datos
	existingWorker, err := s.workerRepo.FindWorkerByNie(worker.Nie)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe el trabajador, esto está bien
		} else {
			return errors.New("error al verificar la existencia del trabajador")
		}
	}
	if existingWorker != nil && existingWorker.ID != "" {
		return errors.New("el trabajador ya existe")
	}

	// Iniciamos la transacción
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("error al iniciar la transaccion")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generamos un nombre de usuario unico para el trabajador
	username, err := generateUniqueUsername(tx, s.userRepo, worker.Name)
	if err != nil {
		tx.Rollback()
		return errors.New("error al generar el nombre de usuario")
	}

	// Generamos la contraseña (PIN) del trabajador
	password := strings.ToLower(worker.Nie[:4])
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		tx.Rollback()
		return errors.New("error al generar la contraseña")
	}

	// Creamos el usuario del trabajador
	user := &models.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: hashedPassword,
		Role:     "worker",
	}

	// Guardamos el usuario en la tabla de usuarios
	if err := s.userRepo.CreateUser(tx, user); err != nil {
		tx.Rollback()
		return errors.New("error al guardar el usuario en la tabla")
	}

	// Generamos el ID del trabajador
	worker.ID = uuid.New().String()
	worker.UserID = user.ID

	// Guardamos el trabajador en la tabla de trabajadores
	if err := s.workerRepo.CreateWorker(tx, worker); err != nil {
		tx.Rollback()
		return errors.New("error al guardar el trabajador en la tabla")
	}

	// Confirmamos la transaccion
	if err := tx.Commit().Error; err != nil {
		return errors.New("error al confirmar la transaccion")
	}

	return nil
}

// DeleteWorker - Elimina un trabajador y su usuario asociado
// --------------------------------------------------------------------
func (s *AdminService) DeleteWorker(workerID string) error {

	// Iniciamos la transaccion
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Buscamos el trabajador por su ID
	worker, err := s.workerRepo.FindWorkerByID(workerID)
	if err != nil {
		tx.Rollback()
		return errors.New("error al buscar el trabajador")
	}
	if worker == nil {
		tx.Rollback()
		return errors.New("el trabajador no existe")
	}

	// Eliminamos el trabajador
	if err := s.workerRepo.DeleteWorker(tx, workerID); err != nil {
		tx.Rollback()
		return errors.New("error al eliminar el trabajador")
	}

	// Eliminamos el usuario asociado al trabajador
	if err := s.userRepo.DeleteUser(tx, worker.UserID); err != nil {
		tx.Rollback()
		return errors.New("error al eliminar el usuario")
	}

	// Confirmamos la transaccion
	if err := tx.Commit().Error; err != nil {
		return errors.New("error al confirmar la transaccion")
	}

	return nil
}
