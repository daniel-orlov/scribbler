package healthchecker

import "go.uber.org/zap"

// KafkaChecker is an interface for checking the health of the Kafka cluster.
type KafkaChecker interface {
	Check() error
}

// ESChecker is an interface for checking the health of the ElasticSearch cluster.
type ESChecker interface {
	Check() error
}

// Service is the Health Checker service.
type Service struct {
	logger *zap.Logger
	kafka  KafkaChecker
	es     ESChecker
}

// NewService creates a new instance of the Health Checker service.
func NewService(logger *zap.Logger, kafka KafkaChecker, es ESChecker) *Service {
	return &Service{logger: logger, kafka: kafka, es: es}
}
