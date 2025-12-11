# Ejemplo 05 - Required By Default

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el comportamiento del flag `--requiredByDefault` que hace que todos los campos sean requeridos por defecto.

## Flag

```bash
--requiredByDefault
```

## Comportamiento

### SIN el flag (default)
```bash
nexs-swag init
```
- **Todos los campos son OPCIONALES** por defecto
- Solo los campos con `binding:"required"` o `validate:"required"` son requeridos

### CON el flag
```bash
nexs-swag init --requiredByDefault
```
- **Todos los campos son REQUERIDOS** por defecto
- Excepciones:
  - Campos con `json:",omitempty"`
  - Campos con `binding:"omitempty"`  
  - Campos que son punteros (`*Type`)

## Ejemplo

```go
type Product struct {
    ID          int      // ✅ Required (con flag)
    Name        string   // ✅ Required (con flag)
    Description string   `json:"description,omitempty"` // ❌ Optional (omitempty)
    Discount    *float64 // ❌ Optional (pointer)
    Category    string   `json:"category" binding:"omitempty"` // ❌ Optional (binding)
}
```

## Schema Generado

### Sin --requiredByDefault
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

### Con --requiredByDefault
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

## Cómo Ejecutar

```bash
chmod +x run.sh
./run.sh
```

## Casos de Uso

**Use `--requiredByDefault` cuando:**
- La API requiere validación estricta
- La mayoría de los campos son obligatorios
- Prefiere opt-out (marcar opcionales) en lugar de opt-in

**NO use cuando:**
- La API tiene muchos campos opcionales
- Prefiere opt-in explícito con `binding:"required"`
- Compatibilidad con clientes existentes
