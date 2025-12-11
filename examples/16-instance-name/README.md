# Example 16 - Instance Name

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the use of `--instanceName` to customize the generated docs variable name.

## Flag

```bash
--instanceName <name>
```

Default: `swagger`

## Usage

```bash
nexs-swag init --instanceName myapi
```

## Generated Code

### Default (swagger)
```go
package docs

import "github.com/swaggo/swag"

const docTemplate = `...`

var SwaggerInfo = &swag.Spec{...}

func init() {
    swag.Register(swag.Name, SwaggerInfo)
}
```

### Custom (myapi)
```go
package docs

import "github.com/swaggo/swag"

const docTemplate = `...`

var MyapiInfo = &swag.Spec{...}

func init() {
    swag.Register("myapi", MyapiInfo)
}
```

## Use Cases

### Multiple APIs in One Project
```bash
# API v1
nexs-swag init --dir ./apiv1 --output ./docs/v1 --instanceName apiv1

# API v2
nexs-swag init --dir ./apiv2 --output ./docs/v2 --instanceName apiv2
```

### Integration
```go
package main

import (
    _ "myapp/docs/v1" // apiv1
    _ "myapp/docs/v2" // apiv2
    
    httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
    // Serve v1
    http.HandleFunc("/swagger/v1/", httpSwagger.Handler(
        httpSwagger.InstanceName("apiv1"),
    ))
    
    // Serve v2
    http.HandleFunc("/swagger/v2/", httpSwagger.Handler(
        httpSwagger.InstanceName("apiv2"),
    ))
    
    http.ListenAndServe(":8080", nil)
}
```

## Benefits

- **Multiple APIs:** Different specifications in same project
- **Versioning:** API v1, v2, v3 simultaneously
- **Isolation:** No conflicts between specifications
- **Organization:** Clear naming for each API

## How to Run

```bash
./run.sh
```

## Common Names

- `api` - Main API
- `apiv1`, `apiv2` - Versioned APIs
- `admin` - Admin API
- `public` - Public API
- `internal` - Internal API

## Important

Instance name must:
- Start with letter
- Contain only letters and numbers
- Be unique in your project
