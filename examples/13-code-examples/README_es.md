# Ejemplo 13 - Code Examples

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de ejemplos de código en múltiples lenguajes.

## Concepto

OpenAPI soporta ejemplos de código para clients en diferentes lenguajes.

## Tag

```go
// @x-codeSamples.file create_user.go
// @x-codeSamples.lang go
// @x-codeSamples.label "Go Example"
```

## Estructura

```
main.go
code_samples/
├── create_user.go      # Go example
├── create_user.py      # Python example
├── create_user.js      # JavaScript example
└── create_user.sh      # cURL example
```

## Uso

```go
// @Summary Create user
// @Description Creates a new user
// @x-codeSamples.file code_samples/create_user.go
// @x-codeSamples.lang go
// @x-codeSamples.label "Go SDK"
// @x-codeSamples.file code_samples/create_user.py
// @x-codeSamples.lang python
// @x-codeSamples.label "Python SDK"
// @x-codeSamples.file code_samples/create_user.js
// @x-codeSamples.lang javascript
// @x-codeSamples.label "JavaScript SDK"
// @x-codeSamples.file code_samples/create_user.sh
// @x-codeSamples.lang shell
// @x-codeSamples.label "cURL"
func CreateUser(c *gin.Context) {}
```

## Cómo Ejecutar

```bash
./run.sh
```

## Ejemplos de Código

### Go
```go
// code_samples/create_user.go
client := sdk.NewClient("api-key")

user := &sdk.User{
    Name:  "John Doe",
    Email: "john@example.com",
}

result, err := client.Users.Create(user)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created user: %s\n", result.ID)
```

### Python
```python
# code_samples/create_user.py
from sdk import Client

client = Client(api_key='api-key')

user = {
    'name': 'John Doe',
    'email': 'john@example.com'
}

result = client.users.create(user)
print(f"Created user: {result['id']}")
```

### JavaScript
```javascript
// code_samples/create_user.js
const SDK = require('@company/sdk');

const client = new SDK.Client('api-key');

const user = {
  name: 'John Doe',
  email: 'john@example.com'
};

const result = await client.users.create(user);
console.log(`Created user: ${result.id}`);
```

### cURL
```bash
# code_samples/create_user.sh
curl -X POST https://api.example.com/users \
  -H "Authorization: Bearer api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

## Lenguajes Soportados

- `go`
- `python`
- `javascript`
- `java`
- `ruby`
- `php`
- `csharp`
- `shell` (bash/curl)
- `typescript`
- `kotlin`
- `swift`

## Output OpenAPI

```yaml
paths:
  /users:
    post:
      summary: Create user
      x-codeSamples:
        - lang: go
          label: Go SDK
          source: |
            client := sdk.NewClient("api-key")
            user := &sdk.User{
              Name: "John Doe",
              Email: "john@example.com",
            }
            result, err := client.Users.Create(user)
        
        - lang: python
          label: Python SDK
          source: |
            client = Client(api_key='api-key')
            user = {'name': 'John Doe', 'email': 'john@example.com'}
            result = client.users.create(user)
```

## Ferramentas que Soportan

### Swagger UI
```yaml
# swagger-config.yaml
plugins:
  - codeSamples
```

### ReDoc
Soporta x-codeSamples nativamente

### Stoplight
Renderiza ejemplos automáticamente

## Tips

### 1. Organización por Endpoint
```
code_samples/
├── users/
│   ├── create.go
│   ├── create.py
│   └── create.js
├── orders/
│   ├── create.go
│   └── list.py
```

### 2. Ejemplos Realistas
```go
// ✅ GOOD - Ejemplo completo
client := NewClient(os.Getenv("API_KEY"))
ctx := context.Background()

user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
}

result, err := client.Users.Create(ctx, user)
if err != nil {
    return fmt.Errorf("create user: %w", err)
}

// ❌ BAD - Muy simplificado
client.Create(user)
```

### 3. Incluir Error Handling
```python
# ✅ GOOD
try:
    result = client.users.create(user)
except APIError as e:
    print(f"Error: {e}")

# ❌ BAD
result = client.users.create(user)
```

## Cuándo Usar

**Use x-codeSamples cuando:**
- Tiene SDK en múltiples lenguajes
- Quiere facilitar integración
- Documentation para desarrolladores
- Ejemplos específicos por endpoint

**NO use cuando:**
- API muy simple (solo REST puro)
- Sin SDKs disponibles
- Ejemplos cambiarían frecuentemente
