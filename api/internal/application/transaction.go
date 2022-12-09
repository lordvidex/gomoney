package application

import (
	"context"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

// ------------------ Create a new transfer transactions ------------------

type CreateTransferParam struct {
	Amount  float64
	ActorID string
	FromID  int64
	ToID    int64
}

type TransferCommand interface {
	Handle(ctx context.Context, transfer CreateTransferParam) (*gomoney.Transaction, error)
}

type createTransferImpl struct {
	srv Service
}

func NewTransferCommand(srv Service) TransferCommand {
	return &createTransferImpl{srv}
}

func (c *createTransferImpl) Handle(ctx context.Context, transfer CreateTransferParam) (*gomoney.Transaction, error) {
	return c.srv.Transfer(ctx, transfer)
}

// ------------------ Create a new deposit transactions ------------------

type DepositParam struct {
	Amount  float64
	ActorID string
	ToID    int64
}

type DepositCommand interface {
	Handle(ctx context.Context, transfer DepositParam) (*gomoney.Transaction, error)
}

type depositImpl struct {
	srv Service
}

func NewDepositCommand(srv Service) DepositCommand {
	return &depositImpl{srv}
}

func (c *depositImpl) Handle(ctx context.Context, param DepositParam) (*gomoney.Transaction, error) {
	return c.srv.Deposit(ctx, param)
}

// ------------------ Create a new withdraw transactions ------------------

type WithdrawParam struct {
	Amount  float64
	ActorID string
	FromID  int64
}

type WithdrawCommand interface {
	Handle(ctx context.Context, param WithdrawParam) (*gomoney.Transaction, error)
}

type withdrawImpl struct {
	srv Service
}

func NewWithdrawCommand(srv Service) WithdrawCommand {
	return &withdrawImpl{srv}
}

func (c *withdrawImpl) Handle(ctx context.Context, param WithdrawParam) (*gomoney.Transaction, error) {
	return c.srv.Withdraw(ctx, param)
}

// ------------------ get transactions for a user's account ------------------

type UserWithAccount struct {
	UserID    string
	AccountID int64
}

type GetTransactionQuery interface {
	Handle(ctx context.Context, param UserWithAccount) (gomoney.TransactionSummary, error)
}

type getTransactionsImpl struct {
	srv Service
}

func NewGetTransactionsQuery(srv Service) GetTransactionQuery {
	return &getTransactionsImpl{srv}
}

func (c *getTransactionsImpl) Handle(ctx context.Context, param UserWithAccount) (gomoney.TransactionSummary, error) {
	return c.srv.GetTransactions(ctx, param)
}

// ------------------ get transactions summary for a user ------------------

type GetTransactionSummaryQuery interface {
	Handle(ctx context.Context, userID string) ([]gomoney.TransactionSummary, error)
}

type getTransactionSummaryImpl struct {
	srv Service
}

func NewGetTransactionsSummaryQuery(srv Service) GetTransactionSummaryQuery {
	return &getTransactionSummaryImpl{srv}
}

func (c *getTransactionSummaryImpl) Handle(ctx context.Context, userID string) ([]gomoney.TransactionSummary, error) {
	return c.srv.GetTransactionSummary(ctx, userID)
}
