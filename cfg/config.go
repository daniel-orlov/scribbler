package cfg

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// Config is the main configuration of the application.
// It is populated from environment variables and defaults.
type Config struct {
	// LogLevel is the log level to use.
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	// KafkaAddress is the address of the Kafka broker.
	KafkaAddress string `envconfig:"KAFKA_ADDRESS" default:"localhost:9092"`
	// KafkaTopic is the name of the Kafka topic to consume from and produce to.
	KafkaTopic string `envconfig:"KAFKA_TOPIC" default:"tweets"`
	// ESAddress is the address of the Elasticsearch instance.
	ESAddress string `envconfig:"ES_ADDRESS" default:"http://localhost:9200"`
	// ESIndex is the name of the Elasticsearch index to write to.
	ESIndex string `envconfig:"ES_INDEX" default:"tweets"`
}

// NewConfig returns a new Config instance, populated with environment variables and defaults.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, errors.Wrap(err, "processing environment variables")
	}

	return cfg, nil
}
