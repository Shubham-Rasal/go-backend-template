package api

import (
	"fmt"

	"github.com/Shubham-Rasal/blog-backend/token"
	"github.com/gofiber/fiber/v2"
)

func Protected(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(errorResponse(fmt.Errorf("authorization header missing")))
		}

		tokenHeader := authHeader[len("Bearer "):]
		payload, err := tokenMaker.VerifyToken(tokenHeader)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(errorResponse(fmt.Errorf("invalid token : %w", err)))
		}

		ctx.Locals("username", payload.Username)
		return ctx.Next()
	}
}
