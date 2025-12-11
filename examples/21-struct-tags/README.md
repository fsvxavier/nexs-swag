# Example 21 - Struct Tags

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the use of struct tags to enrich Swagger documentation.

## Tags Suportadas

### Validação

#### String Validation
```go
type User struct {
    Name     string `minLength:"2" maxLength:"50"`
    Email    string `format:"email"`
    Username string `pattern:"^[a-z0-9]+$"`
}
```

#### Numeric Validation
```go
type Product struct {
    Price    float64 `minimum:"0" maximum:"99999.99"`
    Quantity int     `minimum:"1" exclusiveMinimum:"true"`
    Rating   float64 `minimum:"0" maximum:"5" multipleOf:"0.5"`
}
```

#### Array Validation
```go
type Post struct {
    Tags []string `minItems:"1" maxItems:"10"`
    Images []string `uniqueItems:"true"`
}
```

### Metadata

#### Example Values
```go
type User struct {
    Name  string `example:"John Doe"`
    Age   int    `example:"25"`
    Email string `example:"john@example.com"`
}
```

#### Default Values
```go
type Config struct {
    Timeout int    `default:"30"`
    Enabled bool   `default:"true"`
    Format  string `default:"json"`
}
```

#### Enums
```go
type Order struct {
    Status string `enums:"pending,processing,completed,cancelled"`
    Type   string `enums:"online,store"`
}
```

### Type Control

#### swaggertype
```go
type Model struct {
    // Override type detection
    ID        sql.NullInt64 `swaggertype:"integer"`
    UUID      uuid.UUID     `swaggertype:"string"`
    Metadata  interface{}   `swaggertype:"object"`
    
    // Primitive types
    CreatedAt time.Time `swaggertype:"string" format:"date-time"`
    Data      []byte    `swaggertype:"string" format:"byte"`
}
```

#### swaggerignore
```go
type User struct {
    Name     string `json:"name"`
    Password string `json:"-" swaggerignore:"true"`  // Não aparece na doc
    Internal string `swaggerignore:"true"`           // Também ignorado
}
```

#### format
```go
type User struct {
    Email     string    `format:"email"`
    URL       string    `format:"url"`
    IPv4      string    `format:"ipv4"`
    IPv6      string    `format:"ipv6"`
    Birthday  time.Time `format:"date"`
    CreatedAt time.Time `format:"date-time"`
    Token     string    `format:"uuid"`
}
```

### Extensions

#### Standard Extensions
```go
type Model struct {
    // Nullable field
    MiddleName string `x-nullable:"true"`
    
    // Omit empty
    Description string `x-omitempty:"true"`
    
    // Order in documentation
    ID   int    `x-order:"1"`
    Name string `x-order:"2"`
}
```

#### Custom Extensions
```go
type Product struct {
    // Any x- prefix
    SKU   string  `x-internal-code:"true"`
    Price float64 `x-currency:"USD"`
    Stock int     `x-warehouse:"main"`
}
```

### Other Tags

#### readonly
```go
type User struct {
    ID        int       `readonly:"true"`
    CreatedAt time.Time `readonly:"true"`
    UpdatedAt time.Time `readonly:"true"`
}
```

#### required (via validation tags)
```go
type User struct {
    // Campo obrigatório se tiver validação
    Email string `validate:"required" format:"email"`
}
```

## Combinações Comuns

### Email Field
```go
Email string `json:"email" example:"user@example.com" format:"email" minLength:"5" maxLength:"100"`
```

### Password Field
```go
Password string `json:"password" example:"********" minLength:"8" maxLength:"72" format:"password" x-nullable:"false"`
```

### UUID Field
```go
ID uuid.UUID `json:"id" swaggertype:"string" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
```

### Timestamp Field
```go
CreatedAt time.Time `json:"created_at" swaggertype:"string" format:"date-time" example:"2025-01-01T00:00:00Z" readonly:"true"`
```

### Enum Status
```go
Status string `json:"status" enums:"active,inactive,pending" default:"pending" example:"active"`
```

### Price Field
```go
Price float64 `json:"price" example:"99.99" minimum:"0" maximum:"999999.99" multipleOf:"0.01"`
```

## Como Executar

```bash
./run.sh
```

## OpenAPI Gerado

### Input
```go
type User struct {
    Name  string `json:"name" example:"John" minLength:"2" maxLength:"50"`
    Age   int    `json:"age" example:"25" minimum:"0" maximum:"150"`
    Email string `json:"email" example:"john@example.com" format:"email"`
}
```

### Output
```json
{
  "definitions": {
    "User": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string",
          "minLength": 2,
          "maxLength": 50,
          "example": "John"
        },
        "age": {
          "type": "integer",
          "minimum": 0,
          "maximum": 150,
          "example": 25
        },
        "email": {
          "type": "string",
          "format": "email",
          "example": "john@example.com"
        }
      }
    }
  }
}
```

## Prioridade de Tags

1. **swaggertype** (maior prioridade)
2. **format**
3. **Auto-detection**
4. **Validations**

```go
type Model struct {
    // swaggertype override tudo
    Field1 int `swaggertype:"string"` // Será string, não int
    
    // format adiciona detalhe
    Field2 string `format:"email"` // string com format email
    
    // Auto-detected
    Field3 time.Time // Detectado como string date-time
}
```

## Best Practices

### 1. Sempre use example
```go
✅ Name string `example:"John Doe"`
❌ Name string
```

### 2. Validação adequada
```go
✅ Email string `format:"email" minLength:"5"`
❌ Email string
```

### 3. Ignore sensitive data
```go
✅ Password string `swaggerignore:"true"`
❌ Password string `json:"password"`
```

### 4. Override quando necessário
```go
✅ UUID uuid.UUID `swaggertype:"string" format:"uuid"`
❌ UUID uuid.UUID // Pode gerar object complexo
```

### 5. Use defaults
```go
✅ Timeout int `default:"30"`
❌ Timeout int
```

## Validation vs OpenAPI

⚠️ **Atenção:** Tags de validação são para DOCUMENTAÇÃO, não validação em runtime.

```go
type User struct {
    Age int `minimum:"0" maximum:"150"`
}

// Ainda precisa validar no código!
func CreateUser(u User) error {
    if u.Age < 0 || u.Age > 150 {
        return errors.New("invalid age")
    }
}
```

## Ferramentas Complementares

### go-playground/validator
```go
import "github.com/go-playground/validator/v10"

type User struct {
    Email string `validate:"required,email" json:"email" format:"email"`
    Age   int    `validate:"min=0,max=150" json:"age" minimum:"0" maximum:"150"`
}

validate := validator.New()
err := validate.Struct(user)
```

### ozzo-validation
```go
import validation "github.com/go-ozzo/ozzo-validation/v4"

type User struct {
    Email string `json:"email" format:"email"`
}

func (u User) Validate() error {
    return validation.ValidateStruct(&u,
        validation.Field(&u.Email, validation.Required, is.Email),
    )
}
```

## Recomendação

✅ **Use tags extensively** para documentação rica e auto-documentada.

Mantenha consistência entre:
- Struct tags (documentação)
- Validation tags (runtime)
- Código de validação
