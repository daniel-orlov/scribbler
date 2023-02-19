package messagestore

import (
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

type Store struct {
	logger *zap.Logger
	es     *elasticsearch.Client
}

func NewStore(logger *zap.Logger, es *elasticsearch.Client) *Store {
	return &Store{logger: logger, es: es}
}
