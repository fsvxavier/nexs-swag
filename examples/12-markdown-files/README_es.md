# Ejemplo 12 - Markdown Files

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de archivos Markdown para documentación extendida en OpenAPI.

## Concepto

Puede usar archivos .md para descripciones largas:

```go
// @description.file docs/api-description.md
```

## Uso

```go
// @title Order API
// @description.file api-description.md
// @version 1.0

// @Summary Create order
// @Description.file create-order-description.md
func CreateOrder(c *gin.Context) {}
```

## Estructura

```
main.go
docs/
├── api-description.md         # Description general
├── create-order-description.md  # Description de operation
└── models/
    └── order-model.md         # Description de schema
```

## Ventajas

### 1. Descripciones Largas
```markdown
<!-- api-description.md -->
# Order Management API

Esta API permite gestionar pedidos...

## Características

- CRUD completo
- Validación robusta
- Soporte para bulk operations
```

### 2. Formateo Markdown
- **Negrita** y *itálica*
- Listas con viñetas
- Bloques de código
- Links
- Tablas

### 3. Mantenimiento
- Separación de concerns
- Fácil edición
- Versionamiento independiente
- Reutilización

## Cómo Ejecutar

```bash
./run.sh
```

## Tags Soportados

### General Info
```go
// @description.file api-description.md
// @termsOfService.file terms.md
```

### Operations
```go
// @Summary Quick summary (inline)
// @Description.file operation-details.md
```

### Schemas
```go
type Order struct {
    // @description.file order-description.md
    ID string
}
```

## Ejemplo Completo

### main.go
```go
// @title E-commerce API
// @description.file docs/api-description.md
// @version 2.0.0
// @contact.name Support
// @contact.email support@example.com

// @Summary Create new order
// @Description.file docs/orders/create.md
// @Tags orders
// @Accept json
// @Produce json
// @Param order body Order true "Order object"
// @Success 201 {object} Order
// @Router /orders [post]
func CreateOrder(c *gin.Context) {}
```

### docs/api-description.md
```markdown
# E-commerce API Documentation

Esta API provee endpoints para gestionar un sistema de e-commerce.

## Features

- 🛒 **Orders**: Create, read, update, delete
- 📦 **Products**: Catalog management
- �� **Users**: Authentication and profiles
- 💳 **Payments**: Multiple payment methods

## Rate Limiting

- 1000 requests/hour para usuarios autenticados
- 100 requests/hour para usuarios anónimos

## Authentication

Usa Bearer token:

\`\`\`bash
Authorization: Bearer <token>
\`\`\`
```

### docs/orders/create.md
```markdown
## Create Order Endpoint

Crea un nuevo pedido en el sistema.

### Request Body

El body debe contener:
- Items válidos con quantities > 0
- Shipping address completo
- Payment method válido

### Validations

1. **Items**: Mínimo 1 item
2. **Stock**: Verifica disponibilidad
3. **Payment**: Valida método de pago
4. **Address**: Formato correcto

### Example

\`\`\`json
{
  "items": [
    {"product_id": "123", "quantity": 2}
  ],
  "shipping_address": {
    "street": "123 Main St",
    "city": "New York",
    "zipcode": "10001"
  },
  "payment_method": "credit_card"
}
\`\`\`
```

## Path Resolution

nexs-swag busca archivos .md relativos a:

1. Directorio del archivo .go
2. Directorio actual
3. `--dir` especificado

```bash
main.go → busca en ./
api/handlers.go → busca en api/ luego ./
```

## Tips

### 1. Organización
```
docs/
├── general/
│   ├── api-description.md
│   └── terms.md
├── operations/
│   ├── orders/
│   │   ├── create.md
│   │   └── update.md
│   └── products/
└── schemas/
    ├── order.md
    └── product.md
```

### 2. Templates
Cree templates reutilizables:

```markdown
<!-- template/operation.md -->
## {OPERATION_NAME}

### Description
{DESCRIPTION}

### Validations
{VALIDATIONS}

### Example
{EXAMPLE}
```

### 3. Versionamiento
```
docs/
├── v1/
│   └── api-description.md
└── v2/
    └── api-description.md
```

## Cuándo Usar

**Use archivos .md cuando:**
- Description > 5 líneas
- Necesita formateo Markdown
- Quiere reutilizar descriptions
- Colaboración con technical writers

**Use inline cuando:**
- Descriptions cortas
- Sin formateo especial
- Prototyping rápido
