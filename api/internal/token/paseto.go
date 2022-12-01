package token

import (
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func NewPasetoMaker(symmetricKey []byte) *PasetoMaker {
	return &PasetoMaker{
		symmetricKey: symmetricKey,
		paseto:       paseto.NewV2(),
	}
}

func (maker *PasetoMaker) CreateToken(phone string, duration time.Duration) (string, error) {
	payload, err := NewPayload(phone, duration)

	if err != nil {
		return "", err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, err
}

func (maker *PasetoMaker) VerifyToken(token string) (Payload, error) {
	payload := Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil)
	return payload, err
}
