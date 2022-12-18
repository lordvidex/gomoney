package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/adapters/handlers/response"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

func userFromCtx(ctx *fiber.Ctx) (*core.ApiUser, error) {
	u, ok := ctx.Locals(AuthUserPayload).(*core.ApiUser)
	if !ok {
		return nil, ctx.Status(fiber.StatusUnauthorized).
			JSON(response.ErrM("User is unauthenticated"))
	}
	return u, nil
}

func parseBody(ctx *fiber.Ctx, obj interface{}) error {
	if err := ctx.BodyParser(obj); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).
			JSON(response.ErrM(err.Error()))
		return err
	}
	return nil
}

func parseParams(ctx *fiber.Ctx, obj interface{}) error {
	if err := ctx.ParamsParser(obj); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).
			JSON(response.ErrM(err.Error()))
		return err
	}
	return nil
}

func setCtxBodyError(ctx *fiber.Ctx, err error) error {
	ge := func(ge *gomoney.Error) error {
		msgs := make([]response.Error, len(ge.Messages))
		for i, m := range ge.Messages {
			msgs[i] = response.Err(m, ge.Code)
		}
		return ctx.Status(response.C(ge.Code)).JSON(response.Errs(msgs...))
	}
	switch err := err.(type) {
	case *gomoney.Error:
		return ge(err)
	case gomoney.Error:
		return ge(&err)
	}
	return ctx.Status(fiber.StatusInternalServerError).
		JSON(response.ErrM(err.Error()))
}

func strptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func parseUser(u *gomoney.User) UserDTO {
	return UserDTO{
		ID:    strptr(u.ID.String()),
		Name:  strptr(u.Name),
		Phone: strptr(u.Phone),
	}
}
