# Exemplo 12 - Markdown Files

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra como usar arquivos markdown para descriptions detalhadas.

## Flag

```bash
--markdownFiles <directory>
--md <directory>
```

## Uso

```bash
nexs-swag init --markdownFiles ./docs
```

## Sintaxe

No código Go, use `file(nome.md)` na description:

```go
// @Description file(create-user.md)
// @Router /users [post]
func CreateUser() {}
```

O nexs-swag vai:
1. Ler `docs/create-user.md`
2. Substituir `file(create-user.md)` pelo conteúdo do arquivo
3. Adicionar o conteúdo na description do OpenAPI

## Estrutura

```
12-markdown-files/
├── main.go                   # API com file() references
└── docs/
    ├── create-user.md        # Description detalhada
    └── get-user.md           # Outra description
```

## Benefícios

### 1. Descriptions Detalhadas
```markdown
# Create User

Creates a new user with validation.

## Request Body
- name: string (required)
- email: string (required, valid email)

## Validation Rules
...
```

### 2. Separação de Concerns
- Código Go: lógica da aplicação
- Markdown: documentação detalhada
- Mais fácil de manter

### 3. Reutilização
```go
// Múltiplos endpoints podem usar o mesmo markdown
// @Description file(auth-required.md)
func Endpoint1() {}

// @Description file(auth-required.md)
func Endpoint2() {}
```

### 4. Formatação Rica
- Headers
- Listas
- Code blocks
- Tabelas
- Links

## Exemplo Real

```go
// @Description file(user-endpoints.md)
// @Router /users [post]
func CreateUser() {}
```

**user-endpoints.md:**
```markdown
# User Management

## Authentication Required
All user endpoints require Bearer token.

## Rate Limiting
- 100 requests/minute

## Response Codes
| Code | Description |
|------|-------------|
| 200  | Success |
| 401  | Unauthorized |
| 429  | Rate limit exceeded |
```

## Como Executar

```bash
./run.sh
```

## Comparação

### Sem Markdown
```json
{
  "description": "file(create-user.md)"
}
```

### Com Markdown
```json
{
  "description": "# Create User Endpoint\n\nCreates a new user in the system...\n\n## Request Body\n..."
}
```

## Casos de Uso

- APIs complexas com muita documentação
- Documentação colaborativa (tech writers)
- Versionamento de docs separado do código
- Geração de documentação a partir de Wiki
