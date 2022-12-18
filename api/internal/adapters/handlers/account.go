package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

// type accountResponse struct {
// 	ID          int64            `json:"id"`
// 	Title       string           `json:"title"`
// 	Description string           `json:"description"`
// 	Currency    gomoney.Currency `json:"currency"`
// 	Balance     float64          `json:"balance"`
// }

// GetAccounts godoc
//
//	@Summary		get all user accounts
//	@Description	returns all the accounts for the currently logged in user
//	@Tags			accounts
//	@Produce		json
//	@Success		200	{object}	response.JSON{data=[]gomoney.Account}
//	@Security		bearerAuth
//	@Router			/accounts [get]
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
	// res := make([]accountResponse, len(accounts))
	// for i, acc := range accounts {
	// 	res[i] = accToResponse(acc)
	// }
	return ctx.Status(fiber.StatusOK).JSON(accounts)
}

type createAccountReq struct {
	Title       string           `json:"title" validate:"required,min=5"`
	Description string           `json:"description" validate:"required,min=5"`
	Currency    gomoney.Currency `json:"currency" validate:"required,oneof=USD RUB NGN"`
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

	// validate body request
	if errs := validateStruct(req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
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

// func accToResponse(acc gomoney.Account) accountResponse {
// 	return accountResponse{
// 		ID:          acc.Id,
// 		Title:       acc.Title,
// 		Description: acc.Description,
// 		Currency:    acc.Currency,
// 		Balance:     acc.Balance,
// 	}
// }
