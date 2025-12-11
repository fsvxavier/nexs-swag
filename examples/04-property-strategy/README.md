# Example 04 - Property Naming Strategy

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the different naming strategies for struct fields.

## Flag

```bash
--propertyStrategy <strategy>
-p <strategy>
```

## Available Strategies

### 1. snake_case
```bash
nexs-swag init --propertyStrategy snakecase
```
- `FirstName` → `first_name`
- `LastName` → `last_name`
- `IsActive` → `is_active`

### 2. camelCase (default)
```bash
nexs-swag init --propertyStrategy camelcase
```
- `FirstName` → `firstName`
- `LastName` → `lastName`
- `IsActive` → `isActive`

### 3. PascalCase
```bash
nexs-swag init --propertyStrategy pascalcase
```
- `FirstName` → `FirstName`
- `LastName` → `LastName`
- `IsActive` → `IsActive`

## Important Rules

### ✅ Applied when:
- Field does **NOT** have `json` tag
- Field does **NOT** have `json:"-"`

### ❌ NOT applied when:
- Field has tag `json:"explicit_name"` → uses "explicit_name"
- Field has `json:"-"` → ignored
- Field has `json:",omitempty"` → applies strategy + omitempty

## Example

```go
type User struct {
    UserID    int    `json:"user_id"`      // ✅ ALWAYS "user_id"
    FirstName string                       // ⚠️ Depends on strategy
    LastName  string `json:",omitempty"`   // ⚠️ Strategy + omitempty
    Password  string `json:"-"`            // ❌ Ignored
}
```

With `--propertyStrategy snakecase`:
```json
{
    "user_id": 123,
    "first_name": "John",
    "last_name": "Doe"
}
```

With `--propertyStrategy camelcase`:
```json
{
    "user_id": 123,
    "firstName": "John",
    "lastName": "Doe"
}
```

## How to Run

```bash
chmod +x run.sh
./run.sh
```

## Use Cases

- **snake_case:** Python, Ruby APIs, databases
- **camelCase:** JavaScript APIs, standard JSON
- **PascalCase:** C#, .NET APIs
