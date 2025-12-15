# Example 09 - Parse Dependency

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the use of `--parseDependency` to include types from imported packages.

## Flag

```bash
--parseDependency
--pd
# Or explicitly:
--parseDependency=true
--pd=true
```

> **Note:** Both syntaxes are valid. Use `--parseDependency` or `--pd` (without value), or `--parseDependency=true` (explicit).

## Concept

This example demonstrates how nexs-swag can parse dependencies when you have types defined in separate packages. In this simplified example, we show the concept with a single file, but in real projects you would have:

```
myapp/
├── main.go              # Usa models.Product
└── models/
    └── product.go       # Define Product
```

## Usage

```bash
nexs-swag init --parseDependency
# Or:
nexs-swag init --pd
# Or explicitly:
nexs-swag init --parseDependency=true
```

## How It Works

### WITHOUT --parseDependency
Only types from the current package are included in the documentation.

### WITH --parseDependency
Types from imported packages are also parsed and included.

## Structure in Real Projects

```go
// main.go
package main

import "myapp/models"

// @Success 200 {object} models.Product
func GetProduct() {}

// models/product.go
package models

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

## How to Run

```bash
./run.sh
```

## When to Use

**Use --parseDependency when:**
- Models in separate packages
- Modular structure
- Imported types
- Shared libraries

**NOT needed when:**
- All types in the same package
- Simple API
- No model imports

## Parsing Levels

Combine with `--parseDependencyLevel` to control depth:

```bash
# Level 0: Only main directory
nexs-swag init --parseDependency --parseDependencyLevel 0

# Level 1: + 1 dependency level
nexs-swag init --parseDependency --parseDependencyLevel 1

# Level 2: + 2 levels (default)
nexs-swag init --parseDependency --parseDependencyLevel 2
```

## Performance

⚠️ **Warning:** Parsing many dependencies can be slow.

Optimizations:
```bash
# Only what's necessary
nexs-swag init --parseDependency --parseDependencyLevel 1

# Limit with --exclude
nexs-swag init --parseDependency --exclude "vendor,node_modules"
```
