package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/api/internal/core"
)

type createTransactionReq struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        float64
	Note          string
}

func (server *api.Server) createTransaction(ctx *fiber.Ctx) error {
	var req createTransactionReq
	if err := ctx.BodyParser(&req); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}

	if req.FromAccountID == req.ToAccountID {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "from and to account cannot be the same",
		})
		return application.ErrSimilarAccountTransaction(req.FromAccountID, req.ToAccountID)
	}

	payload := ctx.Locals(payloadHeader).(core.Payload)

	// TODO: Request to gomoney service grpc to verify that both accounts have the same currency
	// If accounts have the same currency, proceed to create transaction
	// If not return an error

	// TODO: Request to gomoney service grpc to verify that the user is the account owner
	// If user is not the account owner, return an error
	// If user is the account owner, continue

	// TODO: Request to gomoney service grpc to create transaction
	// If transaction is created, return transaction's details
	// If error occurs, return an error

	_, _ = req, payload
	return nil
}

type getTransactionsReq struct {
	AccountID int64
}

func (server *api.Server) getTransactions(ctx *fiber.Ctx) error {
	var req getTransactionsReq
	if err := ctx.ParamsParser(&req); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}

	payload := ctx.Locals(payloadHeader).(core.Payload)

	// TODO: Request to gomoney service grpc to verify that the user is the account owner
	// If user is not the account owner, return an error
	// If user is the account owner, continue

	// TODO: Request to gomoney service grpc to get transactions
	// If transactions are found, return transactions' details
	// If error occurs, return an error

	_ = payload
	return nil
}
