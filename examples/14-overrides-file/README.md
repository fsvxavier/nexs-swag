# Exemplo 14 - Overrides File

Demonstra o uso de arquivo `.swaggo` para overrides globais de tipos.

## Flag

```bash
--overridesFile <path>
```

Default: `.swaggo`

## Uso

```bash
nexs-swag init --overridesFile .swaggo
```

## Arquivo .swaggo

```json
{
  "types": {
    "database/sql.NullInt64": "integer",
    "database/sql.NullString": "string",
    "database/sql.NullBool": "boolean",
    "database/sql.NullFloat64": "number",
    "github.com/google/uuid.UUID": "string",
    "github.com/shopspring/decimal.Decimal": "number"
  }
}
```

## Problema

Tipos customizados do Go geram schemas complexos:

### SEM Overrides
```json
{
  "Account": {
    "properties": {
      "id": {
        "type": "object",
        "properties": {
          "Int64": {"type": "integer"},
          "Valid": {"type": "boolean"}
        }
      }
    }
  }
}
```

### COM Overrides
```json
{
  "Account": {
    "properties": {
      "id": {
        "type": "integer"
      }
    }
  }
}
```

## Tipos Comuns

### SQL Nullable Types
```json
{
  "database/sql.NullInt64": "integer",
  "database/sql.NullString": "string",
  "database/sql.NullBool": "boolean",
  "database/sql.NullFloat64": "number",
  "database/sql.NullTime": "string"
}
```

### UUID Libraries
```json
{
  "github.com/google/uuid.UUID": "string",
  "github.com/satori/go.uuid.UUID": "string"
}
```

### Decimal Libraries
```json
{
  "github.com/shopspring/decimal.Decimal": "number",
  "gopkg.in/inf.v0.Dec": "number"
}
```

### Custom Time Types
```json
{
  "myapp/types.Timestamp": "integer",
  "myapp/types.UnixTime": "integer"
}
```

## Alternativa: swaggertype Tag

Sem arquivo global:
```go
type Account struct {
    ID sql.NullInt64 `json:"id" swaggertype:"integer"`
}
```

Com arquivo global:
```go
type Account struct {
    ID sql.NullInt64 `json:"id"` // Automaticamente integer
}
```

## Prioridade

1. **swaggertype tag** (maior prioridade)
2. **Overrides file**
3. **Auto-detection** (menor prioridade)

```go
type User struct {
    ID        sql.NullInt64 `swaggertype:"string"` // string (tag)
    AccountID sql.NullInt64                         // integer (override)
    Balance   float64                               // number (auto)
}
```

## Como Executar

```bash
./run.sh
```

## Benefícios

- **Global:** Um lugar para todos os overrides
- **Reutilizável:** Compartilhar entre projetos
- **Limpo:** Código sem tags repetitivas
- **Versionável:** Commitar junto com código

## Casos de Uso

- Projetos com database/sql
- Bibliotecas de terceiros
- Tipos customizados da empresa
- Padronização de tipos
