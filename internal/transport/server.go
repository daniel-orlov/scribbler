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

type Server struct {
	logger *zap.Logger
	cfg    *cfg.Config
	engine *gin.Engine
}

// EightMb is 8 megabytes.
const EightMb = 1 << 20

// NewServer creates an HTTP server, adds logging and sets timeouts.
func NewServer(logger *zap.Logger, cfg *cfg.Config, engine *gin.Engine) *Server {
	server := &Server{
		logger: logger,
		cfg:    cfg,
		engine: engine,
	}

	server.init()

	return server
}

func (s *Server) init() *http.Server {
	// Adding logging and recovery middleware.
	// These are the global middlewares, which will be called for every request.
	s.engine.Use(ginzap.Ginzap(s.logger, time.RFC3339, true))
	s.engine.Use(ginzap.RecoveryWithZap(s.logger, true))

	// Setting up the server.
	s.engine.MaxMultipartMemory = EightMb
	s.engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	onPort := fmt.Sprint(":", s.cfg.GinPort)

	s.logger.Info("created an http server to listen on port", zap.String("port", onPort))

	return &http.Server{
		Addr:              onPort,
		Handler:           s.engine,
		ReadTimeout:       time.Duration(s.cfg.GinReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.cfg.GinReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.cfg.GinWriteTimeout) * time.Second,
		MaxHeaderBytes:    EightMb,
	}
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	return s.engine.Run(fmt.Sprint(":", s.cfg.GinPort))
}

func (s *Server) RegisterHandlers(handlers ...Handler) {
	root := s.engine.Group("")

	for i := range handlers {
		handlers[i].Init(root, nil)
	}
}
