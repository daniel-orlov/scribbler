package transport

import (
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"scribbler/cfg"
)

// EightMb is 8 megabytes.
const EightMb = 1 << 20

// NewServer creates an HTTP server, adds logging and sets timeouts.
func NewServer(logger *zap.Logger, cfg *cfg.Config, engine *gin.Engine) *http.Server {
	// Adding logging and recovery middleware.
	// These are the global middlewares, which will be called for every request.
	engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger, true))

	onPort := fmt.Sprint(":", cfg.GinPort)

	logger.Info("created an http server to listen on port", zap.String("port", onPort))

	return &http.Server{
		Addr:              onPort,
		Handler:           engine,
		ReadTimeout:       time.Duration(cfg.GinReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.GinReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.GinWriteTimeout) * time.Second,
		MaxHeaderBytes:    EightMb,
	}
}
