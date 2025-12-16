# Example 25: Include Types Complex

This example demonstrates advanced usage of the `--includeTypes` flag with complex scenarios including nested structs, swaggertype overrides, multiple format specifications, and transitive dependencies.

## Features Demonstrated

1. **Nested Struct Dependencies**: OrderRequest → OrderItem → Money → Address
2. **swaggertype Overrides**:
   - `uuid.UUID` → `string` with `format:"uuid"`
   - `time.Time` → `string` with `format:"date-time"`
   - `int64` → `string` for precision (Money.Amount)
   - `interface{}` → `object` for dynamic data
3. **Advanced Format Specifications**:
   - `format:"uuid"`, `format:"date-time"`, `format:"double"`, `format:"int64"`
   - Validation: `minimum`, `maximum`, `minLength`, `maxLength`, `pattern`
4. **Transitive Dependency Resolution**:
   - Referenced struct → nested struct → deeply nested struct
5. **Enum Types**: StatusCode, PaymentMethod (string constants)
6. **Optional Fields**: Pointers, omitempty tags
7. **Multi-file Structure**:
   - `models/order.go` - All model definitions
   - `services/order.go` - Business logic and interfaces
   - `main.go` - API handlers

## File Structure

```
25-include-types-complex/
├── main.go                 # API handlers with OpenAPI annotations
├── models/
│   └── order.go           # Complex model definitions
├── services/
│   └── order.go           # Service interfaces and implementations
├── README.md
└── run.sh
```

## Type Categories and Dependencies

### Referenced Types (Included)
From `@Success`, `@Param`, `@Failure` annotations:

1. **OrderRequest** (referenced in `@Param`)
   - → OrderItem (nested in Items field)
     - → Money (nested in UnitPrice, Subtotal, Discount)
   - → PaymentInfo (nested in Payment field)
     - → PaymentMethod (nested in Method field)
     - → time.Time (converted via swaggertype)
   - → Address (nested in ShippingAddr, BillingAddr)

2. **OrderResponse** (referenced in `@Success`)
   - → uuid.UUID (converted via swaggertype)
   - → OrderItem, PaymentInfo, Address, Money (transitive)
   - → OrderStatus
     - → StatusCode
   - → time.Time (converted via swaggertype)

3. **StatusUpdate** (referenced in `@Param`)
   - → StatusCode

4. **OrderListResponse** (referenced in `@Success`)
   - → []OrderResponse (array, transitive)

5. **ErrorResponse** (referenced in `@Failure`)

### Unreferenced Types (Excluded)
- `UnusedComplexModel` in models/order.go
- `UnusedService` in services/order.go
- `OrderProcessor` interface (unless --includeTypes="interface")

## Usage Examples

### 1. Include All Types (Default)
```bash
nexs-swag init -g main.go -o ./docs
```
Result: All referenced structs included

### 2. Include Only Structs
```bash
nexs-swag init -g main.go -o ./docs --includeTypes="struct"
```
Result: All struct models included, interfaces excluded

### 3. Include Structs and Interfaces
```bash
nexs-swag init -g main.go -o ./docs --includeTypes="struct,interface"
```
Result: Structs + OrderProcessor interface (if referenced)

### 4. Test Dependency Resolution
```bash
# Generate with default (should include Money transitively)
nexs-swag init -g main.go -o ./docs-all

# Check if Money is included even though not directly referenced
grep -o '"Money":' ./docs-all/openapi.json
```

## swaggertype Examples in This Example

### UUID Conversion
```go
OrderID uuid.UUID `swaggertype:"string" format:"uuid"`
```
OpenAPI output:
```json
{
  "order_id": {
    "type": "string",
    "format": "uuid",
    "example": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### Time Conversion
```go
CreatedAt time.Time `swaggertype:"string" format:"date-time"`
```
OpenAPI output:
```json
{
  "created_at": {
    "type": "string",
    "format": "date-time",
    "example": "2025-12-16T10:30:00Z"
  }
}
```

### Integer as String for Precision
```go
Amount int64 `swaggertype:"string" format:"int64" example:"99999"`
```
OpenAPI output:
```json
{
  "amount": {
    "type": "string",
    "format": "int64",
    "example": "99999"
  }
}
```

### Dynamic Object
```go
Metadata interface{} `swaggertype:"object"`
```
OpenAPI output:
```json
{
  "metadata": {
    "type": "object"
  }
}
```

## Expected Schemas

With `--includeTypes="struct"`, the following schemas should be generated:

1. `models.OrderRequest`
2. `models.OrderResponse`
3. `models.OrderItem`
4. `models.PaymentInfo`
5. `models.Address`
6. `models.Money`
7. `models.OrderStatus`
8. `models.StatusCode` (enum)
9. `models.PaymentMethod` (enum)
10. `models.StatusUpdate`
11. `models.OrderListResponse`
12. `models.ErrorResponse`

**NOT included**:
- `models.UnusedComplexModel` (not referenced)
- `services.UnusedService` (not referenced)
- `services.OrderProcessor` (interface, only if --includeTypes includes "interface")

## Validation Features

The example demonstrates various OpenAPI validation keywords:

- **minimum/maximum**: `Latitude`, `Longitude`, `Page`
- **minLength/maxLength**: `State`, `Country`, `Notes`
- **pattern**: `CardLast4` (exactly 4 digits), `PostalCode` (5 digits)
- **minItems**: `Items` (at least 1 item)
- **enum**: StatusCode, PaymentMethod (via const values)

## Running the Example

```bash
# Generate documentation
./run.sh

# View generated schemas
cat docs-all/openapi.json | python3 -m json.tool | grep -A5 '"schemas"'

# Count schemas
grep -c '"type": "object"' docs-all/openapi.json

# Check swaggertype conversions
grep -A2 '"order_id"\|"created_at"\|"amount"' docs-all/openapi.json
```

## Key Takeaways

1. **Transitive Dependencies**: All nested structs are automatically included
2. **swaggertype Override**: Prevents parsing complex external types (uuid.UUID, time.Time)
3. **Precision Control**: Use swaggertype to control how types are represented
4. **Selective Parsing**: Only types reachable from API operations are included
5. **Type Category Filtering**: --includeTypes adds additional granular control
6. **Multi-package Support**: Handles cross-package references correctly
7. **Validation Rich**: OpenAPI validation keywords enhance API documentation

## Testing Scenarios

1. **Verify transitive inclusion**: Money is included via OrderItem
2. **Verify swaggertype works**: uuid.UUID appears as string, not object
3. **Verify exclusion**: UnusedComplexModel not in output
4. **Verify nested arrays**: Items array properly references OrderItem
5. **Verify optional fields**: Discount (*Money) handled correctly
6. **Verify enums**: StatusCode and PaymentMethod have enum values
