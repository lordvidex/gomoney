package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	lgrpc "github.com/lordvidex/gomoney/pkg/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	ErrServiceCall = errors.New("call to service failed")
)

type grpcstatus interface {
	GRPCStatus() *status.Status
}

type service struct {
	ucl lgrpc.UserServiceClient
	acl lgrpc.AccountServiceClient
	tcl lgrpc.TransactionServiceClient
}

func (s service) CreateUser(ctx context.Context, param application.CreateUserParam) (string, error) {
	id, err := s.ucl.CreateUser(ctx, &lgrpc.User{
		Name:  param.Name,
		Phone: param.Phone,
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return "", gomoney.
				Err(gomoney.ErrAlreadyExists).
				WithMessage(err.(grpcstatus).GRPCStatus().Message())
		}
		return "", errors.Wrap(err, ErrServiceCall.Error())
	}
	return id.GetId(), nil
}

func (s service) GetUserByPhone(ctx context.Context, phone string) (*core.ApiUser, error) {
	u, err := s.ucl.GetUserByPhone(ctx, &lgrpc.Phone{Number: phone})
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
	acc, err := s.acl.GetAccounts(ctx, &lgrpc.StringID{Id: ID})
	if err != nil {
		return nil, errors.Wrap(err, ErrServiceCall.Error())
	}
	accs := make([]gomoney.Account, len(acc.Accounts))
	for i, acci := range acc.Accounts {
		accs[i] = *mapPbAccToAcc(acci)
		accs[i].OwnerID = uuid.MustParse(ID)
	}
	return accs, nil
}

func (s service) CreateAccount(ctx context.Context, param application.CreateAccountParam) (int64, error) {
	accID, err := s.acl.CreateAccount(ctx, &lgrpc.ManyAccounts{
		Owner: &lgrpc.StringID{
			Id: param.UserID.String(),
		},
		Accounts: []*lgrpc.Account{
			{
				Title:       param.Title,
				Description: param.Description,
				Balance:     0,
				Currency:    lgrpc.Currency(lgrpc.Currency_value[string(param.Currency)]),
				IsBlocked:   false,
			},
		},
	})
	if err != nil {
		return 0, errors.Wrap(err, ErrServiceCall.Error())
	}
	return accID.GetId(), nil
}

func (s service) DeleteAccount(ctx context.Context, param application.DeleteAccountParam) error {
	_, err := s.acl.DeleteAccount(ctx, &lgrpc.UserWithAccount{
		User:    &lgrpc.StringID{Id: param.UserID},
		Account: &lgrpc.IntID{Id: param.AccountID},
	})
	if err != nil {
		return errors.Wrap(err, ErrServiceCall.Error())
	}
	return nil
}

func (s service) Transfer(ctx context.Context, param application.CreateTransferParam) (*gomoney.Transaction, error) {
	transaction, err := s.tcl.Transfer(ctx, &lgrpc.TransactionParam{
		To:     &param.ToID,
		From:   &param.FromID,
		Amount: param.Amount,
		Actor: &lgrpc.StringID{
			Id: param.ActorID,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, ErrServiceCall.Error())
	}
	return mapPbTxToTx(transaction), nil
}

func (s service) Deposit(ctx context.Context, param application.DepositParam) (*gomoney.Transaction, error) {
	transaction, err := s.tcl.Deposit(ctx, &lgrpc.TransactionParam{
		To:     &param.ToID,
		Amount: param.Amount,
		Actor: &lgrpc.StringID{
			Id: param.ActorID,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, ErrServiceCall.Error())
	}
	return mapPbTxToTx(transaction), nil
}

func (s service) Withdraw(ctx context.Context, param application.WithdrawParam) (*gomoney.Transaction, error) {
	transaction, err := s.tcl.Withdraw(ctx, &lgrpc.TransactionParam{
		From:   &param.FromID,
		Amount: param.Amount,
		Actor: &lgrpc.StringID{
			Id: param.ActorID,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, ErrServiceCall.Error())
	}
	return mapPbTxToTx(transaction), nil
}

func (s service) GetTransactionSummary(ctx context.Context, userID string) ([]gomoney.TransactionSummary, error) {
	transactions, err := s.tcl.GetTransactionSummary(ctx, &lgrpc.StringID{Id: userID})
	if err != nil {
		return nil, errors.Wrap(err, ErrServiceCall.Error())
	}

	accs, err := s.GetAccounts(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, ErrServiceCall.Error())
	}

	m := mapAccountsIntoMap(accs)
	txs := make([]gomoney.TransactionSummary, len(transactions.Transactions))
	for i, tx := range transactions.Transactions {
		txs[i] = mapPbTxSummaryToTxSummary(tx, m)
	}

	return txs, nil
}

func (s service) GetTransactions(ctx context.Context, param application.UserWithAccount) (gomoney.TransactionSummary, error) {
	accTx, err := s.tcl.GetTransactions(ctx, &lgrpc.UserWithAccount{
		User:    &lgrpc.StringID{Id: param.UserID},
		Account: &lgrpc.IntID{Id: param.AccountID},
	})
	if err != nil {
		return gomoney.TransactionSummary{}, errors.Wrap(err, ErrServiceCall.Error())
	}

	accs, err := s.GetAccounts(ctx, param.UserID)
	if err != nil {
		return gomoney.TransactionSummary{}, errors.Wrap(err, ErrServiceCall.Error())
	}

	m := mapAccountsIntoMap(accs)
	return mapPbTxSummaryToTxSummary(accTx, m), nil
}

func mapPbAccToAcc(acc *lgrpc.Account) *gomoney.Account {
	if acc == nil {
		return nil
	}
	return &gomoney.Account{
		Id:          acc.GetId(),
		Title:       acc.GetTitle(),
		Description: acc.GetDescription(),
		Balance:     acc.GetBalance(),
		Currency:    gomoney.Currency(acc.GetCurrency().String()),
		IsBlocked:   acc.GetIsBlocked(),
	}
}

func mapPbTxToTx(t *lgrpc.Transaction) *gomoney.Transaction {
	if t == nil {
		return nil
	}
	return &gomoney.Transaction{
		ID:      uuid.MustParse(t.GetId().GetId()),
		From:    mapPbAccToAcc(t.GetFrom()),
		To:      mapPbAccToAcc(t.GetTo()),
		Type:    gomoney.TransactionType(t.GetType()),
		Amount:  t.GetAmount(),
		Created: t.GetCreatedAt().AsTime(),
	}
}

func mapPbTxSummaryToTxSummary(t *lgrpc.AccountTransactions, m map[int64]gomoney.Account) gomoney.TransactionSummary {
	var accTx gomoney.TransactionSummary
	temp := m[t.Account.Id]
	accTx.Account = &temp
	accTx.Transactions = make([]gomoney.Transaction, len(t.Transactions))
	for i, tx := range t.Transactions {
		accTx.Transactions[i] = *mapPbTxToTx(tx)
	}
	return accTx
}

func mapAccountsIntoMap(accs []gomoney.Account) map[int64]gomoney.Account {
	m := make(map[int64]gomoney.Account)
	for _, acc := range accs {
		m[acc.Id] = acc
	}
	return m
}

func New(conn grpc.ClientConnInterface) application.Service {
	return &service{
		ucl: lgrpc.NewUserServiceClient(conn),
		acl: lgrpc.NewAccountServiceClient(conn),
		tcl: lgrpc.NewTransactionServiceClient(conn),
	}
}
