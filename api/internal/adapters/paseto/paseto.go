package paseto

import (
	"github.com/lordvidex/gomoney/api/internal/core"
	"time"

	"github.com/o1egl/paseto"
)

type Maker struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func New(symmetricKey []byte) *Maker {
	return &Maker{
		symmetricKey: symmetricKey,
		paseto:       paseto.NewV2(),
	}
}

func (m *Maker) CreateToken(payload core.Payload) (string, error) {
	payload.IssuedAt = time.Now()
	payload.ExpireAt = time.Now().Add(m.TokenDuration())
	return m.paseto.Encrypt(m.symmetricKey, payload, nil)
}

func (m *Maker) VerifyToken(token string) (core.Payload, error) {
	payload := core.Payload{}
	err := m.paseto.Decrypt(token, m.symmetricKey, &payload, nil)

	if err = payload.Valid(); err != nil {
		return core.Payload{}, err
	}
	
	return payload, err
}

func (m *Maker) TokenDuration() time.Duration {
	return time.Hour * 24
}
