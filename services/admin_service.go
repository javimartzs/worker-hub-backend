package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/javimartzs/worker-hub-backend/config"
	"github.com/javimartzs/worker-hub-backend/models"
	"github.com/javimartzs/worker-hub-backend/models/dtos"
	"github.com/javimartzs/worker-hub-backend/repositories"
	"github.com/javimartzs/worker-hub-backend/utils"
	"gorm.io/gorm"
)

type AdminService struct {
	userRepo     *repositories.UserRepository
	workerRepo   *repositories.WorkerRepository
	storeRepo    *repositories.StoreRepository
	holidaysRepo *repositories.HolidaysRepository

	db *gorm.DB
}

func NewAdminService(
	userRepo *repositories.UserRepository,
	workerRepo *repositories.WorkerRepository,
	storeRepo *repositories.StoreRepository,
	holidaysRepo *repositories.HolidaysRepository,
	db *gorm.DB) *AdminService {
	return &AdminService{
		userRepo:     userRepo,
		workerRepo:   workerRepo,
		storeRepo:    storeRepo,
		holidaysRepo: holidaysRepo,
		db:           db,
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

// GetAllWorkers - Obtiene todos los trabajadores
// --------------------------------------------------------------------
func (s *AdminService) GetAllWorkers() ([]models.Worker, error) {
	return s.workerRepo.GetAllWorkers()
}

// UpdateWorker - Actualiza un trabajador
// --------------------------------------------------------------------
func (s *AdminService) UpdateWorker(workerID string, worker *models.Worker) error {

	// Validaciones de los campos del trabajador
	if err := utils.ValidateWorkerFields(worker); err != nil {
		return err
	}

	// Llamamos al repositorio para actualizar el trabajador
	return s.workerRepo.UpdateWorker(workerID, worker)
}

// CreateStore - Crea una nueva tienda y su usuario asociado
// --------------------------------------------------------------------
func (s *AdminService) CreateStore(store *models.Store) error {

	// Validaciones de los campos de la tienda
	if err := utils.ValidateStoreFields(store); err != nil {
		return err
	}

	// Comprobamos que la tienda no exista ya en la base de datos
	existingStore, err := s.storeRepo.FindStoreByName(store.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe la tienda, esto está bien
		} else {
			return errors.New("error al verificar la existencia de la tienda")
		}
	}
	if existingStore != nil && existingStore.ID != "" {
		return errors.New("la tienda ya existe")
	}

	// Iniciamos la transaccion
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("error al iniciar la transaccion")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generamos un nombre de usuario unico para la tienda
	username, err := generateUniqueUsername(tx, s.userRepo, store.Name)
	if err != nil {
		tx.Rollback()
		return errors.New("error al generar el nombre de usuario")
	}

	// Generamos la contraseña de la tienda
	password := config.Env.StorePass
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		tx.Rollback()
		return errors.New("error al generar la contraseña")
	}

	// Creamos el usuario de la tienda
	user := &models.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: hashedPassword,
		Role:     "store",
	}

	// Guardamos el usuario en la tabla de usuarios
	if err := s.userRepo.CreateUser(tx, user); err != nil {
		tx.Rollback()
		return errors.New("error al guardar el trabajador en la tablaa")
	}

	// Generamos el ID de la tienda
	store.ID = uuid.New().String()
	store.UserID = user.ID

	// Guardamos la tienda en la tabla de tiendas
	if err := s.storeRepo.CreateStore(tx, store); err != nil {
		tx.Rollback()
		return errors.New("error al guardar la tienda en la tabla")
	}

	// Confirmamos la transaccion
	if err := tx.Commit().Error; err != nil {
		return errors.New("error al confirmar la transaccion")
	}

	return nil
}

// DeleteStore - Elimina una tienda y su usuario asociado
// --------------------------------------------------------------------
func (s *AdminService) DeleteStore(storeID string) error {

	// Iniciamos la transaccion
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("error al iniciar la transaccion")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Buscamos la tienda por su ID
	store, err := s.storeRepo.FindStoreByID(storeID)
	if err != nil {
		tx.Rollback()
		return errors.New("error al buscar la tienda")
	}
	if store == nil {
		tx.Rollback()
		return errors.New("la tienda no existe")
	}

	// Eliminamos la tienda
	if err := s.storeRepo.DeleteStore(tx, storeID); err != nil {
		tx.Rollback()
		return errors.New("error al eliminar la tienda")
	}

	// Eliminamos el usuario asociado a la tienda
	if err := s.userRepo.DeleteUser(tx, store.UserID); err != nil {
		tx.Rollback()
		return errors.New("error al eliminar el usuario")
	}

	// Confirmamos la transaccion
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("error al confirmar la transaccion")
	}

	return nil
}

// GetAllStores - Obtiene todas las tiendas
// --------------------------------------------------------------------
func (s *AdminService) GetAllStores() ([]models.Store, error) {
	return s.storeRepo.GetAllStores()
}

// UpdateStore - Actualiza una tienda
// --------------------------------------------------------------------
func (s *AdminService) UpdateStore(storeID string, store *models.Store) error {

	// Validaciones de los campos de la tienda
	if err := utils.ValidateStoreFields(store); err != nil {
		return err
	}

	// Llamamos al repositorio para actualizar la tienda
	return s.storeRepo.UpdateStore(storeID, store)
}

// CreateHoliday - Crea una nueva vacacion
// --------------------------------------------------------------------
func (s *AdminService) CreateHoliday(holiday *models.Holiday) error {

	// Validaciones de los campos
	if err := utils.ValidateHolidaysFields(holiday); err != nil {
		return err
	}

	// Llamamos al repositorio para crear las vacaciones de un trabajador
	if err := s.holidaysRepo.CreateHoliday(holiday); err != nil {
		return errors.New("error al crear las vacaciones")
	}

	return nil
}

// GetAllHolidays - Obtiene todas las vacaciones
// --------------------------------------------------------------------
func (s *AdminService) GetAllHolidays() ([]models.Holiday, error) {
	return s.holidaysRepo.GetAllHolidays()
}

// GetHolidaysWithWorker - Obtiene todas las vacaciones con el nombre del trabajador
// --------------------------------------------------------------------
func (s *AdminService) GetHolidaysWithWorker() ([]dtos.HolidayWithWorkerName, error) {
	return s.holidaysRepo.GetHolidaysWithWorker()
}

// DeleteHoliday - Elimina una vacacion
// --------------------------------------------------------------------
func (s *AdminService) DeleteHoliday(holidayID string) error {
	return s.holidaysRepo.DeleteHoliday(holidayID)
}

// UpdateHoliday - Actualiza una vacacion
// --------------------------------------------------------------------
func (s *AdminService) UpdateHoliday(holidayID string, holiday *models.Holiday) error {

	// Validaciones de los campos
	if err := utils.ValidateHolidaysFields(holiday); err != nil {
		return err
	}

	// Llamamos al repositorio para actualizar la vacacion
	return s.holidaysRepo.UpdateHoliday(holidayID, holiday)
}
