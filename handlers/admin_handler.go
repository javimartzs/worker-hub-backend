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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID del trabajador requerido",
		})
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

// Handler para obtener todos los trabajadores
// --------------------------------------------------------------------
func (h *AdminHandler) GetAllWorkers(c *gin.Context) {
	workers, err := h.adminService.GetAllWorkers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudieron obtener los trabajadores", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workers": workers})
}

// Handler para actualizar datos de un trabajador
// --------------------------------------------------------------------
func (h *AdminHandler) UpdateWorker(c *gin.Context) {
	workerID := c.Param("id")
	if workerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID del trabajador requerido"})
		return
	}

	var worker models.Worker
	if err := c.ShouldBind(&worker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := h.adminService.UpdateWorker(workerID, &worker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al actualizar el trabajador",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trabajador actualizado correctamente",
		"worker":  worker,
	})
}

// Handler para crear una tienda y su usuario asociado
// --------------------------------------------------------------------
func (h *AdminHandler) CreateStore(c *gin.Context) {

	// Parseamos el cuerpo de la solicitud
	var store models.Store
	if err := c.ShouldBind(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Llamamos al servicio para crear la tienda
	if err := h.adminService.CreateStore(&store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Devolvemos una respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Tienda creada exitosamente",
	})
}

// Handler para eliminar una tienda y su usuario asociado
// --------------------------------------------------------------------
func (h *AdminHandler) DeleteStore(c *gin.Context) {

	storeID := c.Param("id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de la tienda requerido",
		})
		return
	}

	// Llamamos al servicio para eliminar la tienda
	if err := h.adminService.DeleteStore(storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Devolvemos una respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Tienda eliminada exitosamente",
	})
}

// Handler para obtener todas las tiendas
// --------------------------------------------------------------------
func (h *AdminHandler) GetAllStores(c *gin.Context) {
	stores, err := h.adminService.GetAllStores()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudieron obtener las tiendas", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stores": stores,
	})
}

// Handler para actualizar los datos de una tienda
// --------------------------------------------------------------------
func (h *AdminHandler) UpdateStore(c *gin.Context) {
	storeID := c.Param("id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de la tienda requerido"})
		return
	}

	var store models.Store
	if err := c.ShouldBind(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := h.adminService.UpdateStore(storeID, &store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al actualizar la tienda",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tienda actualizada correctamente",
		"store":   store,
	})
}

// Handler para crear una nueva vacacion
// --------------------------------------------------------------------
func (h *AdminHandler) CreateHoliday(c *gin.Context) {

	// Parseamos el cuerpo de la respuesta
	var holiday models.Holiday
	if err := c.ShouldBind(&holiday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Llamamos al servicio para crear la vacacion
	if err := h.adminService.CreateHoliday(&holiday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Devolvemos una respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Vacacion creada correctamente",
	})
}

// Handler para obtener todas las vacaciones
// --------------------------------------------------------------------
func (h *AdminHandler) GetAllHolidays(c *gin.Context) {
	holidays, err := h.adminService.GetAllHolidays()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudieron obtener las vacaciones",
		})
		return
	}

	// Devolvemos una respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"holidays": holidays,
	})
}

// Handler para obtener todas las vacaciones con el nombre del trabajador
// --------------------------------------------------------------------
func (h *AdminHandler) GetHolidaysWithWorker(c *gin.Context) {
	holidays, err := h.adminService.GetHolidaysWithWorker()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudieron obtener las vacaciones",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"holidays": holidays,
	})
}

// Handler para eliminar una vacacion
// --------------------------------------------------------------------
func (h *AdminHandler) DeleteHoliday(c *gin.Context) {
	holidayID := c.Param("id")
	if holidayID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de la vacacion requerido",
		})
		return
	}

	if err := h.adminService.DeleteHoliday(holidayID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vacacion eliminada correctamente",
	})
}

// Handler para actualizar los datos de una vacacion
// --------------------------------------------------------------------
func (h *AdminHandler) UpdateHoliday(c *gin.Context) {
	holidayID := c.Param("id")
	if holidayID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de la vacacion requerido",
		})
		return
	}

	var holiday models.Holiday
	if err := c.ShouldBind(&holiday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if err := h.adminService.UpdateHoliday(holidayID, &holiday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vacacion actualizada correctamente",
		"holiday": holiday,
	})
}
