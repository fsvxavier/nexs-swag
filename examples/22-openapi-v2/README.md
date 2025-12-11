# Example 22 - OpenAPI 2.0 / Swagger 2.0

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

This example demonstrates how to generate both OpenAPI 3.1.0 and Swagger 2.0 specifications from the same Go code.

## Overview

A simple Product API with CRUD operations demonstrating:
- Swagger 2.0 generation (`--openapi-version 2.0`)
- OpenAPI 3.1.0 generation (default)
- Conversion warnings and compatibility

## Running the Example

### Generate OpenAPI 3.1.0 (default)

```bash
nexs-swag init -g main.go -o ./docs
```

### Generate Swagger 2.0

```bash
nexs-swag init -g main.go -o ./docs --openapi-version 2.0
```

### Generate Both Versions

```bash
# Generate OpenAPI 3.1.0
nexs-swag init -g main.go -o ./docs/v3

# Generate Swagger 2.0
nexs-swag init -g main.go -o ./docs/v2 --openapi-version 2.0
```

## Quick Start

```bash
# Install nexs-swag
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest

# Generate documentation
./run.sh

# Run the server
go run main.go
```

## API Endpoints

- `GET /api/v1/products` - List all products
- `POST /api/v1/products` - Create a new product
- `GET /api/v1/products/{id}` - Get a product by ID
- `PUT /api/v1/products/{id}` - Update a product
- `DELETE /api/v1/products/{id}` - Delete a product

## Generated Files

### OpenAPI 3.1.0
- `docs/v3/openapi.json` - JSON specification
- `docs/v3/openapi.yaml` - YAML specification
- `docs/v3/docs.go` - Go embedded specification

### Swagger 2.0
- `docs/v2/swagger.json` - JSON specification
- `docs/v2/swagger.yaml` - YAML specification
- `docs/v2/docs.go` - Go embedded specification

## Key Differences: OpenAPI 3.1.0 vs Swagger 2.0

| Feature | OpenAPI 3.1.0 | Swagger 2.0 |
|---------|---------------|-------------|
| Servers | `servers` array with URLs | `host`, `basePath`, `schemes` |
| Request Body | `requestBody` object | `body` parameter |
| Components | `components` section | `definitions`, `parameters`, `responses` |
| Media Types | `content` with MIME types | `consumes`, `produces` arrays |
| Nullable | `type: ["string", "null"]` | `x-nullable: true` extension |
| Webhooks | Supported | Not supported |
| JSON Schema | Draft 2020-12 | Draft 4 |

## Conversion Notes

When converting from OpenAPI 3.1.0 to Swagger 2.0, nexs-swag will:
- ✅ Convert `servers[0]` to `host`/`basePath`/`schemes`
- ✅ Convert `requestBody` to `body` parameter
- ✅ Convert `components` to `definitions`/`parameters`/`responses`
- ✅ Extract `consumes`/`produces` from operations
- ⚠️ Warn about unsupported features (webhooks, JSON Schema 2020-12, etc.)
- ⚠️ Ignore features not available in Swagger 2.0

## Testing

```bash
# Test GET
curl http://localhost:8080/api/v1/products

# Test POST
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Monitor","description":"4K Monitor","price":399.99,"stock":30}'

# Test GET by ID
curl http://localhost:8080/api/v1/products/1

# Test PUT
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Gaming Laptop","description":"High-end gaming laptop","price":1499.99,"stock":25}'

# Test DELETE
curl -X DELETE http://localhost:8080/api/v1/products/1
```

## View Documentation

### Swagger UI (Swagger 2.0)
```bash
# Use Docker to serve Swagger UI
docker run -p 8081:8080 -e SWAGGER_JSON=/docs/swagger.json \
  -v $(pwd)/docs/v2:/docs swaggerapi/swagger-ui
```

Open: http://localhost:8081

### Swagger Editor (supports both versions)
```bash
# Use Docker to serve Swagger Editor
docker run -p 8082:8080 swaggerapi/swagger-editor
```

Open: http://localhost:8082 and load your JSON file

## Advanced Usage

### Generate multiple formats

```bash
nexs-swag init -g main.go -o ./docs --format json,yaml,go --openapi-version 2.0
```

### With custom instance name

```bash
nexs-swag init -g main.go -o ./docs --instanceName api --openapi-version 2.0
```

### With timestamp

```bash
nexs-swag init -g main.go -o ./docs --generatedTime --openapi-version 2.0
```

## References

- [OpenAPI 3.1.0 Specification](https://spec.openapis.org/oas/v3.1.0)
- [Swagger 2.0 Specification](https://swagger.io/specification/v2/)
- [nexs-swag Documentation](../../README.md)
