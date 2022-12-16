package gomoney

import "strings"

// errors.go contain application defined errors in which all services must
// adhere to and map their errors to.
// for example gomoney.ErrNotFound == 404 (HTTP) == codes.NotFound (GRPC)
//
// It mostly leans towards the HTTP codes for convenience

type ErrType int

func (e ErrType) String() string {
	switch e {
	case ErrNotFound:
		return "entity not found"
	case ErrAlreadyExists:
		return "entity already exists"
	case ErrInvalidInput:
		return "invalid input"
	case ErrInternal:
		return "internal error"
	default:
		return "unknown error"
	}
}

const (
	ErrInternal ErrType = iota
	ErrNotFound 
	ErrAlreadyExists
	ErrInvalidInput
)

type Error struct {
	Code     ErrType
	Messages []string
}


// Err is a factory method to create an Error instance with the code set to 
// ErrInternal and no messages.
func Err() Error {
	return Error{Code: ErrInternal}
}

// Is returns true if the Error instance is equal to the error passed in by
// comparing the code and the list of messages.
//
// In the case where the error passed in is not an Error instance, it will
// compare the error string with the code and the list of messages.
func (e Error) Is(err error) bool {
	if err == nil {
		return false
	}
	switch ex := err.(type) {
	case Error:
		return e.Code == ex.Code && e.Error() == ex.Error()
	case *Error:
		return e.Code == ex.Code && e.Error() == ex.Error()
	default:
		msgs := append(e.Messages, e.Code.String())
		for _, msg := range msgs {
			if strings.Contains(err.Error(), msg) || strings.Contains(msg, err.Error()) {
				return true
			}
		}
		return false
	}
}

func (e Error) IsCode(err error) bool {
	if err == nil {
		return false
	}
	switch ex := err.(type) {
	case Error:
		return e.Code == ex.Code
	case *Error:
		return e.Code == ex.Code
	default:
		return strings.Contains(e.Code.String(), err.Error()) || strings.Contains(err.Error(), e.Code.String())
	}
}

// WithMessage should ALWAYS be used to add message to the error instance.
// This is because the messages field might be nil
func (e Error) WithMessage(message string) Error {
	e.Messages = append(e.Messages, message)
	return e
}

func (e Error) WithCode(code ErrType) Error {
	e.Code = code
	return e
}

func (e Error) Error() string {
	return e.Code.String() + ": [" + strings.Join(e.Messages, ", ") + "]"
}
