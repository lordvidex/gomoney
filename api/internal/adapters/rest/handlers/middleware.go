package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
)

const (
	authorizationHeader = "Authorization"
	// AuthUserPayload is used to encode and decode a *core.ApiUser type
	AuthUserPayload = "authUser"
)

func AuthMiddleware(uc *application.Usecases, ctx *fiber.Ctx) error {
	token := ctx.Get(authorizationHeader)
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing token",
		})
	}

	authUser, err := uc.GetAPIUser.Handle(ctx.UserContext(), token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx.Locals(AuthUserPayload, authUser)
	return ctx.Next()
}
