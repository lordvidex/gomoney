package core

import "github.com/lordvidex/gomoney/pkg/gomoney"

type ApiUser struct {
	gomoney.User
	Password string
}

type ApiUserWithToken struct {
	*ApiUser
	Token string
}
