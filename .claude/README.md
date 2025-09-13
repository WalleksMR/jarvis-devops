# Jarvis DevOps - Nginx Configuration Manager

Um serviço em Go para gerenciamento de configurações do nginx através de uma interface web moderna e intuitiva.

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Funcionalidades](#funcionalidades)
- [Requisitos](#requisitos)
- [Instalação](#instalação)
- [Configuração](#configuração)
- [Uso](#uso)
- [API](#api)
- [Segurança](#segurança)
- [Desenvolvimento](#desenvolvimento)
- [Contribuição](#contribuição)

## 🎯 Visão Geral

O Jarvis DevOps é uma aplicação web desenvolvida em Go que permite o gerenciamento completo das configurações do nginx através de uma interface web moderna. A aplicação utiliza o framework Gin para o backend e Tailwind CSS com Alpine.js para o frontend, proporcionando uma experiência de usuário fluida e responsiva.

## ✨ Funcionalidades

### Principais

- ✅ **Verificação de Instalação**: Detecta automaticamente se o nginx está instalado no sistema
- ✅ **Validação de Configuração**: Executa `nginx -t` para validar arquivos de configuração
- ✅ **Listagem de Arquivos**: Lista todos os arquivos `.conf` no diretório configurado
- ✅ **Editor Web**: Interface de edição com syntax highlighting para nginx
- ✅ **Reload Seguro**: Recarrega configurações do nginx após validação
- ✅ **Restart de Serviço**: Reinicia o serviço nginx com confirmação
- ✅ **Logs em Tempo Real**: Visualização dos logs do nginx via journalctl

### Interface

- 🎨 **Design Moderno**: Interface responsiva com Tailwind CSS
- ⚡ **Carregamento Dinâmico**: Conteúdo carregado via JavaScript/Alpine.js
- 📱 **Mobile Friendly**: Totalmente responsivo para dispositivos móveis
- 🌙 **Suporte a Tema Escuro**: Suporte automático para modo escuro
- ♿ **Acessibilidade**: Suporte a leitores de tela e navegação por teclado

### Segurança

- 🔐 **Autenticação Basic Auth**: Proteção via usuário/senha
- 🛡️ **Validação de Entrada**: Proteção contra path traversal
- 📝 **Backup Automático**: Backup automático antes de modificações
- 🔍 **Validação Prévia**: Sempre valida antes de aplicar mudanças

## 📋 Requisitos

### Sistema

- **Sistema Operacional**: Linux (Ubuntu, CentOS, Debian, etc.)
- **Go**: Versão 1.21 ou superior
- **Nginx**: Instalado e configurado no sistema
- **Systemctl**: Para gerenciamento do serviço nginx

### Permissões

- Acesso de leitura/escrita ao diretório de configurações do nginx
- Permissão para executar comandos nginx
- Permissão para usar systemctl (sudo pode ser necessário)

## 🚀 Instalação

### 1. Clone o Repositório

```bash
git clone <repository-url>
cd jarvis-devops
```

### 2. Instale as Dependências

```bash
go mod tidy
```

### 3. Configure as Variáveis de Ambiente

```bash
cp .env.example .env
# Edite o arquivo .env conforme suas necessidades
```

### 4. Compile a Aplicação

```bash
go build -o jarvis-devops ./cmd/server
```

### 5. Execute

```bash
./jarvis-devops
```

## ⚙️ Configuração

### Variáveis de Ambiente

A aplicação utiliza as seguintes variáveis de ambiente (todas opcionais com valores padrão):

```env
# Configuração do Servidor
SERVER_HOST=0.0.0.0          # Host do servidor web
SERVER_PORT=8080             # Porta do servidor web

# Configuração do Nginx
NGINX_CONFIG_PATH=/etc/nginx/sites-available  # Diretório dos arquivos .conf
NGINX_BINARY=/usr/sbin/nginx                  # Caminho do binário nginx
NGINX_SERVICE_NAME=nginx                      # Nome do serviço systemctl

# Segurança
BASIC_AUTH_USER=admin        # Usuário para autenticação
BASIC_AUTH_PASSWORD=admin123 # Senha para autenticação

# Aplicação
DEBUG=false                  # Modo debug
LOG_LEVEL=info              # Nível de log (debug, info, warn, error)
```

### Arquivo .env

Crie um arquivo `.env` na raiz do projeto:

```env
SERVER_PORT=8080
NGINX_CONFIG_PATH=/etc/nginx/sites-available
BASIC_AUTH_USER=admin
BASIC_AUTH_PASSWORD=sua_senha_segura
DEBUG=false
```

## 📖 Uso

### Interface Web

1. **Acesse a aplicação**: `http://localhost:8080`
2. **Faça login**: Use as credenciais configuradas (padrão: admin/admin123)
3. **Visualize o status**: A página inicial mostra o status atual do nginx
4. **Gerencie configurações**:
   - Clique em "Editar" em qualquer arquivo para abrir o editor
   - Use os botões de ação para validar, recarregar ou reiniciar
   - Visualize logs em tempo real

### Editor de Configuração

- **Syntax Highlighting**: Destaque de sintaxe específico para nginx
- **Atalhos de Teclado**:
  - `Ctrl+S` / `Cmd+S`: Salvar arquivo
  - `Esc`: Sair do modo de edição
- **Validação Automática**: Valide antes de salvar
- **Backup Automático**: Backup criado automaticamente antes de salvar

### Fluxo de Trabalho Recomendado

1. **Edite** o arquivo de configuração
2. **Valide** a configuração usando o botão "Validar"
3. **Salve** as alterações
4. **Recarregue** o nginx para aplicar as mudanças
5. **Monitore** os logs para verificar se tudo está funcionando

## 🔌 API

### Endpoints Disponíveis

#### Status do Nginx

```http
GET /api/status
```

Retorna informações sobre a instalação e status do nginx.

#### Listar Configurações

```http
GET /api/configs
```

Lista todos os arquivos de configuração disponíveis.

#### Obter Configuração

```http
GET /api/config/{filename}
```

Retorna o conteúdo de um arquivo específico.

#### Atualizar Configuração

```http
PUT /api/config/{filename}
Content-Type: application/json

{
  "content": "server { ... }"
}
```

Atualiza o conteúdo de um arquivo de configuração.

#### Validar Configuração

```http
POST /api/validate
```

Valida a configuração atual do nginx.

#### Recarregar Nginx

```http
POST /api/reload
```

Recarrega a configuração do nginx.

#### Reiniciar Nginx

```http
POST /api/restart
```

Reinicia o serviço nginx.

#### Obter Logs

```http
GET /api/logs?lines=50
```

Retorna os logs recentes do nginx.

### Exemplos de Resposta

#### Status

```json
{
  "is_installed": true,
  "is_running": true,
  "version": "1.18.0",
  "config_valid": true,
  "config_error": ""
}
```

#### Erro

```json
{
  "error": "Failed to validate config: nginx: [emerg] unexpected end of file"
}
```

## 🔒 Segurança

### Autenticação

- Basic Auth habilitado por padrão
- Credenciais configuráveis via variáveis de ambiente
- Sessões protegidas

### Proteções

- **Path Traversal**: Proteção contra acesso a arquivos fora do diretório configurado
- **Backup Automático**: Cria backup antes de qualquer modificação
- **Validação Prévia**: Sempre valida configuração antes de aplicar
- **Sanitização**: Validação de entrada em todos os endpoints

### Recomendações

- Use senhas fortes para autenticação
- Execute com usuário com privilégios mínimos necessários
- Configure SSL/TLS em produção
- Monitore logs de acesso

## 🛠️ Desenvolvimento

### Estrutura do Projeto

```
jarvis-devops/
├── cmd/server/             # Ponto de entrada da aplicação
├── internal/
│   ├── config/            # Gerenciamento de configuração
│   ├── handlers/          # Handlers HTTP
│   └── service/           # Lógica de negócio
├── web/
│   ├── static/           # Arquivos estáticos (CSS, JS)
│   └── templates/        # Templates HTML
├── .claude/              # Documentação
├── .env.example          # Exemplo de configuração
├── go.mod                # Dependências Go
└── README.md
```

### Tecnologias Utilizadas

#### Backend

- **Go 1.21+**: Linguagem principal
- **Gin**: Framework web HTTP
- **GoDotEnv**: Carregamento de variáveis de ambiente

#### Frontend

- **Tailwind CSS**: Framework CSS utilitário
- **Alpine.js**: Framework JavaScript reativo
- **CodeMirror**: Editor de código com syntax highlighting

### Executar em Desenvolvimento

```bash
# Instalar dependências
go mod tidy

# Executar com hot reload (se air estiver instalado)
air

# Ou executar diretamente
go run ./cmd/server
```

### Testes

```bash
# Executar testes
go test ./...

# Executar com coverage
go test -cover ./...
```

### Build para Produção

```bash
# Build otimizado
go build -ldflags="-w -s" -o jarvis-devops ./cmd/server

# Build para múltiplas plataformas
GOOS=linux GOARCH=amd64 go build -o jarvis-devops-linux ./cmd/server
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

### Diretrizes

- Mantenha o código limpo e bem documentado
- Adicione testes para novas funcionalidades
- Siga as convenções de código Go
- Atualize a documentação quando necessário

## 📄 Licença

Este projeto está licenciado sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 🆘 Suporte

Para suporte, abra uma issue no GitHub ou entre em contato através do email do projeto.

---

Desenvolvido com ❤️ para simplificar o gerenciamento de configurações nginx.
