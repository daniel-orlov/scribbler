## Scribbler

A twitter fetcher and parser for downstream processing in Elasticsearch.

This application consists of four main components:

- `messagefetcher` service that fetches tweets from Twitter API and passes them to Kafka
- Kafka producer that sends tweets to Kafka
- Kafka consumer that reads tweets from Kafka
- `messagesaver` service that stores tweets using the Elasticsearch

## Project structure

>
>> /cfg - configuration files
>
>> /cmd - entry points for the application
>
>> /deploy - deployment artifacts
>
>> /internal - internal application code
>>
>>> /adapters - adapters for external services
>>
>>> /models - domain models
>>
>>> /services - business logic
>>
>>> /storage - storage adapters
>>
>>> /transport - transport layer
>>
>
>> /pkg - reusable packages
>

## Development

### Prerequisites

- Docker
- Docker Compose
- Go 1.19
- Make
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
