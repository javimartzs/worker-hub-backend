package main

import (
	"github.com/gin-gonic/gin"
	"github.com/javimartzs/worker-hub-backend/config"
	"github.com/javimartzs/worker-hub-backend/database"
	"github.com/javimartzs/worker-hub-backend/handlers"
	"github.com/javimartzs/worker-hub-backend/logger"
	"github.com/javimartzs/worker-hub-backend/repositories"
	"github.com/javimartzs/worker-hub-backend/routes"
	"github.com/javimartzs/worker-hub-backend/services"
)

func main() {
	// Inicializamos el logger
	logger.InitializeLogger()
	defer logger.CleanupLogger()

	// Iniciamos la configuracion de la aplicacion
	config.LoadEnv()
	db := database.ConnectDB()

	// Iniciamos las instancias de los repositorios
	userRepo := repositories.NewUserRepository(db)
	workerRepo := repositories.NewWorkerRepository(db)
	storeRepo := repositories.NewStoreRepository(db)
	holidaysRepo := repositories.NewHolidaysRepository(db)
	timelogRepo := repositories.NewTimelogRepository(db)

	// Iniciamos las instancias de los servicios
	adminService := services.NewAdminService(userRepo, workerRepo, storeRepo, holidaysRepo, timelogRepo, db)

	// Iniciamos las instancias de los handlers
	adminHandler := handlers.NewAdminHandler(adminService)

	// Iniciamos el router de Gin
	router := gin.Default()

	// Configuramos las rutas
	routes.SetupRoutes(router, adminHandler)

	// Iniciamos el servidor
	router.Run(":8080")
	logger.Logger.Info("Starting server on port :8080")
}
