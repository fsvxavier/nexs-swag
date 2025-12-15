# Exemplo 10 - Dependency Level

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra o uso de `--parseDependencyLevel` para controlar profundidade de parsing.

## Flag

```bash
--parseDependencyLevel <0-3>
--pdl <0-3>
```

Default: `0`

Requer: `--parseDependency` (ou `--parseDependency=true`)

## Conceito

Este exemplo mostra types aninhados (nested types) no mesmo arquivo para demonstrar o conceito de níveis de dependência:

```go
type Order struct {
    Items []Item  // Level 1: Order referencia Item
}

type Item struct {
    Metadata Meta  // Level 2: Item referencia Meta
}

type Meta struct {
    CreatedAt string  // Level 3: Type final
}
```

## Níveis

### Level 0 (Default)
Apenas o diretório principal (`--dir`)

```bash
nexs-swag init --parseDependency --parseDependencyLevel 0
```

### Level 1
Principal + 1 nível de dependências

```bash
nexs-swag init --parseDependency --parseDependencyLevel 1
```

### Level 2
Principal + 2 níveis de dependências

```bash
nexs-swag init --parseDependency --parseDependencyLevel 2
```

### Level 3
Principal + 3 níveis de dependências

```bash
nexs-swag init --parseDependency --parseDependencyLevel 3
```

## Estrutura em Projetos Reais

Em projetos com múltiplos packages:

```
main.go
  └── services.Order (Level 1)
        └── models.Item (Level 2)
              └── types.Meta (Level 3)
```

## Comparação

| Level | Parseia | Definitions |
|-------|---------|-------------|
| 0 | main/ | Order apenas |
| 1 | main/ + refs | Order, Item |
| 2 | main/ + refs + refs | Order, Item, Meta |
| 3 | main/ + refs + refs + refs | Todos os types |

## Como Executar

```bash
./run.sh
```

## Quando Usar Cada Nível

### Level 0
```bash
# API simples, types no mesmo package
myapp/
└── main.go  # Todos os types aqui
```

### Level 1
```bash
# Models em subpackage direto
myapp/
├── main.go
└── models/
    └── user.go
```

### Level 2
```bash
# Models com nested types
myapp/
├── main.go
├── services/
│   └── order.go    # Usa models.Item
└── models/
    └── item.go
```

### Level 3
```bash
# Hierarquia profunda
myapp/
├── api/
│   └── handlers.go      # Usa services.Order
├── services/
│   └── order.go         # Usa models.Item
├── models/
│   └── item.go          # Usa types.Meta
└── types/
    └── meta.go
```

## Performance

⚠️ Níveis maiores = parsing mais lento

| Level | Tempo | Arquivos |
|-------|-------|----------|
| 0 | Rápido | ~10 |
| 1 | Normal | ~50 |
| 2 | Lento | ~200 |
| 3 | Muito lento | ~1000+ |

## Otimização

### Combinar com --exclude
```bash
nexs-swag init \
  --parseDependency \
  --parseDependencyLevel 2 \
  --exclude "vendor,testdata,mocks"
```

### Usar --parseGoList
```bash
# Mais rápido para projetos grandes
nexs-swag init \
  --parseDependency \
  --parseDependencyLevel 2 \
  --parseGoList
```

## Recomendações

**Use Level 1 para:**
- Projetos médios
- Models em 1 subpackage
- Performance importante

**Use Level 2 para:**
- Projetos grandes
- Hierarquia moderada
- Balanceamento performance/completude

**Use Level 3 apenas se:**
- Hierarquia muito profunda
- Todas as definitions necessárias
- Performance não é crítica
