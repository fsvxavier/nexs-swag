# X-Visibility with Swagger 2.0 Example

This example demonstrates the `@x-visibility` annotation feature with Swagger 2.0, allowing you to generate separate documentation files for public and private APIs.

## Overview

The `@x-visibility` annotation enables you to:
- Separate public and private endpoints into different Swagger files
- Automatically filter schemas (definitions) based on usage
- Maintain a single codebase with dual documentation output
- Works with Swagger 2.0 and OpenAPI 3.x

## How It Works

### Annotation Syntax

Add `@x-visibility` to your operation comments:

```go
// GetUser godoc
// @Summary      Get user (public)
// @Description  Get user details for public consumption
// @Tags         users
// @Success      200  {object}  UserPublic
// @Router       /users/{id} [get]
// @x-visibility public
func GetUser(c *gin.Context) {
    // handler implementation
}
```

### Visibility Options

- `@x-visibility public` - Endpoint appears only in `swagger_public.json`
- `@x-visibility private` - Endpoint appears only in `swagger_private.json`
- No annotation - Endpoint appears in **both** files (shared endpoint)

## Generated Files

When using `@x-visibility` with Swagger 2.0, nexs-swag generates:

```
docs/
├── swagger_public.json    # Public API specification
├── swagger_private.json   # Private API specification
├── swagger_public.yaml    # Public API (YAML)
├── swagger_private.yaml   # Private API (YAML)
├── docs_public.go         # Public API Go code
└── docs_private.go        # Private API Go code
```

## Running This Example

```bash
# Generate Swagger 2.0 documentation
nexs-swag init --output ./docs --ov 2.0

# Or use the run script
./run.sh

# Verify separation
jq '.paths | keys' docs/swagger_public.json
# Output: ["/users", "/users/{id}"]

jq '.paths | keys' docs/swagger_private.json
# Output: ["/admin/users/{id}", "/users"]

jq '.definitions | keys' docs/swagger_public.json
# Output: ["ErrorResponse", "UserPublic"]

jq '.definitions | keys' docs/swagger_private.json
# Output: ["ErrorResponse", "UserPrivate", "UserPublic"]
```

## Compatibility

The `@x-visibility` feature works with:
- ✅ Swagger 2.0 (this example)
- ✅ OpenAPI 3.0.x
- ✅ OpenAPI 3.1.x
- ✅ OpenAPI 3.2.0

Extensions are preserved during conversion between versions.

## Comparison with OpenAPI 3.x

| Feature | Swagger 2.0 | OpenAPI 3.x |
|---------|-------------|-------------|
| Public/Private Separation | ✅ | ✅ |
| Schema Filtering | ✅ (definitions) | ✅ (components.schemas) |
| Extension Support | ✅ x-visibility | ✅ x-visibility |
| File Names | swagger_*.json | openapi_*.json |

For OpenAPI 3.x version, see [example 26](../26-x-visibility/).

## Notes

- Schemas (definitions) are recursively collected including nested references
- Operations without `@x-visibility` appear in both specs
- Shared schemas (like `ErrorResponse`) are included where needed
- All other Swagger features work normally within each spec
- Extensions are automatically converted when switching between OpenAPI versions
