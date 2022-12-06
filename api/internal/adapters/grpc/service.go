package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	grpc3 "github.com/lordvidex/gomoney/pkg/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	ErrServiceCall = errors.New("call to service failed")
)

type service struct {
	ucl grpc3.UserServiceClient
	acl grpc3.AccountServiceClient
	tcl grpc3.TransactionServiceClient
}

func (s service) CreateUser(ctx context.Context, param application.CreateUserParam) (string, error) {
	id, err := s.ucl.CreateUser(ctx, &grpc3.User{
		Name:  param.Name,
		Phone: param.Phone,
	})
	if err != nil {
		return "", errors.Wrap(err, ErrServiceCall.Error())
	}
	return id.GetId(), nil
}

func (s service) GetUserByPhone(ctx context.Context, phone string) (*core.ApiUser, error) {
	u, err := s.ucl.GetUserByPhone(ctx, &grpc3.Phone{Number: phone})
	if err != nil {
		return nil, errors.Wrap(err, ErrServiceCall.Error())
	}
	return &core.ApiUser{
		User: gomoney.User{
			ID:    uuid.MustParse(u.GetId()),
			Name:  u.Name,
			Phone: u.Phone,
		},
	}, nil
}

func (s service) GetAccounts(ctx context.Context, ID string) ([]gomoney.Account, error) {
	acc, err := s.acl.GetAccounts(ctx, &grpc3.StringID{Id: ID})
	if err != nil {
		return nil, errors.Wrap(err, ErrServiceCall.Error())
	}
	accs := make([]gomoney.Account, len(acc.Accounts))
	for i := 0; i < len(accs); i++ {
		acci := acc.GetAccounts()[i]
		accs[i] = gomoney.Account{
			Id:          acci.GetId(),
			Title:       acci.GetTitle(),
			Description: acci.GetDescription(),
			Balance:     acci.GetBalance(),
			Currency:    gomoney.Currency(acci.GetCurrency().String()),
			IsBlocked:   acci.GetIsBlocked(),
		}
	}
	return accs, nil
}

func (s service) CreateAccount(ctx context.Context, userID string, acc *gomoney.Account) (int64, error) {
	accID, err := s.acl.CreateAccount(ctx, &grpc3.ManyAccounts{
		Owner: &grpc3.StringID{
			Id: userID,
		},
		Accounts: []*grpc3.Account{
			{
				Title:       acc.Title,
				Description: acc.Description,
				Balance:     acc.Balance,
				Currency:    grpc3.Currency(grpc3.Currency_value[string(acc.Currency)]),
				IsBlocked:   acc.IsBlocked,
			},
		},
	})
	if err != nil {
		return 0, errors.Wrap(err, ErrServiceCall.Error())
	}
	return accID.GetId(), nil
}

func (s service) GetAccount(ctx context.Context, accountID int64) (gomoney.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetAccountTransfers(ctx context.Context, accountID int64) ([]gomoney.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetTransfers(ctx context.Context, userID string) ([]gomoney.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) CreateTransfer(ctx context.Context, param application.CreateTransferParam) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func New(conn grpc.ClientConnInterface) application.Service {
	return &service{
		ucl: grpc3.NewUserServiceClient(conn),
		acl: grpc3.NewAccountServiceClient(conn),
		tcl: grpc3.NewTransactionServiceClient(conn),
	}
}
