package application

import "github.com/gofiber/fiber/v2"

const (
	authorizationHeader = "Authorization"
	payloadHeader       = "payload"
)

func (server *Server) authMiddleware(ctx *fiber.Ctx) error {
	token := ctx.Get(authorizationHeader)
	if token == "" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing token",
		})
		return nil
	}

	payload, err := server.maker.VerifyToken(token)
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
		return nil
	}

	ctx.Locals(payloadHeader, payload)
	return ctx.Next()
}
