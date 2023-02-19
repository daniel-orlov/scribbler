package twit

import (
	"context"
	"math/rand"

	"scribbler/cfg"
	"scribbler/internal/models"

	"go.uber.org/zap"
)

type Client struct {
	cfg *cfg.Config
	log *zap.Logger
}

func NewClient(log *zap.Logger) *Client {
	return &Client{log: log}
}

// QueryMessages returns the messages that match the query.
// Note: This implementation uses fake tweet generator for now.
func (c *Client) QueryMessages(_ context.Context, _ *models.MessageFilter) ([]models.Message, error) {
	fakeTweets := generateFakeTweets(rand.Intn(c.cfg.MaxTweetsPerRequest))

	c.log.Info("successfully fetched messages", zap.Int("count", len(fakeTweets)))

	return fakeTweets, nil
}

func generateFakeTweets(count int) []models.Message {
	tweets := make([]models.Message, count, 0)

	for i := 0; i < count; i++ {
		tweets = append(tweets, models.NewFakeMessage())
	}

	return tweets
}
