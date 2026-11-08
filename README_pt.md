# nexs-swag

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

[![Versão Go](https://img.shields.io/badge/Go-1.27.1+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1.0-6BA539?style=flat&logo=openapiinitiative)](https://spec.openapis.org/oas/v3.1.0)
[![Swagger](https://img.shields.io/badge/Swagger-2.0-85EA2D?style=flat&logo=swagger)](https://swagger.io/specification/v2/)
[![Licença](https://img.shields.io/badge/Licença-Apache_2.0-blue.svg)](LICENSE)
[![Cobertura](https://img.shields.io/badge/Cobertura-80.1%25-brightgreen.svg)](/)
[![Exemplos](https://img.shields.io/badge/Exemplos-27-blue.svg)](examples/)

**Gere automaticamente documentação OpenAPI 3.1.0 ou Swagger 2.0 a partir de anotações no código Go.**

nexs-swag converte anotações Go para especificação OpenAPI 3.1.0 ou Swagger 2.0. Foi projetado como uma evolução do [swaggo/swag](https://github.com/swaggo/swag) com suporte completo para a especificação OpenAPI mais recente e compatibilidade total com Swagger 2.0.

## Índice

- [Visão Geral](#visão-geral)
- [Primeiros Passos](#primeiros-passos)
  - [Instalação](#instalação)
  - [Início Rápido](#início-rápido)
- [Frameworks Web Suportados](#frameworks-web-suportados)
- [Como usar com Gin](#como-usar-com-gin)
- [Referência CLI](#referência-cli)
  - [Comando init](#comando-init)
  - [Comando fmt](#comando-fmt)
- [Status de Implementação](#status-de-implementação)
- [Versões OpenAPI](OPENAPI_VERSIONS.md) - Guia completo de todas as versões suportadas
- [Formato de Comentários Declarativos](#formato-de-comentários-declarativos)
  - [Informações Gerais da API](#informações-gerais-da-api)
  - [Operação de API](#operação-de-api)
  - [Tags de Struct](#tags-de-struct)
- [Exemplos](#exemplos)
- [Qualidade e Testes](#qualidade-e-testes)
- [Compatibilidade com swaggo/swag](#compatibilidade-com-swaggoswag)
- [Sobre o Projeto](#sobre-o-projeto)
- [Contribuindo](#contribuindo)
- [Licença](#licença)

## Visão Geral

### Recursos Principais

- ✅ **100% compatível com swaggo/swag** - Substituto direto com todas as anotações e tags
- ✅ **Suporte a múltiplas versões OpenAPI** - Gere v2.0.0, v3.0.x, v3.1.x ou v3.2.x (3.2.0, 3.2.1)
- ✅ **OpenAPI 3.2.x** - Suporte completo para a versão mais recente (método QUERY, streaming, etc)
- ✅ **OpenAPI 3.1.x** - Compatível com JSON Schema 2020-12, webhooks e recursos modernos
- ✅ **OpenAPI 3.0.x** - Todas as versões desde 3.0.0 até 3.0.4
- ✅ **Swagger 2.0** - Compatibilidade total com sistemas legados
- ✅ **Conversão automática** - Conversão entre formatos com avisos para incompatibilidades
- ✅ **20+ atributos de validação** - minimum, maximum, pattern, enum, format e mais
- ✅ **Validação de frameworks** - Suporte nativo para Gin (binding) e go-playground/validator
- ✅ **Headers de resposta** - Documentação completa de headers
- ✅ **Múltiplos tipos de conteúdo** - JSON, XML, YAML, CSV, PDF e tipos MIME customizados
- ✅ **Extensões customizadas** - Suporte completo para x-*
- ✅ **@x-visibility** - Gere documentação pública/privada separada a partir de uma única base de código
- ✅ **80.1% de cobertura de testes** - Pronto para produção com suite de testes abrangente incluindo testes roundtrip
- ✅ **27 exemplos funcionais** - Aprenda com exemplos completos e executáveis

### Por que nexs-swag?

| Recurso | swaggo/swag | nexs-swag |
|---------|-------------|-----------|
| OpenAPI 3.2.x | ❌ | ✅ |
| OpenAPI 3.1.x | ❌ | ✅ |
| OpenAPI 3.0.x | ❌ | ✅ |
| Swagger 2.0 | ✅ | ✅ |
| Múltiplas Versões | ❌ | ✅ (todas do mesmo código) |
| JSON Schema | Draft 4 | Draft 4 + 2020-12 |
| Webhooks | ❌ | ✅ (OpenAPI 3.1+) |
| Headers de Resposta | Limitado | Suporte Completo |
| Suporte a Nullable | `x-nullable` | Nativo + `x-nullable` |
| Cobertura de Testes | ~70% | 80.1% |
| Exemplos | ~10 | 25 |
| Versão Go | 1.19+ | 1.27.1+ |

## Primeiros Passos

### Instalação

#### Usando go install (Recomendado)

```bash
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest
```

Para verificar a instalação:

```bash
nexs-swag --version
```

#### Compilando do Código Fonte

Requer [Go 1.27.1 ou superior](https://go.dev/dl/).

```bash
git clone https://github.com/fsvxavier/nexs-swag.git
cd nexs-swag
go build -o nexs-swag ./cmd/nexs-swag
```

#### Usando Docker

```bash
docker pull ghcr.io/fsvxavier/nexs-swag:latest
docker run --rm -v $(pwd):/app ghcr.io/fsvxavier/nexs-swag:latest init
```

### Início Rápido

#### 1. Adicionar Anotações da API

Adicione anotações gerais da API ao seu `main.go`:

```go
package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
)

// @title           API de Gerenciamento de Usuários
// @version         1.0.0
// @description     Uma API de gerenciamento de usuários com documentação OpenAPI 3.1.0 completa
// @termsOfService  http://swagger.io/terms/

// @contact.name   Suporte da API
// @contact.url    http://www.example.com/suporte
// @contact.email  suporte@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
    r := gin.Default()
    // Configuração da sua aplicação
    r.Run(":8080")
}

// User representa um usuário do sistema
type User struct {
    // ID do usuário (sql.NullInt64 → integer no OpenAPI)
    ID sql.NullInt64 `json:"id" swaggertype:"integer" extensions:"x-primary-key=true"`
    
    // Nome completo (3-100 caracteres obrigatório)
    Name string `json:"name" binding:"required" minLength:"3" maxLength:"100" example:"João Silva"`
    
    // Endereço de email (validado)
    Email string `json:"email" binding:"required,email" format:"email" extensions:"x-unique=true"`
    
    // Senha (oculta da documentação)
    Password string `json:"password" swaggerignore:"true"`
    
    // Status da conta
    Status string `json:"status" enum:"active,inactive,pending" default:"active"`
    
    // Saldo da conta
    Balance float64 `json:"balance" minimum:"0" extensions:"x-currency=BRL"`
}

// CreateUser cria um novo usuário
// @Summary      Criar usuário
// @Description  Cria um novo usuário no sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      User  true  "Dados do usuário"
// @Success      201   {object}  User
// @Header       201   {string}  X-Request-ID  "Identificador da requisição"
// @Header       201   {string}  Location      "URL do recurso do usuário"
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /users [post]
// @Security     ApiKeyAuth
func CreateUser(c *gin.Context) {
    // Implementação
}
```

#### 2. Gerar Documentação

**OpenAPI 3.1.0 (padrão):**

```bash
nexs-swag init
# ou explicitamente
nexs-swag init --openapi-version 3.1
```

**Swagger 2.0:**

```bash
nexs-swag init --openapi-version 2.0
```

**Gerar ambas as versões:**

```bash
# OpenAPI 3.1.0 em ./docs/v3
nexs-swag init --output ./docs/v3 --openapi-version 3.1

# Swagger 2.0 em ./docs/v2
nexs-swag init --output ./docs/v2 --openapi-version 2.0
```

Ou especifique os diretórios:

```bash
nexs-swag init -d ./cmd/api -o ./docs --openapi-version 3.1
```

#### 3. Arquivos Gerados

**OpenAPI 3.1.0 (padrão):**
- **`docs/openapi.json`** - Especificação OpenAPI 3.1.0 em JSON
- **`docs/openapi.yaml`** - Especificação OpenAPI 3.1.0 em YAML
- **`docs/docs.go`** - Arquivo de documentação Go embarcado

**Swagger 2.0 (com `--openapi-version 2.0`):**
- **`docs/swagger.json`** - Especificação Swagger 2.0 em JSON
- **`docs/swagger.yaml`** - Especificação Swagger 2.0 em YAML
- **`docs/docs.go`** - Arquivo de documentação Go embarcado

#### 4. Integrar com Sua Aplicação

Importe o pacote docs gerado:

```go
import _ "seu-modulo/docs"  // Importar docs gerado

func main() {
    r := gin.Default()
    
    // Servir Swagger UI
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    r.Run(":8080")
}
```

Acesse http://localhost:8080/swagger/index.html para ver sua documentação API!

## Frameworks Web Suportados

nexs-swag funciona com todos os frameworks web Go populares através de pacotes middleware swagger:

- [gin](https://github.com/swaggo/gin-swagger) - `github.com/swaggo/gin-swagger`
- [echo](https://github.com/swaggo/echo-swagger) - `github.com/swaggo/echo-swagger`
- [fiber](https://github.com/gofiber/swagger) - `github.com/gofiber/swagger`
- [net/http](https://github.com/swaggo/http-swagger) - `github.com/swaggo/http-swagger`
- [gorilla/mux](https://github.com/swaggo/http-swagger) - `github.com/swaggo/http-swagger`
- [go-chi/chi](https://github.com/swaggo/http-swagger) - `github.com/swaggo/http-swagger`
- [hertz](https://github.com/hertz-contrib/swagger) - `github.com/hertz-contrib/swagger`
- [buffalo](https://github.com/swaggo/buffalo-swagger) - `github.com/swaggo/buffalo-swagger`

## Como usar com Gin

Exemplo completo usando framework Gin. Encontre o código completo em [examples/03-general-info](examples/03-general-info).

**1. Instalar dependências:**

```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

**2. Adicionar informações gerais da API ao `main.go`:**

```go
package main

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    
    _ "seu-projeto/docs"  // Importar docs gerado
)

// @title           API de Exemplo Swagger
// @version         1.0
// @description     Este é um servidor de exemplo com nexs-swag.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Suporte da API
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization

func main() {
    r := gin.Default()
    
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users/:id", GetUser)
        v1.POST("/users", CreateUser)
    }
    
    // Endpoint Swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    r.Run(":8080")
}
```

**3. Adicionar anotações de operação:**

```go
// GetUser recupera um usuário por ID
// @Summary      Buscar usuário por ID
// @Description  Buscar detalhes do usuário pelo seu identificador único
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Usuário"  minimum(1)
// @Success      200  {object}  User
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [get]
// @Security     ApiKeyAuth
func GetUser(c *gin.Context) {
    // Implementação
}
```

**4. Gerar e executar:**

```bash
nexs-swag init
go run main.go
```

Visite http://localhost:8080/swagger/index.html

## Referência CLI

### Comando init

Gera documentação OpenAPI a partir do código fonte.

```bash
nexs-swag init [opções]
```

**Opções Principais:**

| Flag | Curto | Padrão | Descrição |
|------|-------|--------|-----------|
| `--generalInfo` | `-g` | `main.go` | Caminho para arquivo com informações gerais da API |
| `--dir` | `-d` | `./` | Diretórios para analisar (separados por vírgula) |
| `--output` | `-o` | `./docs` | Diretório de saída para arquivos gerados |
| `--outputTypes` | `--ot` | `go,json,yaml` | Tipos de arquivo de saída |
| `--parseDepth` | | `100` | Profundidade de análise de dependência |
| `--parseDependency` | `--pd` | `false` | Analisar arquivos go em dependências |
| `--parseDependencyLevel` | `--pdl` | `0` | 0=desabilitado, 1=modelos, 2=operações, 3=tudo |
| `--parseInternal` | | `false` | Analisar pacotes internos |
| `--parseGoList` | | `true` | Usar `go list` para análise |
| `--propertyStrategy` | `-p` | `camelcase` | Nomenclatura de propriedade: `snakecase`, `camelcase`, `pascalcase` |
| `--requiredByDefault` | | `false` | Marcar todos os campos como obrigatórios |
| `--validate` | | `true` | Validar especificação gerada |
| `--exclude` | | | Excluir diretórios (separados por vírgula) |
| `--tags` | `-t` | | Filtrar por tags (separados por vírgula) |
| `--markdownFiles` | `--md` | | Analisar arquivos markdown para descrições |
| `--codeExampleFiles` | `--cef` | | Analisar arquivos de exemplo de código |
| `--generatedTime` | | `false` | Adicionar timestamp de geração |
| `--instanceName` | | `swagger` | Nome da instância para múltiplos docs |
| `--overridesFile` | | `.swaggo` | Arquivo de overrides de tipo |
| `--templateDelims` | `--td` | `{{,}}` | Delimitadores de template customizados |
| `--collectionFormat` | `--cf` | `csv` | Formato de array padrão |
| `--parseFuncBody` | | `false` | Analisar corpos de função |
| `--includeTypes` | `--it` | `all` | Filtrar tipos a incluir: `struct`, `interface`, `func`, `const`, `type`, `all` |
| `--openapi-version` | `--ov` | `3.1` | Versão OpenAPI: `2.0`, `3.0`, `3.1` |

> **⚠️ Importante: Sintaxe de Flags Booleanas**
>
> Flags booleanas aceitam duas sintaxes válidas:
> - ✅ **Sem valor** (presença = true): `--parseInternal`, `--pd`
> - ✅ **Com sinal de igual**: `--parseInternal=true`, `--pd=false`
> - ❌ **Errado** (separado por espaço): `--parseInternal true`, `--pd true`
>
> A sintaxe separada por espaço não funciona porque o parser CLI trata a palavra após a flag como um argumento posicional separado, não como o valor da flag.

**Exemplos:**

```bash
# Uso básico (OpenAPI 3.1.0)
nexs-swag init

# Gerar Swagger 2.0
nexs-swag init --openapi-version 2.0

# Gerar ambas as versões
nexs-swag init --output ./docs/v3 --openapi-version 3.1
nexs-swag init --output ./docs/v2 --openapi-version 2.0

# Especificar diretórios
nexs-swag init -d ./cmd/api,./internal/handlers -o ./api-docs

# Analisar dependências (nível 1 - apenas modelos)
nexs-swag init --parseDependency --parseDependencyLevel 1
# Ou com sintaxe explícita:
nexs-swag init --parseDependency=true --parseDependencyLevel 1

# Analisar pacotes internos
nexs-swag init --parseInternal
# Ou explicitamente:
nexs-swag init --parseInternal=true

# Saída apenas JSON
nexs-swag init --outputTypes json

# Nomes de propriedade em snake_case
nexs-swag init --propertyStrategy snakecase

# Filtrar por tags
nexs-swag init --tags "users,products"

# Usar descrições em markdown
nexs-swag init --markdownFiles ./docs/api

# Delimitadores de template customizados (evitar conflitos)
nexs-swag init --templateDelims "[[,]]"

# Filtrar tipos a incluir (apenas structs)
nexs-swag init --includeTypes struct

# Filtrar múltiplas categorias de tipos
nexs-swag init --includeTypes "struct,interface"

# Forma curta
nexs-swag init -it struct
```

### Comando fmt

Formata comentários swagger automaticamente.

```bash
nexs-swag fmt [opções]
```

**Opções:**

| Flag | Curto | Padrão | Descrição |
|------|-------|--------|-----------|
| `--dir` | `-d` | `./` | Diretórios para formatar |
| `--exclude` | | | Excluir diretórios |
| `--generalInfo` | `-g` | `main.go` | Arquivo de informações gerais |

**Exemplo:**

```bash
# Formatar diretório atual
nexs-swag fmt

# Formatar diretório específico
nexs-swag fmt -d ./internal/api

# Excluir vendor
nexs-swag fmt --exclude ./vendor
```

## Status de Implementação

### Suporte OpenAPI 3.1.0

✅ **Totalmente Implementado:**
- JSON Schema 2020-12
- Estrutura básica (Info, Servers, Paths, Components)
- Request bodies com múltiplos content types
- Definições de resposta com headers
- Definições de parâmetros (path, query, header, cookie)
- Security schemes (Basic, Bearer, API Key, OAuth2)
- Composição de schemas (allOf, oneOf, anyOf)
- Validação de schemas (min, max, pattern, enum, format)
- Exemplos e descrições
- Documentação externa
- Extensões customizadas (x-*)
- Webhooks
- Tags e agrupamento

### Suporte Swagger 2.0

✅ **Totalmente Compatível:**
- Estrutura básica (Info, Host, BasePath, Paths, Definitions)
- Definições de request/response
- Definições de parâmetros (path, query, header, body, formData)
- Definições de segurança (Basic, API Key, OAuth2)
- Composição de schemas (allOf)
- Validação de schemas (min, max, pattern, enum, format)
- Exemplos e descrições
- Documentação externa
- Extensões customizadas (x-*)
- Tags e agrupamento

⚠️ **Conversão Automática com Avisos:**
- Servers → Host + BasePath (usa a primeira URL de server)
- Webhooks → ⚠️ Não suportado em Swagger 2.0
- Callbacks → ⚠️ Não suportado em Swagger 2.0
- oneOf/anyOf → ⚠️ Suporte limitado (convertido para object)
- propriedade nullable → Usa extensão `x-nullable`

### Compatibilidade com swaggo/swag

✅ **100% Compatível:**
- Todas as anotações (@title, @version, @description, etc.)
- Todas as tags de struct (json, binding, validate, swaggertype, swaggerignore, extensions)
- Todas as flags CLI (28/28 implementadas)
- Comandos: init, fmt
- Type overrides via arquivo .swaggo
- Descrições em Markdown
- Exemplos de código

## Formato de Comentários Declarativos

### Informações Gerais da API

Adicione ao seu `main.go` ou ponto de entrada:

| Anotação | Exemplo | Descrição |
|----------|---------|-----------|
| `@title` | `@title Minha API` | **Obrigatório.** Título da API |
| `@version` | `@version 1.0` | **Obrigatório.** Versão da API |
| `@description` | `@description Esta é minha API` | Descrição da API |
| `@description.markdown` | `@description.markdown` | Carregar descrição de api.md |
| `@termsOfService` | `@termsOfService http://example.com/terms` | URL dos termos de serviço |
| `@contact.name` | `@contact.name Suporte da API` | Nome do contato |
| `@contact.url` | `@contact.url http://example.com` | URL do contato |
| `@contact.email` | `@contact.email support@example.com` | Email do contato |
| `@license.name` | `@license.name Apache 2.0` | **Obrigatório.** Nome da licença |
| `@license.url` | `@license.url http://apache.org/licenses` | URL da licença |
| `@host` | `@host localhost:8080` | Host da API |
| `@BasePath` | `@BasePath /api/v1` | Caminho base |
| `@schemes` | `@schemes http https` | Protocolos de transferência |
| `@accept` | `@accept json xml` | Tipos MIME Accept padrão |
| `@produce` | `@produce json xml` | Tipos MIME Produce padrão |
| `@tag.name` | `@tag.name Users` | Nome da tag |
| `@tag.description` | `@tag.description Operações de usuário` | Descrição da tag |
| `@externalDocs.description` | `@externalDocs.description OpenAPI` | Descrição de docs externos |
| `@externalDocs.url` | `@externalDocs.url https://swagger.io` | URL de docs externos |
| `@x-<nome>` | `@x-custom-info value` | Extensão customizada |

**Anotações Específicas de Versão:**

Ao gerar **Swagger 2.0** (`--openapi-version 2.0`):
- Use anotações `@host`, `@BasePath` e `@schemes`
- Estas são automaticamente convertidas para os campos `host`, `basePath` e `schemes`

Ao gerar **OpenAPI 3.x** (`--openapi-version 3.0` ou `3.1`):
- Use anotação `@server`: `// @server http://localhost:8080/api/v1 Servidor de desenvolvimento`
- Alternativamente, use `@host`, `@BasePath` e `@schemes` que serão convertidos para servers

Ambos os estilos de anotação funcionam com qualquer versão - o conversor lida com a transformação automaticamente.

**Definições de Segurança:**

```go
// Autenticação Basic
// @securityDefinitions.basic BasicAuth

// API Key
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// OAuth2 Application Flow
// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Concede acesso de escrita
// @scope.admin Concede acesso de administrador
```

### Operação de API

Adicione às funções handler:

| Anotação | Exemplo | Descrição |
|----------|---------|-----------|
| `@Summary` | `@Summary Buscar usuário` | Resumo curto |
| `@Description` | `@Description Buscar usuário por ID` | Descrição detalhada |
| `@Description.markdown` | `@Description.markdown details` | Carregar de details.md |
| `@Tags` | `@Tags users,accounts` | Tags da operação |
| `@Accept` | `@Accept json` | Tipo de conteúdo da requisição |
| `@Produce` | `@Produce json,xml` | Tipos de conteúdo da resposta |
| `@Param` | Veja abaixo | Definição de parâmetro |
| `@Success` | `@Success 200 {object} User` | Resposta de sucesso |
| `@Failure` | `@Failure 400 {object} Error` | Resposta de erro |
| `@Header` | `@Header 200 {string} Token` | Header de resposta |
| `@Router` | `@Router /users/{id} [get]` | Caminho e método da rota |
| `@Security` | `@Security ApiKeyAuth` | Requisito de segurança |
| `@Deprecated` | `@Deprecated` | Marcar como deprecated |
| `@x-visibility` | `@x-visibility public` | Separar docs públicas/privadas |
| `@x-<nome>` | `@x-code-samples file.json` | Extensão customizada |

**Sintaxe de Parâmetro:**

```
@Param <nome> <em> <tipo> <obrigatório> <descrição> [atributos]
```

- **nome**: Nome do parâmetro
- **em**: `query`, `path`, `header`, `body`, `formData`
- **tipo**: Tipo de dado (string, int, bool, object, array, file)
- **obrigatório**: `true` ou `false`
- **descrição**: Descrição (entre aspas se contém espaços)
- **atributos**: Atributos de validação opcionais

**Exemplos:**

```go
// Parâmetro de caminho
// @Param id path int true "ID do Usuário" minimum(1) maximum(1000)

// Parâmetro de query com validação
// @Param name query string false "Nome do usuário" minLength(3) maxLength(50)

// Parâmetro de query com enum
// @Param status query string false "Filtro de status" Enums(active,inactive,pending)

// Array de query com formato de coleção
// @Param tags query []string false "Tags" collectionFormat(multi)

// Parâmetro de header
// @Param X-Request-ID header string true "ID da Requisição" format(uuid)

// Parâmetro de body
// @Param user body User true "Objeto do usuário"

// Form data com arquivo
// @Param avatar formData file true "Imagem do avatar"
```

**Sintaxe de Resposta:**

```go
// Resposta simples
// @Success 200 {object} User

// Resposta com descrição
// @Success 201 {object} User "Usuário criado com sucesso"

// Resposta de array
// @Success 200 {array} User "Lista de usuários"

// Resposta primitiva
// @Success 200 {string} string "Mensagem de sucesso"

// Resposta genérica
// @Success 200 {object} Response{data=User} "Resposta do usuário"

// Múltiplos campos de dados
// @Success 200 {object} Response{data=User,meta=Metadata}
```

**Sintaxe de Header:**

```go
// Código de status único
// @Header 200 {string} X-Request-ID "Identificador da requisição"

// Múltiplos códigos de status
// @Header 200,201 {string} Location "URL do recurso"

// Todas as respostas
// @Header all {string} X-API-Version "Versão da API"
```

### Tags de Struct

#### Tags Padrão

```go
type User struct {
    // Serialização JSON
    ID   int    `json:"id"`
    Name string `json:"name,omitempty"`  // omitempty = não obrigatório
    
    // Validação (Gin binding)
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=150"`
    
    // Validação (go-playground/validator)
    UUID  string `json:"uuid" validate:"required,uuid"`
    
    // Atributos OpenAPI
    Price  float64  `json:"price" minimum:"0" maximum:"9999.99"`
    Status string   `json:"status" enum:"active,inactive" default:"active"`
    SKU    string   `json:"sku" pattern:"^[A-Z]{3}-[0-9]{6}$"`
    Items  []string `json:"items" minLength:"1" maxLength:"100"`
    
    // Valor de exemplo
    Bio string `json:"bio" example:"Desenvolvedor de software"`
    
    // Formato
    CreatedAt string `json:"created_at" format:"date-time"`
}
```

#### swaggertype - Override de Tipo

Converter tipos customizados para tipos OpenAPI:

```go
type Account struct {
    // Override sql.NullInt64 para integer
    ID sql.NullInt64 `json:"id" swaggertype:"integer"`
    
    // Tipo de tempo customizado para unix timestamp (integer)
    CreatedAt TimestampTime `json:"created_at" swaggertype:"primitive,integer"`
    
    // Array de bytes para string base64
    Certificate []byte `json:"cert" swaggertype:"string" format:"base64"`
    
    // Array de número customizado
    Coeffs []big.Float `json:"coeffs" swaggertype:"array,number"`
    
    // Tipos customizados aninhados
    Metadata map[string]interface{} `json:"metadata" swaggertype:"object"`
}
```

**Formato:** `swaggertype:"[primitive,]<tipo>"`

- Para tipos primitivos: `swaggertype:"string"`, `swaggertype:"integer"`, `swaggertype:"number"`, `swaggertype:"boolean"`
- Para arrays: `swaggertype:"array,<tipo-elemento>"`
- Para objetos: `swaggertype:"object"`

#### swaggerignore - Ocultar Campos

Excluir campos da documentação (ainda presente no JSON):

```go
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    
    // Presente no JSON, oculto nos docs
    Password string `json:"password" swaggerignore:"true"`
    
    // Campo interno, não no JSON ou docs
    internal string `swaggerignore:"true"`
    
    // Dado sensível
    SSN string `json:"ssn" swaggerignore:"true"`
}
```

#### extensions - Extensões Customizadas

Adicionar metadados customizados com prefixo `x-*`:

```go
type Product struct {
    // Indicador de chave primária
    ID int `json:"id" extensions:"x-primary-key=true"`
    
    // Formatação de moeda
    Price float64 `json:"price" extensions:"x-currency=BRL,x-format=currency"`
    
    // Múltiplas extensões
    Name string `json:"name" extensions:"x-order=1,x-searchable=true,x-filterable=true"`
    
    // Extensão booleana
    Featured bool `json:"featured" extensions:"x-promoted=true"`
    
    // Extensão nullable
    Discount float64 `json:"discount" extensions:"x-nullable"`
}
```

OpenAPI Gerado:

```json
{
  "properties": {
    "id": {
      "type": "integer",
      "x-primary-key": true
    },
    "price": {
      "type": "number",
      "x-currency": "BRL",
      "x-format": "currency"
    }
  }
}
```

## Recursos OpenAPI 3.2.0

nexs-swag oferece suporte completo aos recursos do OpenAPI 3.2.0 e ao patch 3.2.1 (`--openapi-version 3.2.0` ou `3.2.1`), mantendo total compatibilidade com versões anteriores (OpenAPI 2.0, 3.0.x, 3.1.x). A 3.2.1 não adiciona campos, apenas esclarece o texto da 3.2.0; nas duas o nexs-swag emite o `jsonSchemaDialect` padrão do 3.2 (`https://spec.openapis.org/oas/3.2/dialect/2025-09-17`).

### Método HTTP QUERY

O OpenAPI 3.2.0 introduz o método HTTP `QUERY` para consultas seguras com corpo de requisição:

```go
// @Summary      Buscar produtos complexa
// @Description  Buscar produtos usando parâmetros complexos no corpo da requisição
// @Tags         produtos
// @Accept       json
// @Produce      json
// @Param        filtros body ProductFilter true "Critérios de busca"
// @Success      200 {array} Product
// @Router       /products/query [query]
func QueryProducts(c *gin.Context) {}
```

### SecurityScheme Deprecated

Marque esquemas de segurança obsoletos com `@securityDefinitions.*.deprecated`:

```go
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @deprecated true
// @description ⚠️ Este método de autenticação será descontinuado. Use OAuth2 em vez disso.
```

Resultado no OpenAPI:
```yaml
securitySchemes:
  ApiKeyAuth:
    type: apiKey
    name: X-API-Key
    in: header
    deprecated: true
    description: ⚠️ Este método de autenticação será descontinuado. Use OAuth2 em vez disso.
```

### OAuth2 Metadata URL

Para descoberta automática de configuração OAuth2 via `@securityDefinitions.*.oauth2metadataurl`:

```go
// @securityDefinitions.oauth2.application OAuth2Application
// @tokenUrl https://auth.example.com/token
// @oauth2metadataurl https://auth.example.com/.well-known/oauth-authorization-server
```

Resultado no OpenAPI:
```yaml
securitySchemes:
  OAuth2Application:
    type: oauth2
    flows:
      clientCredentials:
        tokenUrl: https://auth.example.com/token
    oauth2MetadataUrl: https://auth.example.com/.well-known/oauth-authorization-server
```

### Device Authorization Flow

Suporte ao OAuth 2.0 Device Authorization Grant (RFC 8628) via `@securityDefinitions.*.deviceAuthorization`:

```go
// @securityDefinitions.oauth2.deviceAuth OAuth2Device
// @deviceAuthorization https://auth.example.com/device https://auth.example.com/token device-code
// @scopes.tv:watch Assistir canais de TV
// @scopes.tv:record Gravar conteúdo
```

Resultado no OpenAPI:
```yaml
securitySchemes:
  OAuth2Device:
    type: oauth2
    flows:
      urn:ietf:params:oauth:grant-type:device_code:
        deviceAuthorizationUrl: https://auth.example.com/device
        tokenUrl: https://auth.example.com/token
        scopes:
          tv:watch: Assistir canais de TV
          tv:record: Gravar conteúdo
```

### Respostas de Streaming

Para respostas SSE (Server-Sent Events) ou streaming, use `@Success {stream}`:

```go
// @Summary      Stream de eventos
// @Description  Recebe atualizações em tempo real de eventos do sistema
// @Tags         eventos
// @Produce      text/event-stream
// @Success      200 {stream} SystemEvent "Stream de eventos em tempo real"
// @Router       /events/stream [get]
func StreamEvents(c *gin.Context) {}
```

Resultado no OpenAPI:
```yaml
responses:
  '200':
    description: Stream de eventos em tempo real
    content:
      text/event-stream:
        itemSchema:
          $ref: '#/components/schemas/SystemEvent'
```

### Webhooks

Documentar webhooks que sua API envia para clientes via `@webhook`:

```go
// @webhook      OrderCreated
// @Description  Webhook enviado quando um novo pedido é criado
// @Tags         webhooks
// @Accept       json
// @Param        order body Order true "Dados do pedido criado"
// @Success      200 {object} WebhookResponse
func DocumentOrderWebhook() {}
```

### Callbacks

Para operações assíncronas com callbacks, use `@Callback`:

```go
// @Summary      Processar pagamento assíncrono
// @Description  Inicia processamento de pagamento e chama URL de callback
// @Tags         pagamentos
// @Accept       json
// @Param        payment body PaymentRequest true "Dados do pagamento"
// @Success      202 {object} PaymentResponse
// @Callback     paymentStatus {$request.body#/callbackUrl} post PaymentStatusCallback
// @Router       /payments/async [post]
func ProcessAsyncPayment(c *gin.Context) {}
```

### Separação por Visibilidade (@x-visibility)

Gere documentação separada para APIs públicas e privadas a partir de uma única base de código.

```go
// GetPublicUser retorna informações públicas do usuário
// @Summary      Buscar usuário (público)
// @Description  Retorna informações do usuário para consumo público
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "ID do Usuário"
// @Success      200  {object}  UserPublic
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [get]
// @x-visibility public
func GetPublicUser(c *gin.Context) {
    c.JSON(200, UserPublic{ID: 1, Name: "João"})
}

// GetAdminUser retorna detalhes completos incluindo dados sensíveis
// @Summary      Buscar usuário (admin)
// @Description  Retorna informações completas do usuário para uso administrativo
// @Tags         admin
// @Produce      json
// @Param        id   path      int  true  "ID do Usuário"
// @Success      200  {object}  UserPrivate
// @Failure      404  {object}  ErrorResponse
// @Router       /admin/users/{id} [get]
// @x-visibility private
func GetAdminUser(c *gin.Context) {
    c.JSON(200, UserPrivate{
        ID:       1,
        Name:     "João",
        Email:    "joao@example.com",
        Password: "hashed",
        Role:     "admin",
    })
}
```

**Opções de Visibilidade:**
- `@x-visibility public` - Endpoint aparece apenas em `openapi_public.json` ou `swagger_public.json`
- `@x-visibility private` - Endpoint aparece apenas em `openapi_private.json` ou `swagger_private.json`
- Sem anotação - Endpoint aparece em **ambas** as especificações (endpoint compartilhado)

**Arquivos Gerados:**
```
docs/
├── openapi_public.json    # Especificação API pública (OpenAPI 3.x)
├── openapi_private.json   # Especificação API privada/admin (OpenAPI 3.x)
├── swagger_public.json    # Especificação API pública (Swagger 2.0)
├── swagger_private.json   # Especificação API privada (Swagger 2.0)
├── openapi_public.yaml
├── openapi_private.yaml
├── docs_public.go
└── docs_private.go
```

**Filtragem de Schemas:**

Os schemas são automaticamente filtrados baseado no uso:
- Spec pública inclui apenas schemas referenciados por endpoints públicos
- Spec privada inclui apenas schemas referenciados por endpoints privados
- Schemas compartilhados (como `ErrorResponse`) aparecem onde necessário
- Dependências recursivas de schemas são coletadas automaticamente

**Casos de Uso:**
- Separar documentação de API estável pública de APIs experimentais privadas
- Ocultar endpoints administrativos internos da documentação pública
- Gerar diferentes SDKs de cliente para APIs públicas vs privadas
- Hospedar documentação separada para diferentes públicos
- Distinguir entre APIs externas e APIs inter-serviços

**Compatibilidade:**
- ✅ Swagger 2.0
- ✅ OpenAPI 3.0.x
- ✅ OpenAPI 3.1.x
- ✅ OpenAPI 3.2.x

Para um exemplo completo, veja [examples/26-x-visibility/](examples/26-x-visibility/) (OpenAPI 3.x) e [examples/27-x-visibility-v2/](examples/27-x-visibility-v2/) (Swagger 2.0).

### Migração 3.1.x → 3.2.0

nexs-swag detecta automaticamente a versão OpenAPI. Para ativar recursos 3.2.0:

1. **Não requer alterações** - recursos são ativados ao usar as anotações
2. **Compatível** - anotações antigas continuam funcionando
3. **Progressivo** - adicione recursos 3.2.0 gradualmente

**Avisos de depreciação** aparecem automaticamente se você usar:
- `@securityDefinitions.*.deprecated true` - mostra badge de descontinuação
- Esquemas obsoletos sem migração - sugestão para atualizar

## Exemplos

nexs-swag inclui 27 exemplos abrangentes e executáveis. Cada exemplo demonstra recursos específicos e inclui um README e script de execução.

### Exemplos Básicos

| Exemplo | Descrição | Recursos Principais |
|---------|-----------|---------------------|
| [01-basic](examples/01-basic) | Uso básico | Configuração mínima, API simples |
| [02-formats](examples/02-formats) | Formatos de saída | Saída JSON, YAML, Go |
| [03-general-info](examples/03-general-info) | Informações gerais da API | Metadados completos da API |

### Recursos Avançados

| Exemplo | Descrição | Recursos Principais |
|---------|-----------|---------------------|
| [04-property-strategy](examples/04-property-strategy) | Estratégias de nomenclatura | Snake_case, camelCase, PascalCase |
| [05-required-default](examples/05-required-default) | Obrigatório por padrão | Auto-require todos os campos |
| [06-exclude](examples/06-exclude) | Excluir diretórios | Filtrar caminhos indesejados |
| [07-tags-filter](examples/07-tags-filter) | Filtragem por tag | Gerar subconjunto de APIs |
| [08-parse-internal](examples/08-parse-internal) | Pacotes internos | Analisar diretório internal/ |
| [09-parse-dependency](examples/09-parse-dependency) | Dependências | Analisar pacotes vendor/go.mod |
| [10-dependency-level](examples/10-dependency-level) | Profundidade de dependência | Controlar nível de análise (0-3) |
| [11-parse-golist](examples/11-parse-golist) | Análise de go list | Usar `go list` para descoberta |

### Recursos de Documentação

| Exemplo | Descrição | Recursos Principais |
|---------|-----------|---------------------|
| [12-markdown-files](examples/12-markdown-files) | Descrições em Markdown | Carregar docs de arquivos .md |
| [13-code-examples](examples/13-code-examples) | Amostras de código | Exemplos em múltiplas linguagens |
| [14-overrides-file](examples/14-overrides-file) | Overrides de tipo | Configuração de arquivo .swaggo |
| [15-generated-time](examples/15-generated-time) | Timestamp de geração | Adicionar data de geração |
| [16-instance-name](examples/16-instance-name) | Múltiplas instâncias | Conjuntos de documentação nomeados |
| [17-template-delims](examples/17-template-delims) | Delimitadores customizados | Evitar conflitos de template |

### Validação e Estrutura

| Exemplo | Descrição | Recursos Principais |
|---------|-----------|---------------------|
| [18-collection-format](examples/18-collection-format) | Formatos de array | CSV, multi, pipes, SSV, TSV |
| [19-parse-func-body](examples/19-parse-func-body) | Corpos de função | Analisar anotações inline |
| [20-fmt-command](examples/20-fmt-command) | Comando de formatação | Auto-formatar comentários |
| [21-struct-tags](examples/21-struct-tags) | Todas as tags de struct | Referência completa de tags |
| [22-openapi-v2](examples/22-openapi-v2) | Versionamento OpenAPI | Swagger 2.0 & OpenAPI 3.1.0 |
| [23-recursive-parsing](examples/23-recursive-parsing) | Análise recursiva | parseInternal, exclude, parseDependency |
| [24-include-types-simple](examples/24-include-types-simple) | Include types (simples) | Uso básico de `--includeTypes` |
| [25-include-types-complex](examples/25-include-types-complex) | Include types (complexo) | `--includeTypes` entre pacotes |
| [26-x-visibility](examples/26-x-visibility) | x-visibility (OpenAPI 3.x) | Separação spec pública/privada |
| [27-x-visibility-v2](examples/27-x-visibility-v2) | x-visibility (Swagger 2.0) | Separação spec pública/privada |

### Executando Exemplos

Cada exemplo inclui um script `run.sh`:

```bash
cd examples/01-basic
./run.sh
```

Ou manualmente (OpenAPI 3.1.0):

```bash
cd examples/01-basic
nexs-swag init -d . -o ./docs
cat docs/openapi.json
```

Ou gerar Swagger 2.0:

```bash
cd examples/01-basic
nexs-swag init -d . -o ./docs --openapi-version 2.0
cat docs/swagger.json
```

### Exemplo: API CRUD Completa

Veja [examples/03-general-info](examples/03-general-info) para uma API CRUD completa com:
- Múltiplos endpoints (GET, POST, PUT, DELETE)
- Modelos de request/response
- Regras de validação
- Respostas de erro
- Esquemas de segurança
- Headers de resposta

## Qualidade e Testes

### Cobertura de Testes

```bash
$ go test ./pkg/... -cover
```

| Pacote | Cobertura | Testes |
|---------|----------|--------|
| pkg/converter | 85.1% | 16 testes (com roundtrip) |
| pkg/format | 95.1% | 15 testes |
| pkg/generator/v2 | 80.3% | 12 testes |
| pkg/generator/v3 | 83.3% | 8 testes |
| pkg/openapi/v2 | 92.0% | 12 testes |
| pkg/openapi/v3 | 88.9% | 10 testes |
| pkg/parser | 84.6% | 195 testes |
| **Geral** | **80.1%** | **320+ testes** |

### Métricas de Qualidade

- ✅ **0 avisos de linter** (golangci-lint com 20+ linters)
- ✅ **0 condições de corrida** (testado com flag `-race`)
- ✅ **27 testes de integração** (exemplos executáveis)
- ✅ **~8.500 linhas de código de teste**
- ✅ **Pronto para produção** (mantido ativamente)
- ✅ **100% compatível com swaggo/swag**
- ✅ **Suporte a múltiplas versões** (OpenAPI 3.1.0 + Swagger 2.0)

### Executando Testes

```bash
# Testes unitários
go test ./pkg/... -v

# Com cobertura
go test ./pkg/... -cover

# Com detecção de race condition
go test ./pkg/... -race

# Pacote específico
go test ./pkg/parser -v

# Executar exemplos
cd examples && for d in */; do cd "$d" && ./run.sh && cd ..; done
```

## Compatibilidade com swaggo/swag

nexs-swag é projetado como um **substituto direto** para swaggo/swag com recursos aprimorados.

### Migração do swaggo/swag

**Nenhuma mudança necessária!** Simplesmente substitua o binário:

```bash
# Ao invés de
go install github.com/swaggo/swag/cmd/swag@latest

# Use
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest

# Os mesmos comandos funcionam
nexs-swag init
nexs-swag fmt
```

### Tabela de Compatibilidade

| Recurso | swaggo/swag | nexs-swag | Notas |
|---------|-------------|-----------|-------|
| Versão OpenAPI | 2.0 | 3.1.0 | Retrocompatível |
| Todas as anotações | ✅ | ✅ | 100% compatível |
| Tags de struct | ✅ | ✅ | swaggertype, swaggerignore, extensions |
| Flags CLI | ✅ | ✅ | Todas as 28 flags suportadas |
| Arquivo .swaggo | ✅ | ✅ | Overrides de tipo |
| Markdown | ✅ | ✅ | Descrições baseadas em arquivo |
| Exemplos de código | ✅ | ✅ | Amostras em múltiplas linguagens |
| Webhooks | ❌ | ✅ | Recurso OpenAPI 3.1 |
| JSON Schema 2020-12 | ❌ | ✅ | Schema moderno |
| Headers de resposta | Limitado | ✅ | Suporte completo |
| Cobertura de testes | ~70% | 80.1% | Maior qualidade |
| Versão Go | 1.19+ | 1.27.1+ | Recursos Go modernos |

### O que é Diferente?

**Aprimorado (retrocompatível):**
- Saída OpenAPI 3.1.0 (vs 2.0)
- Melhor tratamento de nullable
- Mais atributos de validação
- Mensagens de erro melhoradas
- Melhor cobertura de testes

**Mesma API:**
- Todas as flags de linha de comando
- Todas as anotações
- Todas as tags de struct
- Estrutura gerada de docs.go
- Integração com Swagger UI

## Sobre o Projeto

### Estatísticas do Projeto

- **Linhas de Código:** ~5.200 (pkg/ excluindo testes)
- **Código de Teste:** ~8.500 linhas
- **Arquivos Go:** 42 arquivos de implementação
- **Arquivos de Teste:** 29 arquivos de teste
- **Pacotes:** 9 (converter, format, generator, generator/v2, generator/v3, openapi, openapi/v2, openapi/v3, parser)
- **Exemplos:** 27 exemplos completos
- **Cobertura de Testes:** 80.1%
- **Versões OpenAPI:** 4 (Swagger 2.0, OpenAPI 3.0.x, 3.1.x, 3.2.x)
- **Dependências:** 3 dependências diretas
  - urfave/cli/v2 (framework CLI)
  - golang.org/x/tools (análise AST Go)
  - gopkg.in/yaml.v3 (suporte YAML)

### Estrutura do Projeto

```
nexs-swag/
├── cmd/
│   └── nexs-swag/          # Ponto de entrada CLI
├── pkg/
│   ├── converter/          # Conversão de versão (v3 ↔ v2)
│   ├── format/             # Formatação de código
│   ├── generator/          # Geração OpenAPI
│   │   ├── v2/             # Gerador Swagger 2.0
│   │   └── v3/             # Gerador OpenAPI 3.x
│   ├── openapi/            # Modelos OpenAPI
│   │   ├── v2/             # Modelos Swagger 2.0
│   │   └── v3/             # Modelos OpenAPI 3.x
│   └── parser/             # Análise de código Go (AST)
├── examples/               # 27 exemplos
│   ├── 01-basic/
│   ├── 02-formats/
│   └── ...
├── docs/                   # Documentação do projeto
├── README.md               # Versão em inglês
├── README_pt.md            # Este arquivo
├── README_es.md            # Versão em espanhol
└── LICENSE                 # Licença Apache 2.0
```

### Inspiração e Créditos

Este projeto foi inspirado pelo [swaggo/swag](https://github.com/swaggo/swag) e construído para estender suas capacidades com suporte completo ao OpenAPI 3.1.0, mantendo 100% de compatibilidade retroativa.

**Créditos:**
- [swaggo/swag](https://github.com/swaggo/swag) - Gerador Swagger 2.0 original
- [OpenAPI Initiative](https://www.openapis.org/) - Especificação OpenAPI
- [Go Team](https://go.dev/) - Linguagem e ferramentas incríveis
- Todos os contribuidores e a comunidade Go

## Contribuindo

Contribuições são bem-vindas! Por favor, siga estas diretrizes:

### Como Contribuir

1. **Fork** o repositório
2. **Crie** uma branch de feature (`git checkout -b feature/recurso-incrivel`)
3. **Faça** suas mudanças
4. **Adicione** testes para nova funcionalidade
5. **Execute** os testes (`go test ./...`)
6. **Execute** o linter (`golangci-lint run`)
7. **Commit** suas mudanças (`git commit -m 'Adiciona recurso incrível'`)
8. **Push** para a branch (`git push origin feature/recurso-incrivel`)
9. **Abra** um Pull Request

### Configuração de Desenvolvimento

```bash
# Clonar repositório
git clone https://github.com/fsvxavier/nexs-swag.git
cd nexs-swag

# Instalar dependências
go mod download

# Executar testes
go test ./... -v

# Executar linter
golangci-lint run

# Build
go build -o nexs-swag ./cmd/nexs-swag
```

### Reportando Issues

Por favor inclua:
- Versão do Go (`go version`)
- Versão do nexs-swag (`nexs-swag --version`)
- Exemplo reproduzível mínimo
- Comportamento esperado vs real

### Solicitações de Recursos

Abra uma issue com:
- Descrição clara do recurso
- Caso de uso e benefícios
- Implementação proposta (se houver)

## Licença

Este projeto está licenciado sob a **Licença Apache 2.0** - veja o arquivo [LICENSE](LICENSE) para detalhes.

Copyright 2025 Fabricio S. V. Xavier

O nexs-swag inclui partes derivadas do [swaggo/swag](https://github.com/swaggo/swag), Copyright (c) 2017 Eason Lin, utilizadas sob a Licença MIT. Veja o arquivo [NOTICE](NOTICE) para o aviso completo de terceiros.

## Suporte e Comunidade

- **Issues:** [GitHub Issues](https://github.com/fsvxavier/nexs-swag/issues)
- **Discussões:** [GitHub Discussions](https://github.com/fsvxavier/nexs-swag/discussions)
- **Documentação:** [Wiki](https://github.com/fsvxavier/nexs-swag/wiki)
- **Exemplos:** [examples/](examples/)

---

**Feito com ❤️ para a comunidade Go**

[⬆ Voltar ao topo](#nexs-swag)
