package utils

import (
	"errors"
	"time"

	"github.com/javimartzs/worker-hub-backend/models"
)

// Función para validar los campos del trabajador
func ValidateWorkerFields(worker *models.Worker) error {
	if worker.Name == "" || worker.LastName == "" {
		return errors.New("el nombre y apellido son obligatorios")
	}
	if worker.Email == "" {
		return errors.New("el Email del trabajador es obligatorio")
	}
	if worker.Cargo == "" {
		return errors.New("el cargo del trabajador es obligatorio")
	}
	if worker.Nie == "" {
		return errors.New("el NIE del trabajador es obligatorio")
	}
	if worker.Status != "Alta" && worker.Status != "Baja" {
		return errors.New("el estado debe ser Alta o Baja")
	}
	if worker.Prueba != "Si" && worker.Prueba != "No" {
		return errors.New("la prueba debe ser Si o No")
	}
	return nil
}

// Funcion para validar los campos de las tiendas
func ValidateStoreFields(store *models.Store) error {
	if store.Name == "" {
		return errors.New("el nombre de la tienda es obligatorio")
	}
	if store.City == "" {
		return errors.New("la ciudad de la tienda es obligatoria")
	}
	if store.Phone < 100000000 || store.Phone > 999999999 {
		return errors.New("el telefono de la tienda debe tener 9 dígitos")
	}
	if store.Status == "" {
		return errors.New("el estado de la tienda es obligatorio")
	}
	return nil
}

// Funcion para validar los campos de las vacaciones
func ValidateHolidaysFields(holiday *models.Holiday) error {

	if holiday.WorkerID == "" {
		return errors.New("el id del trabajador es obligatorio")
	}
	startDate, err := time.Parse("2006-01-02", holiday.StartDate)
	if err != nil {
		return errors.New("la fecha de inicio no tiene el formato YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", holiday.EndDate)
	if err != nil {
		return errors.New("la fecha de fin no tiene el formato YYYY-MM-DD")
	}
	if endDate.Before(startDate) {
		return errors.New("la fecha de fin no puede ser anterior a la fecha de inicio")
	}
	if holiday.Status != "Pendientes" && holiday.Status != "Disfrutadas" {
		return errors.New("el estado de la vacacion es obligatorio")
	}

	return nil
}

// Funcion para validar los campos de los usuarios
func ValidateUserFields(user *models.User) error {
	if user.Username == "" {
		return errors.New("el nombre del usuario es obligatorio")
	}
	if user.Password == "" {
		return errors.New("la contraseña del usuario es obligatoria")
	}
	if user.Role == "" {
		return errors.New("el rol del usuario es obligatorio")
	}
	return nil
}

// Funcion para validar los campos de los registros horarios
func ValidateTimelogFields(timelog *models.Timelog) error {
	if timelog.WorkerID == "" {
		return errors.New("el id del trabajador es obligatorio")
	}
	if timelog.StoreID == "" {
		return errors.New("el id de la tienda es obligatorio")
	}
	if timelog.InOut != "Entrada" && timelog.InOut != "Salida" {
		return errors.New("el estado de entrada/salida es obligatorio")
	}
	if timelog.Timelog == "" {
		return errors.New("el registro horario es obligatorio")
	}
	return nil
}
