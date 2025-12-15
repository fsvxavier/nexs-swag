# nexs-swag

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

[![Versión Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1.0-6BA539?style=flat&logo=openapiinitiative)](https://spec.openapis.org/oas/v3.1.0)
[![Swagger](https://img.shields.io/badge/Swagger-2.0-85EA2D?style=flat&logo=swagger)](https://swagger.io/specification/v2/)
[![Licencia](https://img.shields.io/badge/Licencia-MIT-blue.svg)](LICENSE)
[![Cobertura](https://img.shields.io/badge/Cobertura-86.1%25-brightgreen.svg)](/))
[![Ejemplos](https://img.shields.io/badge/Ejemplos-22-blue.svg)](examples/)

**Genere automáticamente documentación OpenAPI 3.1.0 o Swagger 2.0 a partir de anotaciones en código Go.**

nexs-swag convierte anotaciones Go a especificación OpenAPI 3.1.0 o Swagger 2.0. Está diseñado como una evolución de [swaggo/swag](https://github.com/swaggo/swag) con soporte completo para la especificación OpenAPI más reciente y compatibilidad total con Swagger 2.0.

## Índice

- [Descripción General](#descripción-general)
- [Primeros Pasos](#primeros-pasos)
  - [Instalación](#instalación)
  - [Inicio Rápido](#inicio-rápido)
- [Frameworks Web Soportados](#frameworks-web-soportados)
- [Cómo usar con Gin](#cómo-usar-con-gin)
- [Referencia CLI](#referencia-cli)
  - [Comando init](#comando-init)
  - [Comando fmt](#comando-fmt)
- [Estado de Implementación](#estado-de-implementación)
- [Formato de Comentarios Declarativos](#formato-de-comentarios-declarativos)
  - [Información General de la API](#información-general-de-la-api)
  - [Operación de API](#operación-de-api)
  - [Tags de Struct](#tags-de-struct)
- [Ejemplos](#ejemplos)
- [Calidad y Pruebas](#calidad-y-pruebas)
- [Compatibilidad con swaggo/swag](#compatibilidad-con-swaggoswag)
- [Sobre el Proyecto](#sobre-el-proyecto)
- [Contribuyendo](#contribuyendo)
- [Licencia](#licencia)

## Descripción General

### Características Principales

- ✅ **100% compatible con swaggo/swag** - Reemplazo directo con todas las anotaciones y tags
- ✅ **Soporte dual de versiones** - Genere OpenAPI 3.1.0 **o** Swagger 2.0 desde las mismas anotaciones
- ✅ **OpenAPI 3.1.0** - Soporte completo para JSON Schema 2020-12, webhooks y características modernas
- ✅ **Swagger 2.0** - Compatibilidad total con sistemas legacy
- ✅ **Conversión automática** - Conversión interna entre formatos con avisos para incompatibilidades
- ✅ **20+ atributos de validación** - minimum, maximum, pattern, enum, format y más
- ✅ **Validación de frameworks** - Soporte nativo para Gin (binding) y go-playground/validator
- ✅ **Headers de respuesta** - Documentación completa de headers
- ✅ **Múltiples tipos de contenido** - JSON, XML, YAML, CSV, PDF y tipos MIME personalizados
- ✅ **Extensiones personalizadas** - Soporte completo para x-*
- ✅ **86.1% de cobertura de pruebas** - Listo para producción con suite de pruebas completa
- ✅ **22 ejemplos funcionales** - Aprenda con ejemplos completos y ejecutables

### ¿Por qué nexs-swag?

| Característica | swaggo/swag | nexs-swag |
|----------------|-------------|-----------||
| OpenAPI 3.1.0 | ❌ | ✅ |
| Swagger 2.0 | ✅ | ✅ |
| Generación Dual | ❌ | ✅ (ambos del mismo código) |
| JSON Schema | Draft 4 | Draft 4 + 2020-12 |
| Webhooks | ❌ | ✅ (OpenAPI 3.1.0) |
| Headers de Respuesta | Limitado | Soporte Completo |
| Soporte Nullable | `x-nullable` | Nativo + `x-nullable` |
| Cobertura de Pruebas | ~70% | 86.1% |
| Ejemplos | ~10 | 22 |
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

#### 1. Agregar Anotaciones de API

Agregue anotaciones generales de API a su `main.go`:

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
    // Configuración de su aplicación
    r.Run(":8080")
}

// User representa un usuario del sistema
type User struct {
    // ID del usuario (sql.NullInt64 → integer en OpenAPI)
    ID sql.NullInt64 `json:"id" swaggertype:"integer" extensions:"x-primary-key=true"`
    
    // Nombre completo (3-100 caracteres requerido)
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
// @Header       201   {string}  X-Request-ID  "Identificador de la solicitud"
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

**OpenAPI 3.1.0 (por defecto):**

```bash
nexs-swag init
```

**Swagger 2.0:**

```bash
nexs-swag init --openapi-version 2.0
```

**Generar ambas versiones:**

```bash
nexs-swag init -o ./docs/v3 --openapi-version 3.1
nexs-swag init -o ./docs/v2 --openapi-version 2.0
```

O especifique los directorios:

```bash
nexs-swag init -d ./cmd/api -o ./docs
```

#### 3. Archivos Generados

**OpenAPI 3.1.0:**
Los siguientes archivos se crearán en su directorio de salida (por defecto: `./docs`):

- **`docs/openapi.json`** - Especificación OpenAPI 3.1.0 en formato JSON
- **`docs/openapi.yaml`** - Especificación OpenAPI 3.1.0 en formato YAML
- **`docs/docs.go`** - Archivo de documentación Go embebido

**Swagger 2.0:**
Cuando use `--openapi-version 2.0`, los archivos generados serán:

- **`docs/swagger.json`** - Especificación Swagger 2.0 en formato JSON
- **`docs/swagger.yaml`** - Especificación Swagger 2.0 en formato YAML
- **`docs/docs.go`** - Archivo de documentación Go embebido

#### 4. Integrar con Su Aplicación

Importe el paquete docs generado:

```go
import _ "su-modulo/docs"  // Importar docs generado

func main() {
    r := gin.Default()
    
    // Servir Swagger UI
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    r.Run(":8080")
}
```

¡Acceda a http://localhost:8080/swagger/index.html para ver su documentación API!

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

Ejemplo completo usando framework Gin. Encuentre el código completo en [examples/03-general-info](examples/03-general-info).

Consulte la [versión en inglés](README.md#how-to-use-with-gin) para detalles completos.

## Referencia CLI

### Comando init

Genera documentación OpenAPI a partir del código fuente.

```bash
nexs-swag init [opciones]
```

**Opciones Principales:**

- `--dir, -d` - Directorios para analizar (por defecto: `./`)
- `--output, -o` - Directorio de salida (por defecto: `./docs`)
- `--outputTypes, --ot` - Tipos de archivo de salida (por defecto: `go,json,yaml`)
- `--openapi-version, --ov` - Versión OpenAPI: `2.0`, `3.0`, `3.1` (por defecto: `3.1`)
- `--parseDependency, --pd` - Analizar dependencias (por defecto: `false`)
- `--parseInternal` - Analizar paquetes internos (por defecto: `false`)
- `--propertyStrategy, -p` - Estrategia de nomenclatura: `snakecase`, `camelcase`, `pascalcase`
- `--validate` - Validar especificación generada (por defecto: `true`)

> **⚠️ Importante: Sintaxis de Flags Booleanas**
>
> Las flags booleanas aceptan dos sintaxis válidas:
> - ✅ **Sin valor** (presencia = true): `--parseInternal`, `--pd`
> - ✅ **Con signo igual**: `--parseInternal=true`, `--pd=false`
> - ❌ **Incorrecto** (separado por espacio): `--parseInternal true`, `--pd true`
>
> La sintaxis separada por espacio no funciona porque el parser CLI trata la palabra después de la flag como un argumento posicional separado, no como el valor de la flag.

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

# Solo salida JSON
nexs-swag init --outputTypes json

# Nombres de propiedad en snake_case
nexs-swag init --propertyStrategy snakecase
```

### Comando fmt

Formatea comentarios swagger automáticamente.

```bash
nexs-swag fmt [opciones]
```

**Ejemplo:**

```bash
# Formatear directorio actual
nexs-swag fmt

# Formatear directorio específico
nexs-swag fmt -d ./internal/api
```

## Formato de Comentarios Declarativos

Para documentación completa de todas las anotaciones, parámetros, tags de struct y ejemplos, consulte la [versión en inglés](README.md#declarative-comments-format).

### Resumen Rápido

**Información General de la API:**
- `@title` - Título de la API (requerido)
- `@version` - Versión de la API (requerido)
- `@description` - Descripción de la API
- `@host` - Host de la API
- `@BasePath` - Ruta base
- `@securityDefinitions.*` - Definiciones de seguridad

**Operación de API:**
- `@Summary` - Resumen corto
- `@Description` - Descripción detallada
- `@Tags` - Tags de la operación
- `@Param` - Definición de parámetro
- `@Success` - Respuesta exitosa
- `@Failure` - Respuesta de error
- `@Router` - Ruta y método de la ruta

**Tags de Struct:**
- `json` - Serialización JSON
- `binding` - Validación Gin
- `validate` - Validación go-playground
- `swaggertype` - Override de tipo
- `swaggerignore` - Ocultar campo
- `extensions` - Extensiones personalizadas

## Ejemplos

nexs-swag incluye 22 ejemplos completos y ejecutables demostrando todas las funcionalidades, incluyendo generación de OpenAPI 3.1.0 y Swagger 2.0. Consulte la [sección de ejemplos](README.md#examples) en la versión en inglés para la lista completa.

### Ejecutando Ejemplos

Cada ejemplo incluye un script `run.sh`:

```bash
cd examples/01-basic
./run.sh
```

## Calidad y Pruebas

### Cobertura de Pruebas

| Paquete | Cobertura | Pruebas |
|---------|-----------|---------||
| pkg/converter | 92.3% | 13 pruebas |
| pkg/format | 95.1% | 15 pruebas |
| pkg/generator | 71.6% | 16 pruebas |
| pkg/generator/v2 | 88.4% | 12 pruebas |
| pkg/generator/v3 | 85.2% | 8 pruebas |
| pkg/openapi | 83.3% | 22 pruebas |
| pkg/openapi/v2 | 89.7% | 12 pruebas |
| pkg/openapi/v3 | 91.5% | 10 pruebas |
| pkg/parser | 82.1% | 192 pruebas |
| **General** | **87.9%** | **300+ pruebas** |

### Métricas de Calidad

- ✅ **0 advertencias de linter** (golangci-lint con 20+ linters)
- ✅ **0 condiciones de carrera** (probado con flag `-race`)
- ✅ **22 pruebas de integración** (ejemplos ejecutables)
- ✅ **~8.500 líneas de código de prueba**
- ✅ **Listo para producción** (mantenido activamente)
- ✅ **100% compatible con swaggo/swag**
- ✅ **Soporte dual de versiones** (OpenAPI 3.1.0 + Swagger 2.0)

## Compatibilidad con swaggo/swag

nexs-swag está diseñado como un **reemplazo directo** para swaggo/swag con características mejoradas.

### Migración desde swaggo/swag

**¡No se necesitan cambios!** Simplemente reemplace el binario:

```bash
# En lugar de
go install github.com/swaggo/swag/cmd/swag@latest

# Use
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest

# Los mismos comandos funcionan
nexs-swag init
nexs-swag fmt
```

## Sobre el Proyecto

### Estadísticas del Proyecto

- **Líneas de Código:** ~5.200 (pkg/ excluyendo pruebas)
- **Código de Prueba:** ~8.500 líneas
- **Archivos Go:** 42 archivos de implementación
- **Archivos de Prueba:** 29 archivos de prueba
- **Paquetes:** 9 (converter, format, generator, generator/v2, generator/v3, openapi, openapi/v2, openapi/v3, parser)
- **Ejemplos:** 22 ejemplos completos
- **Cobertura de Pruebas:** 87.9%
- **Versiones OpenAPI:** 2 (Swagger 2.0 + OpenAPI 3.1.0)
- **Dependencias:** 3 dependencias directas

### Inspiración y Créditos

Este proyecto fue inspirado por [swaggo/swag](https://github.com/swaggo/swag) y construido para extender sus capacidades con soporte completo para OpenAPI 3.1.0, manteniendo 100% de compatibilidad retroactiva.

## Contribuyendo

¡Las contribuciones son bienvenidas! Consulte la [versión en inglés](README.md#contributing) para directrices detalladas.

### Cómo Contribuir

1. **Fork** el repositorio
2. **Cree** una rama de característica (`git checkout -b feature/caracteristica-increible`)
3. **Haga** sus cambios
4. **Agregue** pruebas para nueva funcionalidad
5. **Ejecute** las pruebas (`go test ./...`)
6. **Ejecute** el linter (`golangci-lint run`)
7. **Commit** sus cambios (`git commit -m 'Agrega característica increíble'`)
8. **Push** a la rama (`git push origin feature/caracteristica-increible`)
9. **Abra** un Pull Request

## Licencia

Este proyecto está licenciado bajo la **Licencia MIT** - consulte el archivo [LICENSE](LICENSE) para detalles.

## Soporte y Comunidad

- **Issues:** [GitHub Issues](https://github.com/fsvxavier/nexs-swag/issues)
- **Discusiones:** [GitHub Discussions](https://github.com/fsvxavier/nexs-swag/discussions)
- **Documentación:** [Wiki](https://github.com/fsvxavier/nexs-swag/wiki)
- **Ejemplos:** [examples/](examples/)

---

**Hecho con ❤️ para la comunidad Go**

[⬆ Volver arriba](#nexs-swag)
