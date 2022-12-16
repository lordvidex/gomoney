package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/application"
)

type loginUserReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func Login(uc *application.Usecases, ctx *fiber.Ctx) error {
	// parse request body
	var req loginUserReq
	if err := parseBody(ctx, &req); err != nil {
		return err
	}

	// call app function
	u, err := uc.Login.Handle(ctx.UserContext(),
		application.LoginParam{Phone: req.Phone, Password: req.Password})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  u.User,
		"token": u.Token,
	})
}

type createUserReq struct {
	Phone    string `json:"phone" validate:"required,number"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,alpha"`
}

func Register(uc *application.Usecases, ctx *fiber.Ctx) error {
	var req createUserReq
	if err := parseBody(ctx, &req); err != nil {
		return err
	}

	_, err := uc.CreateUser.Handle(ctx.UserContext(), application.CreateUserParam{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User successfully created",
	})
}
