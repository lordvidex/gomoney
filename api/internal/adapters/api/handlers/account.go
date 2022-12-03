package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type createAccountReq struct {
	Title       string
	Description string
	Currency    gomoney.Currency
}

func userFromCtx(ctx *fiber.Ctx) (*core.ApiUser, error) {
	u, ok := ctx.Locals(AuthUserPayload).(*core.ApiUser)
	if !ok {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User is unauthorized",
		})
	}
	return u, nil
}

func parseBody(ctx *fiber.Ctx, obj interface{}) error {
	if err := ctx.BodyParser(obj); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}
	return nil
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
