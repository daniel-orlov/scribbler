package messagestore

import (
	"context"

	"go.uber.org/zap"

	"scribbler/internal/models"
)

func (s *Store) StoreMessages(ctx context.Context, messages []models.Message) error {
	s.logger.Info("storing messages", zap.Int("count", len(messages)))
	return nil
}
