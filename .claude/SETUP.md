# Guia de Instalação e Setup - Jarvis DevOps

Este guia fornece instruções detalhadas para instalação e configuração do Jarvis DevOps em diferentes ambientes.

## 📋 Pré-requisitos

### Sistema Operacional

- **Linux**: Ubuntu 18.04+, CentOS 7+, Debian 9+, ou distribuições compatíveis
- **Arquitetura**: x86_64 (AMD64)

### Software Necessário

- **Go**: Versão 1.21 ou superior
- **Nginx**: Qualquer versão moderna (1.14+)
- **Git**: Para clone do repositório
- **Systemd**: Para gerenciamento de serviços (padrão na maioria das distros modernas)

### Permissões

O usuário que executará o serviço precisa de:

- Acesso de leitura/escrita ao diretório de configurações do nginx
- Permissão para executar `nginx -t` e `nginx -s reload`
- Permissão para usar `systemctl` (pode requerer sudo)
- Acesso ao `journalctl` para visualizar logs

## 🚀 Instalação

### Método 1: Instalação Manual

#### 1. Instalar Go (se não estiver instalado)

**Ubuntu/Debian:**

```bash
sudo apt update
sudo apt install golang-go
```

**CentOS/RHEL:**

```bash
sudo yum install golang
# ou para versões mais recentes:
sudo dnf install golang
```

**Instalação via tarball (versão mais recente):**

```bash
cd /tmp
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 2. Verificar instalação do Nginx

```bash
nginx -v
systemctl status nginx
```

Se o nginx não estiver instalado:

**Ubuntu/Debian:**

```bash
sudo apt install nginx
sudo systemctl enable nginx
sudo systemctl start nginx
```

**CentOS/RHEL:**

```bash
sudo yum install nginx
# ou:
sudo dnf install nginx
sudo systemctl enable nginx
sudo systemctl start nginx
```

#### 3. Clone e build da aplicação

```bash
# Clone o repositório
git clone <repository-url>
cd jarvis-devops

# Instalar dependências
go mod tidy

# Verificar se tudo compila
go build ./cmd/server

# Build otimizado para produção
go build -ldflags="-w -s" -o jarvis-devops ./cmd/server
```

#### 4. Configuração inicial

```bash
# Copiar arquivo de exemplo
cp .env.example .env

# Editar configurações
nano .env
```

Exemplo de configuração para produção:

```env
SERVER_HOST=127.0.0.1
SERVER_PORT=8080
NGINX_CONFIG_PATH=/etc/nginx/sites-available
NGINX_BINARY=/usr/sbin/nginx
NGINX_SERVICE_NAME=nginx
BASIC_AUTH_USER=admin
BASIC_AUTH_PASSWORD=sua_senha_muito_segura_aqui
DEBUG=false
LOG_LEVEL=info
```

#### 5. Teste inicial

```bash
# Executar em modo de teste
./jarvis-devops

# Em outro terminal, testar
curl -u admin:sua_senha_muito_segura_aqui http://localhost:8080/api/status
```

### Método 2: Instalação como Serviço Systemd

#### 1. Criar usuário dedicado

```bash
# Criar usuário sem shell para segurança
sudo useradd --system --no-create-home --shell /bin/false jarvis-devops

# Adicionar ao grupo necessário para acessar nginx
sudo usermod -a -G nginx jarvis-devops
```

#### 2. Preparar diretórios

```bash
# Criar diretório da aplicação
sudo mkdir -p /opt/jarvis-devops
sudo mkdir -p /etc/jarvis-devops
sudo mkdir -p /var/log/jarvis-devops

# Copiar arquivos
sudo cp jarvis-devops /opt/jarvis-devops/
sudo cp -r web /opt/jarvis-devops/
sudo cp .env /etc/jarvis-devops/

# Ajustar permissões
sudo chown -R jarvis-devops:jarvis-devops /opt/jarvis-devops
sudo chown -R jarvis-devops:jarvis-devops /var/log/jarvis-devops
sudo chmod 640 /etc/jarvis-devops/.env
```

#### 3. Configurar sudoers (para comandos nginx)

```bash
sudo visudo -f /etc/sudoers.d/jarvis-devops
```

Adicionar:

```
jarvis-devops ALL=(ALL) NOPASSWD: /usr/sbin/nginx -t
jarvis-devops ALL=(ALL) NOPASSWD: /usr/sbin/nginx -s reload
jarvis-devops ALL=(ALL) NOPASSWD: /bin/systemctl restart nginx
jarvis-devops ALL=(ALL) NOPASSWD: /bin/systemctl stop nginx
jarvis-devops ALL=(ALL) NOPASSWD: /bin/systemctl start nginx
jarvis-devops ALL=(ALL) NOPASSWD: /bin/journalctl -u nginx *
```

#### 4. Criar service file

```bash
sudo nano /etc/systemd/system/jarvis-devops.service
```

Conteúdo:

```ini
[Unit]
Description=Jarvis DevOps - Nginx Configuration Manager
After=network.target nginx.service
Wants=nginx.service

[Service]
Type=simple
User=jarvis-devops
Group=jarvis-devops
WorkingDirectory=/opt/jarvis-devops
ExecStart=/opt/jarvis-devops/jarvis-devops
EnvironmentFile=/etc/jarvis-devops/.env
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Segurança
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ReadWritePaths=/etc/nginx/sites-available /var/log/jarvis-devops
ProtectHome=true
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true

[Install]
WantedBy=multi-user.target
```

#### 5. Habilitar e iniciar o serviço

```bash
# Recarregar systemd
sudo systemctl daemon-reload

# Habilitar para iniciar automaticamente
sudo systemctl enable jarvis-devops

# Iniciar o serviço
sudo systemctl start jarvis-devops

# Verificar status
sudo systemctl status jarvis-devops

# Verificar logs
sudo journalctl -u jarvis-devops -f
```

## 🔧 Configuração Avançada

### Proxy Reverso com Nginx

Para usar o Jarvis DevOps em produção, configure um proxy reverso:

```bash
sudo nano /etc/nginx/sites-available/jarvis-devops.conf
```

```nginx
server {
    listen 80;
    server_name jarvis-devops.exemplo.com;

    # Redirecionar para HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name jarvis-devops.exemplo.com;

    # Configuração SSL
    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Headers de segurança
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

    # Proxy para aplicação
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support (se necessário no futuro)
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=jarvis:10m rate=10r/m;
    limit_req zone=jarvis burst=5;
}
```

```bash
# Habilitar site
sudo ln -s /etc/nginx/sites-available/jarvis-devops.conf /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Configuração de Firewall

**UFW (Ubuntu):**

```bash
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 80/tcp      # HTTP
sudo ufw allow 443/tcp     # HTTPS
sudo ufw enable
```

**Firewalld (CentOS/RHEL):**

```bash
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --reload
```

### Backup Automático

Criar script de backup:

```bash
sudo nano /opt/jarvis-devops/backup.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/var/backups/jarvis-devops"
DATE=$(date +%Y%m%d_%H%M%S)

# Criar diretório de backup
mkdir -p "$BACKUP_DIR"

# Backup da configuração
tar -czf "$BACKUP_DIR/config_$DATE.tar.gz" /etc/jarvis-devops/

# Backup dos logs
tar -czf "$BACKUP_DIR/logs_$DATE.tar.gz" /var/log/jarvis-devops/

# Manter apenas os últimos 30 backups
find "$BACKUP_DIR" -name "*.tar.gz" -type f -mtime +30 -delete

echo "Backup completed: $DATE"
```

```bash
sudo chmod +x /opt/jarvis-devops/backup.sh

# Adicionar ao crontab
echo "0 2 * * * /opt/jarvis-devops/backup.sh" | sudo crontab -
```

## 🔒 Configuração de Segurança

### 1. Configurar HTTPS

Use Let's Encrypt para certificados gratuitos:

```bash
# Instalar certbot
sudo apt install certbot python3-certbot-nginx

# Obter certificado
sudo certbot --nginx -d jarvis-devops.exemplo.com

# Renovação automática
sudo crontab -e
# Adicionar: 0 12 * * * /usr/bin/certbot renew --quiet
```

### 2. Configurar Fail2Ban

```bash
sudo apt install fail2ban

sudo nano /etc/fail2ban/jail.local
```

```ini
[nginx-http-auth]
enabled = true
filter = nginx-http-auth
port = http,https
logpath = /var/log/nginx/error.log
maxretry = 3
bantime = 3600
```

### 3. Monitoramento de Logs

```bash
# Instalar logrotate para jarvis-devops
sudo nano /etc/logrotate.d/jarvis-devops
```

```
/var/log/jarvis-devops/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 jarvis-devops jarvis-devops
    postrotate
        systemctl reload jarvis-devops
    endscript
}
```

## 🚦 Verificação da Instalação

### 1. Testes de Funcionalidade

```bash
# Verificar se serviço está rodando
sudo systemctl status jarvis-devops

# Testar API
curl -u admin:senha http://localhost:8080/api/status

# Verificar logs
sudo journalctl -u jarvis-devops -n 50

# Testar interface web
curl -I http://localhost:8080/
```

### 2. Checklist de Segurança

- [ ] Senha forte configurada para Basic Auth
- [ ] HTTPS configurado em produção
- [ ] Firewall configurado
- [ ] Logs sendo monitorados
- [ ] Backup automático configurado
- [ ] Fail2Ban configurado
- [ ] Usuário dedicado criado
- [ ] Permissões mínimas configuradas

### 3. Monitoramento

```bash
# Script para monitoramento básico
nano /opt/jarvis-devops/health-check.sh
```

```bash
#!/bin/bash
SERVICE="jarvis-devops"
URL="http://localhost:8080/api/status"
AUTH="admin:sua_senha"

# Verificar se serviço está rodando
if ! systemctl is-active --quiet "$SERVICE"; then
    echo "ERROR: $SERVICE is not running"
    exit 1
fi

# Verificar se API responde
if ! curl -s -u "$AUTH" "$URL" | grep -q "is_installed"; then
    echo "ERROR: API not responding"
    exit 1
fi

echo "OK: $SERVICE is healthy"
```

## 🔧 Solução de Problemas

### Problemas Comuns

#### 1. Erro de Permissão

```bash
# Verificar permissões
ls -la /etc/nginx/sites-available/
sudo chown jarvis-devops:nginx /etc/nginx/sites-available/
sudo chmod 664 /etc/nginx/sites-available/*
```

#### 2. Nginx não encontrado

```bash
# Verificar caminho do nginx
which nginx
# Atualizar NGINX_BINARY no .env
```

#### 3. Porta já em uso

```bash
# Verificar o que está usando a porta
sudo netstat -tlnp | grep :8080
# Alterar SERVER_PORT no .env
```

#### 4. Logs do serviço

```bash
# Ver logs em tempo real
sudo journalctl -u jarvis-devops -f

# Ver logs de erro
sudo journalctl -u jarvis-devops -p err

# Ver logs de sistema
sudo journalctl -xe
```

### Comandos Úteis

```bash
# Reiniciar serviço
sudo systemctl restart jarvis-devops

# Ver configuração carregada
sudo systemctl show jarvis-devops

# Testar configuração do nginx
sudo nginx -t

# Recarregar nginx
sudo systemctl reload nginx

# Ver status de todos os serviços
sudo systemctl list-units --type=service --state=running
```

## 📞 Suporte

Para problemas de instalação:

1. Verificar logs do sistema: `sudo journalctl -xe`
2. Verificar logs da aplicação: `sudo journalctl -u jarvis-devops`
3. Testar conectividade: `curl -v http://localhost:8080`
4. Verificar configurações: `cat /etc/jarvis-devops/.env`

Para suporte adicional, abra uma issue no repositório do projeto.
