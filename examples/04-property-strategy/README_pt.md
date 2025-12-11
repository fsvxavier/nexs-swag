# Exemplo 04 - Property Naming Strategy

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra as diferentes estratégias de naming para campos de struct.

## Flag

```bash
--propertyStrategy <strategy>
-p <strategy>
```

## Estratégias Disponíveis

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

## Regras Importantes

### ✅ Aplicado quando:
- Campo **NÃO** tem tag `json`
- Campo **NÃO** tem `json:"-"`

### ❌ NÃO aplicado quando:
- Campo tem tag `json:"explicit_name"` → usa "explicit_name"
- Campo tem `json:"-"` → ignorado
- Campo tem `json:",omitempty"` → aplica strategy + omitempty

## Exemplo

```go
type User struct {
    UserID    int    `json:"user_id"`      // ✅ SEMPRE "user_id"
    FirstName string                       // ⚠️ Depende da strategy
    LastName  string `json:",omitempty"`   // ⚠️ Strategy + omitempty
    Password  string `json:"-"`            // ❌ Ignorado
}
```

Com `--propertyStrategy snakecase`:
```json
{
    "user_id": 123,
    "first_name": "John",
    "last_name": "Doe"
}
```

Com `--propertyStrategy camelcase`:
```json
{
    "user_id": 123,
    "firstName": "John",
    "lastName": "Doe"
}
```

## Como Executar

```bash
chmod +x run.sh
./run.sh
```

## Casos de Uso

- **snake_case:** APIs Python, Ruby, bancos de dados
- **camelCase:** APIs JavaScript, JSON padrão
- **PascalCase:** APIs C#, .NET
