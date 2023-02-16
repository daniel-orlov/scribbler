package main

import (
	"log"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

const (
	kafkaBootstrapServerAddress = "localhost:9092"
	kafkaTopic                  = "first_topic"
)

func main() {
	// create a new Zap logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("failed to create logger", err)
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Fatal("failed to sync logger", err)
		}
	}(logger)

	// create a new Kafka client
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true

	client, err := sarama.NewClient([]string{kafkaBootstrapServerAddress}, kafkaConfig)
	if err != nil {
		logger.Fatal("failed to create client", zap.Error(err))
	}
	logger.Info("created client", zap.String("client", kafkaBootstrapServerAddress))

	// close the client
	defer func(client sarama.Client) {
		if err = client.Close(); err != nil {
			logger.Error("failed to close client", zap.Error(err))
		}
		logger.Info("closed client")
	}(client)

	// create a new Kafka consumer.
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		logger.Fatal("failed to create consumer", zap.Error(err))
	}
	logger.Info("created consumer")

	// close the consumer
	defer func(consumer sarama.Consumer) {
		// close the consumer
		if err = consumer.Close(); err != nil {
			logger.Error("failed to close consumer", zap.Error(err))
		}
		logger.Info("closed consumer")
	}(consumer)

	// create a new Kafka producer.
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		logger.Fatal("failed to create producer", zap.Error(err))
	}
	logger.Info("created producer")

	// close the producer
	defer func(producer sarama.SyncProducer) {
		if err = producer.Close(); err != nil {
			logger.Error("failed to close producer", zap.Error(err))
		}
		logger.Info("closed producer")
	}(producer)

	// create a test message
	messageValue := "Hello, World!"
	message := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.StringEncoder(messageValue),
	}

	// send data to Kafka
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		logger.Error("failed to send message", zap.Error(err))
	}
	logger.Info("sent message to topic",
		zap.String("topic", "first_topic"),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.String("message", messageValue),
	)

	// consume data from Kafka
	partitionConsumer, err := consumer.ConsumePartition("first_topic", partition, sarama.OffsetOldest)
	if err != nil {
		logger.Error("failed to consume partition", zap.Error(err))
	}
	logger.Info("consuming partition", zap.String("topic", "first_topic"), zap.Int32("partition", 0))

	// consume messages from Kafka
	for consumedMessage := range partitionConsumer.Messages() {
		logger.Info("received message", zap.String("message", string(consumedMessage.Value)))
	}
}
