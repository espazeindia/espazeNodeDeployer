package api

import (
	"strings"

	"github.com/espazeindia/espazeNodeDeployer/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		// Extract token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization format"})
		}

		token := tokenParts[1]

		// Validate token
		claims, err := auth.ValidateToken(token, jwtSecret)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Store claims in context
		c.Locals("user", claims)
		c.Locals("userId", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

