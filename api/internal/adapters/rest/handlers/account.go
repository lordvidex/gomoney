package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

func GetAccounts(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get the user from the context
	u, err := userFromCtx(ctx)
	if err != nil {
		return err
	}
	accounts, err := uc.GetAccounts.Handle(ctx.UserContext(), application.GetAccountsParam{
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
	Title       string           `json:"title" validate:"required"`
	Description string           `json:"description" validate:"required"`
	Currency    gomoney.Currency `json:"currency" validate:"required"`
}

func CreateAccount(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get the user from the context
	u, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	// parse the request body
	var req createAccountReq
	if err = parseBody(ctx, &req); err != nil {
		return err
	}
	if req.Currency.IsValid() == false {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid currency",
		})
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
