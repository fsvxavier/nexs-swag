# Exemplo 02 - Múltiplos Formatos

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra como gerar documentação em diferentes formatos.

## Flags Utilizadas

- `--format json` - Gera apenas openapi.json
- `--format yaml` - Gera apenas openapi.yaml
- `--format go` - Gera apenas docs.go
- `--format json,yaml,go` - Gera todos os formatos

## Como Executar

```bash
chmod +x run.sh
./run.sh
```

## Formatos Disponíveis

### JSON (`--format json`)
```bash
nexs-swag init --format json
# Gera: docs/openapi.json
```

### YAML (`--format yaml`)
```bash
nexs-swag init --format yaml
# Gera: docs/openapi.yaml
```

### Go (`--format go`)
```bash
nexs-swag init --format go
# Gera: docs/docs.go
```

### Múltiplos
```bash
nexs-swag init --format json,yaml
# Gera: docs/openapi.json + docs/openapi.yaml
```

## Casos de Uso

- **JSON only:** Para servir via HTTP endpoint
- **YAML only:** Para documentação legível por humanos
- **Go only:** Para embedding na aplicação
- **Todos:** Para máxima compatibilidade
