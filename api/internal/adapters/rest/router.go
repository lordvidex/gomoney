package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/adapters/rest/handlers"
	"github.com/lordvidex/gomoney/api/internal/application"
)

type router struct {
	app *application.Usecases
	f   *fiber.App
}

type UseCaseHandler func(app *application.Usecases, ctx *fiber.Ctx) error

func (h *router) wrap(uc UseCaseHandler) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return uc(h.app, ctx)
	}
}

func (h *router) setupRoutes() {
	api := h.f.Group("/api")
	auth := h.wrap(handlers.AuthMiddleware)

	// Unauthenticated routes
	api.Post("/login", h.wrap(handlers.Login))
	api.Post("/register", h.wrap(handlers.Register))

	// Authenticated EndPoints

	// - Accounts EndPoint
	api.Get("/accounts", auth, h.wrap(handlers.GetAccounts))
	api.Post("/accounts", auth, h.wrap(handlers.CreateAccount))
	//auth.Get("/account:id", h)

	// - Transactions EndPoint
	api.Post("/transactions/transfer", auth, h.wrap(handlers.CreateTransfers))
	api.Post("/transactions/deposit", auth, h.wrap(handlers.CreateDeposit))
	api.Post("/transactions/withdraw", auth, h.wrap(handlers.CreateWithdraw))
	api.Get("/transactions/:id", auth, h.wrap(handlers.GetAccountTransactions))
	api.Get("/transactions/", auth, h.wrap(handlers.GetTransactions))
}

func (h *router) Listen() error {
	return h.f.Listen(":8080")
}

func New(app *application.Usecases) *router {
	f := fiber.New()
	r := &router{app, f}
	r.setupRoutes()
	return r
}
