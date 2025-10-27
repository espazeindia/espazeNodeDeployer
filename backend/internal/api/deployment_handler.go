package api

import (
	"github.com/espazeindia/espazeNodeDeployer/internal/domain/entities"
	"github.com/espazeindia/espazeNodeDeployer/internal/k8s"
	"github.com/espazeindia/espazeNodeDeployer/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetupDeploymentRoutes(router fiber.Router, deploymentUC usecase.DeploymentUseCase, jwtSecret string) {
	deployments := router.Group("/deployments", AuthMiddleware(jwtSecret))

	deployments.Post("/", func(c *fiber.Ctx) error {
		userID := c.Locals("userId").(string)
		userObjID, _ := primitive.ObjectIDFromHex(userID)

		var req entities.DeploymentRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Get node ID from request or use current node
		nodeID, err := primitive.ObjectIDFromHex(c.Query("nodeId"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Node ID is required"})
		}

		// Get GitHub token from header
		githubToken := c.Get("X-GitHub-Token")
		if githubToken == "" {
			return c.Status(400).JSON(fiber.Map{"error": "GitHub token is required"})
		}

		deployment, err := deploymentUC.CreateDeployment(c.Context(), userObjID, nodeID, &req, githubToken)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(201).JSON(deployment)
	})

	deployments.Get("/", func(c *fiber.Ctx) error {
		userID := c.Locals("userId").(string)
		userObjID, _ := primitive.ObjectIDFromHex(userID)

		deployments, err := deploymentUC.GetDeploymentsByUser(c.Context(), userObjID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(deployments)
	})

	deployments.Get("/node/:nodeId", func(c *fiber.Ctx) error {
		nodeID, err := primitive.ObjectIDFromHex(c.Params("nodeId"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid node ID"})
		}

		deployments, err := deploymentUC.GetDeploymentsByNode(c.Context(), nodeID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(deployments)
	})

	deployments.Get("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid deployment ID"})
		}

		deployment, err := deploymentUC.GetDeployment(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(deployment)
	})

	deployments.Put("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid deployment ID"})
		}

		var req entities.DeploymentUpdateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := deploymentUC.UpdateDeployment(c.Context(), id, &req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Deployment updated successfully"})
	})

	deployments.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid deployment ID"})
		}

		// Note: In production, pass the correct k8s client based on node
		// For now, we'll skip this implementation detail
		if err := deploymentUC.DeleteDeployment(c.Context(), id, nil); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Deployment deleted successfully"})
	})

	deployments.Post("/:id/restart", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid deployment ID"})
		}

		if err := deploymentUC.RestartDeployment(c.Context(), id, nil); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Deployment restarted successfully"})
	})

	deployments.Post("/:id/scale", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid deployment ID"})
		}

		var req struct {
			Replicas int32 `json:"replicas"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := deploymentUC.ScaleDeployment(c.Context(), id, req.Replicas, nil); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Deployment scaled successfully"})
	})

	deployments.Get("/stats", func(c *fiber.Ctx) error {
		var nodeID *primitive.ObjectID
		if nodeIDStr := c.Query("nodeId"); nodeIDStr != "" {
			id, err := primitive.ObjectIDFromHex(nodeIDStr)
			if err == nil {
				nodeID = &id
			}
		}

		stats, err := deploymentUC.GetDeploymentStats(c.Context(), nodeID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(stats)
	})
}

