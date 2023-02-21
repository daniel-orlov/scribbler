package messages

import "github.com/gin-gonic/gin"

const (
	RootEndpoint = "/messages"
)

// Re-run this if you change the interface below.
//go:generate mockgen -source ./handler.go -package=mock -destination=./mock/usecase_mock.go

type UseCase interface {
}

type Handler struct {
	uc UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Init(root *gin.RouterGroup, middlewares []func() gin.HandlerFunc) {
	messages := root.Group(RootEndpoint)

	for _, middleware := range middlewares {
		messages.Use(middleware())
	}

	messages.GET("/search", h.SearchMessages)
}
