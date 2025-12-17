# X-Visibility Example

This example demonstrates the `@x-visibility` annotation feature that allows you to generate separate OpenAPI documentation files for public and private APIs.

## Overview

The `@x-visibility` annotation enables you to:
- Separate public and private endpoints into different OpenAPI files
- Automatically filter schemas based on usage
- Maintain a single codebase with dual documentation output

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

- `@x-visibility public` - Endpoint appears only in `openapi_public.json`
- `@x-visibility private` - Endpoint appears only in `openapi_private.json`
- No annotation - Endpoint appears in **both** files (shared endpoint)

## Generated Files

When using `@x-visibility`, nexs-swag generates:

```
docs/
├── openapi_public.json    # Public API specification
├── openapi_private.json   # Private API specification
├── openapi_public.yaml    # Public API (YAML)
├── openapi_private.yaml   # Private API (YAML)
├── docs_public.go         # Public API Go code
└── docs_private.go        # Private API Go code
```

## Example Structure

```go
// Public endpoint - for external consumers
// @x-visibility public
func GetUser(c *gin.Context) {
    c.JSON(200, UserPublic{ID: 1, Name: "John"})
}

// Private endpoint - for internal/admin use
// @x-visibility private
func GetUserAdmin(c *gin.Context) {
    c.JSON(200, UserPrivate{
        ID: 1, 
        Name: "John",
        Email: "john@example.com",
        Password: "hashed",
        Role: "admin",
    })
}

// Shared endpoint - available to both
func CreateUser(c *gin.Context) {
    c.JSON(201, UserPublic{ID: 2, Name: "Jane"})
}
```

## Schema Filtering

Schemas are automatically filtered based on usage:

- **Public spec**: Only includes schemas referenced by public operations (`UserPublic`, `ErrorResponse`)
- **Private spec**: Only includes schemas referenced by private operations (`UserPrivate`, `ErrorResponse`)
- Shared schemas appear in both if used by operations without visibility annotations

## Running This Example

```bash
# Generate documentation
nexs-swag init --output ./docs --ov 3.1

# Verify separation
jq '.paths | keys' docs/openapi_public.json
# Output: ["/users", "/users/{id}"]

jq '.paths | keys' docs/openapi_private.json
# Output: ["/admin/users/{id}", "/users"]

jq '.components.schemas | keys' docs/openapi_public.json
# Output: ["ErrorResponse", "UserPublic"]

jq '.components.schemas | keys' docs/openapi_private.json
# Output: ["ErrorResponse", "UserPrivate", "UserPublic"]
```

## Use Cases

1. **API Versioning**: Separate stable public APIs from experimental private APIs
2. **Security**: Hide internal admin endpoints from public documentation
3. **Client Libraries**: Generate different client SDKs for public vs private APIs
4. **Documentation Sites**: Host separate documentation for different audiences
5. **Microservices**: Distinguish between external APIs and inter-service APIs

## Benefits

- **Single Source of Truth**: Maintain all APIs in one codebase
- **Automatic Schema Management**: No manual schema duplication
- **Type Safety**: Same Go types ensure consistency
- **Selective Exposure**: Control what information is publicly visible
- **Flexible Deployment**: Choose which spec to publish based on context

## Notes

- Schemas are recursively collected including nested references
- Operations without `@x-visibility` appear in both specs
- Shared schemas (like `ErrorResponse`) are included where needed
- All other OpenAPI features work normally within each spec
