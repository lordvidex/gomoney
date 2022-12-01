package application

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/token"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type createAccountReq struct {
	Title       string
	Description string
	Currency    gomoney.Currency
}

func (server *Server) createAccount(ctx *fiber.Ctx) error {
	var req createAccountReq
	if err := ctx.BodyParser(&req); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}

	payload := ctx.Locals(payloadHeader).(token.Payload)

	_, _ = req, payload
	// TODO: Request to gomoney service grpc to create account
	// If account is created, return account's details
	// If error occurs, return an error

	return nil
}

func (server *Server) getAccounts(ctx *fiber.Ctx) error {

	payload := ctx.Locals(payloadHeader).(token.Payload)

	_ = payload
	// TODO: Request to gomoney service grpc to get all accounts
	// If account is created, return account's details
	// If error occurs, return an error

	return nil
}
