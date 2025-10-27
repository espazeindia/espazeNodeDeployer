package api

import (
	"github.com/espaze/espazeNodeDeployer/internal/domain/entities"
	"github.com/espaze/espazeNodeDeployer/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetupNodeRoutes(router fiber.Router, nodeUC usecase.NodeUseCase, jwtSecret string) {
	nodes := router.Group("/nodes")

	nodes.Post("/register", func(c *fiber.Ctx) error {
		var req entities.NodeRegistrationRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		node, err := nodeUC.RegisterNode(c.Context(), &req)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(201).JSON(node)
	})

	nodes.Get("/current", func(c *fiber.Ctx) error {
		nodeInfo, err := nodeUC.GetCurrentNodeInfo()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(nodeInfo)
	})

	nodes.Get("/", AuthMiddleware(jwtSecret), func(c *fiber.Ctx) error {
		filters := make(map[string]interface{})
		
		if status := c.Query("status"); status != "" {
			filters["status"] = status
		}

		nodes, err := nodeUC.GetAllNodes(c.Context(), filters)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(nodes)
	})

	nodes.Get("/:id", AuthMiddleware(jwtSecret), func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid node ID"})
		}

		node, err := nodeUC.GetNode(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(node)
	})

	nodes.Put("/:id", AuthMiddleware(jwtSecret), func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid node ID"})
		}

		var req entities.NodeUpdateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := nodeUC.UpdateNode(c.Context(), id, &req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Node updated successfully"})
	})

	nodes.Delete("/:id", AuthMiddleware(jwtSecret), func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid node ID"})
		}

		if err := nodeUC.DeleteNode(c.Context(), id); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Node deleted successfully"})
	})

	nodes.Post("/:id/heartbeat", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid node ID"})
		}

		if err := nodeUC.Heartbeat(c.Context(), id); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Heartbeat recorded"})
	})

	nodes.Get("/stats", AuthMiddleware(jwtSecret), func(c *fiber.Ctx) error {
		stats, err := nodeUC.GetNodeStats(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(stats)
	})
}

