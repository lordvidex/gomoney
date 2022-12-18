package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type JSON struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty" swaggerignore:"true"`
	Error   []Error     `json:"error,omitempty" swaggerignore:"true"`
}

type Error struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func Success(data interface{}) JSON {
	return JSON{Success: true, Data: data}
}

func Errs(err ...Error) JSON {
	return JSON{Success: false, Error: err}
}

// ErrM creates a single error with message err
// and returns a JSON response with Success set to false.
func ErrM(err string) JSON {
	return JSON{Success: false, Error: []Error{{Message: err}}}
}

// Err is the standard constructor for Error struct that
// maps the error type to
func Err(msg string, code gomoney.ErrType) Error {
	return Error{
		Message: msg,
		Error:   code.String(),
		Code:    C(code),
	}
}

// C maps the error type to HTTP status code.
func C(t gomoney.ErrType) int {
	switch t {
	case gomoney.ErrNotFound:
		return fiber.StatusNotFound
	case gomoney.ErrInvalidInput:
		return fiber.StatusBadRequest
	case gomoney.ErrInternal:
		return fiber.StatusInternalServerError
	case gomoney.ErrAlreadyExists:
		return fiber.StatusConflict
	}
	return fiber.StatusInternalServerError
}
