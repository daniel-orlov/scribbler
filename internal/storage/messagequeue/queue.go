package messagequeue

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"scribbler/cfg"
)

type Queue struct {
	cfg      *cfg.Config
	logger   *zap.Logger
	consumer sarama.Consumer
	producer sarama.SyncProducer
}

func NewQueue(cfg *cfg.Config, logger *zap.Logger, consumer sarama.Consumer, producer sarama.SyncProducer) *Queue {
	return &Queue{cfg: cfg, logger: logger, consumer: consumer, producer: producer}
}
