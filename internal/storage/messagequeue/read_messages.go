package messagequeue

import (
	"context"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"scribbler/internal/models"
)

func (q *Queue) ReadMessages(ctx context.Context) ([]models.Message, error) {
	partition := int32(0)

	partitionConsumer, err := q.consumer.ConsumePartition(q.cfg.KafkaTopic, partition, sarama.OffsetOldest)
	if err != nil {
		q.logger.Error(
			"failed to consume partition",
			zap.Error(err),
			zap.String("topic", q.cfg.KafkaTopic),
			zap.Int32("partition", partition),
		)
	}

	q.logger.Info("consuming partition",
		zap.String("topic", q.cfg.KafkaTopic),
		zap.Int32("partition", partition),
	)

	// consume messages from Kafka
	for consumedMessage := range partitionConsumer.Messages() {
		q.logger.Info("received message", zap.String("message", string(consumedMessage.Value)))
	}

	// close the partition consumer
	defer func(partitionConsumer sarama.PartitionConsumer) {
		if err = partitionConsumer.Close(); err != nil {
			q.logger.Error("failed to close partition consumer", zap.Error(err))
		}
		q.logger.Info("closed partition consumer")
	}(partitionConsumer)

	return nil, nil
}
