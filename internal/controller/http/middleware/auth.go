package middleware

import (
	"strings"

	"github.com/andreyxaxa/auth-svc/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

func Auth(uc usecase.Sessions) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing 'Authorization' header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid 'Authorization' header format")
		}

		userID, err := uc.Me(ctx.Context(), parts[1])
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}

		ctx.Locals("user_id", userID)

		return ctx.Next()
	}
}
