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

func (payload *Payload) Valid() error {
	if payload.ExpireAt.Before(time.Now()) {
		return ErrTokenExpired
	}
	return nil
}
