package paseto

import (
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPasetoMaker(t *testing.T) {
	symmetricKey := []byte("12345678901234567890123456789012")
	paseto := New(symmetricKey)

	tests := []struct {
		name    string
		payload core.Payload
	}{
		{name: "OK", payload: core.Payload{
			Phone: "79600313041",
		},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := paseto.CreateToken(tt.payload)
			require.NoError(t, err)
			require.NotEmpty(t, token)

			payload, err := paseto.VerifyToken(token)
			require.NoError(t, err)
			require.NotEmpty(t, payload)
			require.Equal(t, tt.payload.Phone, payload.Phone)
			require.True(t, payload.ExpireAt.After(payload.IssuedAt))
		})
	}
}
