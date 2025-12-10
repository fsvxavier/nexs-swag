# nexs-swag

**Gerador de documentação OpenAPI 3.1.x para Go** - Compatível com swaggo/swag + Recursos Avançados

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1.0-6BA539?style=flat&logo=openapiinitiative)](https://spec.openapis.org/oas/v3.1.0)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Examples](https://img.shields.io/badge/Examples-21-blue.svg)](examples/)

## 🚀 Visão Geral

O **nexs-swag** é um gerador de documentação OpenAPI 3.1.x completo para aplicações Go, criado como evolução do swaggo/swag com suporte total à especificação mais recente do OpenAPI.

### ✨ Diferenciais

- ✅ **100% compatível com swaggo/swag** - Suporta todas as annotations e tags
- ✅ **OpenAPI 3.1.0** - JSON Schema 2020-12, webhooks, e recursos modernos
- ✅ **Tags swaggo/swag** - `swaggertype`, `swaggerignore`, `extensions`
- ✅ **20+ atributos de validação** - minimum, maximum, pattern, enum, etc
- ✅ **Validação de frameworks** - Gin (binding), go-playground/validator
- ✅ **Response headers** - Documentação completa de headers
- ✅ **Múltiplos content-types** - JSON, XML, CSV, PDF, etc
- ✅ **Extensões customizadas** - Suporte completo a x-*
- ✅ **86.1% de cobertura de testes** - Testado em produção

## 📦 Instalação

```bash
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest
```

Ou clone e compile:

```bash
git clone https://github.com/fsvxavier/nexs-swag.git
cd nexs-swag
go build -o nexs-swag ./cmd/nexs-swag
```

## 🎯 Uso Rápido

### 1. Adicione annotations ao seu código

```go
package main

import "database/sql"

// @title API de Exemplo
// @version 1.0
// @description API demonstrando todas as funcionalidades do nexs-swag
// @host localhost:8080
// @BasePath /api/v1
func main() {
    // Sua aplicação
}

// User representa um usuário do sistema
type User struct {
    // ID do usuário (sql.NullInt64 → integer)
    ID sql.NullInt64 `json:"id" swaggertype:"integer" extensions:"x-primary-key=true"`
    
    // Nome completo
    Name string `json:"name" binding:"required" minLength:"3" maxLength:"100" example:"John Doe"`
    
    // Email (validado)
    Email string `json:"email" binding:"required,email" format:"email" extensions:"x-unique=true"`
    
    // Senha (oculta da documentação)
    Password string `json:"password" swaggerignore:"true"`
    
    // Status da conta
    Status string `json:"status" enum:"active,inactive,pending" default:"active"`
    
    // Saldo da conta
    Balance float64 `json:"balance" minimum:"0" extensions:"x-currency=USD,x-format=currency"`
}

// createUser cria um novo usuário
// @Summary Criar usuário
// @Description Cria um novo usuário no sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "Dados do usuário"
// @Success 201 {object} User
// @Header 201 {string} X-Request-ID "ID da requisição"
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func createUser() {}
```

### 2. Gere a documentação

```bash
nexs-swag init -d ./cmd/api -o ./docs
```

### 3. Arquivos gerados

- `docs/openapi.json` - OpenAPI 3.1.0 em JSON
- `docs/openapi.yaml` - OpenAPI 3.1.0 em YAML
- `docs/docs.go` - Documentação embarcada em Go

## 🏷️ Tags de Struct - Compatibilidade swaggo/swag

### swaggertype - Override de Tipos

Converta tipos customizados para tipos OpenAPI:

```go
type Account struct {
    // sql.NullInt64 → integer
    ID sql.NullInt64 `json:"id" swaggertype:"integer"`
    
    // TimestampTime → integer (unix timestamp)
    CreatedAt TimestampTime `json:"created_at" swaggertype:"primitive,integer"`
    
    // []byte → string com base64
    Certificate []byte `json:"cert" swaggertype:"string" format:"base64"`
    
    // []big.Float → array de numbers
    Coeffs []big.Float `json:"coeffs" swaggertype:"array,number"`
}
```

### swaggerignore - Ocultar Campos

Oculte campos da documentação sem afetar JSON:

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    
    // Aparece no JSON, oculto na doc
    Password string `json:"password" swaggerignore:"true"`
    
    // Campo interno
    Internal string `swaggerignore:"true"`
}
```

### extensions - Extensões Customizadas

Adicione metadados customizados (x-*):

```go
type Product struct {
    ID    int     `json:"id" extensions:"x-primary-key=true"`
    Price float64 `json:"price" extensions:"x-currency=USD,x-format=currency"`
    Name  string  `json:"name" extensions:"x-order=1,x-searchable=true"`
}
```

## 📋 Validações Suportadas

### Tags de Validação (Gin - binding)

```go
type CreateRequest struct {
    Name  string  `binding:"required,min=3,max=100"`
    Email string  `binding:"required,email"`
    Age   int     `binding:"gte=0,lte=150"`
    Price float64 `binding:"gt=0"`
}
```

### Tags de Validação (go-playground/validator)

```go
type User struct {
    UUID  string `validate:"required,uuid"`
    Email string `validate:"required,email"`
    Date  string `validate:"datetime=2006-01-02"`
}
```

### Tags Customizadas OpenAPI

```go
type Product struct {
    SKU   string  `json:"sku" pattern:"^[A-Z]{3}-[0-9]{6}$" example:"ABC-123456"`
    Price float64 `json:"price" minimum:"0" maximum:"9999.99" example:"99.99"`
    Tags  []string `json:"tags" minLength:"1" maxLength:"50"`
}
```

## 🎨 Annotations de Operação

### Parâmetros Avançados

```go
// @Param id path int true "User ID" minimum(1) maximum(1000) example(123)
// @Param name query string false "Name" minLength(3) pattern(^[a-z]+$)
// @Param status query string false "Status" enum(active,inactive) default(active)
// @Param tags query []string false "Tags" collectionFormat(multi)
```

### Response Headers

```go
// @Header 200 {string} X-Request-ID "Request identifier"
// @Header 200 {int} X-Rate-Limit "Rate limit"
// @Header 201 {string} Location "Resource location"
```

### Múltiplos Content-Types

```go
// @Accept json,xml,yaml
// @Produce json,xml,csv,pdf
// @Success 200 {object} User
```

## 📚 Exemplos Completos (21 exemplos)

Todos os exemplos estão em [`examples/`](examples/) e incluem:

**Básicos:**
- `01-basic` - Uso básico do nexs-swag
- `02-formats` - Múltiplos formatos (JSON, YAML, Go)
- `03-general-info` - Arquivo de informações gerais

**Avançados:**
- `04-property-strategy` - Estratégias de naming
- `05-required-default` - Campos required por padrão
- `06-exclude` - Exclusão de diretórios
- `07-tags-filter` - Filtro por tags
- `08-parse-internal` - Parse de packages internos
- `09-parse-dependency` - Parse de dependências
- `10-dependency-level` - Níveis de parse (0-3)
- `11-parse-golist` - Parse via go list
- `12-markdown-files` - Markdown como descrições
- `13-code-examples` - Code samples em 23+ linguagens
- `14-overrides-file` - Arquivo .swaggo de overrides
- `15-generated-time` - Timestamp na documentação
- `16-instance-name` - Nome customizado da instância
- `17-template-delims` - Delimitadores customizados
- `18-collection-format` - Formatos de array
- `19-parse-func-body` - Parse de anotações em funções
- `20-fmt-command` - Formatação de comentários swagger
- `21-struct-tags` - Demonstração de 18 struct tags

**Executar exemplos:**
```bash
cd examples/01-basic
./run.sh
```

## 🔧 Comandos CLI

### init - Gerar Documentação

```bash
nexs-swag init [opções]
```

**Opções principais:**
- `-d, --dir` - Diretório de código Go (default: "./")
- `-o, --output` - Diretório de saída (default: "./docs")
- `-f, --format` - Formatos: json, yaml, go (default: "json,yaml,go")
- `--validate` - Validar especificação (default: true)
- `--parseDependency` - Parse de dependências (default: false)
- `--parseInternal` - Parse de packages internos (default: false)
- `--parseGoList` - Parse via go list (default: true)

**Exemplos:**

```bash
# Gerar em ./docs
nexs-swag init -d ./cmd/api -o ./docs

# Apenas JSON
nexs-swag init -d . -o ./api-docs -f json

# Com dependências (nível 3 - completo)
nexs-swag init -d . --parseDependency --parseDependencyLevel 3

# Sem validação
nexs-swag init -d . --validate=false
```

### fmt - Formatar Comentários

```bash
nexs-swag fmt [opções]
```

Formata automaticamente comentários swagger usando AST do Go.

**Exemplo:**

```bash
# Formatar diretório atual
nexs-swag fmt -d ./cmd/api
```

## 📈 Qualidade e Testes

- **Cobertura de Testes:** 86.1% (META: 80% ✅)
  - pkg/format: 95.1%
  - pkg/generator: 84.6%
  - pkg/openapi: 83.3%
  - pkg/parser: 81.5%
- **Arquivos de Teste:** 21 arquivos, ~5.000 linhas
- **Testes Integração:** 21 exemplos funcionais
- **Race Conditions:** Zero (testado com -race)
- **CI/CD:** Pronto para integração contínua

## 📖 Documentação Completa

- [README.md](README.md) - Este arquivo (visão geral e início rápido)
- [INSTALL.md](INSTALL.md) - Guia completo de instalação
- [PENDENCIAS.md](PENDENCIAS.md) - Status do projeto e roadmap
- [examples/README.md](examples/README.md) - Guia de exemplos

## 🎯 Compatibilidade

### OpenAPI 3.1.0
- ✅ JSON Schema 2020-12
- ✅ Webhooks
- ✅ Composition (allOf, oneOf, anyOf)
- ✅ Nullable via type array
- ✅ Const e prefixItems

### swaggo/swag (100% compatível)
- ✅ Todas as annotations (@Summary, @Param, @Success, etc)
- ✅ Tags de struct (json, binding, validate)
- ✅ swaggertype, swaggerignore, extensions
- ✅ Atributos de parâmetros (minimum, enum, pattern, etc)
- ✅ Response headers
- ✅ 28/28 flags CLI implementadas
- ✅ Comandos init e fmt

## 📊 Estatísticas do Projeto

- **Linhas de código:** ~3.854 (pkg/, excluindo testes)
- **Arquivos Go:** 33 arquivos de implementação
- **Arquivos de teste:** 21 arquivos (~5.000 linhas)
- **Packages:** 4 (format, generator, openapi, parser)
- **Exemplos:** 21 exemplos funcionais
- **Cobertura de testes:** 86.1% (META: 80% ✅)
- **Status:** ✅ Pronto para produção

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor:

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🙏 Agradecimentos

- [swaggo/swag](https://github.com/swaggo/swag) - Inspiração e compatibilidade
- [OpenAPI Initiative](https://www.openapis.org/) - Especificação OpenAPI
- Comunidade Go

## 📞 Suporte

- Issues: [GitHub Issues](https://github.com/fsvxavier/nexs-swag/issues)
- Documentação: [Wiki](https://github.com/fsvxavier/nexs-swag/wiki)

---

**Desenvolvido com ❤️ para a comunidade Go**