package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/password"
)

type loginUserReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type createUserReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func Login(uc *application.Usecases, ctx *fiber.Ctx) error {
	// parse dto
	var req loginUserReq
	if err := parseBody(ctx, &req); err != nil {
		return err
	}
	// call app function
	u, err := uc.Login.Handle(ctx.UserContext(),
		application.LoginParam{Phone: req.Phone, Password: req.Password})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": u.Token,
	})
}

func Register(uc *application.Usecases, ctx *fiber.Ctx) error {
	var req createUserReq
	if err := parseBody(ctx, &req); err != nil {
		return err
	}

	hashPassword, err := password.CreatePasswordHash(req.Password)
	if err != nil {
		return err
	}

	_, err = uc.CreateUser.Handle(ctx.UserContext(), application.CreateUserParam{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: hashPassword,
	})

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User successfully created",
	})
}
