# Example 08 - Parse Internal Packages

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the use of `--parseInternal` to include `internal/` packages.

## Flag

```bash
--parseInternal
```

## Comportamento

### SEM flag (default)
```bash
nexs-swag init
```
- **Ignora** diretórios `internal/`
- Apenas APIs públicas são documentadas

### COM flag
```bash
nexs-swag init --parseInternal
```
- **Inclui** diretórios `internal/`
- APIs internas também são documentadas

## Estrutura

```
08-parse-internal/
├── main.go              # ✅ Sempre parseado
└── internal/
    └── config.go        # ⚠️ Apenas com --parseInternal
```

## Por que usar?

### Convenção Go
Em Go, diretórios `internal/` têm significado especial:
- Código em `internal/` só pode ser importado por pacotes pai
- É uma convenção para código privado/interno

### Casos de Uso

**NÃO usar --parseInternal quando:**
- Documentação pública para consumidores externos
- APIs de biblioteca/pacote público
- Cliente não deve conhecer detalhes internos

**Usar --parseInternal quando:**
- Documentação interna da equipe
- Microserviços internos
- APIs administrativas
- Debugging e desenvolvimento

## Exemplo

```go
// projeto/
// ├── api/
// │   └── public.go        ✅ Sempre documentado
// └── internal/
//     ├── auth/
//     │   └── auth.go      ⚠️ Apenas com flag
//     └── db/
//         └── queries.go   ⚠️ Apenas com flag
```

### Documentação Pública
```bash
nexs-swag init --output ./docs/public
# Apenas api/public.go
```

### Documentação Completa
```bash
nexs-swag init --parseInternal --output ./docs/internal
# api/public.go + internal/auth + internal/db
```

## Como Executar

```bash
./run.sh
```

## Resultado

- **Sem flag:** 1 schema (User), 1 endpoint (/users)
- **Com flag:** 2 schemas (User + Config), 2 endpoints (/users + /internal/config)
