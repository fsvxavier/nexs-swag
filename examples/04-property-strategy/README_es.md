# Ejemplo 04 - Property Naming Strategy

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra las diferentes estrategias de nomenclatura para campos de struct.

## Flag

```bash
--propertyStrategy <strategy>
-p <strategy>
```

## Estrategias Disponibles

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

## Reglas Importantes

### ✅ Aplicado cuando:
- El campo **NO** tiene tag `json`
- El campo **NO** tiene `json:"-"`

### ❌ NO aplicado cuando:
- El campo tiene tag `json:"explicit_name"` → usa "explicit_name"
- El campo tiene `json:"-"` → ignorado
- El campo tiene `json:",omitempty"` → aplica strategy + omitempty

## Ejemplo

```go
type User struct {
    UserID    int    `json:"user_id"`      // ✅ SIEMPRE "user_id"
    FirstName string                       // ⚠️ Depende de la estrategia
    LastName  string `json:",omitempty"`   // ⚠️ Strategy + omitempty
    Password  string `json:"-"`            // ❌ Ignorado
}
```

Con `--propertyStrategy snakecase`:
```json
{
    "user_id": 123,
    "first_name": "John",
    "last_name": "Doe"
}
```

Con `--propertyStrategy camelcase`:
```json
{
    "user_id": 123,
    "firstName": "John",
    "lastName": "Doe"
}
```

## Cómo Ejecutar

```bash
chmod +x run.sh
./run.sh
```

## Casos de Uso

- **snake_case:** APIs Python, Ruby, bases de datos
- **camelCase:** APIs JavaScript, JSON estándar
- **PascalCase:** APIs C#, .NET
