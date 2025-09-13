# Jarvis DevOps - Nginx Configuration Manager

This document provides essential guidance for AI coding agents working with the Jarvis DevOps codebase.

## Architecture Overview

Jarvis DevOps is a Go web application that provides an interface for managing Nginx configurations. The architecture follows a layered pattern:

1. **Entry Point (`cmd/server/main.go`)**: Initializes components and starts the HTTP server
2. **Config Layer (`internal/config/config.go`)**: Manages environment-based configuration 
3. **Service Layer (`internal/service/nginx.go`)**: Contains core business logic for nginx management
4. **Handler Layer (`internal/handlers/handlers.go`)**: HTTP handlers and API endpoints
5. **Assets Layer (`internal/assets/embed.go`)**: Embedded static files and templates
6. **Web UI**: HTML templates and JS/CSS assets (embedded in binary)

Key data flows:
- HTTP requests → Gin router → Handlers → Service → File system / Nginx commands
- Configuration is loaded once at startup from environment variables or `.env` file
- User edits are validated before being written to nginx config files
- Service layer integrates with system commands via `os/exec` for nginx operations
- Static assets are served from embedded filesystem (no external dependencies)

## Build System

### Embedded Assets Architecture

The application uses Go's `//go:embed` directive to include all static assets in the binary:

```
internal/assets/
├── embed.go              # Embedded asset configuration
└── web/
    ├── static/
    │   ├── css/style.css
    │   └── js/ (app.js, editor.js)
    └── templates/ (*.html)
```

**Key Benefits:**
- Single binary deployment
- No external file dependencies
- Simplified distribution and Docker builds
- Cross-platform compatibility

### Build Commands

```bash
# Development mode
go run ./cmd/server/main.go

# Quick build
make build

# Optimized build (recommended for production)
make build-optimized
./build.sh

# Cross-platform builds
GOOS=windows GOARCH=amd64 ./build.sh
GOOS=darwin GOARCH=amd64 ./build.sh

# Build for multiple platforms
make build-all
```

### Build Optimization

The build system includes several optimizations:
- **`-ldflags="-w -s"`**: Removes debug info, reduces binary size (~30% smaller)
- **`-trimpath`**: Removes absolute paths from binary
- **`CGO_ENABLED=0`**: Creates static binary with no C dependencies
- **Embedded assets**: All web files included in binary
- **UPX compression**: Optional further compression (install `upx` package)

Typical binary sizes:
- Standard build: ~18MB
- Optimized build: ~12MB
- With UPX compression: ~4-6MB

## Development Workflow

### Environment Setup

Create a `.env` file in the project root:

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
NGINX_CONFIG_PATH=/etc/nginx/sites-available
NGINX_BINARY=/usr/sbin/nginx
NGINX_SERVICE_NAME=nginx
BASIC_AUTH_USER=admin
BASIC_AUTH_PASSWORD=admin123
DEBUG=true
LOG_LEVEL=info
```

### Deployment

**Single Binary Deployment:**
1. Build: `./build.sh` or `make build-optimized`
2. Copy binary to target server: `scp jarvis-devops-linux-amd64 user@server:/opt/jarvis/`
3. Set permissions: `chmod +x jarvis-devops-linux-amd64`
4. Run: `./jarvis-devops-linux-amd64`

No additional files needed - all assets are embedded!

## Project-Specific Conventions

1. **Error Handling**: Service layer returns detailed errors, handler layer transforms them into HTTP responses
2. **Path Validation**: All file operations validate paths to prevent traversal attacks
   - Files with `..` or `/` in their names are rejected to prevent directory traversal
3. **Backup System**: Automatic backups are created before modifying any nginx config file
   - Backups are stored with timestamp suffixes: `{filename}.backup.{timestamp}`
4. **HTTP Authentication**: All routes (except health check) require Basic Auth
5. **Validation Before Action**: The application always validates nginx configurations before applying changes
6. **Frontend Notifications**: All actions trigger user notifications through the UI feedback system
7. **Embedded Assets**: All static files are embedded in the binary for simplified deployment

## Key Components and Integration Points

### Asset Management (`internal/assets/embed.go`)

The asset system provides:
- `GetStaticFS()`: Returns embedded static files as http.FileSystem
- `GetTemplates()`: Returns parsed HTML templates from embedded files
- `SetupRoutes(router)`: Configures Gin router to use embedded assets

**Important**: When modifying web assets, they must be updated in `internal/assets/web/` directory to be included in the binary.

### Nginx Service Integration

The `NginxService` in `internal/service/nginx.go` encapsulates all nginx interactions:
- Validates config files using `nginx -t`
- Manages nginx via systemctl commands
- Fetches logs via journalctl
- Reads/writes config files with path validation

Key service methods:
- `CheckInstallation()`: Verifies nginx installation status
- `ValidateConfig()`: Validates nginx configuration syntax
- `ListConfigFiles()`: Lists available configuration files
- `ReadConfigFile()`: Reads a specific configuration file
- `WriteConfigFile()`: Updates a configuration file (with automatic backup)
- `ReloadNginx()`: Reloads nginx configuration
- `RestartNginx()`: Restarts the nginx service
- `GetLogs()`: Fetches recent nginx logs

### Web UI Integration

- Frontend built with Tailwind CSS + Alpine.js
- Configuration editor uses CodeMirror with nginx syntax highlighting
- API calls are made from JavaScript to update content dynamically
- **All assets are embedded** - no external file dependencies

## Common Tasks

### Adding a New API Endpoint

1. Create a handler function in `internal/handlers/handlers.go`
2. Register it in the `RegisterRoutes` method
3. If needed, add corresponding business logic to `internal/service/nginx.go`

### Modifying the UI

1. Edit templates in `internal/assets/web/templates/`
2. Update corresponding JavaScript in `internal/assets/web/static/js/`
3. For editor functionality, modify `editor.js`
4. **Important**: Assets are embedded at build time - rebuild binary to see changes

### Adding Configuration Options

1. Add the new field to the `Config` struct in `internal/config/config.go`
2. Set a default value in the `Load()` function
3. Update the environment variable documentation

### Working with Embedded Assets

When modifying web assets:
1. Edit files in `internal/assets/web/` (this is the only web directory)
2. Test changes: `go run ./cmd/server/main.go`
3. Build production binary: `make build-optimized`
4. Deploy the new binary

**Important**: The original `web/` directory has been removed. All assets are now in `internal/assets/web/` and embedded in the binary.

## API Integration

The application provides a comprehensive REST API for automation:

### Key API Endpoints

- `GET /health`: Health check (no auth required)
- `GET /api/status`: Get nginx installation and runtime status
- `GET /api/configs`: List available configuration files
- `GET /api/config/{filename}`: Get contents of a specific config file
- `PUT /api/config/{filename}`: Update a config file
- `POST /api/validate`: Validate nginx configuration
- `POST /api/reload`: Reload nginx configuration
- `POST /api/restart`: Restart nginx service
- `GET /api/logs`: Get nginx logs

All endpoints (except `/health`) use Basic Auth and return responses in JSON format. Errors follow a consistent pattern with an `error` field containing the error message.

## Testing and Validation

Always validate nginx configurations before applying changes using the `ValidateConfig()` method from `NginxService`. The application includes built-in validation before any nginx reload or restart.

### Recommended Testing Flow

1. Make changes to configuration files
2. Validate the changes using `nginx -t` (handled by `ValidateConfig()`)
3. If valid, reload nginx to apply changes
4. Check logs to ensure the changes were applied correctly

### Build Testing

1. Test development mode: `go run ./cmd/server/main.go`
2. Test embedded build: `make build && ./jarvis-devops`
3. Verify assets work without external `web/` directory
4. Test binary portability by copying to different location

## Production Deployment Best Practices

1. **Use optimized builds**: `./build.sh` for smallest binaries
2. **Single binary deployment**: No need to copy web assets separately
3. **Environment configuration**: Use `.env` file or environment variables
4. **Service management**: Create systemd service for production
5. **Security**: Use strong Basic Auth credentials
6. **Backup**: Automatic nginx config backups are created before changes
7. **Monitoring**: Use `/health` endpoint for health checks

