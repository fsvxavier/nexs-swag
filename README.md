# nexs-swag

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1.0-6BA539?style=flat&logo=openapiinitiative)](https://spec.openapis.org/oas/v3.1.0)
[![Swagger](https://img.shields.io/badge/Swagger-2.0-85EA2D?style=flat&logo=swagger)](https://swagger.io/specification/v2/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Coverage](https://img.shields.io/badge/Coverage-86.1%25-brightgreen.svg)](/)
[![Examples](https://img.shields.io/badge/Examples-22-blue.svg)](examples/)

**Automatically generate OpenAPI 3.1.0 or Swagger 2.0 documentation from Go source code annotations.**

nexs-swag converts Go annotations to OpenAPI 3.1.0 or Swagger 2.0 Specification. It is designed as an evolution of [swaggo/swag](https://github.com/swaggo/swag) with full support for the latest OpenAPI specification and complete backward compatibility with Swagger 2.0.

## Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
- [Supported Web Frameworks](#supported-web-frameworks)
- [How to use with Gin](#how-to-use-with-gin)
- [CLI Reference](#cli-reference)
  - [init Command](#init-command)
  - [fmt Command](#fmt-command)
- [Implementation Status](#implementation-status)
- [Declarative Comments Format](#declarative-comments-format)
  - [General API Info](#general-api-info)
  - [API Operation](#api-operation)
  - [Struct Tags](#struct-tags)
- [Examples](#examples)
- [Quality & Testing](#quality--testing)
- [swaggo/swag Compatibility](#swaggoswag-compatibility)
- [About the Project](#about-the-project)
- [Contributing](#contributing)
- [License](#license)

## Overview

### Key Features

- ✅ **100% swaggo/swag compatible** - Drop-in replacement with all annotations and tags
- ✅ **Dual version support** - Generate OpenAPI 3.1.0 **or** Swagger 2.0 from the same annotations
- ✅ **OpenAPI 3.1.0** - Full support for JSON Schema 2020-12, webhooks, and modern features
- ✅ **Swagger 2.0** - Complete backward compatibility with legacy systems
- ✅ **Automatic conversion** - Internal conversion between formats with warnings for incompatibilities
- ✅ **20+ validation attributes** - minimum, maximum, pattern, enum, format, and more
- ✅ **Framework validation** - Native support for Gin (binding) and go-playground/validator
- ✅ **Response headers** - Complete header documentation
- ✅ **Multiple content types** - JSON, XML, YAML, CSV, PDF, and custom MIME types
- ✅ **Custom extensions** - Full x-* extension support
- ✅ **86.1% test coverage** - Production-ready with comprehensive test suite
- ✅ **22 working examples** - Learn from complete, runnable examples

### Why nexs-swag?

| Feature | swaggo/swag | nexs-swag |
|---------|-------------|-----------|
| OpenAPI 3.1.0 | ❌ | ✅ |
| Swagger 2.0 | ✅ | ✅ |
| Dual Generation | ❌ | ✅ (both from same code) |
| JSON Schema | Draft 4 | Draft 4 + 2020-12 |
| Webhooks | ❌ | ✅ (OpenAPI 3.1.0) |
| Response Headers | Limited | Full Support |
| Nullable Support | `x-nullable` | Native + `x-nullable` |
| Test Coverage | ~70% | 86.1% |
| Examples | ~10 | 22 |
| Go Version | 1.19+ | 1.23+ |

## Getting Started

### Installation

#### Using go install (Recommended)

```bash
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest
```

To verify installation:

```bash
nexs-swag --version
```

#### Building from Source

Requires [Go 1.23 or newer](https://go.dev/dl/).

```bash
git clone https://github.com/fsvxavier/nexs-swag.git
cd nexs-swag
go build -o nexs-swag ./cmd/nexs-swag
```

#### Using Docker

```bash
docker pull ghcr.io/fsvxavier/nexs-swag:latest
docker run --rm -v $(pwd):/app ghcr.io/fsvxavier/nexs-swag:latest init
```

### Quick Start

#### 1. Add API Annotations

Add general API annotations to your `main.go`:

```go
package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
)

// @title           User Management API
// @version         1.0.0
// @description     A user management API with complete OpenAPI 3.1.0 documentation
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
    r := gin.Default()
    // Your application setup
    r.Run(":8080")
}

// User represents a system user
type User struct {
    // User ID (sql.NullInt64 → integer in OpenAPI)
    ID sql.NullInt64 `json:"id" swaggertype:"integer" extensions:"x-primary-key=true"`
    
    // Full name (3-100 characters required)
    Name string `json:"name" binding:"required" minLength:"3" maxLength:"100" example:"John Doe"`
    
    // Email address (validated)
    Email string `json:"email" binding:"required,email" format:"email" extensions:"x-unique=true"`
    
    // Password (hidden from documentation)
    Password string `json:"password" swaggerignore:"true"`
    
    // Account status
    Status string `json:"status" enum:"active,inactive,pending" default:"active"`
    
    // Account balance
    Balance float64 `json:"balance" minimum:"0" extensions:"x-currency=USD"`
}

// CreateUser creates a new user
// @Summary      Create user
// @Description  Create a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      User  true  "User data"
// @Success      201   {object}  User
// @Header       201   {string}  X-Request-ID  "Request identifier"
// @Header       201   {string}  Location      "User resource URL"
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /users [post]
// @Security     ApiKeyAuth
func CreateUser(c *gin.Context) {
    // Implementation
}
```

#### 2. Generate Documentation

**OpenAPI 3.1.0 (default):**

```bash
nexs-swag init
# or explicitly
nexs-swag init --openapi-version 3.1
```

**Swagger 2.0:**

```bash
nexs-swag init --openapi-version 2.0
```

**Generate both versions:**

```bash
# OpenAPI 3.1.0 in ./docs/v3
nexs-swag init --output ./docs/v3 --openapi-version 3.1

# Swagger 2.0 in ./docs/v2
nexs-swag init --output ./docs/v2 --openapi-version 2.0
```

Or specify directories:

```bash
nexs-swag init -d ./cmd/api -o ./docs --openapi-version 3.1
```

#### 3. Generated Files

**OpenAPI 3.1.0 (default):**
- **`docs/openapi.json`** - OpenAPI 3.1.0 specification in JSON
- **`docs/openapi.yaml`** - OpenAPI 3.1.0 specification in YAML
- **`docs/docs.go`** - Embedded Go documentation file

**Swagger 2.0 (with `--openapi-version 2.0`):**
- **`docs/swagger.json`** - Swagger 2.0 specification in JSON
- **`docs/swagger.yaml`** - Swagger 2.0 specification in YAML
- **`docs/docs.go`** - Embedded Go documentation file

#### 4. Integrate with Your Application

Import the generated docs package:

```go
import _ "your-module/docs"  // Import generated docs

func main() {
    r := gin.Default()
    
    // Serve Swagger UI
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    r.Run(":8080")
}
```

Browse to http://localhost:8080/swagger/index.html to see your API documentation!

## Supported Web Frameworks

nexs-swag works with all popular Go web frameworks through swagger middleware packages:

- [gin](https://github.com/swaggo/gin-swagger) - `github.com/swaggo/gin-swagger`
- [echo](https://github.com/swaggo/echo-swagger) - `github.com/swaggo/echo-swagger`
- [fiber](https://github.com/gofiber/swagger) - `github.com/gofiber/swagger`
- [net/http](https://github.com/swaggo/http-swagger) - `github.com/swaggo/http-swagger`
- [gorilla/mux](https://github.com/swaggo/http-swagger) - `github.com/swaggo/http-swagger`
- [go-chi/chi](https://github.com/swaggo/http-swagger) - `github.com/swaggo/http-swagger`
- [hertz](https://github.com/hertz-contrib/swagger) - `github.com/hertz-contrib/swagger`
- [buffalo](https://github.com/swaggo/buffalo-swagger) - `github.com/swaggo/buffalo-swagger`

## How to use with Gin

Complete example using Gin framework. Find the full source in [examples/03-general-info](examples/03-general-info).

**1. Install dependencies:**

```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

**2. Add general API info to `main.go`:**

```go
package main

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    
    _ "your-project/docs"  // Import generated docs
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server with nexs-swag.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
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
    
    // Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    r.Run(":8080")
}
```

**3. Add operation annotations:**

```go
// GetUser retrieves a user by ID
// @Summary      Get user by ID
// @Description  Get user details by their unique identifier
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"  minimum(1)
// @Success      200  {object}  User
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [get]
// @Security     ApiKeyAuth
func GetUser(c *gin.Context) {
    // Implementation
}
```

**4. Generate and run:**

```bash
nexs-swag init
go run main.go
```

Visit http://localhost:8080/swagger/index.html

## CLI Reference

### init Command

Generate OpenAPI documentation from source code.

```bash
nexs-swag init [options]
```

**Main Options:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--generalInfo` | `-g` | `main.go` | Path to file with general API info |
| `--dir` | `-d` | `./` | Directories to parse (comma-separated) |
| `--output` | `-o` | `./docs` | Output directory for generated files |
| `--outputTypes` | `--ot` | `go,json,yaml` | Output file types |
| `--parseDepth` | | `100` | Dependency parse depth |
| `--parseDependency` | `--pd` | `false` | Parse go files in dependencies |
| `--parseDependencyLevel` | `--pdl` | `0` | 0=disabled, 1=models, 2=operations, 3=all |
| `--parseInternal` | | `false` | Parse internal packages |
| `--parseGoList` | | `true` | Use `go list` for parsing |
| `--propertyStrategy` | `-p` | `camelcase` | Property naming: `snakecase`, `camelcase`, `pascalcase` |
| `--requiredByDefault` | | `false` | Mark all fields as required |
| `--validate` | | `true` | Validate generated spec |
| `--exclude` | | | Exclude directories (comma-separated) |
| `--tags` | `-t` | | Filter by tags (comma-separated) |
| `--markdownFiles` | `--md` | | Parse markdown files for descriptions |
| `--codeExampleFiles` | `--cef` | | Parse code example files |
| `--generatedTime` | | `false` | Add generation timestamp |
| `--instanceName` | | `swagger` | Instance name for multiple docs |
| `--overridesFile` | | `.swaggo` | Type overrides file |
| `--templateDelims` | `--td` | `{{,}}` | Custom template delimiters |
| `--collectionFormat` | `--cf` | `csv` | Default array format |
| `--parseFuncBody` | | `false` | Parse function bodies |
| `--openapi-version` | `--ov` | `3.1` | OpenAPI version: `2.0`, `3.0`, `3.1` |

**Examples:**

```bash
# Basic usage (OpenAPI 3.1.0)
nexs-swag init

# Generate Swagger 2.0
nexs-swag init --openapi-version 2.0

# Generate both versions
nexs-swag init --output ./docs/v3 --openapi-version 3.1
nexs-swag init --output ./docs/v2 --openapi-version 2.0

# Specify directories
nexs-swag init -d ./cmd/api,./internal/handlers -o ./api-docs

# Parse dependencies (level 1 - models only)
nexs-swag init --parseDependency --parseDependencyLevel 1

# Parse internal packages
nexs-swag init --parseInternal

# JSON output only
nexs-swag init --outputTypes json

# Snake case property names
nexs-swag init --propertyStrategy snakecase

# Filter by tags
nexs-swag init --tags "users,products"

# Use markdown descriptions
nexs-swag init --markdownFiles ./docs/api

# Custom template delimiters (avoid conflicts)
nexs-swag init --templateDelims "[[,]]"
```

### fmt Command

Format swagger comments automatically.

```bash
nexs-swag fmt [options]
```

**Options:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | `-d` | `./` | Directories to format |
| `--exclude` | | | Exclude directories |
| `--generalInfo` | `-g` | `main.go` | General info file |

**Example:**

```bash
# Format current directory
nexs-swag fmt

# Format specific directory
nexs-swag fmt -d ./internal/api

# Exclude vendor
nexs-swag fmt --exclude ./vendor
```

## Implementation Status

### OpenAPI 3.1.0 Support

✅ **Fully Implemented:**
- JSON Schema 2020-12
- Basic structure (Info, Servers, Paths, Components)
- Request bodies with multiple content types
- Response definitions with headers
- Parameter definitions (path, query, header, cookie)
- Security schemes (Basic, Bearer, API Key, OAuth2)
- Schema composition (allOf, oneOf, anyOf)
- Schema validation (min, max, pattern, enum, format)
- Examples and descriptions
- External documentation
- Custom extensions (x-*)
- Webhooks
- Tags and grouping

### Swagger 2.0 Support

✅ **Fully Compatible:**
- Basic structure (Info, Host, BasePath, Paths, Definitions)
- Request/response definitions
- Parameter definitions (path, query, header, body, formData)
- Security definitions (Basic, API Key, OAuth2)
- Schema composition (allOf)
- Schema validation (min, max, pattern, enum, format)
- Examples and descriptions
- External documentation
- Custom extensions (x-*)
- Tags and grouping

⚠️ **Automatic Conversion with Warnings:**
- Servers → Host + BasePath (uses first server URL)
- Webhooks → ⚠️ Not supported in Swagger 2.0
- Callbacks → ⚠️ Not supported in Swagger 2.0
- oneOf/anyOf → ⚠️ Limited support (converted to object)
- nullable property → Uses `x-nullable` extension

### swaggo/swag Compatibility

✅ **100% Compatible:**
- All annotations (@title, @version, @description, etc.)
- All struct tags (json, binding, validate, swaggertype, swaggerignore, extensions)
- All CLI flags (28/28 implemented)
- Commands: init, fmt
- Type overrides via .swaggo file
- Markdown descriptions
- Code examples

## Declarative Comments Format

### General API Info

Add to your `main.go` or entry point:

| Annotation | Example | Description |
|------------|---------|-------------|
| `@title` | `@title My API` | **Required.** API title |
| `@version` | `@version 1.0` | **Required.** API version |
| `@description` | `@description This is my API` | API description |
| `@description.markdown` | `@description.markdown` | Load description from api.md |
| `@termsOfService` | `@termsOfService http://example.com/terms` | Terms of service URL |
| `@contact.name` | `@contact.name API Support` | Contact name |
| `@contact.url` | `@contact.url http://example.com` | Contact URL |
| `@contact.email` | `@contact.email support@example.com` | Contact email |
| `@license.name` | `@license.name Apache 2.0` | **Required.** License name |
| `@license.url` | `@license.url http://apache.org/licenses` | License URL |
| `@host` | `@host localhost:8080` | API host |
| `@BasePath` | `@BasePath /api/v1` | Base path |
| `@schemes` | `@schemes http https` | Transfer protocols |
| `@accept` | `@accept json xml` | Default Accept MIME types |
| `@produce` | `@produce json xml` | Default Produce MIME types |
| `@tag.name` | `@tag.name Users` | Tag name |
| `@tag.description` | `@tag.description User operations` | Tag description |
| `@externalDocs.description` | `@externalDocs.description OpenAPI` | External docs description |
| `@externalDocs.url` | `@externalDocs.url https://swagger.io` | External docs URL |
| `@x-<name>` | `@x-custom-info value` | Custom extension |

**Version-Specific Annotations:**

When generating **Swagger 2.0** (`--openapi-version 2.0`):
- Use `@host`, `@BasePath`, and `@schemes` annotations
- These are automatically converted to the `host`, `basePath`, and `schemes` fields

When generating **OpenAPI 3.x** (`--openapi-version 3.0` or `3.1`):
- Use `@server` annotation: `// @server http://localhost:8080/api/v1 Development server`
- Alternatively, use `@host`, `@BasePath`, and `@schemes` which will be converted to servers

Both annotation styles work with either version - the converter handles the transformation automatically.

**Security Definitions:**

```go
// Basic Authentication
// @securityDefinitions.basic BasicAuth

// API Key
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// OAuth2 Application Flow
// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants admin access
```

### API Operation

Add to handler functions:

| Annotation | Example | Description |
|------------|---------|-------------|
| `@Summary` | `@Summary Get user` | Short summary |
| `@Description` | `@Description Get user by ID` | Detailed description |
| `@Description.markdown` | `@Description.markdown details` | Load from details.md |
| `@Tags` | `@Tags users,accounts` | Operation tags |
| `@Accept` | `@Accept json` | Request content type |
| `@Produce` | `@Produce json,xml` | Response content types |
| `@Param` | See below | Parameter definition |
| `@Success` | `@Success 200 {object} User` | Success response |
| `@Failure` | `@Failure 400 {object} Error` | Error response |
| `@Header` | `@Header 200 {string} Token` | Response header |
| `@Router` | `@Router /users/{id} [get]` | Route path and method |
| `@Security` | `@Security ApiKeyAuth` | Security requirement |
| `@Deprecated` | `@Deprecated` | Mark as deprecated |
| `@x-<name>` | `@x-code-samples file.json` | Custom extension |

**Parameter Syntax:**

```
@Param <name> <in> <type> <required> <description> [attributes]
```

- **name**: Parameter name
- **in**: `query`, `path`, `header`, `body`, `formData`
- **type**: Data type (string, int, bool, object, array, file)
- **required**: `true` or `false`
- **description**: Description (in quotes if contains spaces)
- **attributes**: Optional validation attributes

**Examples:**

```go
// Path parameter
// @Param id path int true "User ID" minimum(1) maximum(1000)

// Query parameter with validation
// @Param name query string false "User name" minLength(3) maxLength(50)

// Query parameter with enum
// @Param status query string false "Status filter" Enums(active,inactive,pending)

// Query array with collection format
// @Param tags query []string false "Tags" collectionFormat(multi)

// Header parameter
// @Param X-Request-ID header string true "Request ID" format(uuid)

// Body parameter
// @Param user body User true "User object"

// Form data with file
// @Param avatar formData file true "Avatar image"
```

**Response Syntax:**

```go
// Simple response
// @Success 200 {object} User

// Response with description
// @Success 201 {object} User "User created successfully"

// Array response
// @Success 200 {array} User "List of users"

// Primitive response
// @Success 200 {string} string "Success message"

// Generic response
// @Success 200 {object} Response{data=User} "User response"

// Multiple data fields
// @Success 200 {object} Response{data=User,meta=Metadata}
```

**Header Syntax:**

```go
// Single status code
// @Header 200 {string} X-Request-ID "Request identifier"

// Multiple status codes
// @Header 200,201 {string} Location "Resource URL"

// All responses
// @Header all {string} X-API-Version "API version"
```

### Struct Tags

#### Standard Tags

```go
type User struct {
    // JSON serialization
    ID   int    `json:"id"`
    Name string `json:"name,omitempty"`  // omitempty = not required
    
    // Validation (Gin binding)
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=150"`
    
    // Validation (go-playground/validator)
    UUID  string `json:"uuid" validate:"required,uuid"`
    
    // OpenAPI attributes
    Price  float64  `json:"price" minimum:"0" maximum:"9999.99"`
    Status string   `json:"status" enum:"active,inactive" default:"active"`
    SKU    string   `json:"sku" pattern:"^[A-Z]{3}-[0-9]{6}$"`
    Items  []string `json:"items" minLength:"1" maxLength:"100"`
    
    // Example value
    Bio string `json:"bio" example:"Software developer"`
    
    // Format
    CreatedAt string `json:"created_at" format:"date-time"`
}
```

#### swaggertype - Type Override

Convert custom types to OpenAPI types:

```go
type Account struct {
    // Override sql.NullInt64 to integer
    ID sql.NullInt64 `json:"id" swaggertype:"integer"`
    
    // Custom time type to unix timestamp (integer)
    CreatedAt TimestampTime `json:"created_at" swaggertype:"primitive,integer"`
    
    // Byte array to base64 string
    Certificate []byte `json:"cert" swaggertype:"string" format:"base64"`
    
    // Custom number array
    Coeffs []big.Float `json:"coeffs" swaggertype:"array,number"`
    
    // Nested custom types
    Metadata map[string]interface{} `json:"metadata" swaggertype:"object"`
}
```

**Format:** `swaggertype:"[primitive,]<type>"`

- For primitive types: `swaggertype:"string"`, `swaggertype:"integer"`, `swaggertype:"number"`, `swaggertype:"boolean"`
- For arrays: `swaggertype:"array,<element-type>"`
- For objects: `swaggertype:"object"`

#### swaggerignore - Hide Fields

Exclude fields from documentation (still present in JSON):

```go
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    
    // Present in JSON, hidden in docs
    Password string `json:"password" swaggerignore:"true"`
    
    // Internal field, not in JSON or docs
    internal string `swaggerignore:"true"`
    
    // Sensitive data
    SSN string `json:"ssn" swaggerignore:"true"`
}
```

#### extensions - Custom Extensions

Add custom metadata with `x-*` prefix:

```go
type Product struct {
    // Primary key indicator
    ID int `json:"id" extensions:"x-primary-key=true"`
    
    // Currency formatting
    Price float64 `json:"price" extensions:"x-currency=USD,x-format=currency"`
    
    // Multiple extensions
    Name string `json:"name" extensions:"x-order=1,x-searchable=true,x-filterable=true"`
    
    // Boolean extension
    Featured bool `json:"featured" extensions:"x-promoted=true"`
    
    // Nullable extension
    Discount float64 `json:"discount" extensions:"x-nullable"`
}
```

Generated OpenAPI:

```json
{
  "properties": {
    "id": {
      "type": "integer",
      "x-primary-key": true
    },
    "price": {
      "type": "number",
      "x-currency": "USD",
      "x-format": "currency"
    }
  }
}
```

## Examples

nexs-swag includes 21 comprehensive, runnable examples. Each example demonstrates specific features and includes a README and run script.

### Basic Examples

| Example | Description | Key Features |
|---------|-------------|--------------|
| [01-basic](examples/01-basic) | Basic usage | Minimal setup, simple API |
| [02-formats](examples/02-formats) | Output formats | JSON, YAML, Go output |
| [03-general-info](examples/03-general-info) | General API info | Complete API metadata |

### Advanced Features

| Example | Description | Key Features |
|---------|-------------|--------------|
| [04-property-strategy](examples/04-property-strategy) | Naming strategies | Snake_case, camelCase, PascalCase |
| [05-required-default](examples/05-required-default) | Required by default | Auto-require all fields |
| [06-exclude](examples/06-exclude) | Exclude directories | Filter unwanted paths |
| [07-tags-filter](examples/07-tags-filter) | Tag filtering | Generate subset of APIs |
| [08-parse-internal](examples/08-parse-internal) | Internal packages | Parse internal/ directory |
| [09-parse-dependency](examples/09-parse-dependency) | Dependencies | Parse vendor/go.mod packages |
| [10-dependency-level](examples/10-dependency-level) | Dependency depth | Control parsing level (0-3) |
| [11-parse-golist](examples/11-parse-golist) | Go list parsing | Use `go list` for discovery |

### Documentation Features

| Example | Description | Key Features |
|---------|-------------|--------------|
| [12-markdown-files](examples/12-markdown-files) | Markdown descriptions | Load docs from .md files |
| [13-code-examples](examples/13-code-examples) | Code samples | Multi-language examples |
| [14-overrides-file](examples/14-overrides-file) | Type overrides | .swaggo file configuration |
| [15-generated-time](examples/15-generated-time) | Generation timestamp | Add generation date |
| [16-instance-name](examples/16-instance-name) | Multiple instances | Named documentation sets |
| [17-template-delims](examples/17-template-delims) | Custom delimiters | Avoid template conflicts |

### Validation & Structure

| Example | Description | Key Features |
|---------|-------------|--------------|
| [18-collection-format](examples/18-collection-format) | Array formats | CSV, multi, pipes, SSV, TSV |
| [19-parse-func-body](examples/19-parse-func-body) | Function bodies | Parse inline annotations |
| [20-fmt-command](examples/20-fmt-command) | Format command | Auto-format comments |
| [21-struct-tags](examples/21-struct-tags) | All struct tags | Complete tag reference |

### Running Examples

Each example includes a `run.sh` script:

```bash
cd examples/01-basic
./run.sh
```

Or manually (OpenAPI 3.1.0):

```bash
cd examples/01-basic
nexs-swag init -d . -o ./docs
cat docs/openapi.json
```

Or generate Swagger 2.0:

```bash
cd examples/01-basic
nexs-swag init -d . -o ./docs --openapi-version 2.0
cat docs/swagger.json
```

### Example: Complete CRUD API

See [examples/03-general-info](examples/03-general-info) for a complete CRUD API with:
- Multiple endpoints (GET, POST, PUT, DELETE)
- Request/response models
- Validation rules
- Error responses
- Security schemes
- Response headers

## Quality & Testing

### Test Coverage

```bash
$ go test ./pkg/... -cover
```

| Package | Coverage | Tests |
|---------|----------|-------|
| pkg/converter | 92.3% | 13 tests |
| pkg/format | 95.1% | 15 tests |
| pkg/generator | 71.6% | 16 tests |
| pkg/generator/v2 | 88.4% | 12 tests |
| pkg/generator/v3 | 85.2% | 8 tests |
| pkg/openapi | 83.3% | 22 tests |
| pkg/openapi/v2 | 89.7% | 12 tests |
| pkg/openapi/v3 | 91.5% | 10 tests |
| pkg/parser | 82.1% | 192 tests |
| **Overall** | **87.9%** | **300+ tests** |

### Quality Metrics

- ✅ **0 linter warnings** (golangci-lint with 20+ linters)
- ✅ **0 race conditions** (tested with `-race` flag)
- ✅ **22 integration tests** (runnable examples)
- ✅ **~8,500 lines of test code**
- ✅ **Production-ready** (actively maintained)
- ✅ **100% swaggo/swag compatible**
- ✅ **Dual-version support** (OpenAPI 3.1.0 + Swagger 2.0)

### Running Tests

```bash
# Unit tests
go test ./pkg/... -v

# With coverage
go test ./pkg/... -cover

# With race detection
go test ./pkg/... -race

# Specific package
go test ./pkg/parser -v

# Run examples
cd examples && for d in */; do cd "$d" && ./run.sh && cd ..; done
```

## swaggo/swag Compatibility

nexs-swag is designed as a **drop-in replacement** for swaggo/swag with enhanced features.

### Migration from swaggo/swag

**No changes required!** Simply replace the binary:

```bash
# Instead of
go install github.com/swaggo/swag/cmd/swag@latest

# Use
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest

# Same commands work
nexs-swag init
nexs-swag fmt
```

### Compatibility Table

| Feature | swaggo/swag | nexs-swag | Notes |
|---------|-------------|-----------|-------|
| OpenAPI Version | 2.0 | 3.1.0 | Backward compatible |
| All annotations | ✅ | ✅ | 100% compatible |
| Struct tags | ✅ | ✅ | swaggertype, swaggerignore, extensions |
| CLI flags | ✅ | ✅ | All 28 flags supported |
| .swaggo file | ✅ | ✅ | Type overrides |
| Markdown | ✅ | ✅ | File-based descriptions |
| Code examples | ✅ | ✅ | Multi-language samples |
| Webhooks | ❌ | ✅ | OpenAPI 3.1 feature |
| JSON Schema 2020-12 | ❌ | ✅ | Modern schema |
| Response headers | Limited | ✅ | Full support |
| Test coverage | ~70% | 86.1% | Higher quality |
| Go version | 1.19+ | 1.23+ | Modern Go features |

### What's Different?

**Enhanced (backward compatible):**
- OpenAPI 3.1.0 output (vs 2.0)
- Better nullable handling
- More validation attributes
- Improved error messages
- Better test coverage

**Same API:**
- All command-line flags
- All annotations
- All struct tags
- Generated docs.go structure
- Swagger UI integration

## About the Project

### Project Statistics

- **Lines of Code:** ~5,200 (pkg/ excluding tests)
- **Test Code:** ~8,500 lines
- **Go Files:** 42 implementation files
- **Test Files:** 29 test files
- **Packages:** 9 (converter, format, generator, generator/v2, generator/v3, openapi, openapi/v2, openapi/v3, parser)
- **Examples:** 22 complete examples
- **Test Coverage:** 87.9%
- **OpenAPI Versions:** 2 (Swagger 2.0 + OpenAPI 3.1.0)
- **Dependencies:** 3 direct dependencies
  - urfave/cli/v2 (CLI framework)
  - golang.org/x/tools (Go AST parsing)
  - gopkg.in/yaml.v3 (YAML support)

### Project Structure

```
nexs-swag/
├── cmd/
│   └── nexs-swag/          # CLI entry point
├── pkg/
│   ├── converter/          # Version conversion (v3 ↔ v2)
│   ├── format/             # Code formatting
│   ├── generator/          # OpenAPI generation
│   │   ├── v2/             # Swagger 2.0 generator
│   │   └── v3/             # OpenAPI 3.x generator
│   ├── openapi/            # OpenAPI models
│   │   ├── v2/             # Swagger 2.0 models
│   │   └── v3/             # OpenAPI 3.x models
│   └── parser/             # Go code parsing (AST)
├── examples/               # 22 examples
│   ├── 01-basic/
│   ├── 02-formats/
│   └── ...
├── docs/                   # Project documentation
├── README.md               # This file
├── README_pt.md            # Portuguese version
├── README_es.md            # Spanish version
└── LICENSE                 # MIT License
```

### Inspiration & Credits

This project was inspired by [swaggo/swag](https://github.com/swaggo/swag) and built to extend its capabilities with full OpenAPI 3.1.0 support while maintaining 100% backward compatibility.

**Credits:**
- [swaggo/swag](https://github.com/swaggo/swag) - Original Swagger 2.0 generator
- [OpenAPI Initiative](https://www.openapis.org/) - OpenAPI Specification
- [Go Team](https://go.dev/) - Amazing language and tools
- All contributors and the Go community

## Contributing

Contributions are welcome! Please follow these guidelines:

### How to Contribute

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Make** your changes
4. **Add** tests for new functionality
5. **Run** tests (`go test ./...`)
6. **Run** linter (`golangci-lint run`)
7. **Commit** your changes (`git commit -m 'Add amazing feature'`)
8. **Push** to the branch (`git push origin feature/amazing-feature`)
9. **Open** a Pull Request

### Development Setup

```bash
# Clone repository
git clone https://github.com/fsvxavier/nexs-swag.git
cd nexs-swag

# Install dependencies
go mod download

# Run tests
go test ./... -v

# Run linter
golangci-lint run

# Build
go build -o nexs-swag ./cmd/nexs-swag
```

### Reporting Issues

Please include:
- Go version (`go version`)
- nexs-swag version (`nexs-swag --version`)
- Minimal reproducible example
- Expected vs actual behavior

### Feature Requests

Open an issue with:
- Clear description of the feature
- Use case and benefits
- Proposed implementation (if any)

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2024 Fabricio Xavier

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## Support & Community

- **Issues:** [GitHub Issues](https://github.com/fsvxavier/nexs-swag/issues)
- **Discussions:** [GitHub Discussions](https://github.com/fsvxavier/nexs-swag/discussions)
- **Documentation:** [Wiki](https://github.com/fsvxavier/nexs-swag/wiki)
- **Examples:** [examples/](examples/)

---

**Made with ❤️ for the Go community**

[⬆ Back to top](#nexs-swag)

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
