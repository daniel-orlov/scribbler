package transport

import (
	"github.com/gin-gonic/gin"

	"scribbler/cfg"
)

// HandlerInitializer is an interface for abstracting handler initialization from a router.
type HandlerInitializer interface {
	Init(rg *gin.RouterGroup, middlewares []func() gin.HandlerFunc)
}

// Router holds all the handlers.
// Init method should be called for its initialization before use.
type Router struct {
	handlers []HandlerInitializer
	// isInitialized checks if router is initialized.
	isInitialized bool
}

func NewRouter(handlers []HandlerInitializer) *Router {
	return &Router{handlers: handlers}
}

// Init initializes router by adding all handlers and applying the middlewares.
func (r *Router) Init(cfg *cfg.Config, middlewares []func() gin.HandlerFunc) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	ginEngine := gin.New()

	// Create a root group
	root := ginEngine.Group("")

	// Adding handlers and initializing them
	for i := range r.handlers {
		r.handlers[i].Init(root, middlewares)
	}

	r.isInitialized = true

	return ginEngine
}
