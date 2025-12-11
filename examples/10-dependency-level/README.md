# Example 10 - Dependency Level

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the use of `--parseDependencyLevel` to control parsing depth.

## Flag

```bash
--parseDependencyLevel <0-3>
--pdl <0-3>
```

Default: `0`

Requires: `--parseDependency`

## Concept

This example shows nested types in the same file to demonstrate the concept of dependency levels:

```go
type Order struct {
    Items []Item  // Level 1: Order references Item
}

type Item struct {
    Metadata Meta  // Level 2: Item references Meta
}

type Meta struct {
    CreatedAt string  // Level 3: Final type
}
```

## Levels

### Level 0 (Default)
Only the main directory (`--dir`)

```bash
nexs-swag init --parseDependency --parseDependencyLevel 0
```

### Level 1
Main + 1 dependency level

```bash
nexs-swag init --parseDependency --parseDependencyLevel 1
```

### Level 2
Main + 2 dependency levels

```bash
nexs-swag init --parseDependency --parseDependencyLevel 2
```

### Level 3
Main + 3 dependency levels

```bash
nexs-swag init --parseDependency --parseDependencyLevel 3
```

## Structure in Real Projects

In projects with multiple packages:

```
main.go
  └── services.Order (Level 1)
        └── models.Item (Level 2)
              └── types.Meta (Level 3)
```

## Comparison

| Level | Parses | Definitions |
|-------|--------|-------------|
| 0 | main/ | Order only |
| 1 | main/ + refs | Order, Item |
| 2 | main/ + refs + refs | Order, Item, Meta |
| 3 | main/ + refs + refs + refs | All types |

## How to Run

```bash
./run.sh
```

## When to Use Each Level

### Level 0
```bash
# Simple API, types in same package
myapp/
└── main.go  # All types here
```

### Level 1
```bash
# Models in direct subpackage
myapp/
├── main.go
└── models/
    └── user.go
```

### Level 2
```bash
# Models with nested types
myapp/
├── main.go
├── services/
│   └── order.go    # Uses models.Item
└── models/
    └── item.go
```

### Level 3
```bash
# Deep hierarchy
myapp/
├── api/
│   └── handlers.go      # Uses services.Order
├── services/
│   └── order.go         # Uses models.Item
├── models/
│   └── item.go          # Uses types.Meta
└── types/
    └── meta.go
```

## Performance

⚠️ Higher levels = slower parsing

| Level | Time | Files |
|-------|------|-------|
| 0 | Fast | ~10 |
| 1 | Normal | ~50 |
| 2 | Slow | ~200 |
| 3 | Very slow | ~1000+ |

## Optimization

### Combine with --exclude
```bash
nexs-swag init \
  --parseDependency \
  --parseDependencyLevel 2 \
  --exclude "vendor,testdata,mocks"
```

### Use --parseGoList
```bash
# Faster for large projects
nexs-swag init \
  --parseDependency \
  --parseDependencyLevel 2 \
  --parseGoList
```

## Recommendations

**Use Level 1 for:**
- Medium projects
- Models in 1 subpackage
- Performance matters

**Use Level 2 for:**
- Large projects
- Moderate hierarchy
- Balance performance/completeness

**Use Level 3 only if:**
- Very deep hierarchy
- All definitions needed
- Performance not critical
