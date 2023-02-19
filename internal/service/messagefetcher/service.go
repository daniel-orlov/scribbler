package messagefetcher

import (
	"context"

	"scribbler/internal/models"

	"go.uber.org/zap"
)

type MessageSource interface {
	QueryMessages(ctx context.Context, filter *models.MessageFilter) ([]models.Message, error)
}

type QueueWriter interface {
	PostMessages(ctx context.Context, messages []models.Message) error
}

type Service struct {
	log       *zap.Logger
	msgSource MessageSource
	qWriter   QueueWriter
}

// NewService creates a new instance of the Message Fetcher service.
func NewService(log *zap.Logger, msgSource MessageSource, qWriter QueueWriter) *Service {
	return &Service{log: log, msgSource: msgSource, qWriter: qWriter}
}
