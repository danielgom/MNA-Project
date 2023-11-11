// Package logger is the core logger which is going to be used by our application
package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Initialise starts the global logger instance.
func Initialise() {
	logConfig := zap.NewProductionConfig()
	logConfig.Encoding = "console"
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatalln("could not start logger", err.Error())
	}

	zap.ReplaceGlobals(logger)
}

// Info wraps info and zap fields.
func Info(msg string, tags ...zap.Field) {
	zap.L().Info(msg, tags...)
}

// Error wraps error and adds a zap field for the error.
func Error(msg, err string, tags ...zap.Field) {
	if err != "" {
		tags = append(tags, zap.String("error", err))
	}

	zap.L().Error(msg, tags...)
}
