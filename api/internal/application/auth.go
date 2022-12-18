package application

import (
	"context"
	"github.com/lordvidex/gomoney/api/internal/core"
	g "github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/pkg/errors"
)

// -- errors --

var (
	ErrUserNotFound   = g.Err().WithCode(g.ErrNotFound).WithMessage("user does not exist")
	ErrAssigningToken = g.Err().WithCode(g.ErrInternal).WithMessage("error assigning token")
	ErrInvalidLogin   = g.Err().WithCode(g.ErrInvalidInput).WithMessage("invalid login credentials")
	ErrInvalidToken   = g.Err().WithCode(g.ErrInvalidInput).WithMessage("token is invalid")
	ErrUserDeleted    = g.Err().WithCode(g.ErrNotFound).WithMessage("user has been deleted")
)

// ----------- User login Command ------------

type LoginParam struct {
	Phone    string
	Password string
}

type LoginCommand interface {
	Handle(context.Context, LoginParam) (*core.ApiUserWithToken, error)
}

type loginCommandImpl struct {
	ur  UserRepository
	th  TokenHelper
	srv Service
	ph  PasswordHasher
}

func NewLoginCommand(ur UserRepository, th TokenHelper, srv Service, ph PasswordHasher) LoginCommand {
	return &loginCommandImpl{ur, th, srv, ph}
}

func (l *loginCommandImpl) Handle(ctx context.Context, param LoginParam) (*core.ApiUserWithToken, error) {
	user, err := l.ur.GetUserFromPhone(ctx, param.Phone)
	if err != nil {
		return nil, err
	}
	if err = l.ph.CheckPasswordHash(user.Password, param.Password); err != nil {
		return nil, ErrInvalidLogin
	}

	// Get the user from the service to confirm
	user, err = l.srv.GetUserByPhone(ctx, user.Phone)
	if err != nil {
		return nil, ErrUserDeleted
	}

	// Create user token
	token, err := l.th.CreateToken(core.Payload{Phone: user.Phone})
	if err != nil {
		return nil, ErrAssigningToken
	}

	return &core.ApiUserWithToken{ApiUser: user, Token: token}, nil
}

// ----------- User login Command ------------

type APIUserQuery interface {
	Handle(context.Context, string) (*core.ApiUser, error)
}

type apiUserQueryImpl struct {
	repo UserRepository
	ser  Service
	tok  TokenHelper
}

func NewAPIUserQuery(repo UserRepository, ser Service, t TokenHelper) APIUserQuery {
	return &apiUserQueryImpl{repo, ser, t}
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
