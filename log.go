package fuego

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger exposes basic logging methods.
// This is provided as a sheer convenience.
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
}

// init initialises zap's global loggers (i.e. zap.L() and zap.S().
// Log level is retrieved from environment variable LOB_LEVEL, defaulting to "info" otherwise.
// To activate this facility, add a blank import to this package.
func init() {
	level := os.Getenv("FUEGO_LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	fmt.Println("LEVEL=", level)

	initialiseZapGlobals(level)
}

// Option is a function to provide the Logger creation with extra initialisation options.
type Option func(*zap.Config)

func initialiseZapGlobals(level string) {
	zapLevel := zap.InfoLevel
	_ = zapLevel.UnmarshalText([]byte(level))

	lc := zap.NewProductionConfig()
	lc.Level = zap.NewAtomicLevelAt(zapLevel)
	lc.OutputPaths = []string{"stdout"}
	lc.ErrorOutputPaths = []string{"stderr"}
	lc.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	logger, err := lc.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
}
