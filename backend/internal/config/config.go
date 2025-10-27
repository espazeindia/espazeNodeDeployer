package config

import (
	"os"
	"strings"
)

type Config struct {
	// Server
	Port string
	Env  string

	// Database
	MongoURI     string
	DatabaseName string

	// JWT
	JWTSecret string
	JWTExpiry string

	// GitHub
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string

	// Kubernetes
	KubeConfig       string
	DefaultNamespace string

	// CORS
	AllowedOrigins string

	// Default Deployment Settings
	DefaultMemoryLimit   string
	DefaultMemoryRequest string
	DefaultCPULimit      string
	DefaultCPURequest    string
	DefaultReplicas      string

	// Observability
	EnableMetrics bool
	MetricsPort   string
}

func Load() *Config {
	return &Config{
		Port:                 getEnv("PORT", "8080"),
		Env:                  getEnv("ENV", "development"),
		MongoURI:             getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName:         getEnv("DATABASE_NAME", "espaze_node_deployer"),
		JWTSecret:            getEnv("JWT_SECRET", "change-this-secret-in-production"),
		JWTExpiry:            getEnv("JWT_EXPIRY", "24h"),
		GitHubClientID:       getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret:   getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURL:    getEnv("GITHUB_REDIRECT_URL", "http://localhost:5173/auth/callback"),
		KubeConfig:           getEnv("KUBECONFIG", os.Getenv("HOME")+"/.kube/config"),
		DefaultNamespace:     getEnv("DEFAULT_NAMESPACE", "espaze-node-deployer-apps"),
		AllowedOrigins:       getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000"),
		DefaultMemoryLimit:   getEnv("DEFAULT_MEMORY_LIMIT", "512Mi"),
		DefaultMemoryRequest: getEnv("DEFAULT_MEMORY_REQUEST", "256Mi"),
		DefaultCPULimit:      getEnv("DEFAULT_CPU_LIMIT", "500m"),
		DefaultCPURequest:    getEnv("DEFAULT_CPU_REQUEST", "250m"),
		DefaultReplicas:      getEnv("DEFAULT_REPLICAS", "2"),
		EnableMetrics:        getEnv("ENABLE_METRICS", "true") == "true",
		MetricsPort:          getEnv("METRICS_PORT", "9090"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.TrimSpace(value)
}

