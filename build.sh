#!/bin/bash

# Script de build otimizado para Jarvis DevOps
# Gera um binário único com todos os assets embutidos

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir com cor
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Configurações
APP_NAME="jarvis-devops"
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.0")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GOOS=${GOOS:-linux}
GOARCH=${GOARCH:-amd64}

# Nome do binário final
BINARY_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
if [ "$GOOS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

print_status "Iniciando build do Jarvis DevOps..."
print_status "Versão: ${VERSION}"
print_status "Commit: ${COMMIT}"
print_status "Plataforma: ${GOOS}/${GOARCH}"
print_status "Binário: ${BINARY_NAME}"

# Verificar se os assets existem
if [ ! -d "internal/assets/web" ]; then
    print_error "Diretório de assets não encontrado: internal/assets/web"
    exit 1
fi

# Limpar builds anteriores
print_status "Limpando builds anteriores..."
rm -f ${APP_NAME}-* jarvis-devops-embedded

# Build flags para otimização
LDFLAGS="-w -s"
LDFLAGS="${LDFLAGS} -X main.Version=${VERSION}"
LDFLAGS="${LDFLAGS} -X main.GitCommit=${COMMIT}"
LDFLAGS="${LDFLAGS} -X main.BuildTime=${BUILD_TIME}"

# Build do binário
print_status "Compilando binário..."
CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
    -ldflags="${LDFLAGS}" \
    -trimpath \
    -o ${BINARY_NAME} \
    ./cmd/server

# Verificar se o build foi bem-sucedido
if [ ! -f "${BINARY_NAME}" ]; then
    print_error "Falha na compilação do binário"
    exit 1
fi

# Mostrar informações do binário
BINARY_SIZE=$(du -h ${BINARY_NAME} | cut -f1)
print_success "Build concluído com sucesso!"
print_status "Arquivo: ${BINARY_NAME}"
print_status "Tamanho: ${BINARY_SIZE}"

# Verificar se UPX está disponível para compressão
if command -v upx &> /dev/null; then
    read -p "Deseja comprimir o binário com UPX? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Comprimindo binário com UPX..."
        upx --best --lzma ${BINARY_NAME}
        COMPRESSED_SIZE=$(du -h ${BINARY_NAME} | cut -f1)
        print_success "Binário comprimido! Novo tamanho: ${COMPRESSED_SIZE}"
    fi
else
    print_warning "UPX não encontrado. Para binários menores, instale UPX: sudo apt install upx"
fi

# Criar arquivo de informações
cat > ${BINARY_NAME}.info << EOF
Jarvis DevOps - Build Information
================================
Version: ${VERSION}
Git Commit: ${COMMIT}
Build Time: ${BUILD_TIME}
Platform: ${GOOS}/${GOARCH}
Binary Size: ${BINARY_SIZE}
Go Version: $(go version)

Dependencies Embedded:
- HTML Templates (web/templates)
- Static Assets (web/static/css, web/static/js)

Deployment:
1. Copy ${BINARY_NAME} to target server
2. Set environment variables (see .env.example)
3. Run: ./${BINARY_NAME}

No external dependencies required!
EOF

print_success "Build completo!"
print_status "Arquivo de informações criado: ${BINARY_NAME}.info"
print_status ""
print_status "Para testar o binário:"
print_status "  ./${BINARY_NAME}"
print_status ""
print_status "Para build cross-platform:"
print_status "  GOOS=windows GOARCH=amd64 ./build.sh"
print_status "  GOOS=darwin GOARCH=amd64 ./build.sh"