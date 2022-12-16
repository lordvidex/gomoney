package core

import (
	"encoding/json"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type ApiUser struct {
	gomoney.User
	Password string
}

type ApiUserWithToken struct {
	*ApiUser
	Token string
}

func (au *ApiUser) MarshalBinary() ([]byte, error) {
	return json.Marshal(au)
}

func (au *ApiUser) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, au)
}

func (aut *ApiUserWithToken) MarshalBinary() ([]byte, error) {
	return json.Marshal(aut)
}

func (aut *ApiUserWithToken) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, aut)
}
