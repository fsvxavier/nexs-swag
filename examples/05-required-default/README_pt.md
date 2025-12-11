# Exemplo 05 - Required By Default

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra o comportamento da flag `--requiredByDefault` que torna todos os campos required por padrão.

## Flag

```bash
--requiredByDefault
```

## Comportamento

### SEM a flag (default)
```bash
nexs-swag init
```
- **Todos os campos são OPTIONAL** por padrão
- Apenas campos com `binding:"required"` ou `validate:"required"` são required

### COM a flag
```bash
nexs-swag init --requiredByDefault
```
- **Todos os campos são REQUIRED** por padrão
- Exceções:
  - Campos com `json:",omitempty"`
  - Campos com `binding:"omitempty"`  
  - Campos que são ponteiros (`*Type`)

## Exemplo

```go
type Product struct {
    ID          int      // ✅ Required (com flag)
    Name        string   // ✅ Required (com flag)
    Description string   `json:"description,omitempty"` // ❌ Optional (omitempty)
    Discount    *float64 // ❌ Optional (pointer)
    Category    string   `json:"category" binding:"omitempty"` // ❌ Optional (binding)
}
```

## Schema Gerado

### Sem --requiredByDefault
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

### Com --requiredByDefault
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

## Como Executar

```bash
chmod +x run.sh
./run.sh
```

## Casos de Uso

**Use `--requiredByDefault` quando:**
- API exige validação rigorosa
- Maioria dos campos são obrigatórios
- Prefere opt-out (marcar opcionais) ao invés de opt-in

**NÃO use quando:**
- API tem muitos campos opcionais
- Prefere opt-in explícito com `binding:"required"`
- Compatibilidade com clientes existentes
