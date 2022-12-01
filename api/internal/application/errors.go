package application

import (
	"fmt"
)

var (
	ErrSimilarAccountTransaction = func(from, to int64) error {
		return fmt.Errorf("from account and to account CANNOT be the same, from: %d, to: %d", from, to)
	}
)
