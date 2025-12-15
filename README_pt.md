# nexs-swag

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

[![Versão Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1.0-6BA539?style=flat&logo=openapiinitiative)](https://spec.openapis.org/oas/v3.1.0)
[![Swagger](https://img.shields.io/badge/Swagger-2.0-85EA2D?style=flat&logo=swagger)](https://swagger.io/specification/v2/)
[![Licença](https://img.shields.io/badge/Licença-MIT-blue.svg)](LICENSE)
[![Cobertura](https://img.shields.io/badge/Cobertura-86.1%25-brightgreen.svg)](/))
[![Exemplos](https://img.shields.io/badge/Exemplos-22-blue.svg)](examples/)

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
- ✅ **Suporte dual de versões** - Gere OpenAPI 3.1.0 **ou** Swagger 2.0 a partir das mesmas anotações
- ✅ **OpenAPI 3.1.0** - Suporte completo para JSON Schema 2020-12, webhooks e recursos modernos
- ✅ **Swagger 2.0** - Compatibilidade total com sistemas legados
- ✅ **Conversão automática** - Conversão interna entre formatos com avisos para incompatibilidades
- ✅ **20+ atributos de validação** - minimum, maximum, pattern, enum, format e mais
- ✅ **Validação de frameworks** - Suporte nativo para Gin (binding) e go-playground/validator
- ✅ **Headers de resposta** - Documentação completa de headers
- ✅ **Múltiplos tipos de conteúdo** - JSON, XML, YAML, CSV, PDF e tipos MIME customizados
- ✅ **Extensões customizadas** - Suporte completo para x-*
- ✅ **86.1% de cobertura de testes** - Pronto para produção com suite de testes abrangente
- ✅ **22 exemplos funcionais** - Aprenda com exemplos completos e executáveis

### Por que nexs-swag?

| Recurso | swaggo/swag | nexs-swag |
|---------|-------------|-----------||
| OpenAPI 3.1.0 | ❌ | ✅ |
| Swagger 2.0 | ✅ | ✅ |
| Geração Dual | ❌ | ✅ (ambos do mesmo código) |
| JSON Schema | Draft 4 | Draft 4 + 2020-12 |
| Webhooks | ❌ | ✅ (OpenAPI 3.1.0) |
| Headers de Resposta | Limitado | Suporte Completo |
| Suporte a Nullable | `x-nullable` | Nativo + `x-nullable` |
| Cobertura de Testes | ~70% | 86.1% |
| Exemplos | ~10 | 22 |
| Versão Go | 1.19+ | 1.23+ |

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

Requer [Go 1.23 ou superior](https://go.dev/dl/).

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
```

**Swagger 2.0:**

```bash
nexs-swag init --openapi-version 2.0
```

**Gerar ambas as versões:**

```bash
nexs-swag init -o ./docs/v3 --openapi-version 3.1
nexs-swag init -o ./docs/v2 --openapi-version 2.0
```

Ou especifique os diretórios:

```bash
nexs-swag init -d ./cmd/api -o ./docs
```

#### 3. Arquivos Gerados

**OpenAPI 3.1.0:**
Os seguintes arquivos serão criados no seu diretório de saída (padrão: `./docs`):

- **`docs/openapi.json`** - Especificação OpenAPI 3.1.0 em formato JSON
- **`docs/openapi.yaml`** - Especificação OpenAPI 3.1.0 em formato YAML
- **`docs/docs.go`** - Arquivo de documentação Go embarcado

**Swagger 2.0:**
Quando usar `--openapi-version 2.0`, os arquivos gerados serão:

- **`docs/swagger.json`** - Especificação Swagger 2.0 em formato JSON
- **`docs/swagger.yaml`** - Especificação Swagger 2.0 em formato YAML
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

Veja a [versão em inglês](README.md#how-to-use-with-gin) para detalhes completos.

## Referência CLI

### Comando init

Gera documentação OpenAPI a partir do código fonte.

```bash
nexs-swag init [opções]
```

**Opções Principais:**

- `--dir, -d` - Diretórios para analisar (padrão: `./`)
- `--output, -o` - Diretório de saída (padrão: `./docs`)
- `--outputTypes, --ot` - Tipos de arquivo de saída (padrão: `go,json,yaml`)
- `--openapi-version, --ov` - Versão OpenAPI: `2.0`, `3.0`, `3.1` (padrão: `3.1`)
- `--parseDependency, --pd` - Analisar dependências (padrão: `false`)
- `--parseInternal` - Analisar pacotes internos (padrão: `false`)
- `--propertyStrategy, -p` - Estratégia de nomenclatura: `snakecase`, `camelcase`, `pascalcase`
- `--validate` - Validar especificação gerada (padrão: `true`)

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

# Apenas saída JSON
nexs-swag init --outputTypes json

# Nomes de propriedade em snake_case
nexs-swag init --propertyStrategy snakecase
```

### Comando fmt

Formata comentários swagger automaticamente.

```bash
nexs-swag fmt [opções]
```

**Exemplo:**

```bash
# Formatar diretório atual
nexs-swag fmt

# Formatar diretório específico
nexs-swag fmt -d ./internal/api
```

## Formato de Comentários Declarativos

Para documentação completa de todas as anotações, parâmetros, tags de struct e exemplos, consulte a [versão em inglês](README.md#declarative-comments-format).

### Resumo Rápido

**Informações Gerais da API:**
- `@title` - Título da API (obrigatório)
- `@version` - Versão da API (obrigatório)
- `@description` - Descrição da API
- `@host` - Host da API
- `@BasePath` - Caminho base
- `@securityDefinitions.*` - Definições de segurança

**Operação de API:**
- `@Summary` - Resumo curto
- `@Description` - Descrição detalhada
- `@Tags` - Tags da operação
- `@Param` - Definição de parâmetro
- `@Success` - Resposta de sucesso
- `@Failure` - Resposta de erro
- `@Router` - Caminho e método da rota

**Tags de Struct:**
- `json` - Serialização JSON
- `binding` - Validação Gin
- `validate` - Validação go-playground
- `swaggertype` - Override de tipo
- `swaggerignore` - Ocultar campo
- `extensions` - Extensões customizadas

## Exemplos

nexs-swag inclui 22 exemplos abrangentes e executáveis demonstrando todas as funcionalidades, incluindo geração de OpenAPI 3.1.0 e Swagger 2.0. Veja a [seção de exemplos](README.md#examples) na versão em inglês para a lista completa.

### Executando Exemplos

Cada exemplo inclui um script `run.sh`:

```bash
cd examples/01-basic
./run.sh
```

## Qualidade e Testes

### Cobertura de Testes

| Pacote | Cobertura | Testes |
|---------|----------|--------|
| pkg/converter | 92.3% | 13 testes |
| pkg/format | 95.1% | 15 testes |
| pkg/generator | 71.6% | 16 testes |
| pkg/generator/v2 | 88.4% | 12 testes |
| pkg/generator/v3 | 85.2% | 8 testes |
| pkg/openapi | 83.3% | 22 testes |
| pkg/openapi/v2 | 89.7% | 12 testes |
| pkg/openapi/v3 | 91.5% | 10 testes |
| pkg/parser | 82.1% | 192 testes |
| **Geral** | **87.9%** | **300+ testes** |

### Métricas de Qualidade

- ✅ **0 avisos de linter** (golangci-lint com 20+ linters)
- ✅ **0 condições de corrida** (testado com flag `-race`)
- ✅ **22 testes de integração** (exemplos executáveis)
- ✅ **~8.500 linhas de código de teste**
- ✅ **Pronto para produção** (mantido ativamente)
- ✅ **100% compatível com swaggo/swag**
- ✅ **Suporte dual de versões** (OpenAPI 3.1.0 + Swagger 2.0)

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

## Sobre o Projeto

### Estatísticas do Projeto

- **Linhas de Código:** ~5.200 (pkg/ excluindo testes)
- **Código de Teste:** ~8.500 linhas
- **Arquivos Go:** 42 arquivos de implementação
- **Arquivos de Teste:** 29 arquivos de teste
- **Pacotes:** 9 (converter, format, generator, generator/v2, generator/v3, openapi, openapi/v2, openapi/v3, parser)
- **Exemplos:** 22 exemplos completos
- **Cobertura de Testes:** 87.9%
- **Versões OpenAPI:** 2 (Swagger 2.0 + OpenAPI 3.1.0)
- **Dependências:** 3 dependências diretas

### Inspiração e Créditos

Este projeto foi inspirado pelo [swaggo/swag](https://github.com/swaggo/swag) e construído para estender suas capacidades com suporte completo ao OpenAPI 3.1.0, mantendo 100% de compatibilidade retroativa.

## Contribuindo

Contribuições são bem-vindas! Veja a [versão em inglês](README.md#contributing) para diretrizes detalhadas.

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

## Licença

Este projeto está licenciado sob a **Licença MIT** - veja o arquivo [LICENSE](LICENSE) para detalhes.

## Suporte e Comunidade

- **Issues:** [GitHub Issues](https://github.com/fsvxavier/nexs-swag/issues)
- **Discussões:** [GitHub Discussions](https://github.com/fsvxavier/nexs-swag/discussions)
- **Documentação:** [Wiki](https://github.com/fsvxavier/nexs-swag/wiki)
- **Exemplos:** [examples/](examples/)

---

**Feito com ❤️ para a comunidade Go**

[⬆ Voltar ao topo](#nexs-swag)
