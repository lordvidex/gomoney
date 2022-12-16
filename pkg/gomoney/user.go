package gomoney

import (
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID
	Name  string
	Phone string
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// func (User) UnmarshalBinary(data []byte) error
