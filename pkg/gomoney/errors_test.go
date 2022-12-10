package gomoney

import (
	"errors"
	"testing"
)

func TestError_AnyOf(t *testing.T) {
	tests := []struct {
		name string
		err1 Error
		err2 error
		want bool
	}{
		{
			"ErrNotFound is not ErrInvalidInput",
			ErrNotFound,
			ErrInvalidInput,
			false,
		},
		{
			"ErrInvalidInput in string type is equal to ErrInvalidInput",
			ErrInvalidInput,
			errors.New("invalid input"),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err1.AnyOf(tt.err2); got != tt.want {
				t.Errorf("AnyOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
