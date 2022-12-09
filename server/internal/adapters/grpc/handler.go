// handler.go contains the implementation of the gRPC handlers
package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/lordvidex/gomoney/pkg/grpc"
	app "github.com/lordvidex/gomoney/server/internal/application"
)

type Handler struct {
	application *app.UseCases
	pb.UnimplementedAccountServiceServer
	pb.UnimplementedUserServiceServer
	pb.UnimplementedTransactionServiceServer
}

func NewHandler(application *app.UseCases) *Handler {
	return &Handler{application: application}
}

func (h *Handler) CreateUser(ctx context.Context, u *pb.User) (*pb.StringID, error) {
	id, err := h.application.CreateUser.Handle(ctx, app.CreateUserArg{
		Phone: u.Phone,
		Name:  u.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.StringID{Id: id.String()}, nil
}

func (h *Handler) GetUserByPhone(ctx context.Context, p *pb.Phone) (*pb.User, error) {
	user, err := h.application.GetUser.Handle(ctx, app.GetUserQueryArg{
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
	accs, err := h.application.GetAccountsForUser.Handle(ctx, app.GetUserAccountsParam{
		UserID: uid,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	accounts := make([]*pb.Account, len(accs))
	for i, acc := range accs {
		accounts[i] = mapAcctToPb(&acc)
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
	accid, err := h.application.CreateAccount.Handle(ctx, app.CreateAccountArg{
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

func (h *Handler) DeleteAccount(ctx context.Context, param *pb.UserWithAccount) (*emptypb.Empty, error) {
	id := param.GetUser().GetId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id, please provide uuid")
	}
	err = h.application.DeleteAccount.Handle(ctx, app.DeleteAccountArg{
		ActorID:   uid,
		AccountID: param.GetAccount().GetId(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) Transfer(ctx context.Context, param *pb.TransactionParam) (*pb.Transaction, error) {
	if param.From == nil || param.To == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid transaction param: from and to must be provided")
	}
	tid := param.GetActor().GetId()
	actorUUID, err := uuid.Parse(tid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actor id: expected uuid")
	}
	t, err := h.application.Transfer.Handle(ctx, app.TransferArg{
		FromAccountID: param.GetFrom(),
		ToAccountID:   param.GetTo(),
		Amount:        param.GetAmount(),
		ActorID:       actorUUID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return mapTxToPb(t), nil
}

func (h *Handler) Deposit(ctx context.Context, param *pb.TransactionParam) (*pb.Transaction, error) {
	if param.To == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid transaction param: to must be provided")
	}
	tid := param.GetActor().GetId()
	actorUUID, err := uuid.Parse(tid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actor id: expected uuid")
	}
	t, err := h.application.Deposit.Handle(ctx, app.DepositArg{
		AccountID: param.GetTo(),
		Amount:    param.GetAmount(),
		ActorID:   actorUUID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return mapTxToPb(t), nil
}

func (h *Handler) Withdraw(ctx context.Context, param *pb.TransactionParam) (*pb.Transaction, error) {
	if param.From == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid transaction param: from must be provided")
	}
	tid := param.GetActor().GetId()
	actorUUID, err := uuid.Parse(tid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actor id: expected uuid")
	}
	t, err := h.application.Withdrawal.Handle(ctx, app.WithdrawArg{
		AccountID: param.GetFrom(),
		Amount:    param.GetAmount(),
		ActorID:   actorUUID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return mapTxToPb(t), nil
}

func (h *Handler) GetTransactionSummary(ctx context.Context, param *pb.StringID) (*pb.ManyAccountTransactions, error) {
	id := param.GetId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid transaction id: expected uuid")
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid transaction id: expected uuid")
	}
	summaries, err := h.application.GetTransactionSummaryForUser.Handle(ctx, app.TransactionSummaryArg{ActorID: uid})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	summary := make([]*pb.AccountTransactions, len(summaries))
	for i, s := range summaries {
		summary[i] = mapSummaryToPb(s)
	}
	return &pb.ManyAccountTransactions{
		Transactions: summary,
	}, nil

}
func (h *Handler) GetTransactions(ctx context.Context, param *pb.UserWithAccount) (*pb.AccountTransactions, error) {
	userid := param.GetUser().GetId()
	if userid == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid user id: expected uuid")
	}
	uid, err := uuid.Parse(userid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id: expected uuid")
	}
	accid := param.GetAccount().GetId()
	tx, err := h.application.GetTransactionsInAccount.Handle(ctx, app.TransactionsArg{
		AccountID: accid,
		ActorID:   uid,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := make([]*pb.Transaction, len(tx))
	for i, t := range tx {
		result[i] = mapTxToPb(&t)
	}
	return &pb.AccountTransactions{
		Transactions: result,
		Account:      param.GetAccount(),
	}, nil
}

func mapSummaryToPb(s gomoney.TransactionSummary) *pb.AccountTransactions {
	tx := make([]*pb.Transaction, len(s.Transactions))
	for i, t := range s.Transactions {
		tx[i] = mapTxToPb(&t)
	}
	return &pb.AccountTransactions{
		Account:      &pb.IntID{Id: s.Account.Id},
		Transactions: tx,
	}
}

func mapAcctToPb(acc *gomoney.Account) *pb.Account {
	if acc == nil {
		return nil
	}
	return &pb.Account{
		Id:          acc.Id,
		Title:       acc.Title,
		Description: acc.Description,
		Balance:     acc.Balance,
		Currency:    pb.Currency(pb.Currency_value[string(acc.Currency)]),
		IsBlocked:   acc.IsBlocked,
	}
}

func mapTxToPb(t *gomoney.Transaction) *pb.Transaction {
	if t == nil {
		return nil
	}
	return &pb.Transaction{
		Id:        &pb.StringID{Id: t.ID.String()},
		Amount:    t.Amount,
		From:      mapAcctToPb(t.From),
		To:        mapAcctToPb(t.To),
		CreatedAt: timestamppb.New(t.Created),
		Type:      mapTxTypeToPb(t.Type),
	}
}

func mapTxTypeToPb(t gomoney.TransactionType) pb.TransactionType {
	switch t {
	case gomoney.Transfer:
		return pb.TransactionType_TRANSFER
	case gomoney.Deposit:
		return pb.TransactionType_DEPOSIT
	case gomoney.Withdrawal:
		return pb.TransactionType_WITHDRAWAL
	}
	return -1
}
