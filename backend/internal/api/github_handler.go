package api

import (
	"strconv"

	"github.com/espazeindia/espazeNodeDeployer/internal/domain/entities"
	"github.com/espazeindia/espazeNodeDeployer/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetupGitHubRoutes(router fiber.Router, githubUC usecase.GitHubUseCase, jwtSecret string) {
	github := router.Group("/github", AuthMiddleware(jwtSecret))

	github.Post("/token", func(c *fiber.Ctx) error {
		userID := c.Locals("userId").(string)
		userObjID, _ := primitive.ObjectIDFromHex(userID)

		var req entities.GitHubTokenRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := githubUC.SaveToken(c.Context(), userObjID, req.Token); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "GitHub token saved successfully"})
	})

	github.Get("/user", func(c *fiber.Ctx) error {
		userID := c.Locals("userId").(string)
		userObjID, _ := primitive.ObjectIDFromHex(userID)

		user, err := githubUC.GetUser(c.Context(), userObjID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(user)
	})

	github.Get("/repos", func(c *fiber.Ctx) error {
		userID := c.Locals("userId").(string)
		userObjID, _ := primitive.ObjectIDFromHex(userID)

		page, _ := strconv.Atoi(c.Query("page", "1"))
		perPage, _ := strconv.Atoi(c.Query("perPage", "30"))

		repos, pagination, err := githubUC.GetRepositories(c.Context(), userObjID, page, perPage)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"repositories": repos,
			"pagination":   pagination,
		})
	})

	github.Get("/repos/:owner/:repo", func(c *fiber.Ctx) error {
		userID := c.Locals("userId").(string)
		userObjID, _ := primitive.ObjectIDFromHex(userID)

		owner := c.Params("owner")
		repo := c.Params("repo")

		repository, err := githubUC.GetRepository(c.Context(), userObjID, owner, repo)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(repository)
	})

	github.Get("/repos/:owner/:repo/branches", func(c *fiber.Ctx) error {
		userID := c.Locals("userId").(string)
		userObjID, _ := primitive.ObjectIDFromHex(userID)

		owner := c.Params("owner")
		repo := c.Params("repo")

		branches, err := githubUC.GetBranches(c.Context(), userObjID, owner, repo)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(branches)
	})

	github.Get("/search", func(c *fiber.Ctx) error {
		userID := c.Locals("userId").(string)
		userObjID, _ := primitive.ObjectIDFromHex(userID)

		query := c.Query("q")
		if query == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Query parameter 'q' is required"})
		}

		page, _ := strconv.Atoi(c.Query("page", "1"))
		perPage, _ := strconv.Atoi(c.Query("perPage", "30"))

		repos, err := githubUC.SearchRepositories(c.Context(), userObjID, query, page, perPage)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(repos)
	})
}

