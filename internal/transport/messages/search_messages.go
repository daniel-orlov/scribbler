package messages

import (
	"context"

	"scribbler/internal/models"
)

func (h *Handler) SearchMessages(ctx context.Context, query string) ([]models.Message, error) {
	return nil, nil
}
