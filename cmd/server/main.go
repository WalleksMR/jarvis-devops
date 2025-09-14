package main

import (
	"io/fs"
	"log"
	"net/http"
	"time"

	"jarvis-devops/internal/assets"
	"jarvis-devops/internal/config"
	"jarvis-devops/internal/handlers"
	"jarvis-devops/internal/service"

	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup embedded assets from React build
	assets.SetupRoutes(router)

	// Setup basic auth middleware - only for API routes, not for the React app
	authorized := router.Group("/api", gin.BasicAuth(gin.Accounts{
		cfg.BasicAuthUser: cfg.BasicAuthPassword,
	}))

	// Register routes with auth
	handler.RegisterRoutes(authorized)

	// Handle any routes not matched by the API to support React Router
	router.NoRoute(func(c *gin.Context) {
		// Skip API routes in the NoRoute handler
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[0:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		// For all other routes, serve index.html to support client-side routing
		c.Header("Content-Type", "text/html; charset=utf-8")
		content, err := fs.ReadFile(assets.DistFiles, "dist/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading index.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	})

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
		log.Printf("CORS: Enabled for specific development origins")
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
