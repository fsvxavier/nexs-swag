# Examples - nexs-swag

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Este directorio contiene ejemplos de uso para cada flag y funcionalidad de nexs-swag.

## Prerrequisitos

Instale nexs-swag antes de ejecutar los ejemplos:

```bash
# Desde el directorio raíz del proyecto
cd ..
go install ./cmd/nexs-swag

# O use el script de instalación
./install.sh

# Verificar instalación
nexs-swag --version
```

## Estructura

Cada subdirectorio contiene un ejemplo específico con:
- `main.go` - Código Go con anotaciones Swagger
- `README.md` - Instrucciones detalladas de uso (🌍 Disponible en 3 idiomas)
- `run.sh` - Script para ejecutar el ejemplo

## Lista de Ejemplos

### Básicos (01-08)
- [01-basic](./01-basic) - Uso básico con `--dir` y `--output`
- [02-formats](./02-formats) - Múltiples formatos con `--format`
- [03-general-info](./03-general-info) - Archivo específico con `--generalInfo`
- [04-property-strategy](./04-property-strategy) - `--propertyStrategy` (snake_case, camelCase, PascalCase)
- [05-required-default](./05-required-default) - `--requiredByDefault`
- [06-exclude](./06-exclude) - `--exclude` para excluir directorios
- [07-tags-filter](./07-tags-filter) - `--tags` para filtrar por tags
- [08-parse-internal](./08-parse-internal) - `--parseInternal`

### Dependencias (09-11)
- [09-parse-dependency](./09-parse-dependency) - `--parseDependency`
- [10-dependency-level](./10-dependency-level) - `--parseDependencyLevel` (0-3)
- [11-parse-golist](./11-parse-golist) - `--parseGoList`

### Contenido Externo (12-14)
- [12-markdown-files](./12-markdown-files) - `--markdownFiles`
- [13-code-examples](./13-code-examples) - `--codeExampleFilesDir`
- [14-overrides-file](./14-overrides-file) - `--overridesFile`

### Configuración (15-18)
- [15-generated-time](./15-generated-time) - `--generatedTime`
- [16-instance-name](./16-instance-name) - `--instanceName`
- [17-template-delims](./17-template-delims) - `--templateDelims`
- [18-collection-format](./18-collection-format) - `--collectionFormat`

### Avanzados (19-22)
- [19-parse-func-body](./19-parse-func-body) - `--parseFuncBody`
- [20-fmt-command](./20-fmt-command) - Comando `fmt`
- [21-struct-tags](./21-struct-tags) - swaggertype, swaggerignore, extensions
- [22-openapi-v2](./22-openapi-v2) - `--openapi-version` (Swagger 2.0 / OpenAPI 3.1.0)

## Cómo Usar

### Ejecutar un ejemplo específico

```bash
cd 01-basic
./run.sh
```

### Ejecutar manualmente

```bash
cd 01-basic
nexs-swag init --dir . --output ./docs
```

### Ejecutar todos los ejemplos

```bash
for dir in */; do
    echo "=== Ejecutando $dir ==="
    cd "$dir"
    ./run.sh
    cd ..
    echo ""
done
```

## Estructura de Cada Ejemplo

```
XX-nombre-ejemplo/
├── main.go          # Servidor HTTP con anotaciones
├── run.sh           # Script de demostración
└── README.md        # Documentación completa (🌍 3 idiomas)
```

## Consejos

### Ver documentación generada

```bash
# OpenAPI 3.1.0 (por defecto)
cat docs/openapi.json | jq
cat docs/openapi.yaml

# Swagger 2.0 (si se generó con --openapi-version 2.0)
cat docs/swagger.json | jq
cat docs/swagger.yaml

# Docs Go
cat docs/docs.go
```

### Servir con Swagger UI

```bash
# Instalar swagger ui
docker run -p 8080:8080 \
  -e SWAGGER_JSON=/docs/openapi.json \
  -v $(pwd)/docs:/docs \
  swaggerapi/swagger-ui

# Acceder: http://localhost:8080
```

### Integrar en proyectos

```go
package main

import (
    "net/http"
    
    _ "myapp/docs"  // Importar docs generados
    
    httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
    http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Solución de Problemas

### Error "nexs-swag: command not found"

```bash
# Verificar instalación
which nexs-swag

# Si no está instalado
cd ..
go install ./cmd/nexs-swag

# Verificar si $GOPATH/bin está en PATH
echo $PATH | grep $(go env GOPATH)/bin

# Agregar a PATH si es necesario
export PATH=$PATH:$(go env GOPATH)/bin
```

### Error al generar documentación

```bash
# Verificar si el código compila
go build .

# Ejecutar con más detalles
nexs-swag init --dir . --output ./docs --debug
```

### Limpiar documentación anterior

```bash
rm -rf docs docs-*
```

## Recursos

- [Documentación Completa](../INSTALL.md)
- [swaggo/swag - Documentación Original](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
