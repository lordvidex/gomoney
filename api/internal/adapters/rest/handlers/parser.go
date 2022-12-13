package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/core"
)

func userFromCtx(ctx *fiber.Ctx) (*core.ApiUser, error) {
	u, ok := ctx.Locals(AuthUserPayload).(*core.ApiUser)
	if !ok {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User is unauthorized",
		})
	}
	return u, nil
}

func parseBody(ctx *fiber.Ctx, obj interface{}) error {
	if err := ctx.BodyParser(obj); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}
	return nil
}

func parseUri(ctx *fiber.Ctx, obj interface{}) error {
	if err := ctx.ParamsParser(&obj); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}
	return nil
}

func setCtxBodyError(ctx *fiber.Ctx, err error) error {
	return ctx.JSON(fiber.Map{
		"error": err.Error(),
	})
}
