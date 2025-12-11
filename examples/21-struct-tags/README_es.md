# Ejemplo 21 - Struct Tags

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de struct tags para configurar properties en OpenAPI.

## Concepto

Go struct tags controlan cómo las structs se convierten a JSON/YAML en OpenAPI.

## Tags Soportados

### 1. json Tag
```go
type User struct {
    ID        int    `json:"id"`                    // Nombre del campo
    Name      string `json:"name"`                  // Nombre del campo
    Email     string `json:"email,omitempty"`       // Omite si vacío
    Password  string `json:"-"`                     // Ignora este campo
    CreatedAt string `json:"created_at,omitempty"`  // Snake case
}
```

### 2. binding Tag (Gin Validation)
```go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`              // Requerido
    Email    string `json:"email" binding:"required,email"`       // Email válido
    Age      int    `json:"age" binding:"gte=0,lte=120"`         // Range
    Password string `json:"password" binding:"required,min=8"`    // Min length
}
```

### 3. example Tag
```go
type User struct {
    ID    int    `json:"id" example:"123"`
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john@example.com"`
}
```

### 4. format Tag
```go
type User struct {
    Email     string    `json:"email" format:"email"`
    URL       string    `json:"url" format:"uri"`
    CreatedAt time.Time `json:"created_at" format:"date-time"`
    BirthDate string    `json:"birth_date" format:"date"`
}
```

### 5. enums Tag
```go
type Order struct {
    Status string `json:"status" enums:"pending,approved,rejected"`
}
```

### 6. validate Tag
```go
type Product struct {
    Price float64 `json:"price" validate:"gt=0"`
    Stock int     `json:"stock" validate:"gte=0"`
}
```

## Cómo Ejecutar

```bash
./run.sh
```

## Ejemplos Completos

### Request Body
```go
type CreateUserRequest struct {
    // Name es requerido, min 3 chars
    Name string `json:"name" binding:"required,min=3" example:"John Doe"`
    
    // Email es requerido y debe ser válido
    Email string `json:"email" binding:"required,email" format:"email" example:"john@example.com"`
    
    // Age opcional, pero si se proporciona: 0-120
    Age *int `json:"age,omitempty" binding:"omitempty,gte=0,lte=120" example:"30"`
    
    // Password requerido, min 8 chars, no aparece en spec
    Password string `json:"-" binding:"required,min=8"`
}
```

### Response Body
```go
type User struct {
    ID        int       `json:"id" example:"123"`
    Name      string    `json:"name" example:"John Doe"`
    Email     string    `json:"email" format:"email" example:"john@example.com"`
    Age       int       `json:"age" example:"30"`
    Status    string    `json:"status" enums:"active,inactive,banned" example:"active"`
    CreatedAt time.Time `json:"created_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
    UpdatedAt time.Time `json:"updated_at,omitempty" format:"date-time"`
}
```

### Nested Structs
```go
type Order struct {
    ID         int           `json:"id" example:"456"`
    User       User          `json:"user"`                                    // Nested
    Items      []OrderItem   `json:"items"`                                   // Array
    Total      float64       `json:"total" example:"99.99"`
    Status     string        `json:"status" enums:"pending,paid,shipped"`
    Metadata   OrderMetadata `json:"metadata,omitempty"`                     // Optional nested
}

type OrderItem struct {
    ProductID int     `json:"product_id" example:"789"`
    Quantity  int     `json:"quantity" example:"2" binding:"required,gte=1"`
    Price     float64 `json:"price" example:"49.99"`
}

type OrderMetadata struct {
    Notes      string            `json:"notes,omitempty"`
    CustomData map[string]string `json:"custom_data,omitempty"`
}
```

## OpenAPI Output

### Input
```go
type User struct {
    ID    int    `json:"id" example:"123"`
    Name  string `json:"name" binding:"required" example:"John Doe"`
    Email string `json:"email" format:"email" example:"john@example.com"`
}
```

### Output
```yaml
components:
  schemas:
    User:
      type: object
      required:
        - name
      properties:
        id:
          type: integer
          example: 123
        name:
          type: string
          example: John Doe
        email:
          type: string
          format: email
          example: john@example.com
```

## Binding Tags (Gin)

### Validation Rules
```go
type Product struct {
    // Required
    Name string `json:"name" binding:"required"`
    
    // Min/Max length
    Description string `json:"description" binding:"min=10,max=500"`
    
    // Numeric range
    Price float64 `json:"price" binding:"required,gt=0"`
    Stock int     `json:"stock" binding:"gte=0"`
    
    // Email format
    ContactEmail string `json:"contact_email" binding:"omitempty,email"`
    
    // URL format
    Website string `json:"website" binding:"omitempty,url"`
    
    // One of
    Category string `json:"category" binding:"required,oneof=electronics clothing food"`
    
    // UUID
    SKU string `json:"sku" binding:"uuid"`
    
    // Custom validation
    Code string `json:"code" binding:"required,alphanum,len=8"`
}
```

### Arrays
```go
type Order struct {
    // Array required, min 1 item
    Items []OrderItem `json:"items" binding:"required,min=1,dive"`
    
    // Each item validated
    Tags []string `json:"tags" binding:"dive,required,min=3"`
}
```

## omitempty

```go
type User struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Email *string `json:"email,omitempty"`  // Pointer + omitempty
    Age   int     `json:"age,omitempty"`    // Omite si zero value (0)
}
```

### Behavior
```go
user1 := User{ID: 1, Name: "John", Email: nil, Age: 0}
// JSON: {"id":1,"name":"John"}
// email e age omitidos

user2 := User{ID: 2, Name: "Jane", Email: strPtr("jane@example.com"), Age: 30}
// JSON: {"id":2,"name":"Jane","email":"jane@example.com","age":30}
```

## Ignore Fields

```go
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Password string `json:"-"`              // ❌ NUNCA en JSON
    Token    string `json:"-"`              // ❌ NUNCA en JSON
    internal string                         // ❌ Lowercase = private
}
```

## Tips

### 1. Required vs Optional
```go
// ✅ GOOD - Pointer para opcional
type User struct {
    Name  string  `json:"name" binding:"required"`      // Required
    Email *string `json:"email,omitempty"`              // Optional
}

// ❌ BAD - String para opcional
type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email,omitempty"`  // "" = empty, no null
}
```

### 2. Default Values
```go
type CreateOrderRequest struct {
    Items    []OrderItem `json:"items" binding:"required"`
    Priority string      `json:"priority,omitempty" default:"normal" enums:"low,normal,high"`
    Notes    string      `json:"notes,omitempty"`
}
```

### 3. Multiple Formats
```go
type Contact struct {
    Email string `json:"email" binding:"required,email" format:"email"`
    Phone string `json:"phone" binding:"omitempty" format:"phone" example:"+1-555-0123"`
    URL   string `json:"url" binding:"omitempty,url" format:"uri"`
}
```

### 4. Complex Validation
```go
type Payment struct {
    Amount   float64 `json:"amount" binding:"required,gt=0"`
    Currency string  `json:"currency" binding:"required,len=3" example:"USD"`
    Method   string  `json:"method" binding:"required,oneof=credit_card debit_card paypal"`
    CardNum  string  `json:"card_number,omitempty" binding:"omitempty,credit_card"`
}
```

## Custom Validators

```go
// Register custom validator
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("custom", func(fl validator.FieldLevel) bool {
        // Custom logic
        return true
    })
}

type Product struct {
    Code string `json:"code" binding:"custom"`
}
```

## Best Practices

### 1. Consistencia
```go
// ✅ GOOD - Consistente
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// ❌ BAD - Inconsistente
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"Email"`        // Uppercase
    CreatedAt time.Time `json:"createdAt"`    // camelCase
}
```

### 2. Ejemplos Realistas
```go
// ✅ GOOD
type User struct {
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john.doe@example.com"`
    Age   int    `json:"age" example:"30"`
}

// ❌ BAD
type User struct {
    Name  string `json:"name" example:"test"`
    Email string `json:"email" example:"test@test.com"`
    Age   int    `json:"age" example:"1"`
}
```

### 3. Documentar Validations
```go
type CreateUserRequest struct {
    // Name is required and must be between 3 and 50 characters
    Name string `json:"name" binding:"required,min=3,max=50"`
    
    // Email is required and must be a valid email format
    Email string `json:"email" binding:"required,email"`
    
    // Age is optional but must be between 18 and 120 if provided
    Age *int `json:"age,omitempty" binding:"omitempty,gte=18,lte=120"`
}
```

## Cuándo Usar

**Use struct tags para:**
- ✅ Definir nombres de campos JSON
- ✅ Marcar campos required/optional
- ✅ Validación básica
- ✅ Ejemplos en spec
- ✅ Formatos standard

**NO use para:**
- ❌ Lógica de negocio compleja
- ❌ Validations interdependientes
- ❌ Mensajes de error custom
