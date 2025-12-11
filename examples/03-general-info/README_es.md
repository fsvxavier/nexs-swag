# Ejemplo 03 - General Info File

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de `--generalInfo` para especificar qué archivo contiene las anotaciones generales de la API.

## Problema

Cuando tienes múltiples archivos Go, el parser puede encontrar anotaciones de información general (@title, @version) en varios lugares, causando conflictos.

## Solución

Usa `--generalInfo` para especificar exactamente qué archivo contiene la información general:

```bash
nexs-swag init --generalInfo main.go
```

## Estructura

```
03-general-info/
├── main.go       # ✅ TIENE @title, @version, @host, etc
├── products.go   # ❌ Solo endpoints de productos
├── orders.go     # ❌ Solo endpoints de órdenes
└── run.sh
```

## Regla

- **Archivo de Info General:** Debe tener @title, @version, @host, @BasePath
- **Otros Archivos:** Deben tener SOLO endpoints (@Router, @Summary, etc)

## Cómo Ejecutar

```bash
chmod +x run.sh
./run.sh
```

## Comparación

### Sin --generalInfo
```bash
nexs-swag init --dir .
# Puede generar error si encuentra @title en múltiples archivos
```

### Con --generalInfo
```bash
nexs-swag init --dir . --generalInfo main.go
# ✅ Correcto: solo main.go es parseado para info general
# ✅ products.go y orders.go proporcionan solo endpoints
```

## Beneficios

1. **Evita conflictos:** Un único lugar para la info de la API
2. **Más rápido:** El parser no necesita verificar todos los archivos para info general
3. **Organización:** Separa responsabilidades (info general vs endpoints)
