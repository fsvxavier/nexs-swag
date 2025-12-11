# Exemplo 20 - Fmt Command

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra o comando `fmt` para formatar annotations Swagger.

## Comando

```bash
nexs-swag fmt [flags]
```

## Flags

```bash
--dir <directory>         # Diretório a formatar (default: .)
--exclude <pattern>       # Padrões a excluir
--ext <extensions>        # Extensões (default: .go)
```

## Uso Básico

```bash
# Formatar diretório atual
nexs-swag fmt

# Formatar diretório específico
nexs-swag fmt --dir ./api

# Excluir diretórios
nexs-swag fmt --exclude "vendor,testdata"
```

## O Que Faz

### Antes
```go
// GetUser returns a user
// @Summary Get user
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} Error
// @Router /users/{id} [get]
func GetUser() {}
```

### Depois
```go
// GetUser returns a user
// @Summary      Get user
// @Description  Get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  User
// @Failure      404  {object}  Error
// @Router       /users/{id} [get]
func GetUser() {}
```

## Como Funciona

1. **Detecta annotations Swagger**
```go
// @Summary
// @Description
// @Tags
// @Param
// @Success
// @Failure
// @Router
// etc.
```

2. **Alinha colunas**
```
// @Param   name  query  string  true  "Description"
           ^^^^  ^^^^^  ^^^^^^  ^^^^  ^^^^^^^^^^^^^
           col1  col2   col3    col4  col5
```

3. **Preserva outros comentários**
```go
// Este comentário não muda
// TODO: implementar validação
// @Summary Get data  <- Só este é formatado
func GetData() {}
```

## Como Executar

```bash
./run.sh
```

## Casos de Uso

### 1. Pre-commit Hook
```bash
#!/bin/bash
# .git/hooks/pre-commit
nexs-swag fmt
git add -u
```

### 2. CI/CD Validation
```bash
# .github/workflows/docs.yml
- name: Check swagger format
  run: |
    nexs-swag fmt
    git diff --exit-code || (echo "Run nexs-swag fmt" && exit 1)
```

### 3. Editor Integration
```json
// VSCode settings.json
{
  "emeraldwalk.runonsave": {
    "commands": [
      {
        "match": "\\.go$",
        "cmd": "nexs-swag fmt --dir ${fileDirname}"
      }
    ]
  }
}
```

### 4. Make Target
```makefile
.PHONY: fmt-swagger
fmt-swagger:
	nexs-swag fmt
	
.PHONY: fmt
fmt: fmt-swagger
	go fmt ./...
```

## Annotations Suportadas

```go
@Summary
@Description
@Tags
@Accept
@Produce
@Param
@Success
@Failure
@Header
@Router
@Security
@Deprecated
@ID
@x-codeSamples
```

## Comparação com gofmt

| Ferramenta | Escopo |
|------------|--------|
| `go fmt` | Formata código Go |
| `nexs-swag fmt` | Formata annotations Swagger |

Use ambos:
```bash
# Formatar tudo
go fmt ./...
nexs-swag fmt
```

## Opções Avançadas

### Processar apenas arquivos alterados
```bash
# Git
git diff --name-only | grep '\.go$' | xargs -I {} nexs-swag fmt --dir $(dirname {})

# Todos os arquivos staged
git diff --cached --name-only | grep '\.go$' | xargs -I {} nexs-swag fmt --dir $(dirname {})
```

### Formatar e gerar docs
```bash
nexs-swag fmt && nexs-swag init
```

### Verificar sem modificar (dry-run)
```bash
# Copiar arquivos primeiro
cp -r api api.backup
nexs-swag fmt --dir api
diff -r api api.backup
rm -rf api.backup
```

## Benefícios

✅ **Legibilidade:** Código mais limpo e organizado
✅ **Consistência:** Estilo uniforme em todo o projeto
✅ **Manutenibilidade:** Mais fácil de revisar e atualizar
✅ **Colaboração:** Time segue mesmo padrão
✅ **Automação:** Integra com ferramentas de desenvolvimento

## Limitações

⚠️ **Modifica arquivos:** Sempre commite antes de rodar
⚠️ **Não valida:** Apenas formata, não checa erros
⚠️ **Comentários especiais:** Pode afetar formatação customizada

## Best Practices

1. **Sempre use controle de versão**
```bash
git status  # Verificar antes
nexs-swag fmt
git diff    # Revisar mudanças
```

2. **Automatize no workflow**
```bash
# Pre-commit ou CI/CD
nexs-swag fmt --dir .
```

3. **Combine com linters**
```bash
nexs-swag fmt
golangci-lint run
```

4. **Documente no README**
```markdown
## Development

Format Swagger annotations:
\`\`\`bash
nexs-swag fmt
\`\`\`
```

## Recomendação

✅ **Use nexs-swag fmt regularmente**

Integre no workflow de desenvolvimento para manter documentação consistente e legível.
