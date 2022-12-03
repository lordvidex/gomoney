package application

import (
	"errors"
	"fmt"
)

var (
	ErrSimilarAccountTransaction = func(from, to int64) error {
		return fmt.Errorf("from account and to account CANNOT be the same, from: %d, to: %d", from, to)
	}

	ErrInvalidLogin   = errors.New("invalid login credentials")
	ErrInvalidToken   = errors.New("token is invalid")
	ErrAssigningToken = errors.New("error assigning token")
)
