package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	ErrServiceCall = errors.New("call to service failed")
)

type service struct {
	ucl UserServiceClient
	acl AccountServiceClient
}

func (s service) CreateUser(ctx context.Context, param application.CreateUserParam) (string, error) {
	id, err := s.ucl.CreateUser(ctx, &User{
		Name:  param.Name,
		Phone: param.Phone,
	})
	if err != nil {
		return "", errors.Wrap(err, ErrServiceCall.Error())
	}
	return id.GetId(), nil
}

func (s service) GetUserByPhone(ctx context.Context, phone string) (*core.ApiUser, error) {
	u, err := s.ucl.GetUserByPhone(ctx, &Phone{Number: phone})
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
	acc, err := s.acl.GetAccounts(ctx, &StringID{Id: ID})
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
	accID, err := s.acl.CreateAccount(ctx, &ManyAccounts{
		Owner: &StringID{
			Id: userID,
		},
		Accounts: []*Account{
			{
				Title:       acc.Title,
				Description: acc.Description,
				Balance:     acc.Balance,
				Currency:    Currency(Currency_value[string(acc.Currency)]),
				IsBlocked:   acc.IsBlocked,
			},
		},
	})
	if err != nil {
		return 0, errors.Wrap(err, ErrServiceCall.Error())
	}
	return accID.GetId(), nil
}

func New(conn grpc.ClientConnInterface) application.Service {
	return &service{
		ucl: NewUserServiceClient(conn),
		acl: NewAccountServiceClient(conn),
	}
}
