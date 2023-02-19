package messagesaver

import (
	"context"

	"scribbler/internal/models"

	"go.uber.org/zap"
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
