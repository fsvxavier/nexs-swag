# Ejemplo 02 - Múltiples Formatos

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra cómo generar documentación en diferentes formatos.

## Flags Utilizadas

- `--format json` - Genera solo openapi.json
- `--format yaml` - Genera solo openapi.yaml
- `--format go` - Genera solo docs.go
- `--format json,yaml,go` - Genera todos los formatos

## Cómo Ejecutar

```bash
chmod +x run.sh
./run.sh
```

## Formatos Disponibles

### JSON (`--format json`)
```bash
nexs-swag init --format json
# Genera: docs/openapi.json
```

### YAML (`--format yaml`)
```bash
nexs-swag init --format yaml
# Genera: docs/openapi.yaml
```

### Go (`--format go`)
```bash
nexs-swag init --format go
# Genera: docs/docs.go
```

### Múltiples
```bash
nexs-swag init --format json,yaml
# Genera: docs/openapi.json + docs/openapi.yaml
```

## Casos de Uso

- **JSON only:** Para servir vía endpoint HTTP
- **YAML only:** Para documentación legible por humanos
- **Go only:** Para incrustar en la aplicación
- **Todos:** Para máxima compatibilidad
