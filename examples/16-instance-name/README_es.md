# Ejemplo 16 - Instance Name

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de `--instanceName` para múltiples specs en una aplicación.

## Flag

```bash
--instanceName <name>
--in <name>
```

Default: `"swagger"`

## Concepto

Permite generar múltiples specs OpenAPI en la misma aplicación Go.

## Uso

### Default Instance
```bash
nexs-swag init
# Genera: docs/swagger.yaml, docs/swagger.json
```

```go
import _ "myapp/docs"

// @title Default API
func main() {
    r := gin.Default()
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

### Named Instance
```bash
nexs-swag init --instanceName admin
# Genera: docs/admin_swagger.yaml, docs/admin_swagger.json
```

```go
import _ "myapp/docs"

// @title Admin API
func main() {
    r := gin.Default()
    r.GET("/swagger/*any", ginSwagger.WrapHandler(
        swaggerFiles.Handler,
        ginSwagger.InstanceName("admin"),
    ))
}
```

## Cómo Ejecutar

```bash
./run.sh
```

## Casos de Uso

### 1. API Público vs Admin
```bash
# Public API
nexs-swag init --dir ./public --instanceName public

# Admin API  
nexs-swag init --dir ./admin --instanceName admin
```

```go
// Public routes
r.GET("/swagger/public/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.InstanceName("public"),
))

// Admin routes
r.GET("/swagger/admin/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.InstanceName("admin"),
))
```

### 2. API v1 vs v2
```bash
nexs-swag init --dir ./api/v1 --instanceName v1
nexs-swag init --dir ./api/v2 --instanceName v2
```

```go
r.GET("/swagger/v1/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.InstanceName("v1"),
))

r.GET("/swagger/v2/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.InstanceName("v2"),
))
```

### 3. API Interno vs Externo
```bash
nexs-swag init --dir ./internal --instanceName internal
nexs-swag init --dir ./external --instanceName external
```

```go
// Internal (autenticado)
internal.GET("/swagger/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.InstanceName("internal"),
))

// External (público)
external.GET("/swagger/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.InstanceName("external"),
))
```

## Estructura

```
myapp/
├── main.go
├── docs/
│   ├── swagger.yaml            # Default instance
│   ├── swagger.json
│   ├── admin_swagger.yaml      # Admin instance
│   ├── admin_swagger.json
│   ├── internal_swagger.yaml   # Internal instance
│   └── internal_swagger.json
├── public/
│   └── handlers.go
├── admin/
│   └── handlers.go
└── internal/
    └── handlers.go
```

## Código Generado

### docs/docs.go (Default)
```go
package docs

var SwaggerInfo = &swag.Spec{
    Version: "1.0",
    Host: "localhost:8080",
    BasePath: "/api",
    // ...
}

func init() {
    swag.Register(swag.Name, SwaggerInfo)
}
```

### docs/admin_docs.go (Admin Instance)
```go
package docs

var AdminSwaggerInfo = &swag.Spec{
    Version: "1.0",
    Host: "localhost:8080",
    BasePath: "/admin",
    // ...
}

func init() {
    swag.Register("admin", AdminSwaggerInfo)
}
```

## Makefile

```makefile
# Makefile
.PHONY: docs docs-public docs-admin docs-all

docs-public:
nexs-swag init --dir ./public --instanceName public

docs-admin:
nexs-swag init --dir ./admin --instanceName admin

docs-internal:
nexs-swag init --dir ./internal --instanceName internal

docs-all: docs-public docs-admin docs-internal

clean:
rm -rf docs/*swagger*
```

## CI/CD

```yaml
# .github/workflows/docs.yml
name: Generate Docs

on: [push]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Generate Public API docs
        run: nexs-swag init --dir ./public --instanceName public
      
      - name: Generate Admin API docs
        run: nexs-swag init --dir ./admin --instanceName admin
      
      - name: Commit docs
        run: |
          git add docs/
          git commit -m "Update API docs"
          git push
```

## Nombres Válidos

```bash
# ✅ Válidos
--instanceName api
--instanceName admin
--instanceName v1
--instanceName internal_api
--instanceName publicAPI

# ❌ Inválidos (caracteres especiales)
--instanceName "admin api"  # espacio
--instanceName admin-api    # hyphen en algunos contextos
```

## Configuración Diferente

### public.yaml
```yaml
# Configuración para API público
instanceName: public
dir: ./public
output: ./docs
parseDependency: true
```

### admin.yaml
```yaml
# Configuración para API admin
instanceName: admin
dir: ./admin
output: ./docs
parseDependency: true
parseInternal: true
```

```bash
nexs-swag init --config public.yaml
nexs-swag init --config admin.yaml
```

## Tips

### 1. Ambiente por Instance
```go
var specs = map[string]string{
    "public":   os.Getenv("PUBLIC_API_HOST"),
    "admin":    os.Getenv("ADMIN_API_HOST"),
    "internal": os.Getenv("INTERNAL_API_HOST"),
}

for name, host := range specs {
    docs.UpdateHost(name, host)
}
```

### 2. Middleware por Instance
```go
publicGroup := r.Group("/public")
publicGroup.Use(RateLimiter(1000)) // Liberal
publicGroup.GET("/swagger/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.InstanceName("public"),
))

adminGroup := r.Group("/admin")
adminGroup.Use(Auth(), RateLimiter(10000)) // Restrictivo
adminGroup.GET("/swagger/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.InstanceName("admin"),
))
```

### 3. Documentación Diferente
```go
// public/handlers.go
// @title Public API
// @description API para usuarios externos

// admin/handlers.go
// @title Admin API
// @description API para administradores (requiere autenticación)
```

## Recomendaciones

**Use instanceName cuando:**
- Múltiples APIs en una app
- Separación público/privado
- Versionamiento lado a lado
- Audiences diferentes

**NO use cuando:**
- Solo una API
- Versionamiento via path (/v1/, /v2/)
- Complejidad innecesaria
