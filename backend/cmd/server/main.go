package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/espazeindia/espazeNodeDeployer/internal/api"
	"github.com/espazeindia/espazeNodeDeployer/internal/config"
	"github.com/espazeindia/espazeNodeDeployer/internal/github"
	"github.com/espazeindia/espazeNodeDeployer/internal/k8s"
	"github.com/espazeindia/espazeNodeDeployer/internal/repository"
	"github.com/espazeindia/espazeNodeDeployer/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Ping MongoDB to verify connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("✅ Connected to MongoDB successfully")

	db := mongoClient.Database(cfg.DatabaseName)

	// Initialize Kubernetes client
	k8sClient, err := k8s.NewClient(cfg.KubeConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}
	log.Println("✅ Connected to Kubernetes cluster successfully")

	// Initialize GitHub client
	githubClient := github.NewClient(cfg.GitHubClientID, cfg.GitHubClientSecret)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	deploymentRepo := repository.NewDeploymentRepository(db)
	githubTokenRepo := repository.NewGitHubTokenRepository(db)
	nodeRepo := repository.NewNodeRepository(db)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)
	nodeUseCase := usecase.NewNodeUseCase(nodeRepo)
	deploymentUseCase := usecase.NewDeploymentUseCase(deploymentRepo, k8sClient, githubClient, githubTokenRepo)
	githubUseCase := usecase.NewGitHubUseCase(githubClient, githubTokenRepo)
	k8sUseCase := usecase.NewK8sUseCase(k8sClient)
	metricsUseCase := usecase.NewMetricsUseCase(k8sClient)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:           "Espaze Node Deployer API",
		EnablePrintRoutes: cfg.Env == "development",
		ErrorHandler:      customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// API v1 routes
	apiV1 := app.Group("/api/v1")

	// Initialize handlers
	api.SetupAuthRoutes(apiV1, authUseCase, cfg.JWTSecret)
	api.SetupNodeRoutes(apiV1, nodeUseCase, cfg.JWTSecret)
	api.SetupGitHubRoutes(apiV1, githubUseCase, cfg.JWTSecret)
	api.SetupDeploymentRoutes(apiV1, deploymentUseCase, cfg.JWTSecret)
	api.SetupK8sRoutes(apiV1, k8sUseCase, cfg.JWTSecret)
	api.SetupMetricsRoutes(apiV1, metricsUseCase, cfg.JWTSecret)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "The requested resource was not found",
		})
	})

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("🚀 Server starting on http://localhost%s\n", addr)
		log.Printf("📚 API Documentation: http://localhost%s/api/v1\n", addr)
		log.Printf("🏥 Health Check: http://localhost%s/health\n", addr)
		
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   message,
		"status":  code,
		"path":    c.Path(),
		"method":  c.Method(),
		"time":    time.Now(),
	})
}

