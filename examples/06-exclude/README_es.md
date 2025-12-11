# Ejemplo 06 - Exclude Patterns

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra cómo excluir directorios y archivos del parsing.

## Flag

```bash
--exclude pattern1,pattern2,pattern3
```

## Uso

```bash
# Excluir un directorio
nexs-swag init --exclude mock

# Excluir múltiples
nexs-swag init --exclude mock,testdata,vendor

# Excluir con wildcards
nexs-swag init --exclude "*.test.go,*_mock.go"
```

## Exclusiones Automáticas

Siempre excluidos (no es necesario especificar):
- `vendor/` - Dependencias
- `testdata/` - Datos de prueba
- `docs/` - Documentación generada
- `.git/` - Repositorio Git
- `*_test.go` - Archivos de prueba

## Estructura del Ejemplo

```
06-exclude/
├── main.go           # ✅ Será parseado
├── main_test.go      # ❌ Excluido (test)
├── mock/
│   └── mock.go       # ❌ Excluido (con flag)
└── testdata/
    └── data.go       # ❌ Excluido (automático)
```

## Cómo Ejecutar

```bash
./run.sh
```

## Casos de Uso

- **mock:** Código de mocking para pruebas
- **testdata:** Fixtures y datos de prueba
- **vendor:** Dependencias (si usa vendor)
- **examples:** Código de ejemplo
- **internal:** Paquetes internos (use --parseInternal para incluir)
