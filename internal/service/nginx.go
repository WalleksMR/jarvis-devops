package service

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"jarvis-devops/internal/config"
)

// NginxService provides nginx management functionality
type NginxService struct {
	config *config.Config
}

// ConfigFile represents an nginx configuration file
type ConfigFile struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"modified_time"`
	IsActive     bool      `json:"is_active"`
}

// NginxStatus represents the current status of nginx
type NginxStatus struct {
	IsInstalled bool      `json:"is_installed"`
	IsRunning   bool      `json:"is_running"`
	Version     string    `json:"version"`
	ConfigValid bool      `json:"config_valid"`
	ConfigError string    `json:"config_error,omitempty"`
	LastReload  time.Time `json:"last_reload,omitempty"`
}

// NewNginxService creates a new nginx service instance
func NewNginxService(cfg *config.Config) *NginxService {
	return &NginxService{
		config: cfg,
	}
}

// CheckInstallation verifies if nginx is installed and accessible
func (s *NginxService) CheckInstallation() (*NginxStatus, error) {
	status := &NginxStatus{}

	// Check if nginx binary exists
	if _, err := os.Stat(s.config.NginxBinary); os.IsNotExist(err) {
		status.IsInstalled = false
		return status, nil
	}
	status.IsInstalled = true

	// Get nginx version
	if version, err := s.getNginxVersion(); err == nil {
		status.Version = version
	}

	// Check if nginx is running
	status.IsRunning = s.isNginxRunning()

	// Check config validity
	valid, errMsg := s.ValidateConfig()
	status.ConfigValid = valid
	if errMsg != "" {
		status.ConfigError = errMsg
	}

	return status, nil
}

// getNginxVersion gets the nginx version
func (s *NginxService) getNginxVersion() (string, error) {
	cmd := exec.Command(s.config.NginxBinary, "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// nginx outputs version to stderr, extract version number
	versionStr := string(output)
	if idx := strings.Index(versionStr, "nginx/"); idx != -1 {
		version := versionStr[idx+6:]
		if spaceIdx := strings.Index(version, " "); spaceIdx != -1 {
			version = version[:spaceIdx]
		}
		return strings.TrimSpace(version), nil
	}

	return strings.TrimSpace(versionStr), nil
}

// isNginxRunning checks if nginx process is running
func (s *NginxService) isNginxRunning() bool {
	cmd := exec.Command("pgrep", "nginx")
	err := cmd.Run()
	return err == nil
}

// ValidateConfig validates nginx configuration using nginx -t
func (s *NginxService) ValidateConfig() (bool, string) {
	cmd := exec.Command(s.config.NginxBinary, "-t")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, string(output)
	}

	return true, ""
}

// ListConfigFiles lists all configuration files in the configured directory
func (s *NginxService) ListConfigFiles() ([]ConfigFile, error) {
	var configFiles []ConfigFile

	err := filepath.WalkDir(s.config.NginxConfigPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only include .conf files
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".conf") {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		configFile := ConfigFile{
			Name:         d.Name(),
			Path:         path,
			Size:         info.Size(),
			ModifiedTime: info.ModTime(),
			IsActive:     s.isConfigActive(path),
		}

		configFiles = append(configFiles, configFile)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list config files: %w", err)
	}

	return configFiles, nil
}

// isConfigActive checks if a config file is currently active (linked in sites-enabled)
func (s *NginxService) isConfigActive(configPath string) bool {
	// This is a simplified check - in real nginx setups, you'd check if the file
	// is symlinked in sites-enabled directory
	enabledPath := strings.Replace(configPath, "sites-available", "sites-enabled", 1)
	_, err := os.Stat(enabledPath)
	return err == nil
}

// ReadConfigFile reads the content of a configuration file
func (s *NginxService) ReadConfigFile(filename string) (string, error) {
	// Ensure the filename is safe (no path traversal)
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		return "", fmt.Errorf("invalid filename: %s", filename)
	}

	filePath := filepath.Join(s.config.NginxConfigPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", filename)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

// WriteConfigFile writes content to a configuration file
func (s *NginxService) WriteConfigFile(filename, content string) error {
	// Ensure the filename is safe (no path traversal)
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		return fmt.Errorf("invalid filename: %s", filename)
	}

	filePath := filepath.Join(s.config.NginxConfigPath, filename)

	// Create backup of existing file
	if _, err := os.Stat(filePath); err == nil {
		backupPath := filePath + ".backup." + time.Now().Format("20060102-150405")
		if err := s.copyFile(filePath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Write new content
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// copyFile creates a copy of a file
func (s *NginxService) copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ReloadNginx reloads nginx configuration
func (s *NginxService) ReloadNginx() error {
	// First validate config
	if valid, errMsg := s.ValidateConfig(); !valid {
		return fmt.Errorf("config validation failed: %s", errMsg)
	}

	// Reload nginx
	cmd := exec.Command(s.config.NginxBinary, "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx reload failed: %s", string(output))
	}

	return nil
}

// RestartNginx restarts the nginx service
func (s *NginxService) RestartNginx() error {
	// First validate config
	if valid, errMsg := s.ValidateConfig(); !valid {
		return fmt.Errorf("config validation failed: %s", errMsg)
	}

	// Stop nginx
	stopCmd := exec.Command("systemctl", "stop", s.config.NginxServiceName)
	if output, err := stopCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop nginx: %s", string(output))
	}

	// Start nginx
	startCmd := exec.Command("systemctl", "start", s.config.NginxServiceName)
	if output, err := startCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start nginx: %s", string(output))
	}

	return nil
}

// GetLogs returns recent nginx logs
func (s *NginxService) GetLogs(lines int) (string, error) {
	if lines <= 0 {
		lines = 50
	}

	cmd := exec.Command("journalctl", "-u", s.config.NginxServiceName, "-n", fmt.Sprintf("%d", lines), "--no-pager")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}

	return string(output), nil
}
