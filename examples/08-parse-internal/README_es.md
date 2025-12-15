# Ejemplo 08 - Parse Internal Packages

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de `--parseInternal` para incluir paquetes `internal/`.

## Flag

```bash
--parseInternal
# O explícitamente:
--parseInternal=true
```

> **Nota:** Ambas sintaxis son válidas. Use `--parseInternal` (sin valor) o `--parseInternal=true` (explícito). NO use `--parseInternal true` (separado por espacio).

## Comportamiento

### SIN flag (default)
```bash
nexs-swag init
```
- **Ignora** directorios `internal/`
- Solo APIs públicas son documentadas

### CON flag
```bash
nexs-swag init --parseInternal
```
- **Incluye** directorios `internal/`
- APIs internas también son documentadas

## Estructura

```
08-parse-internal/
├── main.go              # ✅ Siempre parseado
└── internal/
    └── config.go        # ⚠️ Solo con --parseInternal
```

## Por qué usar

### Convención Go
En Go, los directorios `internal/` tienen significado especial:
- Código en `internal/` solo puede ser importado por paquetes padre
- Es una convención para código privado/interno

### Casos de Uso

**NO usar --parseInternal cuando:**
- Documentación pública para consumidores externos
- APIs de biblioteca/paquete público
- Cliente no debe conocer detalles internos

**Usar --parseInternal cuando:**
- Documentación interna del equipo
- Microservicios internos
- APIs administrativas
- Debugging y desarrollo

## Ejemplo

```go
// proyecto/
// ├── api/
// │   └── public.go        ✅ Siempre documentado
// └── internal/
//     ├── auth/
//     │   └── auth.go      ⚠️ Solo con flag
//     └── db/
//         └── queries.go   ⚠️ Solo con flag
```

### Documentación Pública
```bash
nexs-swag init --output ./docs/public
# Solo api/public.go
```

### Documentación Completa
```bash
nexs-swag init --parseInternal --output ./docs/internal
# api/public.go + internal/auth + internal/db
```

## Cómo Ejecutar

```bash
./run.sh
```

## Resultado

- **Sin flag:** 1 schema (User), 1 endpoint (/users)
- **Con flag:** 2 schemas (User + Config), 2 endpoints (/users + /internal/config)
