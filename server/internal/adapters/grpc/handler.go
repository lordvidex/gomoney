// handler.go contains the implementation of the gRPC handlers
package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/lordvidex/gomoney/pkg/grpc"
	"github.com/lordvidex/gomoney/server/internal/application"
)

type Handler struct {
	application *application.UseCases
	pb.UnimplementedAccountServiceServer
	pb.UnimplementedUserServiceServer
}

func NewHandler(application *application.UseCases) *Handler {
	return &Handler{application: application}
}

func (h *Handler) CreateUser(ctx context.Context, u *pb.User) (*pb.StringID, error) {
	id, err := h.application.CreateUser.Handle(ctx, application.CreateUserArg{
		Phone: u.Phone,
		Name:  u.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.StringID{Id: id.String()}, nil
}

func (h *Handler) GetUserByPhone(ctx context.Context, p *pb.Phone) (*pb.User, error) {
	user, err := h.application.GetUser.Handle(ctx, application.GetUserQueryArg{
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

func (h *Handler) GetAccounts(ctx context.Context, id *pb.StringID) (*pb.ManyAccounts, error) {
	uid, err := uuid.Parse(id.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	accs, err := h.application.GetAccountsForUser.Handle(ctx, application.GetUserAccountsParam{
		UserID: uid,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	accounts := make([]*pb.Account, len(accs))
	for i, acc := range accs {
		accounts[i] = &pb.Account{
			Id:          acc.Id,
			Title:       acc.Title,
			Description: acc.Description,
			Balance:     acc.Balance,
			Currency:    pb.Currency(pb.Currency_value[string(acc.Currency)]),
			IsBlocked:   acc.IsBlocked,
		}
	}
	return &pb.ManyAccounts{
		Accounts: accounts,
		Owner:    id,
	}, nil
}
func (h *Handler) CreateAccount(ctx context.Context, param *pb.ManyAccounts) (*pb.IntID, error) {
	// account
	if len(param.GetAccounts()) != 1 {
		return nil, status.Error(codes.InvalidArgument, "invalid number of accounts: provide exactly one account")
	}
	acc := param.GetAccounts()[0]

	// user id
	id := param.GetOwner().GetId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid owner id")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id, please provide uuid")
	}
	accid, err := h.application.CreateAccount.Handle(ctx, application.CreateAccountArg{
		UserID:      uid,
		Title:       acc.GetTitle(),
		Description: acc.GetDescription(),
		Currency:    gomoney.Currency(acc.GetCurrency().String()),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.IntID{Id: accid}, nil
}
