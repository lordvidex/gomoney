package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/adapters/handlers/response"
	"github.com/lordvidex/gomoney/api/internal/application"
)

type loginUserReq struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginUserRes struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type UserDTO struct {
	ID    *string `json:"id"`
	Name  *string `json:"name"`
	Phone *string `json:"phone"`
}

// Login godoc
//
//	@Summary		login with phone and password
//	@Description	login with phone and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		loginUserReq	true	"login user request"
//	@Success		200		{object}	response.JSON{data=loginUserRes}
//	@Failure		400,500	{object}	response.JSON{error=[]response.Error}
//	@Router			/login [post]
func Login(uc *application.Usecases, ctx *fiber.Ctx) error {
	// parse request body
	var req loginUserReq
	if err := parseBody(ctx, &req); err != nil {
		return err
	}

	// validate body request
	if errs := validateStruct(req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Errs(errs...))
	}

	// call app function
	u, err := uc.Login.Handle(ctx.UserContext(),
		application.LoginParam{Phone: req.Phone, Password: req.Password})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(
		response.Success(loginUserRes{
			User:  parseUser(&u.User),
			Token: u.Token,
		}))
}

type createUserReq struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
} //	@name	a.createUserReq

type createUserRes struct {
	ID      string `json:"user_id"`
	Message string `json:"message"`
}

// Register documentation
//
//	@Summary	register a new user
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		createUserReq	true	"create user request"
//	@Success	201		{object}	response.JSON{data=createUserRes}
//	@Failure	400,500	{object}	response.JSON{error=[]response.Error}
//	@Router		/register [post]
func Register(uc *application.Usecases, ctx *fiber.Ctx) error {
	// parse request body
	var req createUserReq
	if err := parseBody(ctx, &req); err != nil {
		return err
	}
	// validate body request
	if errs := validateStruct(req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Errs(errs...))
	}

	// call create user service
	uid, err := uc.CreateUser.Handle(ctx.UserContext(), application.CreateUserParam{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return setCtxBodyError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(response.Success(createUserRes{uid, "User successfully created"}))
}
