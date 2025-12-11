# Ejemplo 18 - Collection Format

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de `collectionFormat` para arrays en parámetros.

## Tag

```go
// @Param items query []string false "Items" collectionFormat(csv)
// @Param tags query []string false "Tags" collectionFormat(multi)
```

## Formatos

| Format | Description | Example |
|--------|-------------|---------|
| `csv` | Comma separated | `?ids=1,2,3` |
| `ssv` | Space separated | `?ids=1 2 3` |
| `tsv` | Tab separated | `?ids=1\t2\t3` |
| `pipes` | Pipe separated | `?ids=1\|2\|3` |
| `multi` | Multiple params | `?ids=1&ids=2&ids=3` |

## Default

### OpenAPI 2.0 (Swagger)
Default: `csv`

### OpenAPI 3.0
- Query params: `form` (similar a csv)
- Path/Header: `simple`

## Uso

```go
// CSV (default)
// @Param ids query []int false "IDs" collectionFormat(csv)
// GET /users?ids=1,2,3

// Multi
// @Param ids query []int false "IDs" collectionFormat(multi)  
// GET /users?ids=1&ids=2&ids=3

// Pipes
// @Param status query []string false "Status" collectionFormat(pipes)
// GET /orders?status=pending|approved|shipped

// Space
// @Param tags query []string false "Tags" collectionFormat(ssv)
// GET /posts?tags=golang go swagger
```

## Cómo Ejecutar

```bash
./run.sh
```

## Ejemplos por Tipo

### Query Parameters

#### CSV
```go
// @Summary List users
// @Param ids query []int false "User IDs" collectionFormat(csv)
// @Router /users [get]
func ListUsers(c *gin.Context) {
    // GET /users?ids=1,2,3
    ids := c.QueryArray("ids")  // ["1,2,3"]
    // O
    ids := strings.Split(c.Query("ids"), ",")  // ["1", "2", "3"]
}
```

#### Multi
```go
// @Summary List users  
// @Param ids query []int false "User IDs" collectionFormat(multi)
// @Router /users [get]
func ListUsers(c *gin.Context) {
    // GET /users?ids=1&ids=2&ids=3
    ids := c.QueryArray("ids")  // ["1", "2", "3"]
}
```

#### Pipes
```go
// @Summary Filter orders
// @Param status query []string false "Status" collectionFormat(pipes)
// @Router /orders [get]
func FilterOrders(c *gin.Context) {
    // GET /orders?status=pending|approved|shipped
    statuses := strings.Split(c.Query("status"), "|")
}
```

### OpenAPI Output

#### CSV
```yaml
parameters:
  - name: ids
    in: query
    type: array
    items:
      type: integer
    collectionFormat: csv
# URL: ?ids=1,2,3
```

#### Multi
```yaml
parameters:
  - name: ids
    in: query
    type: array
    items:
      type: integer
    collectionFormat: multi
# URL: ?ids=1&ids=2&ids=3
```

## Parsing no Backend

### Gin Framework

```go
// Multi format
ids := c.QueryArray("ids")
// GET /users?ids=1&ids=2&ids=3
// Result: ["1", "2", "3"]

// CSV format
idsStr := c.Query("ids")
ids := strings.Split(idsStr, ",")
// GET /users?ids=1,2,3
// Result: ["1", "2", "3"]

// Pipes format
statusStr := c.Query("status")
statuses := strings.Split(statusStr, "|")
// GET /orders?status=pending|approved
// Result: ["pending", "approved"]
```

### Echo Framework

```go
// Multi format
ids := c.QueryParams()["ids"]
// ["1", "2", "3"]

// CSV format
idsStr := c.QueryParam("ids")
ids := strings.Split(idsStr, ",")
```

## Comparación

### CSV vs Multi

**CSV Advantages:**
- ✅ URL más corta
- ✅ Fácil de leer
- ✅ Standard en APIs REST

**CSV Disadvantages:**
- ❌ Problemas con valores que contienen `,`
- ❌ Necesita encoding para special chars

**Multi Advantages:**
- ✅ Sin ambiguedad
- ✅ Standard HTTP
- ✅ Fácil parsing

**Multi Disadvantages:**
- ❌ URL más larga
- ❌ Más verboso

### Ejemplo

```bash
# CSV - Compacto
GET /api/users?ids=1,2,3,4,5

# Multi - Verboso
GET /api/users?ids=1&ids=2&ids=3&ids=4&ids=5

# CSV - Problema con comas
GET /api/search?tags=golang,go # ¿1 o 2 tags?

# Multi - Sin ambiguedad
GET /api/search?tags=golang&tags=go
```

## Client Examples

### JavaScript/TypeScript
```typescript
// CSV
const ids = [1, 2, 3];
const url = `/users?ids=${ids.join(',')}`;
// /users?ids=1,2,3

// Multi
const params = new URLSearchParams();
ids.forEach(id => params.append('ids', id));
const url = `/users?${params}`;
// /users?ids=1&ids=2&ids=3
```

### Python
```python
# CSV
ids = [1, 2, 3]
url = f"/users?ids={','.join(map(str, ids))}"
# /users?ids=1,2,3

# Multi
import urllib.parse
params = [('ids', id) for id in ids]
url = f"/users?{urllib.parse.urlencode(params)}"
# /users?ids=1&ids=2&ids=3
```

### cURL
```bash
# CSV
curl "https://api.example.com/users?ids=1,2,3"

# Multi
curl "https://api.example.com/users?ids=1&ids=2&ids=3"

# Pipes
curl "https://api.example.com/orders?status=pending|approved"
```

## OpenAPI 3.0

### Style Parameter
```yaml
# OpenAPI 3.0 usa 'style' em vez de 'collectionFormat'
parameters:
  - name: ids
    in: query
    schema:
      type: array
      items:
        type: integer
    style: form        # Similar a csv
    explode: false     # ?ids=1,2,3

  - name: ids
    in: query
    schema:
      type: array
      items:
        type: integer
    style: form
    explode: true      # ?ids=1&ids=2&ids=3 (similar a multi)
```

## Recomendaciones

**Use CSV cuando:**
- Arrays pequeños
- Valores simples (numbers, simple strings)
- Quiere URLs cortas
- Standard REST API

**Use Multi cuando:**
- Valores complejos
- Necesita precisión
- Standard HTTP forms
- Backend parsea automáticamente

**Use Pipes cuando:**
- Valores pueden contener comas
- Filtros complejos
- Compatibilidad con sistemas legacy

**Use SSV/TSV raramente:**
- Casos muy específicos
- Compatibilidad con sistemas antiguos
