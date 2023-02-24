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

// Handler is an interface for abstracting handler initialization from a server.
type Handler interface {
	Init(rg *gin.RouterGroup, middlewares []func() gin.HandlerFunc)
}

type Server struct {
	logger *zap.Logger
	cfg    *cfg.Config
	engine *gin.Engine
}

// EightMb is 8 megabytes.
const EightMb = 1 << 20

// NewServer creates an HTTP server, adds logging and sets timeouts.
func NewServer(logger *zap.Logger, cfg *cfg.Config) *Server {
	server := &Server{
		logger: logger,
		cfg:    cfg,
		engine: gin.New(),
	}

	server.init()

	return server
}

func (s *Server) init() *http.Server {
	gin.SetMode(s.cfg.GinMode)

	// Adding logging and recovery middleware.
	// These are the global middlewares, which will be called for every request.
	s.engine.Use(ginzap.Ginzap(s.logger, time.RFC3339, true))
	s.engine.Use(ginzap.RecoveryWithZap(s.logger, true))

	// Setting up the server defaults.
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

func (s *Server) RegisterHandlers(handlers map[Handler][]func() gin.HandlerFunc) {
	v1 := s.engine.Group("api/v1")

	for handler, mids := range handlers {
		handler.Init(v1, mids)
	}
}
