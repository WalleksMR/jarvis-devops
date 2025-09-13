# Jarvis DevOps - Nginx Configuration Manager

🚀 **Um serviço moderno em Go para gerenciamento de configurações nginx através de interface web**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Nginx](https://img.shields.io/badge/Nginx-Compatible-green.svg)](https://nginx.org/)
[![Embedded Assets](https://img.shields.io/badge/Assets-Embedded-orange.svg)](internal/assets/)

## 📖 Sobre

O Jarvis DevOps é uma aplicação web desenvolvida em Go que permite o gerenciamento completo das configurações do nginx através de uma interface moderna e intuitiva. Perfeito para administradores de sistema que precisam de uma forma rápida e segura de editar, validar e aplicar configurações nginx.

**🎯 Binário único com assets embutidos** - Todos os arquivos web (HTML, CSS, JS) são incluídos no binário para distribuição simplificada.

## ✨ Funcionalidades Principais

- ✅ **Verificação automática** da instalação do nginx
- ✅ **Validação em tempo real** dos arquivos de configuração
- ✅ **Editor web** com syntax highlighting para nginx
- ✅ **Listagem dinâmica** de arquivos de configuração
- ✅ **Reload e restart seguros** do nginx
- ✅ **Visualização de logs** em tempo real
- ✅ **Interface responsiva** com Tailwind CSS
- ✅ **Distribuição simplificada** com binário único
- ✅ **API REST completa** para automação
- ✅ **Backup automático** antes de modificações
- ✅ **Autenticação Basic Auth** integrada
- ✅ **Binário único** com todos os recursos embutidos

## 🚀 Início Rápido

### 1. Clone e Configure

```bash
git clone <repository-url>
cd jarvis-devops

# Copie e configure as variáveis de ambiente
cp .env.example .env
nano .env
```

### 2. Execute

```bash
# Instale dependências
go mod tidy

# Execute a aplicação
go run ./cmd/server

# Ou compile e execute
go build -o jarvis-devops ./cmd/server
./jarvis-devops
```

### 3. Acesse

Abra http://localhost:8080 no seu navegador e faça login com as credenciais configuradas (padrão: admin/admin123).

### 4. Compilação com Arquivos Estáticos Embutidos

A aplicação suporta a compilação em um único binário autônomo, com todos os arquivos estáticos e templates embutidos:

```bash
# Compile o binário único com todos os recursos
go build -o jarvis-devops ./cmd/server

# Execute o binário autônomo
./jarvis-devops
```

Isso cria um binário único e portátil que não requer a distribuição separada dos arquivos da pasta `/web`.

## 📁 Estrutura do Projeto

```
jarvis-devops/
├── cmd/server/                 # Ponto de entrada da aplicação
│   └── main.go
├── internal/
│   ├── config/                # Gerenciamento de configuração
│   │   └── config.go
│   ├── handlers/              # Handlers HTTP/API
│   │   └── handlers.go
│   └── service/               # Lógica de negócio do nginx
│       └── nginx.go
├── web/
│   ├── static/
│   │   ├── css/              # Estilos CSS customizados
│   │   │   └── style.css
│   │   └── js/               # JavaScript da aplicação
│   │       ├── app.js        # Funcionalidades principais
│   │       └── editor.js     # Editor de configuração
│   └── templates/            # Templates HTML
│       ├── index.html        # Página principal
│       └── editor.html       # Editor de arquivos
├── .claude/                  # Documentação detalhada
│   ├── README.md            # Documentação principal
│   ├── API.md               # Documentação da API
│   └── SETUP.md             # Guia de instalação
├── .env.example             # Exemplo de configuração
├── go.mod                   # Dependências Go
├── go.sum                   # Checksums das dependências
└── README.md                # Este arquivo
```

## ⚙️ Configuração

### Variáveis de Ambiente Principais

```env
# Servidor
SERVER_HOST=0.0.0.0              # Host do servidor
SERVER_PORT=8080                 # Porta do servidor

# Nginx
NGINX_CONFIG_PATH=/etc/nginx/sites-available  # Diretório dos arquivos .conf
NGINX_BINARY=/usr/sbin/nginx                  # Caminho do nginx
NGINX_SERVICE_NAME=nginx                      # Nome do serviço

# Segurança
BASIC_AUTH_USER=admin            # Usuário para login
BASIC_AUTH_PASSWORD=admin123     # Senha (altere em produção!)

# Aplicação
DEBUG=false                      # Modo debug
LOG_LEVEL=info                   # Nível de log
```

Para configuração completa, veja [.env.example](.env.example) e a [documentação de setup](.claude/SETUP.md).

## 🔌 API REST

A aplicação fornece uma API REST completa para automação:

```bash
# Verificar status do nginx
curl -u admin:senha http://localhost:8080/api/status

# Listar arquivos de configuração
curl -u admin:senha http://localhost:8080/api/configs

# Obter conteúdo de um arquivo
curl -u admin:senha http://localhost:8080/api/config/default.conf

# Validar configuração
curl -u admin:senha -X POST http://localhost:8080/api/validate

# Recarregar nginx
curl -u admin:senha -X POST http://localhost:8080/api/reload
```

Para documentação completa da API, veja [API.md](.claude/API.md).

## 🛠️ Tecnologias Utilizadas

### Backend

- **[Go 1.21+](https://golang.org/)** - Linguagem principal
- **[Gin](https://gin-gonic.com/)** - Framework web HTTP
- **[GoDotEnv](https://github.com/joho/godotenv)** - Carregamento de variáveis de ambiente

### Frontend

- **[Tailwind CSS](https://tailwindcss.com/)** - Framework CSS utilitário
- **[Alpine.js](https://alpinejs.dev/)** - Framework JavaScript reativo
- **[CodeMirror](https://codemirror.net/)** - Editor de código com syntax highlighting

## 📖 Documentação

- **[README Principal](.claude/README.md)** - Documentação completa e detalhada
- **[API Documentation](.claude/API.md)** - Todos os endpoints da API REST
- **[Setup Guide](.claude/SETUP.md)** - Guia completo de instalação e configuração

## 🔒 Segurança

- ✅ Autenticação Basic Auth obrigatória
- ✅ Proteção contra path traversal
- ✅ Validação de entrada em todos os endpoints
- ✅ Backup automático antes de modificações
- ✅ Validação prévia antes de aplicar mudanças

**⚠️ Importante**: Altere as credenciais padrão em produção e configure HTTPS.

## 🚦 Status do Projeto

- ✅ **Core Functionality** - Todas as funcionalidades principais implementadas
- ✅ **Web Interface** - Interface moderna e responsiva
- ✅ **API REST** - API completa para automação
- ✅ **Documentation** - Documentação completa
- ✅ **Security** - Implementações básicas de segurança
- 🔄 **Testing** - Testes unitários (em desenvolvimento)
- 🔄 **CI/CD** - Pipeline de integração contínua (planejado)

## 🤝 Contribuindo

Contribuições são bem-vindas! Para contribuir:

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está licenciado sob a Licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🆘 Suporte

- 📖 **Documentação**: Veja a pasta [.claude/](.claude/) para documentação detalhada
- 🐛 **Issues**: Abra uma issue no GitHub para reportar bugs
- 💡 **Features**: Suggira novas funcionalidades através de issues
- 📧 **Contato**: Entre em contato através do GitHub

---

**Desenvolvido com ❤️ para simplificar o gerenciamento de configurações nginx**
