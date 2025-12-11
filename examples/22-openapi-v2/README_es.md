# Ejemplo 22 - OpenAPI v2 (Swagger)

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra la generación de specs OpenAPI 2.0 (Swagger).

## Flag

```bash
--openAPIVersion 2
--oav 2
```

Default: `3` (OpenAPI 3.0)

## Concepto

nexs-swag puede generar specs en formato OpenAPI 2.0 (también conocido como Swagger 2.0) o OpenAPI 3.0.

## Uso

```bash
# OpenAPI 2.0
nexs-swag init --openAPIVersion 2

# OpenAPI 3.0 (default)
nexs-swag init --openAPIVersion 3
```

## Cómo Ejecutar

```bash
./run.sh
```

## Diferencias Principales

### 1. Info Object

#### OpenAPI 2.0
```yaml
swagger: "2.0"
info:
  title: My API
  version: 1.0.0
host: api.example.com
basePath: /v1
schemes:
  - https
```

#### OpenAPI 3.0
```yaml
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
servers:
  - url: https://api.example.com/v1
```

### 2. Security Definitions

#### OpenAPI 2.0
```yaml
securityDefinitions:
  bearerAuth:
    type: apiKey
    name: Authorization
    in: header
  oauth2:
    type: oauth2
    flow: accessCode
    authorizationUrl: https://auth.example.com/oauth/authorize
    tokenUrl: https://auth.example.com/oauth/token
    scopes:
      read: Read access
      write: Write access
```

#### OpenAPI 3.0
```yaml
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
    oauth2:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://auth.example.com/oauth/authorize
          tokenUrl: https://auth.example.com/oauth/token
          scopes:
            read: Read access
            write: Write access
```

### 3. Request Body

#### OpenAPI 2.0
```yaml
paths:
  /users:
    post:
      parameters:
        - in: body
          name: body
          schema:
            $ref: '#/definitions/User'
      responses:
        201:
          description: Created
```

#### OpenAPI 3.0
```yaml
paths:
  /users:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/User'
      responses:
        '201':
          description: Created
```

### 4. Definitions vs Components

#### OpenAPI 2.0
```yaml
definitions:
  User:
    type: object
    properties:
      id:
        type: integer
      name:
        type: string
```

#### OpenAPI 3.0
```yaml
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
```

### 5. Response Examples

#### OpenAPI 2.0
```yaml
responses:
  200:
    description: Success
    schema:
      $ref: '#/definitions/User'
    examples:
      application/json:
        id: 1
        name: John Doe
```

#### OpenAPI 3.0
```yaml
responses:
  '200':
    description: Success
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/User'
        example:
          id: 1
          name: John Doe
```

### 6. Multiple Media Types

#### OpenAPI 2.0
```yaml
consumes:
  - application/json
  - application/xml
produces:
  - application/json
  - application/xml
```

#### OpenAPI 3.0
```yaml
requestBody:
  content:
    application/json:
      schema: {...}
    application/xml:
      schema: {...}
responses:
  '200':
    content:
      application/json:
        schema: {...}
      application/xml:
        schema: {...}
```

## Annotations

### OpenAPI 2.0 Específicas

```go
// @swagger 2.0

// @host api.example.com
// @basePath /v1
// @schemes https http

// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.accessCode OAuth2
// @tokenUrl https://auth.example.com/oauth/token
// @authorizationurl https://auth.example.com/oauth/authorize
// @scope.read Read access
// @scope.write Write access
```

### OpenAPI 3.0 Específicas

```go
// @openapi 3.0.0

// @server.url https://api.example.com/v1
// @server.description Production server

// @securitySchemes.http BearerAuth
// @scheme bearer
// @bearerFormat JWT

// @securitySchemes.oauth2.authorizationCode OAuth2
// @authorizationUrl https://auth.example.com/oauth/authorize
// @tokenUrl https://auth.example.com/oauth/token
// @scope.read Read access
// @scope.write Write access
```

## Cuándo Usar Cada Versión

### Use OpenAPI 2.0 cuando:

#### 1. Legacy Systems
```bash
# Sistema legacy soporta solo Swagger 2.0
nexs-swag init --openAPIVersion 2
```

#### 2. Tooling Requirements
- Código generators antiguos
- Internal tools que requieren v2
- Clientes legacy

#### 3. Compatibilidad
```bash
# Generar ambas versiones
nexs-swag init --openAPIVersion 2 --output docs/v2
nexs-swag init --openAPIVersion 3 --output docs/v3
```

### Use OpenAPI 3.0 cuando:

#### 1. Nuevos Proyectos
```bash
# Default - OpenAPI 3.0
nexs-swag init
```

#### 2. Features Avanzados
- Multiple servers
- Callbacks
- Links
- anyOf/oneOf/not
- Discriminator

#### 3. Modern Tooling
- Swagger UI 3.x+
- ReDoc
- Stoplight
- Postman

## Migration 2.0 → 3.0

### Automated
```bash
# Use API Matic Transformer
curl -X POST \
  https://www.apimatic.io/api/transform \
  -d @swagger-v2.yaml

# Or swagger-converter
npm install -g swagger2openapi
swagger2openapi swagger-v2.yaml -o openapi-v3.yaml
```

### Manual Changes

#### 1. Version
```yaml
# Before
swagger: "2.0"

# After
openapi: 3.0.0
```

#### 2. Host/BasePath → Servers
```yaml
# Before
host: api.example.com
basePath: /v1
schemes: [https]

# After
servers:
  - url: https://api.example.com/v1
```

#### 3. Body Parameters → RequestBody
```yaml
# Before
parameters:
  - in: body
    name: user
    schema:
      $ref: '#/definitions/User'

# After
requestBody:
  content:
    application/json:
      schema:
        $ref: '#/components/schemas/User'
```

#### 4. Definitions → Components/Schemas
```yaml
# Before
definitions:
  User: {...}

# After
components:
  schemas:
    User: {...}
```

## Exemplo Completo

### main.go
```go
// @title E-commerce API
// @version 1.0.0
// @description API para gestión de e-commerce

// OpenAPI 2.0
// @host api.example.com
// @basePath /v1
// @schemes https

// OpenAPI 3.0
// @server.url https://api.example.com/v1
// @server.description Production

// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @Summary Create user
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User object"
// @Success 201 {object} User
// @Router /users [post]
// @Security ApiKeyAuth
func CreateUser(c *gin.Context) {}
```

### Generar Ambas Versiones
```bash
# V2
nexs-swag init --openAPIVersion 2 --output docs/v2

# V3
nexs-swag init --openAPIVersion 3 --output docs/v3
```

### Servir Ambas
```go
r.GET("/swagger/v2/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.URL("/docs/v2/swagger.json"),
))

r.GET("/swagger/v3/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.URL("/docs/v3/swagger.json"),
))
```

## Validation

### OpenAPI 2.0
```bash
# Swagger Editor
# https://editor.swagger.io/

# CLI
npm install -g swagger-cli
swagger-cli validate docs/swagger.yaml
```

### OpenAPI 3.0
```bash
# Swagger Editor (supports both)
# https://editor.swagger.io/

# CLI
npm install -g @apidevtools/swagger-cli
swagger-cli validate docs/openapi.yaml
```

## Recomendaciones

**Para Nuevos Proyectos:**
```bash
# Use OpenAPI 3.0 (default)
nexs-swag init
```

**Para Proyectos Legacy:**
```bash
# Use OpenAPI 2.0 si necesario
nexs-swag init --openAPIVersion 2

# Considere migrar a 3.0 gradualmente
```

**Para Máxima Compatibilidad:**
```bash
# Genere ambas versiones
make docs-v2 docs-v3
```

## Recursos

- [OpenAPI 2.0 Spec](https://swagger.io/specification/v2/)
- [OpenAPI 3.0 Spec](https://swagger.io/specification/)
- [Migration Guide](https://swagger.io/docs/specification/about/)
- [API Matic Transformer](https://www.apimatic.io/transformer/)
