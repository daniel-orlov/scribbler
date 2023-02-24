package logging

import (
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerSingleton *zap.Logger
var once sync.Once

// Logger returns a concurrence-safe singleton logger.
func Logger() *zap.Logger {
	once.Do(func() {
		loggerSingleton = initLogger()
	})

	return loggerSingleton
}

func initLogger() *zap.Logger {
	logger, err := zap.NewDevelopment(
		zap.IncreaseLevel(zapcore.InfoLevel),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic(errors.Wrap(err, "creating logger"))
	}

	return logger
}
