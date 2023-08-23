package utils

import (
	"os"

	"github.com/rs/zerolog"
)

var LogLevels = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
	"panic": zerolog.PanicLevel,
}

func SetupLogger(verbosity string) *zerolog.Logger {
	// Find log level
	logLevel, ok := LogLevels[verbosity]
	if !ok {
		logLevel = zerolog.InfoLevel
	}

	// Return logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(logLevel)
	return &logger
}
