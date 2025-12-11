# Example 17 - Template Delims

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the use of `--templateDelims` to customize template delimiters.

## Flag

```bash
--templateDelims "<left>,<right>"
```

Default: `"{{,}}"`

## Problem

Default delimiters `{{` and `}}` conflict with Go templates:

```go
const template = `
    API Docs: {{.BasePath}}  // ❌ Conflict!
`
```

## Solution

```bash
nexs-swag init --templateDelims "[[,]]"
```

Now in your code:
```go
const template = `
    API Docs: {{.BasePath}}  // ✅ OK! Not parsed by nexs-swag
    Swagger: [[.Version]]     // ✅ Parsed by nexs-swag
`
```

## Usage Examples

### Square Brackets
```bash
nexs-swag init --templateDelims "[[,]]"
```

### Angle Brackets
```bash
nexs-swag init --templateDelims "<<,>>"
```

### Custom
```bash
nexs-swag init --templateDelims "{%,%}"
```

## Generated docs.go

### Default Delimiters
```go
const docTemplate = `{
    "swagger": "2.0",
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}"
}`
```

### Custom Delimiters [[,]]
```go
const docTemplate = `{
    "swagger": "2.0",
    "host": "[[.Host]]",
    "basePath": "[[.BasePath]]"
}`
```

## When to Use

**Use custom delimiters when:**
- Using Go text/template
- Conflicts with other template engines
- Code has `{{` and `}}` in strings
- Embedding documentation in templates

**Keep default when:**
- No conflicts
- Standard project
- Following conventions

## How to Run

```bash
./run.sh
```

## Common Conflicts

### Go Templates
```go
// ❌ Conflict with default {{,}}
tmpl := template.Must(template.New("").Parse(`
    <h1>{{.Title}}</h1>
`))

// ✅ Solution: use different delimiters
// nexs-swag init --templateDelims "[[,]]"
```

### Vue.js/Angular in Comments
```go
// Documentation: {{variable}}  // ❌ Parsed by nexs-swag
// With [[,]]: {{variable}}     // ✅ Ignored by nexs-swag
```

### String Literals
```go
const json = `{"key": "{{value}}"}`  // ❌ May cause issues
// Use [[,]] delimiters to avoid
```

## Best Practices

1. Choose delimiters that won't conflict
2. Document the choice in README
3. Be consistent across the project
4. Add to build scripts/Makefile
