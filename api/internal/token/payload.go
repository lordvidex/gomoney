package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token is expired")
)

type Payload struct {
	ID       uuid.UUID
	Phone    string    `json:"phone"`
	IssuedAt time.Time `json:"issued_at"`
	ExpireAt time.Time `json:"expire_at"`
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
