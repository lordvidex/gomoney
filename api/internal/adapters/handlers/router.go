package handlers

import (
	"log"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/swagger"
	"github.com/lordvidex/gomoney/api/docs"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/config"
)

const (
	mountPrefix = "/gomoney"
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
	api.Post("/accounts/:id", auth, h.wrap(DeleteAccount))

	// - Transactions EndPoint
	api.Post("/transactions/transfer", auth, h.wrap(CreateTransfers))
	api.Post("/transactions/deposit", auth, h.wrap(CreateDeposit))
	api.Post("/transactions/withdraw", auth, h.wrap(CreateWithdraw))

	// Validation is done in the handler
	api.Get("/transactions/:id", auth, h.wrap(GetAccountTransactions))
	api.Get("/transactions/", auth, h.wrap(GetTransactions))
}

func (h *router) setupSwagger(c *config.Config) {
	docs.SwaggerInfo.Host = c.Get("SWAGGER_HOST")
	if docs.SwaggerInfo.Host == "" {
		docs.SwaggerInfo.Host = "localhost:8000"
	}

	docs.SwaggerInfo.BasePath = c.Get("SWAGGER_BASE_PATH")
	if docs.SwaggerInfo.BasePath == "" {
		docs.SwaggerInfo.BasePath = "/gomoney/api"
	}
	h.f.Get("/docs/*", func(ctx *fiber.Ctx) error {
		p := utils.CopyString(ctx.Params("*"))
		switch p {
		case "/", "":
			prefix := strings.ReplaceAll(ctx.Route().Path, "*", "")
			rr := path.Join(prefix, "index.html")
			log.Println("Redirecting to ", rr)
			return ctx.Redirect(rr)
		default:
			return swagger.HandlerDefault(ctx)
		}
	}) // documentations
}

func (h *router) Listen() error {
	return h.f.Listen(":8080")
}

func New(app *application.Usecases, c *config.Config) *router {
	f := fiber.New()
	r := &router{app, f}
	r.setupRoutes()
	r.setupSwagger(c)
	if isProd(c) {
		log.Println("RUNNING IN PRODUCTION MODE: MOUNTING", mountPrefix, "all routes without prefix will be redirected to", mountPrefix, "/<route>")
		r.mount(mountPrefix)
	} else {
		log.Println("RUNNING IN DEVELOPMENT MODE: NOT MOUNTING")
	}
	return r
}

func (r *router) mount(prefix string) {
	if prefix == "" {
		return
	}
	mnt := fiber.New()
	// redirect all requests without prefix to the prefix
	mnt.Use(func(ctx *fiber.Ctx) error {
		if strings.HasPrefix(ctx.OriginalURL(), prefix) {
			return ctx.Next()
		}
		return ctx.Redirect(prefix + ctx.OriginalURL())
	})
	mnt.Mount(prefix, r.f)
	r.f = mnt
}

func isProd(c *config.Config) bool {
	return strings.ToLower(c.Get("APP_ENV")) == "production"
}
