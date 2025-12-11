# Example 05 - Required By Default

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the behavior of the `--requiredByDefault` flag that makes all fields required by default.

## Flag

```bash
--requiredByDefault
```

## Behavior

### WITHOUT the flag (default)
```bash
nexs-swag init
```
- **All fields are OPTIONAL** by default
- Only fields with `binding:"required"` or `validate:"required"` are required

### WITH the flag
```bash
nexs-swag init --requiredByDefault
```
- **All fields are REQUIRED** by default
- Exceptions:
  - Fields with `json:",omitempty"`
  - Fields with `binding:"omitempty"`  
  - Fields that are pointers (`*Type`)

## Example

```go
type Product struct {
    ID          int      // ✅ Required (with flag)
    Name        string   // ✅ Required (with flag)
    Description string   `json:"description,omitempty"` // ❌ Optional (omitempty)
    Discount    *float64 // ❌ Optional (pointer)
    Category    string   `json:"category" binding:"omitempty"` // ❌ Optional (binding)
}
```

## Generated Schema

### Without --requiredByDefault
```json
{
  "Product": {
    "type": "object",
    "properties": {
      "id": {"type": "integer"},
      "name": {"type": "string"},
      "description": {"type": "string"},
      "discount": {"type": "number"},
      "category": {"type": "string"}
    }
  }
}
```

### With --requiredByDefault
```json
{
  "Product": {
    "type": "object",
    "required": ["id", "name"],
    "properties": {
      "id": {"type": "integer"},
      "name": {"type": "string"},
      "description": {"type": "string"},
      "discount": {"type": "number"},
      "category": {"type": "string"}
    }
  }
}
```

## How to Run

```bash
chmod +x run.sh
./run.sh
```

## Use Cases

**Use `--requiredByDefault` when:**
- API requires strict validation
- Most fields are mandatory
- Prefer opt-out (mark optionals) instead of opt-in

**DON'T use when:**
- API has many optional fields
- Prefer explicit opt-in with `binding:"required"`
- Compatibility with existing clients
