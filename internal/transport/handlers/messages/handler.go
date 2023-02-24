package messages

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	RootEndpoint = "/messages"
)

// Re-run this if you change the interface below.
//go:generate mockgen -source ./handler.go -package=mock -destination=./mock/usecase_mock.go

type UseCase interface {
}

type Handler struct {
	logger *zap.Logger
	uc     UseCase
}

func NewHandler(logger *zap.Logger, uc UseCase) *Handler {
	return &Handler{logger: logger, uc: uc}
}

func (h *Handler) Init(root *gin.RouterGroup, middlewares []func() gin.HandlerFunc) {
	messages := root.Group(RootEndpoint)

	for _, middleware := range middlewares {
		messages.Use(middleware())
	}

	messages.GET("/search", h.SearchMessages)
}
