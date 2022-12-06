package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
)

type createTransfer struct {
	FromAccountID int64 `json:"from_account_id" validate:"required,number,min=1"`
	ToAccountID   int64 `json:"to_account_id" validate:"required,number,min=1"`
	Amount        int64 `json:"amount" validate:"required,number,min=1"`
}

func CreateTransfers(uc *application.Usecases, ctx *fiber.Ctx) error {
	user, err := userFromCtx(ctx)

	var req createTransfer
	err = parseBody(ctx, &req)
	if err != nil {
		return err
	}

	_, err = uc.CreateTransfer.Handle(ctx.UserContext(), application.CreateTransferParam{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})

	_ = user
	if err != nil {
		return nil
	}
	return nil
}

func GetTransfers(uc *application.Usecases, ctx *fiber.Ctx) error {
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	transactions, err := uc.ViewTransfers.Handle(ctx.UserContext(), user.ID.String())
	if err != nil {
		_ = ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}

	_ = ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"transactions": transactions,
	})

	return nil
}

type getAccountTransferParam struct {
	AccountID int64 `uri:"account_id" validate:"required,number,min=1"`
}

func GetAccountTransfers(uc *application.Usecases, ctx *fiber.Ctx) error {
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	var req getAccountTransferParam
	err = parseUri(ctx, &req)
	if err != nil {
		return err
	}

	// TODO: check if user owns the account
	//account, err := uc.ViewAccount.Handle(ctx.UserContext(), application.ViewAccountParam{
	//	AccountID: req.AccountID,
	//})
	//if err != nil {
	//	parseDatabaseInternalError(ctx, err)
	//}

	transactions, err := uc.ViewTransfers.Handle(ctx.UserContext(), user.ID.String())
	if err != nil {
		parseDatabaseInternalError(ctx, err)
		return err
	}

	_ = ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"transactions": transactions,
	})

	return nil
}
