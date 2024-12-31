package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/javimartzs/worker-hub-backend/logger"
	"github.com/javimartzs/worker-hub-backend/models"
	"github.com/javimartzs/worker-hub-backend/services"
	"go.uber.org/zap"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// Handler para crear un trabajador y su usuario
// --------------------------------------------------------------------
func (h *AdminHandler) CreateWorker(c *gin.Context) {

	// Parseamos el cuerpo de la solicitud
	var worker models.Worker
	if err := c.ShouldBind(&worker); err != nil {
		logger.Logger.Error("CreateWorker: Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Llamamos al servicio para crear el trabajador y su usuario
	if err := h.adminService.CreateWorker(&worker); err != nil {
		logger.Logger.Error("CreateWorker: Worker creation failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Devolvemos una respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Trabajador creado exitosamente",
	})
}

// Handler para eliminar un trabajador y su usuario asociado
// --------------------------------------------------------------------
func (h *AdminHandler) DeleteWorker(c *gin.Context) {
	workerID := c.Param("id")
	if workerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del trabajador requerido"})
		return
	}

	if err := h.adminService.DeleteWorker(workerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trabajador eliminado exitosamente",
	})
}
