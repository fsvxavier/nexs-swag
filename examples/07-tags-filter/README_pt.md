# Exemplo 07 - Tag Filtering

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra como filtrar endpoints por tags.

## Flag

```bash
--tags tag1,tag2,tag3
-t tag1,tag2
```

## Sintaxe

### Incluir Tags
```bash
# Apenas endpoints com tag "users"
nexs-swag init --tags users

# Endpoints com "users" OU "admin"
nexs-swag init --tags users,admin
```

### Excluir Tags
```bash
# Todos EXCETO "internal"
nexs-swag init --tags '!internal'

# Todos EXCETO "internal" e "deprecated"
nexs-swag init --tags '!internal,!deprecated'
```

### Combinar
```bash
# Apenas "admin" mas NÃO "internal"
nexs-swag init --tags admin,!internal
```

## Endpoints do Exemplo

| Endpoint | Tags | Incluído em |
|----------|------|-------------|
| `GET /users` | users | --tags users |
| `POST /users` | users,admin | --tags users OU --tags admin |
| `DELETE /users/{id}` | admin | --tags admin |
| `GET /internal/config` | internal | --tags internal |

## Exemplos de Filtros

### 1. Documentação Pública (sem internals)
```bash
nexs-swag init --tags '!internal'
```

### 2. Documentação Admin
```bash
nexs-swag init --tags admin
```

### 3. Documentação Completa (exceto deprecated)
```bash
nexs-swag init --tags '!deprecated'
```

### 4. Múltiplas Versões
```bash
# API v1
nexs-swag init --tags v1 --output docs/v1

# API v2
nexs-swag init --tags v2 --output docs/v2
```

## Como Executar

```bash
./run.sh
```

## Casos de Uso

- **Documentação pública vs interna:** Excluir endpoints internos
- **Múltiplas versões:** Gerar docs separadas para v1, v2
- **Por permissão:** user, admin, superadmin
- **Por status:** stable, beta, deprecated
