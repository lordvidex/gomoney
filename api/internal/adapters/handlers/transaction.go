package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/adapters/handlers/response"
	"github.com/lordvidex/gomoney/api/internal/application"
)

type TransactionDTO struct {
	ID        *string     `json:"id"`
	To        *AccountDTO `json:"to"`
	From      *AccountDTO `json:"from"`
	Amount    *float64    `json:"amount"`
	CreatedAt *string     `json:"created_at"`
	Type      *string     `json:"type"`
}

type TransactionSummaryDTO struct {
	Account     *AccountDTO      `json:"account"`
	Transaction []TransactionDTO `json:"transaction"`
}

type createTransfer struct {
	FromAccountID int64   `json:"from_account_id" validate:"required,number,min=1"`
	ToAccountID   int64   `json:"to_account_id" validate:"required,number,min=1,nefield=FromAccountID"`
	Amount        float64 `json:"amount" validate:"required,number,min=1"`
}

// CreateTransfers godoc
//
//	@Summary		transfer a specified amount from a user's account to another account
//	@Description	transfer a specified amount from a user's account to another account
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createTransfer	true	"create transfer request"
//	@Success		200		{object}	response.JSON{data=TransactionDTO}
//	@Failure		400,500	{object}	response.JSON{error=[]response.Error}
//	@Security		bearerAuth
//	@Router			/transactions/transfer [post]
func CreateTransfers(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get auth user from context
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	// Parse request body
	var req createTransfer
	if err = parseBody(ctx, &req); err != nil {
		return err
	}

	// validate body request
	if errs := validateStruct(req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	tx, err := uc.Transfer.Handle(ctx.UserContext(), application.CreateTransferParam{
		ActorID: user.ID.String(),
		FromID:  req.FromAccountID,
		ToID:    req.ToAccountID,
		Amount:  req.Amount,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(response.Success(parseTransaction(tx)))
}

type createDeposit struct {
	ToAccountID int64   `json:"to_account_id" validate:"required,number,min=1"`
	Amount      float64 `json:"amount" validate:"required,number,min=1"`
}

// CreateDeposit godoc
//
//	@Summary		deposit a specified amount to a user's account
//	@Description	deposit a specified amount to a user's account
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createDeposit	true	"create deposit request"
//	@Success		200		{object}	response.JSON{data=TransactionDTO}
//	@Failure		400,500	{object}	response.JSON{error=[]response.Error}
//	@Security		bearerAuth
//	@Router			/transactions/deposit [post]
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

	// validate body request
	if errs := validateStruct(req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	tx, err := uc.Deposit.Handle(ctx.UserContext(), application.DepositParam{
		ActorID: user.ID.String(), ToID: req.ToAccountID, Amount: req.Amount,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(response.Success(parseTransaction(tx)))
}

type createWithdraw struct {
	FromAccountID int64   `json:"from_account_id" validate:"required,number,min=1"`
	Amount        float64 `json:"amount" validate:"required,number,min=1"`
}

// CreateWithdraw godoc
//
//	@Summary		withdraw a specified amount from a user's account
//	@Description	withdraw a specified amount from a user's account
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createWithdraw	true	"create withdraw request"
//	@Success		200		{object}	response.JSON{data=TransactionDTO}
//	@Failure		400,500	{object}	response.JSON{error=[]response.Error}
//	@Security		bearerAuth
//	@Router			/transactions/withdraw [post]
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

	// validate body request
	if errs := validateStruct(req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	tx, err := uc.Withdraw.Handle(ctx.UserContext(), application.WithdrawParam{
		ActorID: user.ID.String(), FromID: req.FromAccountID, Amount: req.Amount,
	})

	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.Success(parseTransaction(tx)))
}

// GetTransactions  godoc
//
//	@Summary		get all accounts transactions for the logged-in user
//	@Description	get all accounts transactions for the logged-in user
//	@Tags			transactions
//	@Produce		json
//	@Success		200		{object}	response.JSON{data=[]TransactionSummaryDTO}
//	@Failure		400,500	{object}	response.JSON{error=[]response.Error}
//	@Security		bearerAuth
//	@Router			/transactions [get]
func GetTransactions(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get auth user from context
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	txs, err := uc.GetTransactionsSummary.Handle(ctx.UserContext(), user.ID.String())
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	txsDTO := make([]TransactionSummaryDTO, len(txs))
	for i, tx := range txs {
		txsDTO[i] = parseTransactionSummary(&tx)
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success(txsDTO))
}

type getAccountTransferParam struct {
	AccountID int64 `params:"id" validate:"required,number,min=1"`
}

// GetAccountTransactions	godoc
//
//	@Summary		get one account transactions for the logged-in user by account id
//	@Description	get one account transactions for the logged-in user by account id
//	@Tags			transactions
//	@Produce		json
//	@Param			id		path		int64	true	"account id"
//	@Success		200		{object}	response.JSON{data=TransactionSummaryDTO}
//	@Failure		400,500	{object}	response.JSON{error=[]response.Error}
//	@Security		bearerAuth
//	@Router			/transactions/{id} [get]
func GetAccountTransactions(uc *application.Usecases, ctx *fiber.Ctx) error {
	// get user from context
	user, err := userFromCtx(ctx)
	if err != nil {
		return err
	}

	// Parse request body
	var req getAccountTransferParam
	err = parseParams(ctx, &req)
	if err != nil {
		return err
	}

	// validate body request
	if errs := validateStruct(req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Errs(errs...))
	}

	transactions, err := uc.GetAccountTransactions.Handle(ctx.UserContext(), application.UserWithAccount{
		UserID: user.ID.String(), AccountID: req.AccountID,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.Success(parseTransactionSummary(&transactions)))
}
