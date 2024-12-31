package utils

import (
	"errors"

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
