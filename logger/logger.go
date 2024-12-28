package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitializeLogger() {

	// Folder where logd will be stored
	logFolder := "logs"
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		err := os.Mkdir(logFolder, os.ModePerm)
		if err != nil {
			panic("Failed to create log folder" + err.Error())
		}
	}

	// Log File path
	logFilePath := filepath.Join(logFolder, "app.log")

	// Configure the zap logger to log into the file
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to open log file" + err.Error())
	}

	writer := zapcore.AddSync(file)

	// Combine file logging and console logging
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, zapcore.InfoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	Logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Logger)
}

// CleanupLogger ensures all logs are flushed before the application exits
func CleanupLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
