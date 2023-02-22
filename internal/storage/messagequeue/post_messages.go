package messagequeue

import (
	"context"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"scribbler/internal/models"
)

func (q *Queue) PostMessages(ctx context.Context, messages []models.Message) error {
	_ = ctx

	producerMessages := make([]*sarama.ProducerMessage, len(messages))
	for i := range messages {
		producerMessages[i] = &sarama.ProducerMessage{
			Topic: q.cfg.KafkaTopic,
			Value: sarama.StringEncoder(messages[i].Text),
		}

	}

	// send data to Kafka
	err := q.producer.SendMessages(producerMessages)
	if err != nil {
		q.logger.Error("failed to send messages", zap.Error(err))
	}

	q.logger.Info("sent messages to topic",
		zap.String("topic", q.cfg.KafkaTopic),
		zap.Int("count", len(messages)),
	)

	return nil
}
