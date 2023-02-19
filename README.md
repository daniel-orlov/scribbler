## Scribbler

A twitter fetcher and parser for downstream processing in Elasticsearch.

This application consists of four main components:

- `messagefetcher` service that fetches tweets from Twitter API and passes them to Kafka
- Kafka producer that sends tweets to Kafka
- Kafka consumer that reads tweets from Kafka
- `messagesaver` service that stores tweets using the Elasticsearch