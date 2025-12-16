# Exemplos de Uso com Diferentes Versões OpenAPI

Este documento demonstra como gerar documentação OpenAPI em diferentes versões usando o nexs-swag.

## Estrutura Básica

Considere uma API Go simples:

```go
package main

import (
    "github.com/gin-gonic/gin"
)

// @title           API de Exemplo
// @version         1.0
// @description     API de demonstração do nexs-swag
// @contact.name    Suporte da API
// @contact.email   support@example.com
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
    r := gin.Default()
    
    // @Summary      Lista usuários
    // @Description  Obtém lista de todos os usuários
    // @Tags         users
    // @Accept       json
    // @Produce      json
    // @Success      200  {array}  User
    // @Router       /users [get]
    r.GET("/users", getUsers)
    
    r.Run(":8080")
}

type User struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"João Silva"`
    Email string `json:"email" example:"joao@example.com"`
}

func getUsers(c *gin.Context) {
    // implementação
}
```

## Gerando OpenAPI 3.2.0 (Mais Recente)

```bash
# Gera OpenAPI 3.2.0 com todos os recursos mais recentes
nexs-swag init --openapi-version 3.2.0

# Ou usando o atalho
nexs-swag init --ov 3.2
```

**Arquivo gerado** (`docs/openapi.json`):
```json
{
  "openapi": "3.2.0",
  "info": {
    "title": "API de Exemplo",
    "description": "API de demonstração do nexs-swag",
    "contact": {
      "name": "Suporte da API",
      "email": "support@example.com"
    },
    "version": "1.0"
  },
  "servers": [
    {
      "url": "http://localhost:8080/api/v1"
    }
  ],
  "paths": {
    "/users": {
      "get": {
        "summary": "Lista usuários",
        "description": "Obtém lista de todos os usuários",
        "tags": ["users"],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/User"
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "User": {
        "type": "object",
        "properties": {
          "id": {
            "type": "integer",
            "example": 1
          },
          "name": {
            "type": "string",
            "example": "João Silva"
          },
          "email": {
            "type": "string",
            "example": "joao@example.com"
          }
        }
      }
    }
  }
}
```

## Gerando OpenAPI 3.1.2

```bash
# Gera OpenAPI 3.1.2 com compatibilidade JSON Schema 2020-12
nexs-swag init --openapi-version 3.1.2

# Ou
nexs-swag init --ov 3.1
```

**Características específicas v3.1:**
- Compatível com JSON Schema Draft 2020-12
- Suporte a `webhooks`
- Campo `jsonSchemaDialect`
- Pode usar `$schema` em schemas

## Gerando OpenAPI 3.0.4

```bash
# Gera OpenAPI 3.0.4 para máxima compatibilidade com ferramentas
nexs-swag init --openapi-version 3.0.4

# Ou
nexs-swag init --ov 3.0
```

**Ideal para:**
- Ferramentas que ainda não suportam OpenAPI 3.1+
- Integrações com sistemas mais antigos
- Compatibilidade ampla

## Gerando Swagger 2.0 (Legado)

```bash
# Gera Swagger 2.0 para sistemas legados
nexs-swag init --openapi-version 2.0

# Ou
nexs-swag init --ov 2
```

**Arquivo gerado** (`docs/swagger.json`):
```json
{
  "swagger": "2.0",
  "info": {
    "title": "API de Exemplo",
    "description": "API de demonstração do nexs-swag",
    "contact": {
      "name": "Suporte da API",
      "email": "support@example.com"
    },
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/api/v1",
  "paths": {
    "/users": {
      "get": {
        "summary": "Lista usuários",
        "description": "Obtém lista de todos os usuários",
        "tags": ["users"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/User"
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "User": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "example": 1
        },
        "name": {
          "type": "string",
          "example": "João Silva"
        },
        "email": {
          "type": "string",
          "example": "joao@example.com"
        }
      }
    }
  }
}
```

## Recursos Avançados por Versão

### OpenAPI 3.2.0 - Streaming e QUERY Method

```go
// @Summary      Stream de eventos
// @Description  Endpoint de Server-Sent Events
// @Tags         events
// @Accept       json
// @Produce      text/event-stream
// @Success      200  {object}  Event
// @Router       /events [get]
func streamEvents(c *gin.Context) {
    // implementação SSE
}

// @Summary      Busca avançada
// @Description  Busca usando método QUERY (OpenAPI 3.2+)
// @Tags         search
// @Accept       json
// @Produce      json
// @Success      200  {array}  SearchResult
// @Router       /search [query]
func querySearch(c *gin.Context) {
    // implementação
}
```

### OpenAPI 3.1.x - Webhooks

```go
// @title           API com Webhooks
// @version         1.0
// @description     Demonstra webhooks do OpenAPI 3.1

// @webhook         orderCreated
// @description     Webhook disparado quando um pedido é criado
// @tags            webhooks
// @accept          json
// @param           order  body  Order  true  "Dados do pedido"
// @success         200
func main() {
    // código da API
}
```

### Comparação de Features

| Feature | v2.0 | v3.0.x | v3.1.x | v3.2.0 |
|---------|------|--------|--------|--------|
| Servers | ❌ | ✅ | ✅ | ✅ |
| Components | ❌ | ✅ | ✅ | ✅ |
| Webhooks | ❌ | ❌ | ✅ | ✅ |
| JSON Schema 2020-12 | ❌ | ❌ | ✅ | ✅ |
| QUERY method | ❌ | ❌ | ❌ | ✅ |
| Streaming (itemSchema) | ❌ | ❌ | ❌ | ✅ |
| OAuth2 Device Flow | ❌ | ❌ | ❌ | ✅ |

## Scripts de Build

### Makefile

```makefile
.PHONY: docs-3.2 docs-3.1 docs-3.0 docs-2.0

# Gera todas as versões
docs-all: docs-3.2 docs-3.1 docs-3.0 docs-2.0

# OpenAPI 3.2.0 (default)
docs-3.2:
	nexs-swag init --ov 3.2 --output ./docs/v3.2

# OpenAPI 3.1.2
docs-3.1:
	nexs-swag init --ov 3.1 --output ./docs/v3.1

# OpenAPI 3.0.4
docs-3.0:
	nexs-swag init --ov 3.0 --output ./docs/v3.0

# Swagger 2.0
docs-2.0:
	nexs-swag init --ov 2.0 --output ./docs/v2.0

# Valida todas as especificações
validate: docs-all
	@echo "Validando especificações geradas..."
	npx @apidevtools/swagger-cli validate docs/v3.2/openapi.yaml
	npx @apidevtools/swagger-cli validate docs/v3.1/openapi.yaml
	npx @apidevtools/swagger-cli validate docs/v3.0/openapi.yaml
	npx @apidevtools/swagger-cli validate docs/v2.0/swagger.yaml
```

### Script Shell

```bash
#!/bin/bash
# generate-docs.sh

set -e

echo "Gerando documentação OpenAPI em múltiplas versões..."

# OpenAPI 3.2.0
echo "📝 Gerando OpenAPI 3.2.0..."
nexs-swag init --ov 3.2.0 --output ./docs/v3.2 --quiet

# OpenAPI 3.1.2
echo "📝 Gerando OpenAPI 3.1.2..."
nexs-swag init --ov 3.1.2 --output ./docs/v3.1 --quiet

# OpenAPI 3.0.4
echo "📝 Gerando OpenAPI 3.0.4..."
nexs-swag init --ov 3.0.4 --output ./docs/v3.0 --quiet

# Swagger 2.0
echo "📝 Gerando Swagger 2.0..."
nexs-swag init --ov 2.0 --output ./docs/v2.0 --quiet

echo "✅ Todas as versões foram geradas com sucesso!"
echo ""
echo "Arquivos gerados:"
ls -lh docs/*/openapi.{json,yaml} docs/*/swagger.{json,yaml} 2>/dev/null || true
```

## GitHub Actions

```yaml
name: Generate OpenAPI Docs

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  generate-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Install nexs-swag
        run: go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest
      
      - name: Generate OpenAPI 3.2.0
        run: nexs-swag init --ov 3.2 --output ./docs/v3.2
      
      - name: Generate OpenAPI 3.1.2
        run: nexs-swag init --ov 3.1 --output ./docs/v3.1
      
      - name: Generate Swagger 2.0
        run: nexs-swag init --ov 2.0 --output ./docs/v2.0
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: openapi-docs
          path: docs/
```

## Escolhendo a Versão Certa

### Use OpenAPI 3.2.0 quando:
- ✅ Você precisa dos recursos mais recentes
- ✅ Quer suporte para streaming e SSE
- ✅ Precisa do método HTTP QUERY
- ✅ Está começando um novo projeto

### Use OpenAPI 3.1.x quando:
- ✅ Precisa de webhooks
- ✅ Quer compatibilidade total com JSON Schema 2020-12
- ✅ Suas ferramentas suportam 3.1 mas não 3.2

### Use OpenAPI 3.0.x quando:
- ✅ Precisa de ampla compatibilidade com ferramentas
- ✅ Trabalha com sistemas que não suportam 3.1+
- ✅ Quer evitar possíveis incompatibilidades

### Use Swagger 2.0 quando:
- ✅ Precisa dar suporte a sistemas legados
- ✅ Ferramentas antigas que não suportam OpenAPI 3.x
- ✅ Clientes que ainda usam Swagger 2.0

## Referências

- [OpenAPI 3.2.0 Specification](https://spec.openapis.org/oas/v3.2.0.html)
- [OpenAPI 3.1.0 Specification](https://spec.openapis.org/oas/v3.1.0.html)
- [OpenAPI 3.0.0 Specification](https://spec.openapis.org/oas/v3.0.0.html)
- [Swagger 2.0 Specification](https://spec.openapis.org/oas/v2.0.html)
- [Documentação completa de versões](OPENAPI_VERSIONS.md)
