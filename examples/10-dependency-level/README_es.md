# Ejemplo 10 - Dependency Level

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de `--parseDependencyLevel` para controlar la profundidad del parsing.

## Flag

```bash
--parseDependencyLevel <0-3>
--pdl <0-3>
```

Default: `0`

Requiere: `--parseDependency` (o `--parseDependency=true`)

## Concepto

Este ejemplo muestra tipos anidados en el mismo archivo para demostrar el concepto de niveles de dependencia:

```go
type Order struct {
    Items []Item  // Level 1: Order referencia Item
}

type Item struct {
    Metadata Meta  // Level 2: Item referencia Meta
}

type Meta struct {
    CreatedAt string  // Level 3: Tipo final
}
```

## Niveles

### Level 0 (Default)
Solo el directorio principal (`--dir`)

```bash
nexs-swag init --parseDependency --parseDependencyLevel 0
```

### Level 1
Principal + 1 nivel de dependencias

```bash
nexs-swag init --parseDependency --parseDependencyLevel 1
```

### Level 2
Principal + 2 niveles de dependencias

```bash
nexs-swag init --parseDependency --parseDependencyLevel 2
```

### Level 3
Principal + 3 niveles de dependencias

```bash
nexs-swag init --parseDependency --parseDependencyLevel 3
```

## Estructura en Proyectos Reales

En proyectos con múltiples packages:

```
main.go
  └── services.Order (Level 1)
        └── models.Item (Level 2)
              └── types.Meta (Level 3)
```

## Comparación

| Level | Parsea | Definitions |
|-------|--------|-------------|
| 0 | main/ | Order solamente |
| 1 | main/ + refs | Order, Item |
| 2 | main/ + refs + refs | Order, Item, Meta |
| 3 | main/ + refs + refs + refs | Todos los tipos |

## Cómo Ejecutar

```bash
./run.sh
```

## Cuándo Usar Cada Nivel

### Level 0
```bash
# API simple, tipos en el mismo package
myapp/
└── main.go  # Todos los tipos aquí
```

### Level 1
```bash
# Models en subpackage directo
myapp/
├── main.go
└── models/
    └── user.go
```

### Level 2
```bash
# Models con tipos anidados
myapp/
├── main.go
├── services/
│   └── order.go    # Usa models.Item
└── models/
    └── item.go
```

### Level 3
```bash
# Jerarquía profunda
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

⚠️ Niveles mayores = parsing más lento

| Level | Tiempo | Archivos |
|-------|--------|----------|
| 0 | Rápido | ~10 |
| 1 | Normal | ~50 |
| 2 | Lento | ~200 |
| 3 | Muy lento | ~1000+ |

## Optimización

### Combinar con --exclude
```bash
nexs-swag init \
  --parseDependency \
  --parseDependencyLevel 2 \
  --exclude "vendor,testdata,mocks"
```

### Usar --parseGoList
```bash
# Más rápido para proyectos grandes
nexs-swag init \
  --parseDependency \
  --parseDependencyLevel 2 \
  --parseGoList
```

## Recomendaciones

**Use Level 1 para:**
- Proyectos medianos
- Models en 1 subpackage
- Performance importa

**Use Level 2 para:**
- Proyectos grandes
- Jerarquía moderada
- Balance performance/completitud

**Use Level 3 solo si:**
- Jerarquía muy profunda
- Todas las definitions necesarias
- Performance no es crítico
