package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/javimartzs/worker-hub-backend/config"
	"github.com/javimartzs/worker-hub-backend/logger"
	"github.com/javimartzs/worker-hub-backend/models"
	"github.com/javimartzs/worker-hub-backend/utils"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect go to postgres
func ConnectDB() *gorm.DB {

	uri := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s",
		config.Env.DBUser,
		config.Env.DBPass,
		config.Env.DBName,
		config.Env.DBHost,
		config.Env.DBPort,
	)

	DB, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		logger.Logger.Error("Failed to connect PostgresDB", zap.Error(err))
	}

	DB.AutoMigrate(
		&models.User{},
		&models.Store{},
		&models.Worker{},
		&models.Holiday{},
		&models.Timelog{},
		&models.Order{},
		&models.WorkShift{},
	)

	createInitialAdmin(DB)

	logger.Logger.Info("Connected to postgres")
	return DB
}

// Create sudo admin
func createInitialAdmin(db *gorm.DB) {
	// Verificar si ya existe un usuario con rol admin
	var admin models.User
	if err := db.Where("role = ?", "admin").First(&admin).Error; err == nil {
		logger.Logger.Info("Admin user already exists, skipping creation")
		return
	}

	// Crear el usuario admin inicial
	hashedPassword, err := utils.HashPassword(config.Env.DBPass) // Hashear la contraseña inicial
	if err != nil {
		logger.Logger.Error("Failed to hash admin password", zap.Error(err))
	}

	admin = models.User{
		ID:       uuid.New().String(),
		Username: "admin",
		Password: hashedPassword,
		Role:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		logger.Logger.Error("Failed to create initial admin user", zap.Error(err))
	}

	logger.Logger.Info("Initial admin user created successfully")
}
