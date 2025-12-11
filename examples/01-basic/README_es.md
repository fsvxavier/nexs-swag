# Ejemplo 01 - Básico

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso básico de nexs-swag con las flags esenciales.

## Flags Utilizadas

- `--dir .` - Directorio con código Go
- `--output ./docs` - Directorio de salida

## Estructura

```
01-basic/
├── main.go      # API simple con 2 endpoints
├── run.sh       # Script de ejecución
└── README.md    # Este archivo
```

**Nota:** Este ejemplo usa el go.mod de la raíz del proyecto.

## Cómo Ejecutar

```bash
./run.sh
```

## Lo que se Genera

1. **docs/openapi.json** - Especificación OpenAPI en JSON
2. **docs/openapi.yaml** - Especificación OpenAPI en YAML
3. **docs/docs.go** - Código Go con la especificación

## API Endpoints

- `GET /api/v1/users/{id}` - Obtener usuario
- `POST /api/v1/users` - Crear usuario

## Probar la API

```bash
# Ejecutar el servidor
go run main.go

# En otra terminal
curl http://localhost:8080/api/v1/users/1
```
