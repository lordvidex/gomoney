package application

import (
	"context"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

// ------------------ Create a new transfer ------------------

type CreateTransferParam struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type CreateTransferCommand interface {
	Handle(ctx context.Context, transfer CreateTransferParam) (int64, error)
}

type createTransferImpl struct {
	srv  Service
	repo Repository
}

func NewCreateTransferCommand(srv Service, repo Repository) CreateTransferCommand {
	return &createTransferImpl{
		srv: srv, repo: repo,
	}
}

func (c *createTransferImpl) Handle(ctx context.Context, transfer CreateTransferParam) (int64, error) {
	return c.srv.CreateTransfer(ctx, transfer)
}

// ------------------ View transfers for an account ------------------

type ViewAccountTransfersQuery interface {
	Handle(ctx context.Context, accountID int64) ([]gomoney.Transaction, error)
}

type viewAccountTransfersImpl struct {
	srv  Service
	repo Repository
}

func NewViewAccountTransfersQuery(srv Service, repo Repository) ViewAccountTransfersQuery {
	return &viewAccountTransfersImpl{
		srv: srv, repo: repo,
	}
}

func (c *viewAccountTransfersImpl) Handle(ctx context.Context, accountID int64) ([]gomoney.Transaction, error) {
	return c.srv.GetAccountTransfers(ctx, accountID)
}

// ------------------ View transfers for a user ------------------

type ViewTransfersQuery interface {
	Handle(ctx context.Context, userID string) ([]gomoney.Transaction, error)
}

type viewTransfersImpl struct {
	srv  Service
	repo Repository
}

func NewViewTransfersQuery(srv Service, repo Repository) ViewTransfersQuery {
	return &viewTransfersImpl{
		srv: srv, repo: repo,
	}
}

func (c *viewTransfersImpl) Handle(ctx context.Context, userID string) ([]gomoney.Transaction, error) {
	return c.srv.GetTransfers(ctx, userID)
}
