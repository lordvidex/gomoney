package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/lordvidex/gomoney/api/docs"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/config"
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
	auth := h.wrap(AuthMiddleware)

	// Unauthenticated routes
	api.Post("/login", h.wrap(Login))
	api.Post("/register", h.wrap(Register))

	// Authenticated EndPoints

	// - Accounts EndPoint
	api.Get("/accounts", auth, h.wrap(GetAccounts))
	api.Post("/accounts", auth, h.wrap(CreateAccount))

	// - Transactions EndPoint
	api.Post("/transactions/transfer", auth, h.wrap(CreateTransfers))
	api.Post("/transactions/deposit", auth, h.wrap(CreateDeposit))
	api.Post("/transactions/withdraw", auth, h.wrap(CreateWithdraw))

	// Validation is done in the handler
	api.Get("/transactions/:id", auth, h.wrap(GetAccountTransactions))
	api.Get("/transactions/", auth, h.wrap(GetTransactions))
}

func (h *router) setupSwagger() {
	config := config.New()

	docs.SwaggerInfo.Host = config.Get("SWAGGER_HOST")
	if docs.SwaggerInfo.Host == "" {
		docs.SwaggerInfo.Host = "localhost:8000"
	}

	docs.SwaggerInfo.BasePath = config.Get("SWAGGER_BASE_PATH")
	if docs.SwaggerInfo.BasePath == "" {
		docs.SwaggerInfo.BasePath = "/api"
	}

	h.f.Get("/docs/*", swagger.HandlerDefault) // documentations
}

func (h *router) Listen() error {
	return h.f.Listen(":8080")
}

func New(app *application.Usecases) *router {
	f := fiber.New()
	r := &router{app, f}
	r.setupRoutes()
	r.setupSwagger()
	return r
}
