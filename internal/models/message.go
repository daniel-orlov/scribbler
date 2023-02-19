package models

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/oklog/ulid"
)

type Message struct {
	ID        ulid.ULID `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UserID    ulid.ULID `json:"user_id"`
	Likes     int       `json:"likes"`
	Reposts   int       `json:"reposts"`
}

func NewFakeMessage() Message {
	return Message{
		ID:        ulid.MustNew(ulid.Now(), nil),
		Text:      gofakeit.Sentence(rand.Intn(15)),
		CreatedAt: time.Now(),
		UserID:    ulid.MustNew(ulid.Now(), nil),
		Likes:     rand.Intn(1500),
		Reposts:   rand.Intn(1500),
	}
}
