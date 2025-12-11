# Example 14 - Overrides File

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the use of `.swaggo` file for global type overrides.

## Flag

```bash
--overridesFile <path>
```

Default: `.swaggo`

## Usage

```bash
nexs-swag init --overridesFile .swaggo
```

## .swaggo File

```json
{
  "types": {
    "database/sql.NullInt64": "integer",
    "database/sql.NullString": "string",
    "database/sql.NullBool": "boolean",
    "database/sql.NullFloat64": "number",
    "github.com/google/uuid.UUID": "string",
    "github.com/shopspring/decimal.Decimal": "number"
  }
}
```

## Problem

Custom Go types generate complex schemas:

### WITHOUT Overrides
```json
{
  "Account": {
    "properties": {
      "id": {
        "type": "object",
        "properties": {
          "Int64": {"type": "integer"},
          "Valid": {"type": "boolean"}
        }
      }
    }
  }
}
```

### WITH Overrides
```json
{
  "Account": {
    "properties": {
      "id": {
        "type": "integer"
      }
    }
  }
}
```

## Common Types

### SQL Nullable Types
```json
{
  "database/sql.NullInt64": "integer",
  "database/sql.NullString": "string",
  "database/sql.NullBool": "boolean",
  "database/sql.NullFloat64": "number",
  "database/sql.NullTime": "string"
}
```

### UUID Libraries
```json
{
  "github.com/google/uuid.UUID": "string",
  "github.com/satori/go.uuid.UUID": "string"
}
```

### Decimal Libraries
```json
{
  "github.com/shopspring/decimal.Decimal": "number",
  "gopkg.in/inf.v0.Dec": "number"
}
```

### Custom Time Types
```json
{
  "myapp/types.Timestamp": "integer",
  "myapp/types.UnixTime": "integer"
}
```

## Alternative: swaggertype Tag

Without global file:
```go
type Account struct {
    ID sql.NullInt64 `json:"id" swaggertype:"integer"`
}
```

With global file:
```go
type Account struct {
    ID sql.NullInt64 `json:"id"` // Automatically integer
}
```

## Priority

1. **swaggertype tag** (highest priority)
2. **Overrides file**
3. **Auto-detection** (lowest priority)

```go
type User struct {
    ID        sql.NullInt64 `swaggertype:"string"` // string (tag)
    AccountID sql.NullInt64                         // integer (override)
    Balance   float64                               // number (auto)
}
```

## How to Run

```bash
./run.sh
```

## Benefits

- **Global:** One place for all overrides
- **Reusable:** Share across projects
- **Clean:** Code without repetitive tags
- **Versionable:** Commit with code

## Use Cases

- Projects with database/sql
- Third-party libraries
- Company custom types
- Type standardization
