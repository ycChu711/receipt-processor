package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the global logger
var Logger = logrus.New()

// InitLogger initializes the logger
func InitLogger() {
	// Set log format
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Set output
	Logger.SetOutput(os.Stdout)

	// Set log level based on environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	Logger.SetLevel(level)
}
