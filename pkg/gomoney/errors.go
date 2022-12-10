package gomoney

// errors.go contain application defined errors in which all services must
// adhere to and map their errors to.
// for example application.ErrNotFound == 404 (HTTP) == codes.NotFound (GRPC)
//
// It mostly leans towards the HTTP codes for convenience

import (
	"strings"
)

type Error struct {
	message string
}

func (e Error) Error() string {
	return e.message
}

func err(message string) Error {
	return Error{message: message}
}

func (e Error) AnyOf(err ...error) bool {
	for _, e2 := range err {
		if e.Is(e2) {
			return true
		}
	}
	return false
}

func (e Error) Is(err error) bool {
	return strings.Contains(e.message, err.Error()) || strings.Contains(err.Error(), e.message)
}

var (
	ErrNotFound      = err("entity not found")
	ErrAlreadyExists = err("entity already exists")
	ErrInvalidInput  = err("invalid input")
)
