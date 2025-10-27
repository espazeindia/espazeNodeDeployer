package api

import (
	"github.com/espazeindia/espazeNodeDeployer/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupMetricsRoutes(router fiber.Router, metricsUC usecase.MetricsUseCase, jwtSecret string) {
	metrics := router.Group("/metrics", AuthMiddleware(jwtSecret))

	metrics.Get("/pods", func(c *fiber.Ctx) error {
		namespace := c.Query("namespace", "espaze-node-deployer-apps")

		podMetrics, err := metricsUC.GetPodMetrics(c.Context(), namespace)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(podMetrics)
	})

	metrics.Get("/cluster", func(c *fiber.Ctx) error {
		clusterMetrics, err := metricsUC.GetClusterMetrics(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(clusterMetrics)
	})

	metrics.Get("/deployments/:namespace/:name", func(c *fiber.Ctx) error {
		namespace := c.Params("namespace")
		name := c.Params("name")

		deploymentMetrics, err := metricsUC.GetDeploymentMetrics(c.Context(), namespace, name)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(deploymentMetrics)
	})
}

