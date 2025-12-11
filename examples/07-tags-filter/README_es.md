# Ejemplo 07 - Tag Filtering

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra cómo filtrar endpoints por tags.

## Flag

```bash
--tags tag1,tag2,tag3
-t tag1,tag2
```

## Sintaxis

### Incluir Tags
```bash
# Solo endpoints con tag "users"
nexs-swag init --tags users

# Endpoints con "users" O "admin"
nexs-swag init --tags users,admin
```

### Excluir Tags
```bash
# Todos EXCEPTO "internal"
nexs-swag init --tags '!internal'

# Todos EXCEPTO "internal" y "deprecated"
nexs-swag init --tags '!internal,!deprecated'
```

### Combinar
```bash
# Solo "admin" pero NO "internal"
nexs-swag init --tags admin,!internal
```

## Endpoints del Ejemplo

| Endpoint | Tags | Incluido en |
|----------|------|-------------|
| `GET /users` | users | --tags users |
| `POST /users` | users,admin | --tags users O --tags admin |
| `DELETE /users/{id}` | admin | --tags admin |
| `GET /internal/config` | internal | --tags internal |

## Ejemplos de Filtros

### 1. Documentación Pública (sin internals)
```bash
nexs-swag init --tags '!internal'
```

### 2. Documentación Admin
```bash
nexs-swag init --tags admin
```

### 3. Documentación Completa (excepto deprecated)
```bash
nexs-swag init --tags '!deprecated'
```

### 4. Múltiples Versiones
```bash
# API v1
nexs-swag init --tags v1 --output docs/v1

# API v2
nexs-swag init --tags v2 --output docs/v2
```

## Cómo Ejecutar

```bash
./run.sh
```

## Casos de Uso

- **Documentación pública vs interna:** Excluir endpoints internos
- **Múltiples versiones:** Generar docs separadas para v1, v2
- **Por permiso:** user, admin, superadmin
- **Por estado:** stable, beta, deprecated
