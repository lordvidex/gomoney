package application

import (
	"errors"
)

var (
	ErrInvalidLogin   = errors.New("invalid login credentials")
	ErrInvalidToken   = errors.New("token is invalid")
	ErrAssigningToken = errors.New("error assigning token")
	ErrUserDeleted    = errors.New("user has been deleted")
)
