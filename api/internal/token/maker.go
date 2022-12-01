package token

import "time"

type Maker interface {
	CreateToken(phone string, duration time.Duration) (string, error)
	VerifyToken(token string) (Payload, error)
}
