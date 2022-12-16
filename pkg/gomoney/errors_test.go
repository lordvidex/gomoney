package gomoney

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Is(t *testing.T) {
	tests := []struct {
		name string
		e    Error
		err  error
		want bool
	}{
		{name: "empty", e: Err(), err: nil, want: false},
		{name: "test error should be false", e: Err(), err: errors.New("test"), want: false},
		{name: "description of code should be true", e: Error{Code: ErrNotFound}, err: errors.New("entity not found"), want: true},
		{
			name: "two same errors",
			e:    Err().WithCode(ErrNotFound).WithMessage("test1").WithMessage("test2"),
			err:  Err().WithCode(ErrNotFound).WithMessage("test1").WithMessage("test2"), want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Is(tt.err); got != tt.want {
				t.Errorf("Error.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorIsCode(t *testing.T) {
	tests := []struct {
		name string
		e    Error
		err  error
		want bool
	}{
		{
			name: "default error should be internal",
			e:    Err(),
			err:  errors.New(ErrInternal.String()),
			want: true,
		},
		{
			name: "default error should be internal and should not match other errs",
			e:    Err(),
			err:  errors.New("test"),
			want: false,
		},
		{
			name: "nil error should be false",
			e:    Err(),
			err:  nil,
			want: false,
		},
		{
			name: "errors should still be the same for IsCode even with different messages",
			e:    Err().WithCode(ErrNotFound),
			err:  Err().WithCode(ErrNotFound).WithMessage("test"),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.e.IsCode(tt.err), tt.want)
		})
	}

}

func TestErrorWithMessage(t *testing.T) {
	err := Err().WithMessage("test")
	err2 := err.WithMessage("test2")
	err3 := err.WithMessage("test3")
	t.Log(err, err2, err3)
	assert.NotEqualf(t, err.Messages, err2.Messages, "errors %v and %v should be different", err.Messages, err2.Messages)
	assert.NotEqualf(t, err2, err3, "errors %v and %v should be different", err2, err3)
	assert.Equal(t, len(Err().Messages), 0)
}
