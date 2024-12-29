package config

import (
	"os"

	"github.com/javimartzs/worker-hub-backend/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var Env Config

type Config struct {
	DBUser   string
	DBPass   string
	DBName   string
	DBHost   string
	DBPort   string
	JwtKey   string
	Username string
	Password string
}

func LoadEnv() {

	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal("Error loading env file", zap.Error(err))
		return
	}

	Env = Config{
		DBUser:   os.Getenv("DB_USER"),
		DBPass:   os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		DBHost:   os.Getenv("DB_HOST"),
		DBPort:   os.Getenv("DB_PORT"),
		JwtKey:   os.Getenv("JWT_KEY"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	logger.Logger.Info("Env file loaded succesfully")
}
