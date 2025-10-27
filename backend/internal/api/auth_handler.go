package api

import (
	"github.com/espazeindia/espazeNodeDeployer/internal/domain/entities"
	"github.com/espazeindia/espazeNodeDeployer/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router, authUC usecase.AuthUseCase, jwtSecret string) {
	auth := router.Group("/auth")

	auth.Post("/register", func(c *fiber.Ctx) error {
		var req entities.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		user, err := authUC.Register(c.Context(), &req)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(201).JSON(fiber.Map{
			"message": "User registered successfully",
			"user":    user,
		})
	})

	auth.Post("/login", func(c *fiber.Ctx) error {
		var req entities.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		response, err := authUC.Login(c.Context(), &req)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(response)
	})

	auth.Get("/validate", AuthMiddleware(jwtSecret), func(c *fiber.Ctx) error {
		user := c.Locals("user")
		return c.JSON(fiber.Map{
			"valid": true,
			"user":  user,
		})
	})
}

