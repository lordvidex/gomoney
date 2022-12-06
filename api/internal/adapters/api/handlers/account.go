package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

func GetAccount(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get the user from the context
	u, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	// get account from repository
	accounts, err := uc.ViewAccounts.Handle(ctx.UserContext(), application.ViewAccountsParam{
		UserID: u.ID,
	})
	if err != nil {
		parseDatabaseInternalError(ctx, err)
		return err
	}
	
	return ctx.Status(fiber.StatusOK).JSON(accounts)
}

func GetAccounts(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get the user from the context
	u, err := userFromCtx(ctx)
	if err != nil {
		return err
	}
	accounts, err := uc.ViewAccounts.Handle(ctx.UserContext(), application.ViewAccountsParam{
		UserID: u.ID,
	})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(accounts)
}

type createAccountReq struct {
	Title       string
	Description string
	Currency    gomoney.Currency
}

func CreateAccount(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get the user from the context
	u, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	var req createAccountReq
	err = parseBody(ctx, &req)
	if err != nil {
		return err
	}
	id, err := uc.CreateAccount.Handle(ctx.UserContext(), application.CreateAccountParam{
		UserID: u.ID,
		Account: gomoney.Account{
			Title:       req.Title,
			Description: req.Description,
			Currency:    req.Currency,
		},
	})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Account successfully created",
		"data": map[string]int64{
			"id": id,
		},
	})
}
