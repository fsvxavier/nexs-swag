# Example 24: Include Types Simple

This example demonstrates the basic usage of the `--includeTypes` flag to filter which Go type categories are included in the OpenAPI specification.

## Features Demonstrated

1. **Struct Filtering**: Use `--includeTypes="struct"` to include only struct definitions
2. **Selective Parsing**: Only structs referenced in API operations are included
3. **Multi-file Support**: Structs and interfaces are in separate files
   - `models/product.go` - Contains struct definitions
   - `interfaces/repository.go` - Contains interface definitions
   - `main.go` - Contains API handlers with annotations
4. **swaggertype Tag**: Convert complex types (like `time.Time`) to simple OpenAPI types
5. **format Tag**: Specify OpenAPI format for better client generation

## File Structure

- `main.go` - API handlers with OpenAPI annotations
- `models/product.go` - Struct definitions (Product, ProductSummary, UnusedModel)
- `interfaces/repository.go` - Interface definitions (ProductRepository, UnusedInterface)

## Type Categories

The example includes:
- **Structs** (in `models/product.go`):
  - `models.Product` - Referenced in `@Success` annotation ✓ (included)
  - `models.ProductSummary` - Referenced in `@Success` annotation ✓ (included)
  - `models.UnusedModel` - Not referenced anywhere ✗ (excluded)
- **Interfaces** (in `interfaces/repository.go`):
  - `interfaces.ProductRepository` - Not referenced in annotations ✗ (excluded)
  - `interfaces.UnusedInterface` - Not referenced anywhere ✗ (excluded)

## Usage Examples

### 1. Include All Types (Default)
```bash
nexs-swag init -g main.go -o ./docs
# or
nexs-swag init -g main.go -o ./docs --includeTypes="all"
```
Result: Both `Product` and `ProductSummary` are included (but not `UnusedModel` since it's not referenced)

### 2. Include Only Structs
```bash
nexs-swag init -g main.go -o ./docs --includeTypes="struct"
```
Result: Both `Product` and `ProductSummary` are included as struct schemas

### 3. Short Form
```bash
nexs-swag init -g main.go -o ./docs -it="struct"
```
Result: Same as above using the short alias `-it`

## swaggertype and format Tags

### swaggertype Tag
Converts complex Go types to simple OpenAPI types:
```go
CreatedAt time.Time `swaggertype:"string" format:"date-time"`
```
- Go type: `time.Time`
- OpenAPI type: `string` with format `date-time`
- Prevents transitive parsing of `time.Time` struct

### format Tag
Adds OpenAPI format specification:
```go
Price float64 `format:"decimal"`
```
- Helps code generators create appropriate client code
- Common formats: `int32`, `int64`, `float`, `double`, `date`, `date-time`, `password`, `byte`, `binary`

```json
{
  "components": {
    "schemas": {
      "models.Product": {
        "type": "object",
        "properties": {
          "id": {"type": "integer", "example": 1},
          "name": {"type": "string", "example": "Laptop"},
          "price": {"type": "number", "format": "decimal", "example": 999.99},
          "created_at": {"type": "string", "format": "date-time", "example": "2025-12-16T10:00:00Z"},
          "category": {"type": "string", "example": "Electronics"},
          "is_available": {"type": "boolean", "example": true}
        }
      },
      "models.ProductSummary": {
        "type": "object",
        "properties": {
          "id": {"type": "integer", "example": 1},
          "name": {"type": "string", "example": "Laptop"}
        }
      }
    }
  }
}
```

Note: `models.UnusedModel` is NOT included because it's never referenced in any API operation.
}
```

Note: `UnusedModel` is NOT included because it's never referenced in any API operation.

## Running the Example

```bash
# Generate documentation
./run.sh

# Check the generated files
cat docs/openapi.json
```

## Key Takeaways

1. **Selective Parsing**: Only types referenced in operations are included (even with `--includeTypes="all"`)
2. **Type Filtering**: `--includeTypes` adds an additional filter on top of selective parsing
3. **Multi-file Support**: Definitions can be in separate files; the parser handles cross-file references
4. **Package Qualification**: Schemas are named with package prefix (e.g., `models.Product`)
5. **swaggertype**: Prevents deep parsing of complex types by converting them to primitives
6. **format**: Enhances API documentation with OpenAPI format specifications
