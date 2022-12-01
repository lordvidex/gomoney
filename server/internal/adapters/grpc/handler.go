// handler.go contains the implementation of the gRPC handlers
package grpc

import (
	context "context"

	"github.com/lordvidex/gomoney/server/internal/application"
)

type handler struct {
	application *application.UseCases
	UnimplementedAccountServiceServer
	UnimplementedUserServiceServer
}

func NewHandler(application *application.UseCases) *handler {
	return &handler{application: application}
}

func (h *handler) CreateUser(_ context.Context, u *User) (*User_ID, error) {
	id, err := h.application.CreateUser.Handle(application.CreateUserArg{
		Phone: u.Phone,
		Name:  u.Name,
	})
	if err != nil {
		return nil, err
	}
	return &User_ID{Id: id.String()}, nil
}

func (h *handler) GetUserByPhone(_ context.Context, p *Phone) (*User, error) {
	user, err := h.application.GetUser.Handle(application.GetUserQueryArg{
		Phone: p.GetNumber(),
	})
	if err != nil {
		return nil, err
	}
	return &User{
		Id:    user.ID.String(),
		Phone: user.Phone,
		Name:  user.Name,
	}, nil
}
