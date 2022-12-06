package gomoney

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountCanTransfer(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		want    bool
	}{
		{name: "Enough balance, valid account", account: Account{Balance: 100, IsBlocked: false}, want: true},
		{name: "Not enough balance, valid account", account: Account{Balance: 0, IsBlocked: false}, want: false},
		{name: "Enough balance, blocked account", account: Account{Balance: 100, IsBlocked: true}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.CanTransfer(10)
			assert.Equal(t, tt.want, got)
		})
	}
}
