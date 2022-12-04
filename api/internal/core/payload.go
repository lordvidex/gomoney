package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenExpired = errors.New("token is expired")
)

type Payload struct {
	ID       uuid.UUID
	Phone    string    `json:"phone"`
	IssuedAt time.Time `json:"iat"`
	ExpireAt time.Time `json:"exp"`
}

func NewPayload(phone string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:       tokenID,
		Phone:    phone,
		IssuedAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}, nil
}

func (payload *Payload) Valid() error {
	if payload.ExpireAt.Before(time.Now()) {
		return ErrTokenExpired
	}
	return nil
}
