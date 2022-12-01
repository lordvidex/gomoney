// handler.go contains the implementation of the gRPC handlers
package grpc

import (
	context "context"

	"github.com/lordvidex/gomoney/server/internal/application"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type handler struct {
	application *application.UseCases
	UnimplementedAccountServiceServer
	UnimplementedUserServiceServer
}

func NewHandler(application *application.UseCases) *handler {
	return &handler{application: application}
}

func (h *handler) CreateUser(_ context.Context, u *User) (*emptypb.Empty, error) {
	err := h.application.CreateUser.Handle(application.CreateUserArg{
		Phone: u.Phone,
		Name:  u.Name,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *handler) GetUser(_ context.Context, u *User) (*User, error) {
	user, err := h.application.GetUser.Handle(application.GetUserQueryArg{
		Phone: u.Phone,
	})
	if err != nil {
		return nil, err
	}
	return &User{
		Phone: user.Phone,
		Name:  user.Name,
	}, nil
}
