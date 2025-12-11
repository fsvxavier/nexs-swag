# Ejemplo 11 - Parse GoList

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de `--parseGoList` para parsing más rápido y preciso.

## Flag

```bash
--parseGoList
```

## Diferencia

### SIN --parseGoList (Manual)
- Recorre directorios manualmente
- Parsea todos los archivos .go
- Más lento en proyectos grandes
- Puede perder dependencias

### CON --parseGoList (Go Tooling)
- Usa `go list -json`
- Información de Go modules
- Más rápido y preciso
- Respeta go.mod

## Uso

```bash
nexs-swag init --parseGoList
```

## Cómo Funciona

### Manual Parsing
```bash
# nexs-swag hace:
1. Walk en todos los directorios
2. Encuentra archivos .go
3. Parsea cada archivo
4. Resuelve imports manualmente
```

### Go List Parsing
```bash
# nexs-swag ejecuta:
go list -json ./...

# Retorna:
{
  "ImportPath": "myapp",
  "Dir": "/path/to/myapp",
  "GoFiles": ["main.go", "handlers.go"],
  "Imports": ["encoding/json", "net/http"],
  "Deps": ["myapp/models"]
}
```

## Performance

### Proyecto Pequeño (< 50 archivos)
```
Manual:  ~500ms
GoList:  ~400ms
Ganancia: 20%
```

### Proyecto Mediano (100-500 archivos)
```
Manual:  ~5s
GoList:  ~2s
Ganancia: 60%
```

### Proyecto Grande (> 1000 archivos)
```
Manual:  ~30s
GoList:  ~8s
Ganancia: 73%
```

## Ventajas

### 1. Velocidad
```bash
# 3x más rápido en proyectos grandes
time nexs-swag init --parseGoList
# real: 2s vs 6s sin flag
```

### 2. Precisión
```bash
# Respeta go.mod replace
# go.mod:
replace github.com/old/pkg => ../local/pkg

# nexs-swag usa el path correcto
```

### 3. Build Tags
```go
// +build linux

package linux

// Solo parseado si build tag coincide
```

### 4. Vendor Detection
```bash
# Detecta vendor/ automáticamente
go mod vendor
nexs-swag init --parseGoList
# Ignora vendor/ si modules están habilitados
```

## Cómo Ejecutar

```bash
./run.sh
```

## Cuándo Usar

**Use --parseGoList cuando:**
- Proyecto con Go modules
- Muchos packages
- Dependencias complejas
- Quiere velocidad

**NO use cuando:**
- Proyecto sin go.mod
- Proyecto basado en GOPATH
- Necesita parsear archivos específicos ignorados por go list

## Requisitos

```bash
# Go 1.11+ (modules)
go mod init myapp

# o GOPATH configurado
export GOPATH=/path/to/gopath
```

## Combinar con Otras Flags

### Con parseDependency
```bash
nexs-swag init \
  --parseGoList \
  --parseDependency \
  --parseDependencyLevel 2
```

### Con exclude
```bash
nexs-swag init \
  --parseGoList \
  --exclude "testdata,mocks"
```
