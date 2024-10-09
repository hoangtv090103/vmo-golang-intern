package middleware

import (
	"ecommerce/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenStr == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		claim, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Store the claims in the context
		c.Locals("claims", claim)

		// If the token is valid, call Next() to continue to the next middleware/handler
		return c.Next()
	}
}

func IsAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claim := c.Locals("claims").(*utils.Claims)

		if claim.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		return c.Next()
	}
}

func IsUserMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claim := c.Locals("claims").(*utils.Claims)

		if claim.Role != "user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		return c.Next()
	}
}
