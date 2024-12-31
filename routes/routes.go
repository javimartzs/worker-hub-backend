package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/javimartzs/worker-hub-backend/handlers"
)

func SetupRoutes(
	router *gin.Engine,
	adminHandler *handlers.AdminHandler,
) {
	apiGroup := router.Group("/api") // Grupo de rutas para la API
	{
		// Rutas para el administrador
		adminGroup := apiGroup.Group("/admin")
		{
			// Rutas de trabajadores
			adminGroup.POST("/workers/create", adminHandler.CreateWorker)
			adminGroup.POST("/workers/delete/:id", adminHandler.DeleteWorker)
			adminGroup.GET("/workers", adminHandler.GetAllWorkers)
			adminGroup.POST("/workers/update/:id", adminHandler.UpdateWorker)
		}
	}
}
