# Makefile para Jarvis DevOps

# Variáveis
APP_NAME := jarvis-devops
BUILD_DIR := build
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS := -w -s -X main.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)

# Arquivos fonte
SOURCES := $(shell find . -type f -name '*.go')

# Default target
.PHONY: all
all: build

# Build binário com assets embutidos
.PHONY: build
build: $(BUILD_DIR)/$(APP_NAME)

$(BUILD_DIR)/$(APP_NAME): $(SOURCES)
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -trimpath -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server
	@echo "✅ Build complete! Binary: $(BUILD_DIR)/$(APP_NAME)"

# Build otimizado usando o script
.PHONY: build-optimized
build-optimized:
	@echo "🚀 Building optimized binary..."
	./build.sh

# Build para múltiplas plataformas
.PHONY: build-all
build-all:
	@echo "🌍 Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 ./build.sh
	GOOS=windows GOARCH=amd64 ./build.sh
	GOOS=darwin GOARCH=amd64 ./build.sh

# Executar em modo desenvolvimento
.PHONY: dev
dev:
	@echo "🔧 Running in development mode..."
	go run ./cmd/server/main.go

# Executar testes
.PHONY: test
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Verificar dependências
.PHONY: deps
deps:
	@echo "📦 Checking dependencies..."
	go mod tidy
	go mod verify

# Limpeza
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(APP_NAME) $(APP_NAME)-* *.info

# Instalar dependências de desenvolvimento
.PHONY: install-dev-deps
install-dev-deps:
	@echo "📦 Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Lint do código
.PHONY: lint
lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Formatar código
.PHONY: fmt
fmt:
	@echo "✨ Formatting code..."
	go fmt ./...

# Mostrar ajuda
.PHONY: help
help:
	@echo "Jarvis DevOps - Comandos disponíveis:"
	@echo ""
	@echo "  build              Compila o binário básico"
	@echo "  build-optimized    Compila binário otimizado (usa build.sh)"
	@echo "  build-all          Compila para múltiplas plataformas"
	@echo "  dev                Executa em modo desenvolvimento"
	@echo "  test               Executa testes"
	@echo "  deps               Verifica e atualiza dependências"
	@echo "  clean              Limpa arquivos de build"
	@echo "  clean-all          Limpeza completa (inclui backups)"
	@echo "  install-dev-deps   Instala dependências de desenvolvimento"
	@echo "  lint               Executa linter de código"
	@echo "  fmt                Formata o código"
	@echo "  check-assets       Verifica se assets embutidos estão presentes"
	@echo "  test-embedded      Testa build sem diretório web/ original"
	@echo "  help               Mostra esta ajuda"
	@echo ""
	@echo "Exemplo de uso:"
	@echo "  make build         # Build simples"
	@echo "  make build-optimized # Build otimizado"
	@echo "  make dev           # Modo desenvolvimento"
	@echo "  make test-embedded # Testar se assets estão embutidos"

# Verificar se todos os assets estão presentes
.PHONY: check-assets
check-assets:
	@echo "🔍 Checking embedded assets..."
	@if [ ! -d "internal/assets/web" ]; then \
		echo "❌ Assets directory not found: internal/assets/web"; \
		exit 1; \
	fi
	@if [ ! -f "internal/assets/embed.go" ]; then \
		echo "❌ Embed file not found: internal/assets/embed.go"; \
		exit 1; \
	fi
	@echo "✅ All assets are present"

# Verificar se build funciona sem web/ original
.PHONY: test-embedded
test-embedded:
	@echo "🧪 Testing if binary works without original web/ directory..."
	@if [ -d "web" ]; then \
		mv web web.backup.test; \
		echo "📦 Temporarily moved web/ to web.backup.test"; \
	fi
	@make build
	@echo "✅ Build successful without original web/ directory"
	@if [ -d "web.backup.test" ]; then \
		mv web.backup.test web; \
		echo "📦 Restored web/ directory"; \
	fi

# Limpeza completa incluindo arquivos de backup
.PHONY: clean-all
clean-all: clean
	@echo "🧹 Deep cleaning..."
	rm -rf web.backup*
	rm -rf /tmp/jarvis-test*
	@echo "✅ Deep clean complete"

# Target que roda antes do build para verificar assets
build: check-assets

.DEFAULT_GOAL := help