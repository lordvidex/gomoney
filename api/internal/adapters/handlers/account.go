package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/adapters/handlers/response"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type AccountDTO struct {
	ID          *int64   `json:"id"`
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Currency    *string  `json:"currency"`
	Balance     *float64 `json:"balance"`
	IsBlocked   *bool    `json:"is_blocked"`
}

// GetAccounts godoc
//
//	@Summary		get all user accounts
//	@Description	returns all the accounts for the currently logged-in user
//	@Tags			accounts
//	@Produce		json
//	@Success		200		{object}	response.JSON{data=[]AccountDTO}
//	@Failure		400,500	{object}	response.JSON{error=[]response.Error}
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
		return setCtxBodyError(ctx, err)
	}
	res := make([]*AccountDTO, len(accounts))
	for i, acc := range accounts {
		res[i] = parseAccount(&acc)
	}
	return ctx.Status(fiber.StatusOK).JSON(accounts)
}

type createAccountReq struct {
	Title       string           `json:"title" validate:"required,min=5"`
	Description string           `json:"description" validate:"required,min=5"`
	Currency    gomoney.Currency `json:"currency" validate:"required,oneof=USD RUB NGN"`
}

type createAccountRes struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

// CreateAccount godoc
//
//	@Summary		creates a new account for the currently logged-in user
//	@Description	creates a new account for the currently logged-in user
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createAccountReq	true	"login user request"
//	@Success		200		{object}	response.JSON{data=createAccountRes}
//	@Failure		400,500	{object}	response.JSON{error=[]response.Error}
//	@Security		bearerAuth
//	@Router			/accounts [post]
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

	// call create account service
	id, err := uc.CreateAccount.Handle(ctx.UserContext(), application.CreateAccountParam{
		UserID: u.ID, Title: req.Title, Description: req.Description, Currency: req.Currency,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.Success(createAccountRes{ID: id, Message: "Account created successfully"}))
}
