package main

import (
	"log"
	"net/http"

	"jarvis-devops/internal/config"
	"jarvis-devops/internal/handlers"
	"jarvis-devops/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on debug setting
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize services
	nginxService := service.NewNginxService(cfg)

	// Initialize handlers
	handler := handlers.NewHandler(nginxService)

	// Create Gin router
	router := gin.Default()

	// Serve static files and templates at engine level
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	// Setup basic auth middleware
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		cfg.BasicAuthUser: cfg.BasicAuthPassword,
	}))

	// Register routes with auth
	handler.RegisterRoutes(authorized)

	// Health check endpoint (without auth for monitoring)
	router.GET("/health", func(c *gin.Context) {
		status, err := nginxService.CheckInstallation()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		if status.IsInstalled {
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"nginx":  status,
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"nginx":  status,
			})
		}
	})

	// Log startup information
	if cfg.Debug {
		log.Printf("Starting Jarvis DevOps in DEBUG mode")
		log.Printf("Server: %s", cfg.GetServerAddress())
		log.Printf("Nginx Config Path: %s", cfg.NginxConfigPath)
		log.Printf("Nginx Binary: %s", cfg.NginxBinary)
		log.Printf("Basic Auth User: %s", cfg.BasicAuthUser)
	} else {
		log.Printf("Starting Jarvis DevOps server on %s", cfg.GetServerAddress())
	}

	// Start server
	log.Printf("🚀 Jarvis DevOps is running on http://%s", cfg.GetServerAddress())
	log.Printf("📖 Documentation: http://%s (login: %s)", cfg.GetServerAddress(), cfg.BasicAuthUser)

	if err := router.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
