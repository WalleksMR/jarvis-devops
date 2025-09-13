# API Documentation - Jarvis DevOps

Esta documentação detalha todos os endpoints da API REST do Jarvis DevOps.

## Base URL

```
http://localhost:8080/api
```

## Autenticação

Todos os endpoints requerem autenticação Basic Auth:

- **Usuário**: Configurado via `BASIC_AUTH_USER` (padrão: `admin`)
- **Senha**: Configurado via `BASIC_AUTH_PASSWORD` (padrão: `admin123`)

## Endpoints

### 1. Status do Nginx

#### `GET /api/status`

Retorna informações sobre a instalação e status atual do nginx.

**Resposta de Sucesso:**

```json
{
  "is_installed": true,
  "is_running": true,
  "version": "1.18.0",
  "config_valid": true,
  "config_error": "",
  "last_reload": "2023-10-01T10:30:00Z"
}
```

**Campos da Resposta:**

- `is_installed` (boolean): Se o nginx está instalado
- `is_running` (boolean): Se o nginx está executando
- `version` (string): Versão do nginx instalado
- `config_valid` (boolean): Se a configuração atual é válida
- `config_error` (string): Erro de configuração, se houver
- `last_reload` (string): Timestamp do último reload (ISO 8601)

---

### 2. Listar Arquivos de Configuração

#### `GET /api/configs`

Lista todos os arquivos `.conf` no diretório configurado.

**Resposta de Sucesso:**

```json
{
  "configs": [
    {
      "name": "default.conf",
      "path": "/etc/nginx/sites-available/default.conf",
      "size": 2048,
      "modified_time": "2023-10-01T09:15:00Z",
      "is_active": true
    },
    {
      "name": "example.com.conf",
      "path": "/etc/nginx/sites-available/example.com.conf",
      "size": 1536,
      "modified_time": "2023-09-30T14:20:00Z",
      "is_active": false
    }
  ]
}
```

**Campos de Cada Arquivo:**

- `name` (string): Nome do arquivo
- `path` (string): Caminho completo do arquivo
- `size` (integer): Tamanho em bytes
- `modified_time` (string): Data de modificação (ISO 8601)
- `is_active` (boolean): Se o arquivo está ativo (linked em sites-enabled)

---

### 3. Obter Conteúdo de Arquivo

#### `GET /api/config/{filename}`

Retorna o conteúdo de um arquivo de configuração específico.

**Parâmetros:**

- `filename` (path): Nome do arquivo (ex: `default.conf`)

**Resposta de Sucesso:**

```json
{
  "filename": "default.conf",
  "content": "server {\n    listen 80;\n    server_name localhost;\n    # ...\n}"
}
```

**Resposta de Erro:**

```json
{
  "error": "File not found: default.conf"
}
```

---

### 4. Atualizar Arquivo de Configuração

#### `PUT /api/config/{filename}`

Atualiza o conteúdo de um arquivo de configuração.

**Parâmetros:**

- `filename` (path): Nome do arquivo

**Corpo da Requisição:**

```json
{
  "content": "server {\n    listen 80;\n    server_name example.com;\n    root /var/www/html;\n}"
}
```

**Resposta de Sucesso:**

```json
{
  "message": "Configuration file updated successfully"
}
```

**Resposta de Erro:**

```json
{
  "error": "Failed to update config file: permission denied"
}
```

**Notas:**

- Um backup automático é criado antes da atualização
- O arquivo de backup tem o formato: `{filename}.backup.{timestamp}`

---

### 5. Validar Configuração

#### `POST /api/validate`

Executa `nginx -t` para validar a configuração atual.

**Resposta de Sucesso (Configuração Válida):**

```json
{
  "valid": true
}
```

**Resposta de Sucesso (Configuração Inválida):**

```json
{
  "valid": false,
  "error": "nginx: [emerg] unexpected end of file, expecting \"}\" in /etc/nginx/sites-available/default.conf:25"
}
```

---

### 6. Recarregar Nginx

#### `POST /api/reload`

Recarrega a configuração do nginx usando `nginx -s reload`.

**Pré-condições:**

- A configuração deve ser válida (executa validação automaticamente)

**Resposta de Sucesso:**

```json
{
  "message": "Nginx reloaded successfully"
}
```

**Resposta de Erro:**

```json
{
  "error": "Config validation failed: nginx: [emerg] unexpected end of file"
}
```

---

### 7. Reiniciar Nginx

#### `POST /api/restart`

Reinicia o serviço nginx usando `systemctl restart nginx`.

**Pré-condições:**

- A configuração deve ser válida (executa validação automaticamente)

**Resposta de Sucesso:**

```json
{
  "message": "Nginx restarted successfully"
}
```

**Resposta de Erro:**

```json
{
  "error": "Failed to restart nginx: Unit nginx.service not found"
}
```

**Nota:** Este endpoint para o serviço completamente antes de reiniciar.

---

### 8. Obter Logs

#### `GET /api/logs`

Retorna logs recentes do nginx via `journalctl`.

**Parâmetros de Query:**

- `lines` (integer, opcional): Número de linhas para retornar (padrão: 50)

**Exemplo:**

```
GET /api/logs?lines=100
```

**Resposta de Sucesso:**

```json
{
  "logs": "Oct 01 10:30:15 server nginx[1234]: Starting nginx...\nOct 01 10:30:16 server nginx[1234]: nginx started successfully"
}
```

**Resposta de Erro:**

```json
{
  "error": "Failed to get logs: journalctl command not found"
}
```

---

## Códigos de Status HTTP

| Código | Descrição                |
| ------ | ------------------------ |
| 200    | Sucesso                  |
| 400    | Requisição inválida      |
| 401    | Não autorizado           |
| 404    | Recurso não encontrado   |
| 500    | Erro interno do servidor |

## Exemplos de Uso

### cURL

#### Obter Status

```bash
curl -u admin:admin123 http://localhost:8080/api/status
```

#### Listar Configurações

```bash
curl -u admin:admin123 http://localhost:8080/api/configs
```

#### Atualizar Arquivo

```bash
curl -u admin:admin123 \
  -X PUT \
  -H "Content-Type: application/json" \
  -d '{"content":"server { listen 80; }"}' \
  http://localhost:8080/api/config/test.conf
```

#### Validar Configuração

```bash
curl -u admin:admin123 \
  -X POST \
  http://localhost:8080/api/validate
```

### JavaScript (Frontend)

```javascript
// Obter status com autenticação
const response = await fetch("/api/status", {
  headers: {
    Authorization: "Basic " + btoa("admin:admin123"),
  },
});
const status = await response.json();

// Atualizar arquivo
const updateResponse = await fetch("/api/config/default.conf", {
  method: "PUT",
  headers: {
    "Content-Type": "application/json",
    Authorization: "Basic " + btoa("admin:admin123"),
  },
  body: JSON.stringify({
    content: "server { listen 80; server_name example.com; }",
  }),
});
```

## Tratamento de Erros

Todos os endpoints retornam erros no formato:

```json
{
  "error": "Descrição detalhada do erro"
}
```

### Erros Comuns

| Erro                       | Descrição                          | Solução                                    |
| -------------------------- | ---------------------------------- | ------------------------------------------ |
| `File not found`           | Arquivo não existe                 | Verificar se o arquivo existe no diretório |
| `Permission denied`        | Sem permissão para acessar arquivo | Verificar permissões do usuário            |
| `Config validation failed` | Configuração nginx inválida        | Corrigir erro de sintaxe no arquivo        |
| `Nginx not found`          | Binário nginx não encontrado       | Verificar caminho do nginx                 |
| `Service not found`        | Serviço nginx não encontrado       | Verificar se nginx está instalado          |

## Rate Limiting

Atualmente não há rate limiting implementado. Em produção, considere usar um reverse proxy com rate limiting.

## Segurança

- Use HTTPS em produção
- Configure senhas fortes para Basic Auth
- Monitore logs de acesso para detectar tentativas de acesso não autorizado
- Execute o serviço com privilégios mínimos necessários
