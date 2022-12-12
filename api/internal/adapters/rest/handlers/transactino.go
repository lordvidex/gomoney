package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
)

type createTransfer struct {
	FromAccountID int64   `json:"from_account_id" validate:"required,number,min=1"`
	ToAccountID   int64   `json:"to_account_id" validate:"required,number,min=1"`
	Amount        float64 `json:"amount" validate:"required,number,min=1"`
}

func CreateTransfers(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get auth user from context
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	// Parse request body
	var req createTransfer
	err = parseBody(ctx, &req)
	if err != nil {
		return err
	}

	_, err = uc.Transfer.Handle(ctx.UserContext(), application.CreateTransferParam{
		ActorID: user.ID.String(),
		FromID:  req.FromAccountID,
		ToID:    req.ToAccountID,
		Amount:  req.Amount,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return nil
}

type createDeposit struct {
	ToAccountID int64   `json:"to_account_id" validate:"required,number,min=1"`
	Amount      float64 `json:"amount" validate:"required,number,min=1"`
}

func CreateDeposit(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get auth user from context
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	// Parse request body
	var req createDeposit
	err = parseBody(ctx, &req)
	if err != nil {
		return err
	}

	_, err = uc.Deposit.Handle(ctx.UserContext(), application.DepositParam{
		ActorID: user.ID.String(),
		ToID:    req.ToAccountID,
		Amount:  req.Amount,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return nil
}

type createWithdraw struct {
	FromAccountID int64   `json:"from_account_id" validate:"required,number,min=1"`
	Amount        float64 `json:"amount" validate:"required,number,min=1"`
}

func CreateWithdraw(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get auth user from context
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	// Parse request body
	var req createWithdraw
	if err = parseBody(ctx, &req); err != nil {
		return err
	}

	_, err = uc.Withdraw.Handle(ctx.UserContext(), application.WithdrawParam{
		ActorID: user.ID.String(),
		FromID:  req.FromAccountID,
		Amount:  req.Amount,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return nil
}

func GetTransactions(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get auth user from context
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	transactions, err := uc.GetTransactionsSummary.Handle(ctx.UserContext(), user.ID.String())
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"transactions": transactions})
}

type getAccountTransferParam struct {
	AccountID int64 `uri:"account_id" validate:"required,number,min=1"`
}

func GetAccountTransactions(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get user from context
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	// Parse request body
	var req getAccountTransferParam
	err = parseUri(ctx, &req)
	if err != nil {
		return err
	}

	transactions, err := uc.GetAccountTransactions.Handle(ctx.UserContext(), application.UserWithAccount{
		UserID:    user.ID.String(),
		AccountID: req.AccountID,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"transactions": transactions})
}
