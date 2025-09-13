# ⚠️ IMPORTANTE: Esta pasta não é mais usada em produção

Os arquivos desta pasta são mantidos apenas para referência durante o desenvolvimento.

## 🔄 Nova Estrutura

Os assets estáticos agora são **embutidos no binário** e estão localizados em:

```
internal/assets/web/
├── static/
│   ├── css/style.css
│   └── js/ (app.js, editor.js)
└── templates/ (*.html)
```

## 📝 Para Modificar Assets

1. **NÃO edite arquivos nesta pasta (`web/`)**
2. **Edite os arquivos em `internal/assets/web/`**
3. **Faça rebuild do binário** para ver as mudanças

## 🚀 Build e Deploy

```bash
# Desenvolvimento (usa arquivos do filesystem)
go run ./cmd/server/main.go

# Produção (usa assets embutidos)
make build-optimized
./jarvis-devops-linux-amd64
```

## 🗑️ Esta pasta pode ser removida?

Sim, mas recomendamos manter durante a transição para garantir que nada foi esquecido.

Para testar se tudo está funcionando sem esta pasta:
```bash
mv web web.backup
make build && ./jarvis-devops
# Se funcionar, pode remover web.backup
```