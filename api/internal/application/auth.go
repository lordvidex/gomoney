package application

import (
	"context"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/pkg/errors"
)

type LoginParam struct {
	Phone    string
	Password string
}

type LoginCommand interface {
	Handle(context.Context, LoginParam) (*core.ApiUserWithToken, error)
}

type loginCommandImpl struct {
	repo Repository
	tok  TokenHelper
	srv  Service
	ph   PasswordHasher
}

func NewLoginCommand(repo Repository, tok TokenHelper, srv Service, ph PasswordHasher) LoginCommand {
	return &loginCommandImpl{repo, tok, srv, ph}
}

func (l *loginCommandImpl) Handle(ctx context.Context, param LoginParam) (*core.ApiUserWithToken, error) {
	user, err := l.repo.GetUserFromPhone(ctx, param.Phone)
	if err != nil {
		return nil, err
	}
	if err = l.ph.CheckPasswordHash(user.Password, param.Password); err != nil {
		return nil, ErrInvalidLogin
	}
	// get the user from the service to confirm
	user, err = l.srv.GetUserByPhone(ctx, user.Phone)
	if err != nil {
		return nil, ErrUserDeleted
	}
	token, err := l.tok.CreateToken(core.Payload{Phone: user.Phone})
	if err != nil {
		return nil, ErrAssigningToken
	}
	return &core.ApiUserWithToken{ApiUser: user, Token: token}, nil
}

type APIUserQuery interface {
	Handle(context.Context, string) (*core.ApiUser, error)
}

type apiUserQueryImpl struct {
	repo Repository
	tok  TokenHelper
}

func NewAPIUserQuery(repo Repository, t TokenHelper) APIUserQuery {
	return &apiUserQueryImpl{repo, t}
}

func (a *apiUserQueryImpl) Handle(ctx context.Context, token string) (*core.ApiUser, error) {
	payload, err := a.tok.VerifyToken(token)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidToken.Error())
	}
	user, err := a.repo.GetUserFromPhone(ctx, payload.Phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}
