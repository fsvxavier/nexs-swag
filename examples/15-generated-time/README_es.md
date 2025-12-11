# Ejemplo 15 - Generated Time

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el control de timestamp de generación en la spec.

## Flag

```bash
--generatedTime
```

Default: `false`

## Concepto

Por defecto, nexs-swag NO incluye timestamp de generación en la spec.

## Uso

### SIN --generatedTime (Default)
```bash
nexs-swag init
```

```yaml
# docs/swagger.yaml
info:
  title: My API
  version: 1.0.0
# No timestamp
```

### CON --generatedTime
```bash
nexs-swag init --generatedTime
```

```yaml
# docs/swagger.yaml
info:
  title: My API
  version: 1.0.0
  x-generated-at: 2025-12-11T14:30:00Z
```

## Cómo Ejecutar

```bash
./run.sh
```

## Cuándo Usar

### Use --generatedTime cuando:

#### 1. Debugging
```bash
# Saber cuándo fue generada la spec
nexs-swag init --generatedTime

# Check timestamp
grep x-generated-at docs/swagger.yaml
# x-generated-at: 2025-12-11T14:30:00Z
```

#### 2. Version Control
```yaml
# Commit message automático
x-generated-at: 2025-12-11T14:30:00Z

# Git hook puede verificar si spec está actualizada
```

#### 3. CI/CD Tracking
```bash
# Build pipeline
echo "Generating spec..."
nexs-swag init --generatedTime
echo "Generated at: $(grep x-generated-at docs/swagger.yaml)"
```

#### 4. Cache Invalidation
```go
// Servidor puede usar timestamp para invalidar cache
type SpecMetadata struct {
    GeneratedAt time.Time `json:"x-generated-at"`
}

// Client verifica si spec cambió
if spec.GeneratedAt.After(cachedSpec.GeneratedAt) {
    // Reload spec
}
```

### NO use --generatedTime cuando:

#### 1. Version Control Noise
```bash
# Cada generation = nuevo timestamp = diff en git
git diff docs/swagger.yaml
- x-generated-at: 2025-12-11T14:30:00Z
+ x-generated-at: 2025-12-11T14:31:00Z

# Spec no cambió, solo timestamp!
```

#### 2. Deterministic Builds
```bash
# Build reproducible: mismo input = mismo output
nexs-swag init # Sin timestamp
md5sum docs/swagger.yaml
# Always: abc123...

nexs-swag init --generatedTime
md5sum docs/swagger.yaml
# Different cada vez!
```

#### 3. Git Conflicts
```bash
# Merge conflicts innecesarios
<<<<<<< HEAD
x-generated-at: 2025-12-11T14:30:00Z
=======
x-generated-at: 2025-12-11T14:35:00Z
>>>>>>> feature-branch
```

## Best Practices

### 1. Excluir de Git
```bash
# .gitignore
docs/swagger.yaml
docs/swagger.json

# Generar en CI/CD
```

### 2. Separate Metadata
```yaml
# spec.yaml - committed
info:
  title: My API
  version: 1.0.0

# metadata.yaml - generated, not committed
generated_at: 2025-12-11T14:30:00Z
generated_by: nexs-swag v1.0.0
```

### 3. CI/CD Strategy
```yaml
# .github/workflows/docs.yml
name: Generate Docs

on: [push]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Generate spec
        run: nexs-swag init --generatedTime
      
      - name: Upload artifact
        uses: actions/upload-artifact@v2
        with:
          name: swagger-spec
          path: docs/
```

## Formato Timestamp

### ISO 8601
```yaml
x-generated-at: 2025-12-11T14:30:00Z
```

### Components
```
2025-12-11T14:30:00Z
│    │  │ │  │  │ │
│    │  │ │  │  │ └─ UTC
│    │  │ │  │  └─── Seconds
│    │  │ │  └────── Minutes
│    │  │ └───────── Hours
│    │  └──────────── Day
│    └─────────────── Month
└──────────────────── Year
```

## Alternativas

### 1. Version en Info
```go
// @version 1.0.0-20251211
```

```yaml
info:
  version: 1.0.0-20251211
  # Sin x-generated-at
```

### 2. Git Commit
```bash
# Incluir git hash
git rev-parse HEAD > docs/version.txt

# O en description
GIT_HASH=$(git rev-parse --short HEAD)
sed -i "s/VERSION/1.0.0-$GIT_HASH/" docs/swagger.yaml
```

### 3. Build Info
```go
// @version 1.0.0
// @x-build-time 2025-12-11T14:30:00Z
// @x-build-commit abc123
```

## Recomendación

**Default (sin --generatedTime):**
- ✅ Sin noise en git
- ✅ Builds determinísticos
- ✅ Menos conflicts
- ✅ Más limpio

**Use --generatedTime solo si:**
- Realmente necesita timestamp
- Spec no está en git
- Debug/tracking es crítico
