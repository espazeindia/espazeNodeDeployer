package api

import (
	"strconv"

	"github.com/espazeindia/espazeNodeDeployer/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupK8sRoutes(router fiber.Router, k8sUC usecase.K8sUseCase, jwtSecret string) {
	k8s := router.Group("/k8s", AuthMiddleware(jwtSecret))

	k8s.Get("/cluster/info", func(c *fiber.Ctx) error {
		info, err := k8sUC.GetClusterInfo(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(info)
	})

	k8s.Get("/namespaces", func(c *fiber.Ctx) error {
		namespaces, err := k8sUC.GetNamespaces(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(namespaces)
	})

	k8s.Get("/pods", func(c *fiber.Ctx) error {
		namespace := c.Query("namespace", "espaze-node-deployer-apps")

		pods, err := k8sUC.GetPods(c.Context(), namespace)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(pods)
	})

	k8s.Get("/pods/:namespace/:name", func(c *fiber.Ctx) error {
		namespace := c.Params("namespace")
		name := c.Params("name")

		pod, err := k8sUC.GetPod(c.Context(), namespace, name)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(pod)
	})

	k8s.Get("/pods/:namespace/:name/logs", func(c *fiber.Ctx) error {
		namespace := c.Params("namespace")
		name := c.Params("name")
		tailLines, _ := strconv.ParseInt(c.Query("tail", "100"), 10, 64)

		logs, err := k8sUC.GetPodLogs(c.Context(), namespace, name, tailLines)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"logs": logs})
	})

	k8s.Get("/services", func(c *fiber.Ctx) error {
		namespace := c.Query("namespace", "espaze-node-deployer-apps")

		services, err := k8sUC.GetServices(c.Context(), namespace)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(services)
	})

	k8s.Get("/nodes", func(c *fiber.Ctx) error {
		nodes, err := k8sUC.GetNodes(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(nodes)
	})

	k8s.Get("/events", func(c *fiber.Ctx) error {
		namespace := c.Query("namespace", "espaze-node-deployer-apps")

		events, err := k8sUC.GetEvents(c.Context(), namespace)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(events)
	})
}

