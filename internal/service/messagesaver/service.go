package messagesaver

import (
	"context"

	"go.uber.org/zap"

	"scribbler/internal/models"
)

type MessageStore interface {
	StoreMessages(ctx context.Context, messages []models.Message) error
}

type QueueReader interface {
	ReadMessages(ctx context.Context) ([]models.Message, error)
}

type Service struct {
	logger   *zap.Logger
	msgStore MessageStore
	qReader  QueueReader
}

func NewService(logger *zap.Logger, msgStore MessageStore, qReader QueueReader) *Service {
	return &Service{logger: logger, msgStore: msgStore, qReader: qReader}
}
