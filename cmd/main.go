package main

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"scribbler/cfg"
	"scribbler/internal/adapters/twit"
	"scribbler/internal/service/messagefetcher"
	"scribbler/internal/service/messagesaver"
	"scribbler/internal/storage/messagequeue"
	"scribbler/internal/storage/messagestore"
	"scribbler/internal/transport"
	"scribbler/internal/transport/handlers/messages"
	"scribbler/internal/transport/middlewares/loggermw"
	"scribbler/pkg/logging"
)

func main() {
	//------------------------------------------------------------------------------//
	//                           	     CONFIG     	                            //
	//------------------------------------------------------------------------------//
	config, err := cfg.NewConfig()
	if err != nil {
		log.Fatal("failed to create config", err)
	}

	//------------------------------------------------------------------------------//
	//                           	     LOGGER     	                            //
	//------------------------------------------------------------------------------//
	logger := logging.Logger()
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Fatal("failed to sync logger", err)
		}
	}(logger)

	//------------------------------------------------------------------------------//
	//                           	    CLIENTS     	                            //
	//------------------------------------------------------------------------------//
	// KAFKA
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true

	kafkaClient, err := sarama.NewClient([]string{config.KafkaAddress}, kafkaConfig)
	if err != nil {
		logger.Fatal("failed to create kafka client", zap.Error(err))
	}
	logger.Info("created kafka client", zap.String("client", config.KafkaAddress))

	// close the kafka client
	defer func(kafkaClient sarama.Client) {
		if err = kafkaClient.Close(); err != nil {
			logger.Error("failed to close kafka client", zap.Error(err))
		}
		logger.Info("closed kafka client")
	}(kafkaClient)

	// TWITTER
	twitterClient := twit.NewClient(logger)

	// ELASTICSEARCH
	elasticsearchClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		logger.Fatal("failed to create client", zap.Error(err))
	}
	logger.Info("created Elasticsearch client", zap.String("client", config.ESAddress))

	//------------------------------------------------------------------------------//
	//                           	    STORAGES     	                            //
	//------------------------------------------------------------------------------//
	msgStore := messagestore.NewStore(logger, elasticsearchClient)

	//------------------------------------------------------------------------------//
	//                           	     QUEUES     	                            //
	//------------------------------------------------------------------------------//
	// Create a new Kafka consumer.
	consumer, err := sarama.NewConsumerFromClient(kafkaClient)
	if err != nil {
		logger.Fatal("failed to create consumer", zap.Error(err))
	}
	logger.Info("created Kafka consumer")

	// Close the consumer
	defer func(consumer sarama.Consumer) {
		// close the consumer
		if err = consumer.Close(); err != nil {
			logger.Error("failed to close consumer", zap.Error(err))
		}
		logger.Info("closed Kafka consumer")
	}(consumer)

	// Create a new Kafka producer.
	producer, err := sarama.NewSyncProducerFromClient(kafkaClient)
	if err != nil {
		logger.Fatal("failed to create producer", zap.Error(err))
	}
	logger.Info("created Kafka producer")

	// Close the producer
	defer func(producer sarama.SyncProducer) {
		if err = producer.Close(); err != nil {
			logger.Error("failed to close producer", zap.Error(err))
		}
		logger.Info("closed Kafka producer")
	}(producer)

	msgQueue := messagequeue.NewQueue(config, logger, consumer, producer)

	//------------------------------------------------------------------------------//
	//                           	    SERVICES     	                            //
	//------------------------------------------------------------------------------//
	msgFetcherSvc := messagefetcher.NewService(logger, twitterClient, msgQueue)
	msgSaverSvc := messagesaver.NewService(logger, msgStore, msgQueue)
	_ = msgSaverSvc

	//------------------------------------------------------------------------------//
	//                           	    HANDLERS     	                            //
	//------------------------------------------------------------------------------//
	server := transport.NewServer(logger, config)

	// Handlers.
	messagesH := messages.NewHandler(logger, msgFetcherSvc)

	// Middlewares.
	globalMiddlewares := []func() gin.HandlerFunc{
		loggermw.Logger,
	}

	// Register the handlers.
	server.RegisterHandlers(map[transport.Handler][]func() gin.HandlerFunc{
		messagesH: globalMiddlewares,
	})

	// Run the server.
	logger.Fatal("server failed to run", zap.Error(server.Run()))
}
