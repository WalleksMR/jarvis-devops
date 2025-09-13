package handlers

import (
	"net/http"
	"strconv"

	"jarvis-devops/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler contains the service dependencies
type Handler struct {
	nginxService *service.NginxService
}

// NewHandler creates a new handler instance
func NewHandler(nginxService *service.NginxService) *Handler {
	return &Handler{
		nginxService: nginxService,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// API routes are now directly under the /api prefix with authentication
	router.GET("/status", h.getNginxStatus)
	router.GET("/configs", h.listConfigFiles)
	router.GET("/config/:filename", h.getConfigFile)
	router.PUT("/config/:filename", h.updateConfigFile)
	router.POST("/validate", h.validateConfig)
	router.POST("/reload", h.reloadNginx)
	router.POST("/restart", h.restartNginx)
	router.GET("/logs", h.getNginxLogs)
}

// Note: Frontend routing is now handled by React Router

// getNginxStatus returns the current nginx status
func (h *Handler) getNginxStatus(c *gin.Context) {
	status, err := h.nginxService.CheckInstallation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check nginx status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// listConfigFiles returns a list of all configuration files
func (h *Handler) listConfigFiles(c *gin.Context) {
	configs, err := h.nginxService.ListConfigFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list config files: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"configs": configs,
	})
}

// getConfigFile returns the content of a specific configuration file
func (h *Handler) getConfigFile(c *gin.Context) {
	filename := c.Param("filename")

	content, err := h.nginxService.ReadConfigFile(filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to read config file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"content":  content,
	})
}

// updateConfigFile updates a configuration file
func (h *Handler) updateConfigFile(c *gin.Context) {
	filename := c.Param("filename")

	var request struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	err := h.nginxService.WriteConfigFile(filename, request.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update config file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration file updated successfully",
	})
}

// validateConfig validates the nginx configuration
func (h *Handler) validateConfig(c *gin.Context) {
	valid, errMsg := h.nginxService.ValidateConfig()

	response := gin.H{
		"valid": valid,
	}

	if errMsg != "" {
		response["error"] = errMsg
	}

	statusCode := http.StatusOK
	if !valid {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, response)
}

// reloadNginx reloads the nginx configuration
func (h *Handler) reloadNginx(c *gin.Context) {
	err := h.nginxService.ReloadNginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reload nginx: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Nginx reloaded successfully",
	})
}

// restartNginx restarts the nginx service
func (h *Handler) restartNginx(c *gin.Context) {
	err := h.nginxService.RestartNginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to restart nginx: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Nginx restarted successfully",
	})
}

// getNginxLogs returns nginx logs
func (h *Handler) getNginxLogs(c *gin.Context) {
	linesStr := c.DefaultQuery("lines", "50")
	lines, err := strconv.Atoi(linesStr)
	if err != nil {
		lines = 50
	}

	logs, err := h.nginxService.GetLogs(lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get logs: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}
