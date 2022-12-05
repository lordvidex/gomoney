// handler.go contains the implementation of the gRPC handlers
package grpc

import (
	context "context"

	pb "github.com/lordvidex/gomoney/pkg/grpc"
	"github.com/lordvidex/gomoney/server/internal/application"
)

type handler struct {
	application *application.UseCases
	pb.UnimplementedAccountServiceServer
	pb.UnimplementedUserServiceServer
}

func NewHandler(application *application.UseCases) *handler {
	return &handler{application: application}
}

func (h *handler) CreateUser(_ context.Context, u *pb.User) (*pb.StringID, error) {
	id, err := h.application.CreateUser.Handle(application.CreateUserArg{
		Phone: u.Phone,
		Name:  u.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.StringID{Id: id.String()}, nil
}

func (h *handler) GetUserByPhone(_ context.Context, p *pb.Phone) (*pb.User, error) {
	user, err := h.application.GetUser.Handle(application.GetUserQueryArg{
		Phone: p.GetNumber(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Id:    user.ID.String(),
		Phone: user.Phone,
		Name:  user.Name,
	}, nil
}
