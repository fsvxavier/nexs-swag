# Example 02 - Multiple Formats

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates how to generate documentation in different formats.

## Flags Used

- `--format json` - Generates only openapi.json
- `--format yaml` - Generates only openapi.yaml
- `--format go` - Generates only docs.go
- `--format json,yaml,go` - Generates all formats

## How to Run

```bash
chmod +x run.sh
./run.sh
```

## Available Formats

### JSON (`--format json`)
```bash
nexs-swag init --format json
# Generates: docs/openapi.json
```

### YAML (`--format yaml`)
```bash
nexs-swag init --format yaml
# Generates: docs/openapi.yaml
```

### Go (`--format go`)
```bash
nexs-swag init --format go
# Generates: docs/docs.go
```

### Multiple
```bash
nexs-swag init --format json,yaml
# Generates: docs/openapi.json + docs/openapi.yaml
```

## Use Cases

- **JSON only:** To serve via HTTP endpoint
- **YAML only:** For human-readable documentation
- **Go only:** For embedding in application
- **All:** For maximum compatibility
