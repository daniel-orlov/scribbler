package main

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"scribbler/cfg"
	"scribbler/internal/adapters/twit"
	"scribbler/internal/service/messagefetcher"
	"scribbler/internal/service/messagesaver"
	"scribbler/internal/storage/messagequeue"
	"scribbler/internal/storage/messagestore"
	"scribbler/internal/transport"
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
	logger, err := zap.NewDevelopment(
		zap.IncreaseLevel(zapcore.InfoLevel),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		log.Fatal("failed to create logger", err)
	}
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
		logger.Fatal("failed to create client", zap.Error(err))
	}
	logger.Info("created client", zap.String("client", config.KafkaAddress))

	// close the kafka client
	defer func(kafkaClient sarama.Client) {
		if err = kafkaClient.Close(); err != nil {
			logger.Error("failed to close client", zap.Error(err))
		}
		logger.Info("closed client")
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
	// create a new Kafka consumer.
	consumer, err := sarama.NewConsumerFromClient(kafkaClient)
	if err != nil {
		logger.Fatal("failed to create consumer", zap.Error(err))
	}
	logger.Info("created Kafka consumer")

	// close the consumer
	defer func(consumer sarama.Consumer) {
		// close the consumer
		if err = consumer.Close(); err != nil {
			logger.Error("failed to close consumer", zap.Error(err))
		}
		logger.Info("closed Kafka consumer")
	}(consumer)

	// create a new Kafka producer.
	producer, err := sarama.NewSyncProducerFromClient(kafkaClient)
	if err != nil {
		logger.Fatal("failed to create producer", zap.Error(err))
	}
	logger.Info("created Kafka producer")

	// close the producer
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

	//------------------------------------------------------------------------------//
	//                           	    HANDLERS     	                            //
	//------------------------------------------------------------------------------//
	// create a new HTTP server using Gin.
	engine := gin.New()
	server := transport.NewServer(logger, config, engine)

	// register the handlers.
	server.RegisterHandlers(msgFetcherSvc, msgSaverSvc)

}
