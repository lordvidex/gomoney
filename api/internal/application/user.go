package application

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/token"
)

type loginUserReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (s *Server) login(ctx *fiber.Ctx) error {
	var req loginUserReq
	if err := ctx.BodyParser(&req); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}

	// TODO: Request to gomoney service grpc to authenticate user's credentials
	// If credentials are valid, generate a token and return it to the user
	// If credentials are invalid, return an error

	duration := 24 * time.Hour // 1 Day
	token, err := s.maker.CreateToken(req.Phone, duration)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": token,
	})
	return nil
}

func (server *Server) getUser(ctx *fiber.Ctx) error {
	payload := ctx.Locals(payloadHeader).(*token.Payload)

	_ = payload

	// TODO: Request to gomoney service grpc to get user's details
	// If user is found, return user's details
	// If error occurs, return an error

	return nil
}
