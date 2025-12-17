# nexs-swag

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

[![Versión Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1.0-6BA539?style=flat&logo=openapiinitiative)](https://spec.openapis.org/oas/v3.1.0)
[![Swagger](https://img.shields.io/badge/Swagger-2.0-85EA2D?style=flat&logo=swagger)](https://swagger.io/specification/v2/)
[![Licencia](https://img.shields.io/badge/Licencia-MIT-blue.svg)](LICENSE)
[![Cobertura](https://img.shields.io/badge/Cobertura-80.1%25-brightgreen.svg)](/)
[![Ejemplos](https://img.shields.io/badge/Ejemplos-27-blue.svg)](examples/)

**Genera automáticamente documentación OpenAPI 3.1.0 o Swagger 2.0 a partir de anotaciones en código Go.**

nexs-swag convierte anotaciones Go en especificación OpenAPI 3.1.0 o Swagger 2.0. Fue diseñado como una evolución de [swaggo/swag](https://github.com/swaggo/swag) con soporte completo para la especificación OpenAPI más reciente y compatibilidad total con Swagger 2.0.

## Índice

- [Visión General](#visión-general)
- [Primeros Pasos](#primeros-pasos)
  - [Instalación](#instalación)
  - [Inicio Rápido](#inicio-rápido)
- [Frameworks Web Soportados](#frameworks-web-soportados)
- [Cómo usar con Gin](#cómo-usar-con-gin)
- [Referencia CLI](#referencia-cli)
  - [Comando init](#comando-init)
  - [Comando fmt](#comando-fmt)
- [Estado de Implementación](#estado-de-implementación)
- [Versiones OpenAPI](OPENAPI_VERSIONS.md) - Guía completa de todas las versiones soportadas
- [Formato de Comentarios Declarativos](#formato-de-comentarios-declarativos)
  - [Información General de la API](#información-general-de-la-api)
  - [Operación de API](#operación-de-api)
  - [Tags de Struct](#tags-de-struct)
- [Ejemplos](#ejemplos)
- [Calidad y Pruebas](#calidad-y-pruebas)
- [Compatibilidad con swaggo/swag](#compatibilidad-con-swaggoswag)
- [Acerca del Proyecto](#acerca-del-proyecto)
- [Contribuyendo](#contribuyendo)
- [Licencia](#licencia)

## Visión General

### Características Principales

- ✅ **100% compatible con swaggo/swag** - Sustituto directo con todas las anotaciones y tags
- ✅ **Soporte a múltiples versiones OpenAPI** - Genera v2.0.0, v3.0.x, v3.1.x o v3.2.0
- ✅ **OpenAPI 3.2.0** - Soporte completo para la versión más reciente (método QUERY, streaming, etc)
- ✅ **OpenAPI 3.1.x** - Compatible con JSON Schema 2020-12, webhooks y características modernas
- ✅ **OpenAPI 3.0.x** - Todas las versiones desde 3.0.0 hasta 3.0.4
- ✅ **Swagger 2.0** - Compatibilidad total con sistemas legados
- ✅ **Conversión automática** - Conversión entre formatos con avisos para incompatibilidades
- ✅ **20+ atributos de validación** - minimum, maximum, pattern, enum, format y más
- ✅ **Validación de frameworks** - Soporte nativo para Gin (binding) y go-playground/validator
- ✅ **Headers de respuesta** - Documentación completa de headers
- ✅ **Múltiples tipos de contenido** - JSON, XML, YAML, CSV, PDF y tipos MIME personalizados
- ✅ **Extensiones personalizadas** - Soporte completo para x-*
- ✅ **@x-visibility** - Genera documentación pública/privada separada desde una única base de código
- ✅ **80.1% de cobertura de pruebas** - Listo para producción con suite de pruebas integral incluyendo pruebas roundtrip
- ✅ **27 ejemplos funcionales** - Aprende con ejemplos completos y ejecutables

### ¿Por qué nexs-swag?

| Característica | swaggo/swag | nexs-swag |
|----------------|-------------|-----------|
| OpenAPI 3.2.0 | ❌ | ✅ |
| OpenAPI 3.1.x | ❌ | ✅ |
| OpenAPI 3.0.x | ❌ | ✅ |
| Swagger 2.0 | ✅ | ✅ |
| Múltiples Versiones | ❌ | ✅ (todas del mismo código) |
| JSON Schema | Draft 4 | Draft 4 + 2020-12 |
| Webhooks | ❌ | ✅ (OpenAPI 3.1+) |
| Headers de Respuesta | Limitado | Soporte Completo |
| Soporte a Nullable | `x-nullable` | Nativo + `x-nullable` |
| Cobertura de Pruebas | ~70% | 80.1% |
| Ejemplos | ~10 | 25 |
| Versión Go | 1.19+ | 1.23+ |

## Primeros Pasos

### Instalación

#### Usando go install (Recomendado)

```bash
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest
```

Para verificar la instalación:

```bash
nexs-swag --version
```

#### Compilando desde el Código Fuente

Requiere [Go 1.23 o superior](https://go.dev/dl/).

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

### Inicio Rápido

#### 1. Agregar Anotaciones de la API

Agrega anotaciones generales de la API a tu `main.go`:

```go
package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
)

// @title           API de Gestión de Usuarios
// @version         1.0.0
// @description     Una API de gestión de usuarios con documentación OpenAPI 3.1.0 completa
// @termsOfService  http://swagger.io/terms/

// @contact.name   Soporte de la API
// @contact.url    http://www.example.com/soporte
// @contact.email  soporte@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
    r := gin.Default()
    // Configuración de tu aplicación
    r.Run(":8080")
}

// User representa un usuario del sistema
type User struct {
    // ID del usuario (sql.NullInt64 → integer en OpenAPI)
    ID sql.NullInt64 `json:"id" swaggertype:"integer" extensions:"x-primary-key=true"`
    
    // Nombre completo (3-100 caracteres obligatorio)
    Name string `json:"name" binding:"required" minLength:"3" maxLength:"100" example:"Juan Silva"`
    
    // Dirección de correo electrónico (validado)
    Email string `json:"email" binding:"required,email" format:"email" extensions:"x-unique=true"`
    
    // Contraseña (oculta de la documentación)
    Password string `json:"password" swaggerignore:"true"`
    
    // Estado de la cuenta
    Status string `json:"status" enum:"active,inactive,pending" default:"active"`
    
    // Saldo de la cuenta
    Balance float64 `json:"balance" minimum:"0" extensions:"x-currency=USD"`
}

// CreateUser crea un nuevo usuario
// @Summary      Crear usuario
// @Description  Crea un nuevo usuario en el sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      User  true  "Datos del usuario"
// @Success      201   {object}  User
// @Header       201   {string}  X-Request-ID  "Identificador de la petición"
// @Header       201   {string}  Location      "URL del recurso del usuario"
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /users [post]
// @Security     ApiKeyAuth
func CreateUser(c *gin.Context) {
    // Implementación
}
```

#### 2. Generar Documentación

**OpenAPI 3.1.0 (predeterminado):**

```bash
nexs-swag init
# o explícitamente
nexs-swag init --openapi-version 3.1
```

**Swagger 2.0:**

```bash
nexs-swag init --openapi-version 2.0
```

**Generar ambas versiones:**

```bash
# OpenAPI 3.1.0 en ./docs/v3
nexs-swag init --output ./docs/v3 --openapi-version 3.1

# Swagger 2.0 en ./docs/v2
nexs-swag init --output ./docs/v2 --openapi-version 2.0
```

O especifica los directorios:

```bash
nexs-swag init -d ./cmd/api -o ./docs --openapi-version 3.1
```

#### 3. Archivos Generados

**OpenAPI 3.1.0 (predeterminado):**
- **`docs/openapi.json`** - Especificación OpenAPI 3.1.0 en JSON
- **`docs/openapi.yaml`** - Especificación OpenAPI 3.1.0 en YAML
- **`docs/docs.go`** - Archivo de documentación Go embebido

**Swagger 2.0 (con `--openapi-version 2.0`):**
- **`docs/swagger.json`** - Especificación Swagger 2.0 en JSON
- **`docs/swagger.yaml`** - Especificación Swagger 2.0 en YAML
- **`docs/docs.go`** - Archivo de documentación Go embebido

#### 4. Integrar con Tu Aplicación

Importa el paquete docs generado:

```go
import _ "tu-modulo/docs"  // Importar docs generado

func main() {
    r := gin.Default()
    
    // Servir Swagger UI
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    r.Run(":8080")
}
```

¡Accede a http://localhost:8080/swagger/index.html para ver tu documentación API!

## Frameworks Web Soportados

nexs-swag funciona con todos los frameworks web Go populares a través de paquetes middleware swagger:

- [gin](https://github.com/swaggo/gin-swagger) - `github.com/swaggo/gin-swagger`
- [echo](https://github.com/swaggo/echo-swagger) - `github.com/swaggo/echo-swagger`
- [fiber](https://github.com/gofiber/swagger) - `github.com/gofiber/swagger`
- [net/http](https://github.com/swaggo/http-swagger) - `github.com/swaggo/http-swagger`
- [gorilla/mux](https://github.com/swaggo/http-swagger) - `github.com/swaggo/http-swagger`
- [go-chi/chi](https://github.com/swaggo/http-swagger) - `github.com/swaggo/http-swagger`
- [hertz](https://github.com/hertz-contrib/swagger) - `github.com/hertz-contrib/swagger`
- [buffalo](https://github.com/swaggo/buffalo-swagger) - `github.com/swaggo/buffalo-swagger`

## Cómo usar con Gin

Ejemplo completo usando framework Gin. Encuentra el código completo en [examples/03-general-info](examples/03-general-info).

**1. Instalar dependencias:**

```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

**2. Agregar información general de la API a `main.go`:**

```go
package main

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    
    _ "tu-proyecto/docs"  // Importar docs generado
)

// @title           API de Ejemplo Swagger
// @version         1.0
// @description     Este es un servidor de ejemplo con nexs-swag.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Soporte de la API
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

**3. Agregar anotaciones de operación:**

```go
// GetUser recupera un usuario por ID
// @Summary      Buscar usuario por ID
// @Description  Buscar detalles del usuario por su identificador único
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del Usuario"  minimum(1)
// @Success      200  {object}  User
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [get]
// @Security     ApiKeyAuth
func GetUser(c *gin.Context) {
    // Implementación
}
```

**4. Generar y ejecutar:**

```bash
nexs-swag init
go run main.go
```

Visita http://localhost:8080/swagger/index.html

## Referencia CLI

### Comando init

Genera documentación OpenAPI a partir del código fuente.

```bash
nexs-swag init [opciones]
```

**Opciones Principales:**

| Flag | Corto | Predeterminado | Descripción |
|------|-------|----------------|-------------|
| `--generalInfo` | `-g` | `main.go` | Ruta al archivo con información general de la API |
| `--dir` | `-d` | `./` | Directorios para analizar (separados por coma) |
| `--output` | `-o` | `./docs` | Directorio de salida para archivos generados |
| `--outputTypes` | `--ot` | `go,json,yaml` | Tipos de archivo de salida |
| `--parseDepth` | | `100` | Profundidad de análisis de dependencia |
| `--parseDependency` | `--pd` | `false` | Analizar archivos go en dependencias |
| `--parseDependencyLevel` | `--pdl` | `0` | 0=deshabilitado, 1=modelos, 2=operaciones, 3=todo |
| `--parseInternal` | | `false` | Analizar paquetes internos |
| `--parseGoList` | | `true` | Usar `go list` para análisis |
| `--propertyStrategy` | `-p` | `camelcase` | Nomenclatura de propiedad: `snakecase`, `camelcase`, `pascalcase` |
| `--requiredByDefault` | | `false` | Marcar todos los campos como obligatorios |
| `--validate` | | `true` | Validar especificación generada |
| `--exclude` | | | Excluir directorios (separados por coma) |
| `--tags` | `-t` | | Filtrar por tags (separados por coma) |
| `--markdownFiles` | `--md` | | Analizar archivos markdown para descripciones |
| `--codeExampleFiles` | `--cef` | | Analizar archivos de ejemplo de código |
| `--generatedTime` | | `false` | Agregar marca de tiempo de generación |
| `--instanceName` | | `swagger` | Nombre de instancia para múltiples docs |
| `--overridesFile` | | `.swaggo` | Archivo de overrides de tipo |
| `--templateDelims` | `--td` | `{{,}}` | Delimitadores de plantilla personalizados |
| `--collectionFormat` | `--cf` | `csv` | Formato de array predeterminado |
| `--parseFuncBody` | | `false` | Analizar cuerpos de función |
| `--includeTypes` | `--it` | `all` | Filtrar tipos a incluir: `struct`, `interface`, `func`, `const`, `type`, `all` |
| `--openapi-version` | `--ov` | `3.1` | Versión OpenAPI: `2.0`, `3.0`, `3.1` |

> **⚠️ Importante: Sintaxis de Flags Booleanos**
>
> Los flags booleanos aceptan dos sintaxis válidas:
> - ✅ **Sin valor** (presencia = true): `--parseInternal`, `--pd`
> - ✅ **Con signo de igual**: `--parseInternal=true`, `--pd=false`
> - ❌ **Incorrecto** (separado por espacio): `--parseInternal true`, `--pd true`
>
> La sintaxis separada por espacio no funciona porque el parser CLI trata la palabra después del flag como un argumento posicional separado, no como el valor del flag.

**Ejemplos:**

```bash
# Uso básico (OpenAPI 3.1.0)
nexs-swag init

# Generar Swagger 2.0
nexs-swag init --openapi-version 2.0

# Generar ambas versiones
nexs-swag init --output ./docs/v3 --openapi-version 3.1
nexs-swag init --output ./docs/v2 --openapi-version 2.0

# Especificar directorios
nexs-swag init -d ./cmd/api,./internal/handlers -o ./api-docs

# Analizar dependencias (nivel 1 - solo modelos)
nexs-swag init --parseDependency --parseDependencyLevel 1
# O con sintaxis explícita:
nexs-swag init --parseDependency=true --parseDependencyLevel 1

# Analizar paquetes internos
nexs-swag init --parseInternal
# O explícitamente:
nexs-swag init --parseInternal=true

# Salida solo JSON
nexs-swag init --outputTypes json

# Nombres de propiedad en snake_case
nexs-swag init --propertyStrategy snakecase

# Filtrar por tags
nexs-swag init --tags "users,products"

# Usar descripciones en markdown
nexs-swag init --markdownFiles ./docs/api

# Delimitadores de plantilla personalizados (evitar conflictos)
nexs-swag init --templateDelims "[[,]]"

# Filtrar tipos a incluir (solo structs)
nexs-swag init --includeTypes struct

# Filtrar múltiples categorías de tipos
nexs-swag init --includeTypes "struct,interface"

# Forma corta
nexs-swag init -it struct
```

### Comando fmt

Formatea comentarios swagger automáticamente.

```bash
nexs-swag fmt [opciones]
```

**Opciones:**

| Flag | Corto | Predeterminado | Descripción |
|------|-------|----------------|-------------|
| `--dir` | `-d` | `./` | Directorios para formatear |
| `--exclude` | | | Excluir directorios |
| `--generalInfo` | `-g` | `main.go` | Archivo de información general |

**Ejemplo:**

```bash
# Formatear directorio actual
nexs-swag fmt

# Formatear directorio específico
nexs-swag fmt -d ./internal/api

# Excluir vendor
nexs-swag fmt --exclude ./vendor
```

## Estado de Implementación

### Soporte OpenAPI 3.1.0

✅ **Totalmente Implementado:**
- JSON Schema 2020-12
- Estructura básica (Info, Servers, Paths, Components)
- Request bodies con múltiples content types
- Definiciones de respuesta con headers
- Definiciones de parámetros (path, query, header, cookie)
- Security schemes (Basic, Bearer, API Key, OAuth2)
- Composición de schemas (allOf, oneOf, anyOf)
- Validación de schemas (min, max, pattern, enum, format)
- Ejemplos y descripciones
- Documentación externa
- Extensiones personalizadas (x-*)
- Webhooks
- Tags y agrupamiento

### Soporte Swagger 2.0

✅ **Totalmente Compatible:**
- Estructura básica (Info, Host, BasePath, Paths, Definitions)
- Definiciones de request/response
- Definiciones de parámetros (path, query, header, body, formData)
- Definiciones de seguridad (Basic, API Key, OAuth2)
- Composición de schemas (allOf)
- Validación de schemas (min, max, pattern, enum, format)
- Ejemplos y descripciones
- Documentación externa
- Extensiones personalizadas (x-*)
- Tags y agrupamiento

⚠️ **Conversión Automática con Avisos:**
- Servers → Host + BasePath (usa la primera URL de server)
- Webhooks → ⚠️ No soportado en Swagger 2.0
- Callbacks → ⚠️ No soportado en Swagger 2.0
- oneOf/anyOf → ⚠️ Soporte limitado (convertido a object)
- propiedad nullable → Usa extensión `x-nullable`

### Compatibilidad con swaggo/swag

✅ **100% Compatible:**
- Todas las anotaciones (@title, @version, @description, etc.)
- Todas las tags de struct (json, binding, validate, swaggertype, swaggerignore, extensions)
- Todos los flags CLI (28/28 implementados)
- Comandos: init, fmt
- Type overrides vía archivo .swaggo
- Descripciones en Markdown
- Ejemplos de código

## Formato de Comentarios Declarativos

### Información General de la API

Agrega a tu `main.go` o punto de entrada:

| Anotación | Ejemplo | Descripción |
|-----------|---------|-------------|
| `@title` | `@title Mi API` | **Obligatorio.** Título de la API |
| `@version` | `@version 1.0` | **Obligatorio.** Versión de la API |
| `@description` | `@description Esta es mi API` | Descripción de la API |
| `@description.markdown` | `@description.markdown` | Cargar descripción de api.md |
| `@termsOfService` | `@termsOfService http://example.com/terms` | URL de los términos de servicio |
| `@contact.name` | `@contact.name Soporte de la API` | Nombre del contacto |
| `@contact.url` | `@contact.url http://example.com` | URL del contacto |
| `@contact.email` | `@contact.email support@example.com` | Email del contacto |
| `@license.name` | `@license.name Apache 2.0` | **Obligatorio.** Nombre de la licencia |
| `@license.url` | `@license.url http://apache.org/licenses` | URL de la licencia |
| `@host` | `@host localhost:8080` | Host de la API |
| `@BasePath` | `@BasePath /api/v1` | Ruta base |
| `@schemes` | `@schemes http https` | Protocolos de transferencia |
| `@accept` | `@accept json xml` | Tipos MIME Accept predeterminados |
| `@produce` | `@produce json xml` | Tipos MIME Produce predeterminados |
| `@tag.name` | `@tag.name Users` | Nombre de la tag |
| `@tag.description` | `@tag.description Operaciones de usuario` | Descripción de la tag |
| `@externalDocs.description` | `@externalDocs.description OpenAPI` | Descripción de docs externos |
| `@externalDocs.url` | `@externalDocs.url https://swagger.io` | URL de docs externos |
| `@x-<nombre>` | `@x-custom-info value` | Extensión personalizada |

**Anotaciones Específicas de Versión:**

Al generar **Swagger 2.0** (`--openapi-version 2.0`):
- Usa anotaciones `@host`, `@BasePath` y `@schemes`
- Estas son automáticamente convertidas a los campos `host`, `basePath` y `schemes`

Al generar **OpenAPI 3.x** (`--openapi-version 3.0` o `3.1`):
- Usa anotación `@server`: `// @server http://localhost:8080/api/v1 Servidor de desarrollo`
- Alternativamente, usa `@host`, `@BasePath` y `@schemes` que serán convertidos a servers

Ambos estilos de anotación funcionan con cualquier versión - el conversor maneja la transformación automáticamente.

**Definiciones de Seguridad:**

```go
// Autenticación Basic
// @securityDefinitions.basic BasicAuth

// API Key
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// OAuth2 Application Flow
// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Concede acceso de escritura
// @scope.admin Concede acceso de administrador
```

### Operación de API

Agrega a las funciones handler:

| Anotación | Ejemplo | Descripción |
|-----------|---------|-------------|
| `@Summary` | `@Summary Buscar usuario` | Resumen corto |
| `@Description` | `@Description Buscar usuario por ID` | Descripción detallada |
| `@Description.markdown` | `@Description.markdown details` | Cargar de details.md |
| `@Tags` | `@Tags users,accounts` | Tags de la operación |
| `@Accept` | `@Accept json` | Tipo de contenido de la petición |
| `@Produce` | `@Produce json,xml` | Tipos de contenido de la respuesta |
| `@Param` | Ver abajo | Definición de parámetro |
| `@Success` | `@Success 200 {object} User` | Respuesta de éxito |
| `@Failure` | `@Failure 400 {object} Error` | Respuesta de error |
| `@Header` | `@Header 200 {string} Token` | Header de respuesta |
| `@Router` | `@Router /users/{id} [get]` | Ruta y método de la ruta |
| `@Security` | `@Security ApiKeyAuth` | Requisito de seguridad |
| `@Deprecated` | `@Deprecated` | Marcar como deprecated |
| `@x-visibility` | `@x-visibility public` | Separar docs públicas/privadas |
| `@x-<nombre>` | `@x-code-samples file.json` | Extensión personalizada |

**Sintaxis de Parámetro:**

```
@Param <nombre> <en> <tipo> <obligatorio> <descripción> [atributos]
```

- **nombre**: Nombre del parámetro
- **en**: `query`, `path`, `header`, `body`, `formData`
- **tipo**: Tipo de dato (string, int, bool, object, array, file)
- **obligatorio**: `true` o `false`
- **descripción**: Descripción (entre comillas si contiene espacios)
- **atributos**: Atributos de validación opcionales

**Ejemplos:**

```go
// Parámetro de ruta
// @Param id path int true "ID del Usuario" minimum(1) maximum(1000)

// Parámetro de query con validación
// @Param name query string false "Nombre del usuario" minLength(3) maxLength(50)

// Parámetro de query con enum
// @Param status query string false "Filtro de estado" Enums(active,inactive,pending)

// Array de query con formato de colección
// @Param tags query []string false "Tags" collectionFormat(multi)

// Parámetro de header
// @Param X-Request-ID header string true "ID de la Petición" format(uuid)

// Parámetro de body
// @Param user body User true "Objeto del usuario"

// Form data con archivo
// @Param avatar formData file true "Imagen del avatar"
```

**Sintaxis de Respuesta:**

```go
// Respuesta simple
// @Success 200 {object} User

// Respuesta con descripción
// @Success 201 {object} User "Usuario creado con éxito"

// Respuesta de array
// @Success 200 {array} User "Lista de usuarios"

// Respuesta primitiva
// @Success 200 {string} string "Mensaje de éxito"

// Respuesta genérica
// @Success 200 {object} Response{data=User} "Respuesta del usuario"

// Múltiples campos de datos
// @Success 200 {object} Response{data=User,meta=Metadata}
```

**Sintaxis de Header:**

```go
// Código de estado único
// @Header 200 {string} X-Request-ID "Identificador de la petición"

// Múltiples códigos de estado
// @Header 200,201 {string} Location "URL del recurso"

// Todas las respuestas
// @Header all {string} X-API-Version "Versión de la API"
```

### Tags de Struct

#### Tags Estándar

```go
type User struct {
    // Serialización JSON
    ID   int    `json:"id"`
    Name string `json:"name,omitempty"`  // omitempty = no obligatorio
    
    // Validación (Gin binding)
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=150"`
    
    // Validación (go-playground/validator)
    UUID  string `json:"uuid" validate:"required,uuid"`
    
    // Atributos OpenAPI
    Price  float64  `json:"price" minimum:"0" maximum:"9999.99"`
    Status string   `json:"status" enum:"active,inactive" default:"active"`
    SKU    string   `json:"sku" pattern:"^[A-Z]{3}-[0-9]{6}$"`
    Items  []string `json:"items" minLength:"1" maxLength:"100"`
    
    // Valor de ejemplo
    Bio string `json:"bio" example:"Desarrollador de software"`
    
    // Formato
    CreatedAt string `json:"created_at" format:"date-time"`
}
```

#### swaggertype - Override de Tipo

Convertir tipos personalizados a tipos OpenAPI:

```go
type Account struct {
    // Override sql.NullInt64 a integer
    ID sql.NullInt64 `json:"id" swaggertype:"integer"`
    
    // Tipo de tiempo personalizado a unix timestamp (integer)
    CreatedAt TimestampTime `json:"created_at" swaggertype:"primitive,integer"`
    
    // Array de bytes a string base64
    Certificate []byte `json:"cert" swaggertype:"string" format:"base64"`
    
    // Array de número personalizado
    Coeffs []big.Float `json:"coeffs" swaggertype:"array,number"`
    
    // Tipos personalizados anidados
    Metadata map[string]interface{} `json:"metadata" swaggertype:"object"`
}
```

**Formato:** `swaggertype:"[primitive,]<tipo>"`

- Para tipos primitivos: `swaggertype:"string"`, `swaggertype:"integer"`, `swaggertype:"number"`, `swaggertype:"boolean"`
- Para arrays: `swaggertype:"array,<tipo-elemento>"`
- Para objetos: `swaggertype:"object"`

#### swaggerignore - Ocultar Campos

Excluir campos de la documentación (todavía presente en el JSON):

```go
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    
    // Presente en el JSON, oculto en los docs
    Password string `json:"password" swaggerignore:"true"`
    
    // Campo interno, no en JSON o docs
    internal string `swaggerignore:"true"`
    
    // Dato sensible
    SSN string `json:"ssn" swaggerignore:"true"`
}
```

#### extensions - Extensiones Personalizadas

Agregar metadatos personalizados con prefijo `x-*`:

```go
type Product struct {
    // Indicador de clave primaria
    ID int `json:"id" extensions:"x-primary-key=true"`
    
    // Formato de moneda
    Price float64 `json:"price" extensions:"x-currency=USD,x-format=currency"`
    
    // Múltiples extensiones
    Name string `json:"name" extensions:"x-order=1,x-searchable=true,x-filterable=true"`
    
    // Extensión booleana
    Featured bool `json:"featured" extensions:"x-promoted=true"`
    
    // Extensión nullable
    Discount float64 `json:"discount" extensions:"x-nullable"`
}
```

OpenAPI Generado:

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

## Características OpenAPI 3.2.0

nexs-swag ofrece soporte completo para las características de OpenAPI 3.2.0, manteniendo total compatibilidad con versiones anteriores (OpenAPI 2.0, 3.0.x, 3.1.x).

### Método HTTP QUERY

OpenAPI 3.2.0 introduce el método HTTP `QUERY` para consultas seguras con cuerpo de petición:

```go
// @Summary      Búsqueda compleja de productos
// @Description  Buscar productos usando parámetros complejos en el cuerpo de la petición
// @Tags         productos
// @Accept       json
// @Produce      json
// @Param        filtros body ProductFilter true "Criterios de búsqueda"
// @Success      200 {array} Product
// @Router       /products/query [query]
func QueryProducts(c *gin.Context) {}
```

### SecurityScheme Deprecated

Marque esquemas de seguridad obsoletos con `@securityDefinitions.*.deprecated`:

```go
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @deprecated true
// @description ⚠️ Este método de autenticación será descontinuado. Use OAuth2 en su lugar.
```

Resultado en OpenAPI:
```yaml
securitySchemes:
  ApiKeyAuth:
    type: apiKey
    name: X-API-Key
    in: header
    deprecated: true
    description: ⚠️ Este método de autenticación será descontinuado. Use OAuth2 en su lugar.
```

### OAuth2 Metadata URL

Para descubrimiento automático de configuración OAuth2 vía `@securityDefinitions.*.oauth2metadataurl`:

```go
// @securityDefinitions.oauth2.application OAuth2Application
// @tokenUrl https://auth.example.com/token
// @oauth2metadataurl https://auth.example.com/.well-known/oauth-authorization-server
```

Resultado en OpenAPI:
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

Soporte para OAuth 2.0 Device Authorization Grant (RFC 8628) vía `@securityDefinitions.*.deviceAuthorization`:

```go
// @securityDefinitions.oauth2.deviceAuth OAuth2Device
// @deviceAuthorization https://auth.example.com/device https://auth.example.com/token device-code
// @scopes.tv:watch Ver canales de TV
// @scopes.tv:record Grabar contenido
```

Resultado en OpenAPI:
```yaml
securitySchemes:
  OAuth2Device:
    type: oauth2
    flows:
      urn:ietf:params:oauth:grant-type:device_code:
        deviceAuthorizationUrl: https://auth.example.com/device
        tokenUrl: https://auth.example.com/token
        scopes:
          tv:watch: Ver canales de TV
          tv:record: Grabar contenido
```

### Respuestas de Streaming

Para respuestas SSE (Server-Sent Events) o streaming, use `@Success {stream}`:

```go
// @Summary      Stream de eventos
// @Description  Recibe actualizaciones en tiempo real de eventos del sistema
// @Tags         eventos
// @Produce      text/event-stream
// @Success      200 {stream} SystemEvent "Stream de eventos en tiempo real"
// @Router       /events/stream [get]
func StreamEvents(c *gin.Context) {}
```

Resultado en OpenAPI:
```yaml
responses:
  '200':
    description: Stream de eventos en tiempo real
    content:
      text/event-stream:
        itemSchema:
          $ref: '#/components/schemas/SystemEvent'
```

### Webhooks

Documentar webhooks que su API envía a clientes vía `@webhook`:

```go
// @webhook      OrderCreated
// @Description  Webhook enviado cuando se crea una nueva orden
// @Tags         webhooks
// @Accept       json
// @Param        order body Order true "Datos de la orden creada"
// @Success      200 {object} WebhookResponse
func DocumentOrderWebhook() {}
```

### Callbacks

Para operaciones asíncronas con callbacks, use `@Callback`:

```go
// @Summary      Procesar pago asíncrono
// @Description  Inicia procesamiento de pago y llama URL de callback
// @Tags         pagos
// @Accept       json
// @Param        payment body PaymentRequest true "Datos del pago"
// @Success      202 {object} PaymentResponse
// @Callback     paymentStatus {$request.body#/callbackUrl} post PaymentStatusCallback
// @Router       /payments/async [post]
func ProcessAsyncPayment(c *gin.Context) {}
```

### Separación por Visibilidad (@x-visibility)

Genera documentación separada para APIs públicas y privadas desde una única base de código.

```go
// GetPublicUser retorna información pública del usuario
// @Summary      Obtener usuario (público)
// @Description  Retorna información del usuario para consumo público
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "ID del Usuario"
// @Success      200  {object}  UserPublic
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [get]
// @x-visibility public
func GetPublicUser(c *gin.Context) {
    c.JSON(200, UserPublic{ID: 1, Name: "Juan"})
}

// GetAdminUser retorna detalles completos incluyendo datos sensibles
// @Summary      Obtener usuario (admin)
// @Description  Retorna información completa del usuario para uso administrativo
// @Tags         admin
// @Produce      json
// @Param        id   path      int  true  "ID del Usuario"
// @Success      200  {object}  UserPrivate
// @Failure      404  {object}  ErrorResponse
// @Router       /admin/users/{id} [get]
// @x-visibility private
func GetAdminUser(c *gin.Context) {
    c.JSON(200, UserPrivate{
        ID:       1,
        Name:     "Juan",
        Email:    "juan@example.com",
        Password: "hashed",
        Role:     "admin",
    })
}
```

**Opciones de Visibilidad:**
- `@x-visibility public` - Endpoint aparece solo en `openapi_public.json` o `swagger_public.json`
- `@x-visibility private` - Endpoint aparece solo en `openapi_private.json` o `swagger_private.json`
- Sin anotación - Endpoint aparece en **ambas** especificaciones (endpoint compartido)

**Archivos Generados:**
```
docs/
├── openapi_public.json    # Especificación API pública (OpenAPI 3.x)
├── openapi_private.json   # Especificación API privada/admin (OpenAPI 3.x)
├── swagger_public.json    # Especificación API pública (Swagger 2.0)
├── swagger_private.json   # Especificación API privada (Swagger 2.0)
├── openapi_public.yaml
├── openapi_private.yaml
├── docs_public.go
└── docs_private.go
```

**Filtrado de Schemas:**

Los schemas se filtran automáticamente basado en el uso:
- Spec pública incluye solo schemas referenciados por endpoints públicos
- Spec privada incluye solo schemas referenciados por endpoints privados
- Schemas compartidos (como `ErrorResponse`) aparecen donde sea necesario
- Dependencias recursivas de schemas se recopilan automáticamente

**Casos de Uso:**
- Separar documentación de API estable pública de APIs experimentales privadas
- Ocultar endpoints administrativos internos de la documentación pública
- Generar diferentes SDKs de cliente para APIs públicas vs privadas
- Alojar documentación separada para diferentes audiencias
- Distinguir entre APIs externas y APIs inter-servicios

**Compatibilidad:**
- ✅ Swagger 2.0
- ✅ OpenAPI 3.0.x
- ✅ OpenAPI 3.1.x
- ✅ OpenAPI 3.2.0

Para un ejemplo completo, vea [examples/26-x-visibility/](examples/26-x-visibility/) (OpenAPI 3.x) y [examples/27-x-visibility-v2/](examples/27-x-visibility-v2/) (Swagger 2.0).

### Migración 3.1.x → 3.2.0

nexs-swag detecta automáticamente la versión OpenAPI. Para activar características 3.2.0:

1. **No requiere cambios** - las características se activan al usar las anotaciones
2. **Compatible** - las anotaciones antiguas continúan funcionando
3. **Progresivo** - agregue características 3.2.0 gradualmente

**Avisos de depreciación** aparecen automáticamente si usa:
- `@securityDefinitions.*.deprecated true` - muestra badge de descontinuación
- Esquemas obsoletos sin migración - sugerencia para actualizar

## Ejemplos

nexs-swag incluye 21 ejemplos completos y ejecutables. Cada ejemplo demuestra características específicas e incluye un README y script de ejecución.

### Ejemplos Básicos

| Ejemplo | Descripción | Características Principales |
|---------|-------------|----------------------------|
| [01-basic](examples/01-basic) | Uso básico | Configuración mínima, API simple |
| [02-formats](examples/02-formats) | Formatos de salida | Salida JSON, YAML, Go |
| [03-general-info](examples/03-general-info) | Información general de la API | Metadatos completos de la API |

### Características Avanzadas

| Ejemplo | Descripción | Características Principales |
|---------|-------------|----------------------------|
| [04-property-strategy](examples/04-property-strategy) | Estrategias de nomenclatura | Snake_case, camelCase, PascalCase |
| [05-required-default](examples/05-required-default) | Obligatorio por defecto | Auto-require todos los campos |
| [06-exclude](examples/06-exclude) | Excluir directorios | Filtrar rutas no deseadas |
| [07-tags-filter](examples/07-tags-filter) | Filtrado por tag | Generar subconjunto de APIs |
| [08-parse-internal](examples/08-parse-internal) | Paquetes internos | Analizar directorio internal/ |
| [09-parse-dependency](examples/09-parse-dependency) | Dependencias | Analizar paquetes vendor/go.mod |
| [10-dependency-level](examples/10-dependency-level) | Profundidad de dependencia | Controlar nivel de análisis (0-3) |
| [11-parse-golist](examples/11-parse-golist) | Análisis de go list | Usar `go list` para descubrimiento |

### Características de Documentación

| Ejemplo | Descripción | Características Principales |
|---------|-------------|----------------------------|
| [12-markdown-files](examples/12-markdown-files) | Descripciones en Markdown | Cargar docs de archivos .md |
| [13-code-examples](examples/13-code-examples) | Muestras de código | Ejemplos en múltiples lenguajes |
| [14-overrides-file](examples/14-overrides-file) | Overrides de tipo | Configuración de archivo .swaggo |
| [15-generated-time](examples/15-generated-time) | Marca de tiempo de generación | Agregar fecha de generación |
| [16-instance-name](examples/16-instance-name) | Múltiples instancias | Conjuntos de documentación nombrados |
| [17-template-delims](examples/17-template-delims) | Delimitadores personalizados | Evitar conflictos de plantilla |

### Validación y Estructura

| Ejemplo | Descripción | Características Principales |
|---------|-------------|----------------------------|
| [18-collection-format](examples/18-collection-format) | Formatos de array | CSV, multi, pipes, SSV, TSV |
| [19-parse-func-body](examples/19-parse-func-body) | Cuerpos de función | Analizar anotaciones inline |
| [20-fmt-command](examples/20-fmt-command) | Comando de formato | Auto-formatear comentarios |
| [21-struct-tags](examples/21-struct-tags) | Todas las tags de struct | Referencia completa de tags |
| [22-openapi-v2](examples/22-openapi-v2) | Versionado OpenAPI | Swagger 2.0 & OpenAPI 3.1.0 |
| [23-recursive-parsing](examples/23-recursive-parsing) | Análisis recursivo | parseInternal, exclude, parseDependency |

### Ejecutando Ejemplos

Cada ejemplo incluye un script `run.sh`:

```bash
cd examples/01-basic
./run.sh
```

O manualmente (OpenAPI 3.1.0):

```bash
cd examples/01-basic
nexs-swag init -d . -o ./docs
cat docs/openapi.json
```

O generar Swagger 2.0:

```bash
cd examples/01-basic
nexs-swag init -d . -o ./docs --openapi-version 2.0
cat docs/swagger.json
```

### Ejemplo: API CRUD Completa

Consulta [examples/03-general-info](examples/03-general-info) para una API CRUD completa con:
- Múltiples endpoints (GET, POST, PUT, DELETE)
- Modelos de request/response
- Reglas de validación
- Respuestas de error
- Esquemas de seguridad
- Headers de respuesta

## Calidad y Pruebas

### Cobertura de Pruebas

```bash
$ go test ./pkg/... -cover
```

| Paquete | Cobertura | Pruebas |
|---------|-----------|----------|
| pkg/converter | 85.1% | 16 pruebas (con roundtrip) |
| pkg/format | 95.1% | 15 pruebas |
| pkg/generator/v2 | 80.3% | 12 pruebas |
| pkg/generator/v3 | 83.3% | 8 pruebas |
| pkg/openapi/v2 | 92.0% | 12 pruebas |
| pkg/openapi/v3 | 88.9% | 10 pruebas |
| pkg/parser | 84.6% | 195 pruebas |
| **General** | **80.1%** | **320+ pruebas** |

### Métricas de Calidad

- ✅ **0 avisos de linter** (golangci-lint con 20+ linters)
- ✅ **0 condiciones de carrera** (probado con flag `-race`)
- ✅ **22 pruebas de integración** (ejemplos ejecutables)
- ✅ **~8.500 líneas de código de prueba**
- ✅ **Listo para producción** (mantenido activamente)
- ✅ **100% compatible con swaggo/swag**
- ✅ **Soporte a múltiples versiones** (OpenAPI 3.1.0 + Swagger 2.0)

### Ejecutando Pruebas

```bash
# Pruebas unitarias
go test ./pkg/... -v

# Con cobertura
go test ./pkg/... -cover

# Con detección de race condition
go test ./pkg/... -race

# Paquete específico
go test ./pkg/parser -v

# Ejecutar ejemplos
cd examples && for d in */; do cd "$d" && ./run.sh && cd ..; done
```

## Compatibilidad con swaggo/swag

nexs-swag está diseñado como un **sustituto directo** para swaggo/swag con características mejoradas.

### Migración desde swaggo/swag

**¡Ningún cambio necesario!** Simplemente reemplaza el binario:

```bash
# En lugar de
go install github.com/swaggo/swag/cmd/swag@latest

# Usa
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest

# Los mismos comandos funcionan
nexs-swag init
nexs-swag fmt
```

### Tabla de Compatibilidad

| Característica | swaggo/swag | nexs-swag | Notas |
|----------------|-------------|-----------|-------|
| Versión OpenAPI | 2.0 | 3.1.0 | Retrocompatible |
| Todas las anotaciones | ✅ | ✅ | 100% compatible |
| Tags de struct | ✅ | ✅ | swaggertype, swaggerignore, extensions |
| Flags CLI | ✅ | ✅ | Todos los 28 flags soportados |
| Archivo .swaggo | ✅ | ✅ | Overrides de tipo |
| Markdown | ✅ | ✅ | Descripciones basadas en archivo |
| Ejemplos de código | ✅ | ✅ | Muestras en múltiples lenguajes |
| Webhooks | ❌ | ✅ | Característica OpenAPI 3.1 |
| JSON Schema 2020-12 | ❌ | ✅ | Schema moderno |
| Headers de respuesta | Limitado | ✅ | Soporte completo |
| Cobertura de pruebas | ~70% | 80.1% | Mayor calidad |
| Versión Go | 1.19+ | 1.23+ | Características Go modernas |

### ¿Qué es Diferente?

**Mejorado (retrocompatible):**
- Salida OpenAPI 3.1.0 (vs 2.0)
- Mejor manejo de nullable
- Más atributos de validación
- Mensajes de error mejorados
- Mejor cobertura de pruebas

**Misma API:**
- Todos los flags de línea de comandos
- Todas las anotaciones
- Todas las tags de struct
- Estructura generada de docs.go
- Integración con Swagger UI

## Acerca del Proyecto

### Estadísticas del Proyecto

- **Líneas de Código:** ~5.200 (pkg/ excluyendo pruebas)
- **Código de Prueba:** ~8.500 líneas
- **Archivos Go:** 42 archivos de implementación
- **Archivos de Prueba:** 29 archivos de prueba
- **Paquetes:** 9 (converter, format, generator, generator/v2, generator/v3, openapi, openapi/v2, openapi/v3, parser)
- **Ejemplos:** 23 ejemplos completos
- **Cobertura de Pruebas:** 80.1%
- **Versiones OpenAPI:** 4 (Swagger 2.0, OpenAPI 3.0.x, 3.1.x, 3.2.0)
- **Dependencias:** 3 dependencias directas
  - urfave/cli/v2 (framework CLI)
  - golang.org/x/tools (análisis AST Go)
  - gopkg.in/yaml.v3 (soporte YAML)

### Estructura del Proyecto

```
nexs-swag/
├── cmd/
│   └── nexs-swag/          # Punto de entrada CLI
├── pkg/
│   ├── converter/          # Conversión de versión (v3 ↔ v2)
│   ├── format/             # Formateo de código
│   ├── generator/          # Generación OpenAPI
│   │   ├── v2/             # Generador Swagger 2.0
│   │   └── v3/             # Generador OpenAPI 3.x
│   ├── openapi/            # Modelos OpenAPI
│   │   ├── v2/             # Modelos Swagger 2.0
│   │   └── v3/             # Modelos OpenAPI 3.x
│   └── parser/             # Análisis de código Go (AST)
├── examples/               # 22 ejemplos
│   ├── 01-basic/
│   ├── 02-formats/
│   └── ...
├── docs/                   # Documentación del proyecto
├── README.md               # Versión en inglés
├── README_pt.md            # Versión en portugués
├── README_es.md            # Este archivo
└── LICENSE                 # Licencia MIT
```

### Inspiración y Créditos

Este proyecto fue inspirado por [swaggo/swag](https://github.com/swaggo/swag) y construido para extender sus capacidades con soporte completo a OpenAPI 3.1.0, manteniendo 100% de compatibilidad retroactiva.

**Créditos:**
- [swaggo/swag](https://github.com/swaggo/swag) - Generador Swagger 2.0 original
- [OpenAPI Initiative](https://www.openapis.org/) - Especificación OpenAPI
- [Go Team](https://go.dev/) - Lenguaje y herramientas increíbles
- Todos los contribuyentes y la comunidad Go

## Contribuyendo

¡Las contribuciones son bienvenidas! Por favor, sigue estas directrices:

### Cómo Contribuir

1. **Fork** el repositorio
2. **Crea** una rama de característica (`git checkout -b feature/caracteristica-increible`)
3. **Haz** tus cambios
4. **Agrega** pruebas para nueva funcionalidad
5. **Ejecuta** las pruebas (`go test ./...`)
6. **Ejecuta** el linter (`golangci-lint run`)
7. **Commit** tus cambios (`git commit -m 'Agrega característica increíble'`)
8. **Push** a la rama (`git push origin feature/caracteristica-increible`)
9. **Abre** un Pull Request

### Configuración de Desarrollo

```bash
# Clonar repositorio
git clone https://github.com/fsvxavier/nexs-swag.git
cd nexs-swag

# Instalar dependencias
go mod download

# Ejecutar pruebas
go test ./... -v

# Ejecutar linter
golangci-lint run

# Build
go build -o nexs-swag ./cmd/nexs-swag
```

### Reportando Issues

Por favor incluye:
- Versión de Go (`go version`)
- Versión de nexs-swag (`nexs-swag --version`)
- Ejemplo reproducible mínimo
- Comportamiento esperado vs real

### Solicitudes de Características

Abre una issue con:
- Descripción clara de la característica
- Caso de uso y beneficios
- Implementación propuesta (si existe)

## Licencia

Este proyecto está licenciado bajo la **Licencia MIT** - consulta el archivo [LICENSE](LICENSE) para detalles.

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

## Soporte y Comunidad

- **Issues:** [GitHub Issues](https://github.com/fsvxavier/nexs-swag/issues)
- **Discusiones:** [GitHub Discussions](https://github.com/fsvxavier/nexs-swag/discussions)
- **Documentación:** [Wiki](https://github.com/fsvxavier/nexs-swag/wiki)
- **Ejemplos:** [examples/](examples/)

---

**Hecho con ❤️ para la comunidad Go**

[⬆ Volver arriba](#nexs-swag)
