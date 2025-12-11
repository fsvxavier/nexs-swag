# Example 07 - Tag Filtering

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates how to filter endpoints by tags.

## Flag

```bash
--tags tag1,tag2,tag3
-t tag1,tag2
```

## Syntax

### Include Tags
```bash
# Only endpoints with "users" tag
nexs-swag init --tags users

# Endpoints with "users" OR "admin"
nexs-swag init --tags users,admin
```

### Exclude Tags
```bash
# All EXCEPT "internal"
nexs-swag init --tags '!internal'

# All EXCEPT "internal" and "deprecated"
nexs-swag init --tags '!internal,!deprecated'
```

### Combine
```bash
# Only "admin" but NOT "internal"
nexs-swag init --tags admin,!internal
```

## Example Endpoints

| Endpoint | Tags | Included in |
|----------|------|-------------|
| `GET /users` | users | --tags users |
| `POST /users` | users,admin | --tags users OR --tags admin |
| `DELETE /users/{id}` | admin | --tags admin |
| `GET /internal/config` | internal | --tags internal |

## Filter Examples

### 1. Public Documentation (without internals)
```bash
nexs-swag init --tags '!internal'
```

### 2. Admin Documentation
```bash
nexs-swag init --tags admin
```

### 3. Complete Documentation (except deprecated)
```bash
nexs-swag init --tags '!deprecated'
```

### 4. Multiple Versions
```bash
# API v1
nexs-swag init --tags v1 --output docs/v1

# API v2
nexs-swag init --tags v2 --output docs/v2
```

## How to Run

```bash
./run.sh
```

## Use Cases

- **Public vs internal documentation:** Exclude internal endpoints
- **Multiple versions:** Generate separate docs for v1, v2
- **By permission:** user, admin, superadmin
- **By status:** stable, beta, deprecated
