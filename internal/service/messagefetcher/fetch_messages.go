package messagefetcher

import (
	"context"

	"scribbler/internal/models"

	"github.com/pkg/errors"
)

func (s *Service) FetchMessages(ctx context.Context, query *models.MessageFilter) error {
	messages, err := s.msgSource.QueryMessages(ctx, query)
	if err != nil {
		return errors.Wrap(err, "fetching messages from source")
	}

	err = s.qWriter.PostMessages(ctx, messages)
	if err != nil {
		return errors.Wrap(err, "posting message to queue")
	}

	return nil
}
