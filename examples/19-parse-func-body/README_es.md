# Ejemplo 19 - Parse Func Body

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de `--parseFuncBody` para extraer tipos del body de funciones.

## Flag

```bash
--parseFuncBody
--pfb
```

O con valor explícito:
```bash
--parseFuncBody=true
--pfb=true
```

> **Nota:** Ambas sintaxis son válidas. Use `--parseFuncBody` o `--pfb` (sin valor), o `--parseFuncBody=true` (explícito). NO use `--parseFuncBody true` (separado por espacio).

Default: `false`

## Concepto

Por defecto, nexs-swag solo parsea comentarios y signatures. Con `--parseFuncBody`, también parsea el body de las funciones para extraer tipos.

## Problema

### SIN --parseFuncBody
```go
// @Summary Create user
// @Success 200 {object} ??? // ¿Qué tipo?
func CreateUser(c *gin.Context) {
    result := &Response{
        User: &User{
            ID:   1,
            Name: "John",
        },
    }
    c.JSON(200, result)  // Tipo no detectado
}
```

### CON --parseFuncBody
```go
// @Summary Create user
// @Success 200 {object} auto  // Detecta automáticamente
func CreateUser(c *gin.Context) {
    result := &Response{  // ← nexs-swag detecta este tipo
        User: &User{
            ID:   1,
            Name: "John",
        },
    }
    c.JSON(200, result)
}
```

## Uso

```bash
nexs-swag init --parseFuncBody
```

## Cómo Ejecutar

```bash
./run.sh
```

## Detección Automática

nexs-swag detecta tipos en:

### 1. Return Statements
```go
func GetUser(c *gin.Context) {
    user := &User{ID: 1, Name: "John"}
    return c.JSON(200, user)  // Detecta User
}
```

### 2. Variable Declarations
```go
func ListUsers(c *gin.Context) {
    result := &UserList{  // Detecta UserList
        Users: []User{},
        Total: 0,
    }
    c.JSON(200, result)
}
```

### 3. Struct Literals
```go
func CreateOrder(c *gin.Context) {
    c.JSON(201, &OrderResponse{  // Detecta OrderResponse
        Order: order,
        Message: "Created",
    })
}
```

### 4. Function Calls
```go
func GetReport(c *gin.Context) {
    report := generateReport()  // Detecta tipo de retorno
    c.JSON(200, report)
}

func generateReport() *Report {
    return &Report{/* ... */}
}
```

## Ventajas

### 1. Menos Redundancia
```go
// SIN --parseFuncBody - Redundante
// @Success 200 {object} Response  // Manual
func GetData(c *gin.Context) {
    c.JSON(200, &Response{})  // Repite el tipo
}

// CON --parseFuncBody - DRY
// @Success 200 {object} auto  // Automático
func GetData(c *gin.Context) {
    c.JSON(200, &Response{})  // Single source of truth
}
```

### 2. Mantiene Sincronización
```go
// Cambió el tipo de retorno
func GetUser(c *gin.Context) {
    // Antes: c.JSON(200, &User{})
    c.JSON(200, &UserDetailed{})  // Cambió
    // nexs-swag detecta automáticamente el nuevo tipo
}
```

### 3. Tipos Complejos
```go
func GetDashboard(c *gin.Context) {
    // Tipo complejo anidado
    c.JSON(200, &DashboardResponse{
        User: &UserInfo{},
        Stats: &Statistics{},
        RecentOrders: []Order{},
    })
    // Detecta DashboardResponse y todas sus dependencias
}
```

## Limitaciones

### 1. Complejidad Performance
```bash
# SIN --parseFuncBody
Parse time: ~1s

# CON --parseFuncBody  
Parse time: ~5s

# ⚠️ Más lento para proyectos grandes
```

### 2. Tipos Dinámicos
```go
// ❌ NO detecta - tipo dinámico
func GetData(c *gin.Context) {
    var result interface{}
    if condition {
        result = &TypeA{}
    } else {
        result = &TypeB{}
    }
    c.JSON(200, result)
}
```

### 3. External Functions
```go
// ❌ NO detecta - función externa
func GetUser(c *gin.Context) {
    user := externalPkg.FetchUser()
    c.JSON(200, user)  // Tipo de retorno en otro package
}
```

### 4. Ambiguedad
```go
// ❌ NO detecta - múltiples returns
func GetData(c *gin.Context) {
    data1 := &TypeA{}
    data2 := &TypeB{}
    
    c.JSON(200, data1)
    c.JSON(200, data2)  // ¿Cuál es el correcto?
}
```

## Best Practices

### 1. Combine con Comments
```go
// ✅ GOOD - Especifica cuando necesario
// @Success 200 {object} UserResponse
// @Success 404 {object} ErrorResponse
func GetUser(c *gin.Context) {
    if user == nil {
        c.JSON(404, &ErrorResponse{})
        return
    }
    c.JSON(200, &UserResponse{})
}
```

### 2. Simple Functions
```go
// ✅ GOOD - Función simple, un tipo
func GetUser(c *gin.Context) {
    user := service.GetUser(id)
    c.JSON(200, user)  // Detecta automáticamente
}

// ❌ BAD - Función compleja, múltiples tipos
func HandleRequest(c *gin.Context) {
    // Mucha lógica...
    // Varios returns...
    // Difícil de detectar
}
```

### 3. Named Returns
```go
// ✅ GOOD - Tipo explícito
func buildResponse() *UserResponse {
    return &UserResponse{/* ... */}
}

func GetUser(c *gin.Context) {
    c.JSON(200, buildResponse())  // Detecta UserResponse
}
```

## Cuándo Usar

**Use --parseFuncBody cuando:**
- Muchos endpoints
- Quiere reducir redundancia
- Tipos cambian frecuentemente
- Performance no es crítico

**NO use cuando:**
- Proyecto grande (performance)
- Tipos dinámicos/complejos
- Prefiere explícito
- CI/CD debe ser rápido

## Alternativas

### 1. Type Aliases
```go
// Define un tipo
type UserResponse = User

// @Success 200 {object} UserResponse
func GetUser(c *gin.Context) {
    c.JSON(200, &User{})
}
```

### 2. Helper Functions
```go
func success(c *gin.Context, data interface{}) {
    c.JSON(200, data)
}

// @Success 200 {object} User
func GetUser(c *gin.Context) {
    success(c, &User{})
}
```

### 3. Code Generation
```go
//go:generate nexs-swag-helper generate

// El helper genera comments automáticamente
func GetUser(c *gin.Context) {
    c.JSON(200, &User{})
}
```

## Performance Tuning

### 1. Combine con --exclude
```bash
nexs-swag init \
  --parseFuncBody \
  --exclude "testdata,mocks,vendor"
```

### 2. Limite Scope
```bash
# Solo directorios específicos
nexs-swag init \
  --parseFuncBody \
  --dir ./api/handlers
```

### 3. Use Cache
```bash
# Algunos tools soportan cache
nexs-swag init \
  --parseFuncBody \
  --cache-dir .swag-cache
```

## Recomendación

**Para la mayoría de proyectos:**
```bash
# No use --parseFuncBody
nexs-swag init

# Sea explícito en comments
// @Success 200 {object} User
```

**Para proyectos pequeños/medianos con tipos estables:**
```bash
# Use --parseFuncBody para conveniencia
nexs-swag init --parseFuncBody
```
