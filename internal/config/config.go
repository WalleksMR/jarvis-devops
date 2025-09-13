package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration settings for the application
type Config struct {
	// Server configuration
	ServerHost string
	ServerPort string

	// Nginx configuration
	NginxConfigPath  string // Directory where nginx config files are located
	NginxBinary      string // Path to nginx binary
	NginxServiceName string // Service name for systemctl commands

	// Security
	BasicAuthUser     string
	BasicAuthPassword string

	// Application settings
	Debug    bool
	LogLevel string
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		ServerHost:        getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		NginxConfigPath:   getEnv("NGINX_CONFIG_PATH", "/etc/nginx/sites-available"),
		NginxBinary:       getEnv("NGINX_BINARY", "/usr/sbin/nginx"),
		NginxServiceName:  getEnv("NGINX_SERVICE_NAME", "nginx"),
		BasicAuthUser:     getEnv("BASIC_AUTH_USER", "admin"),
		BasicAuthPassword: getEnv("BASIC_AUTH_PASSWORD", "admin123"),
		Debug:             getEnvBool("DEBUG", false),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
	}

	// Validate required configuration
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate performs basic validation of configuration values
func (c *Config) validate() error {
	// Check if nginx config path exists
	if _, err := os.Stat(c.NginxConfigPath); os.IsNotExist(err) {
		log.Printf("Warning: Nginx config path does not exist: %s", c.NginxConfigPath)
	}

	// Check if nginx binary exists
	if _, err := os.Stat(c.NginxBinary); os.IsNotExist(err) {
		log.Printf("Warning: Nginx binary not found at: %s", c.NginxBinary)
	}

	return nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvBool gets a boolean environment variable with a fallback value
func getEnvBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return fallback
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.ServerHost + ":" + c.ServerPort
}
