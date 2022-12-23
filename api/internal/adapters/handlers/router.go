package handlers

import (
	"github.com/lordvidex/gomoney/api"
	"github.com/lordvidex/gomoney/api/docs"
	"log"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/swagger"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/config"
	"github.com/mvrilo/go-redoc"
)

const (
	mountPrefix = "/gomoney"
)

type Router struct {
	app *application.Usecases
	f   *fiber.App
}

type UseCaseHandler func(app *application.Usecases, ctx *fiber.Ctx) error

func (r *Router) wrap(uc UseCaseHandler) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return uc(r.app, ctx)
	}
}

func (r *Router) setupRoutes() {

	api := r.f.Group("/api")
	auth := r.wrap(AuthMiddleware)

	// Unauthenticated routes
	api.Post("/login", r.wrap(Login))
	api.Post("/register", r.wrap(Register))

	// Authenticated EndPoints

	// - Accounts EndPoint
	api.Get("/accounts", auth, r.wrap(GetAccounts))
	api.Post("/accounts", auth, r.wrap(CreateAccount))

	// - Transactions EndPoint
	api.Post("/transactions/transfer", auth, r.wrap(CreateTransfers))
	api.Post("/transactions/deposit", auth, r.wrap(CreateDeposit))
	api.Post("/transactions/withdraw", auth, r.wrap(CreateWithdraw))

	// Validation is done in the handler
	api.Get("/transactions/:id", auth, r.wrap(GetAccountTransactions))
	api.Get("/transactions/", auth, r.wrap(GetTransactions))
}

func (r *Router) setupSwagger(c *config.Config) {
	docs.SwaggerInfo.Host = c.Get("SWAGGER_HOST")
	if docs.SwaggerInfo.Host == "" {
		docs.SwaggerInfo.Host = "localhost:8000"
	}

	docs.SwaggerInfo.BasePath = c.Get("SWAGGER_BASE_PATH")
	if docs.SwaggerInfo.BasePath == "" {
		docs.SwaggerInfo.BasePath = "/api"
	}

	r.f.Get("/docs/*", func(ctx *fiber.Ctx) error {
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

func (r *Router) setupRedoc() {
	doc := redoc.Redoc{
		SpecPath:    "doc.json",
		Title:       "GoMoni Redoc Documentation",
		Description: `This is the documentation for the GoMoni API server by Redoc`,
	}
	r.f.Get("/redoc/*", func(ctx *fiber.Ctx) error {
		p := utils.CopyString(ctx.Params("*"))
		switch {
		case p == "/" || p == "":
			html, err := doc.Body()
			if err != nil {
				return err
			}
			ctx.Set("content-type", "text/html")
			return ctx.Send(html)
		case strings.HasSuffix(p, "doc.json"):
			docJSON, err := api.Docs.ReadFile("docs/swagger.json")
			if err != nil {
				return err
			}
			ctx.Set("content-type", "application/json")
			ctx.Status(200).Write(docJSON)
		}
		return nil
	})

}

func (r *Router) Listen() error {
	return r.f.Listen(":8080")
}

func New(app *application.Usecases, c *config.Config) *Router {
	f := fiber.New()
	r := &Router{app, f}
	r.setupRoutes()
	r.setupSwagger(c)
	r.setupRedoc()
	if isProd(c) {
		log.Println("RUNNING IN PRODUCTION MODE: MOUNTING", mountPrefix, "all routes without prefix will be redirected to", mountPrefix, "/<route>")
		r.mount(mountPrefix)
	} else {
		log.Println("RUNNING IN DEVELOPMENT MODE: NOT MOUNTING")
	}
	return r
}

func (r *Router) mount(prefix string) {
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
