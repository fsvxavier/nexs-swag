# Ejemplo 09 - Parse Dependency

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de `--parseDependency` para incluir tipos de paquetes importados.

## Flag

```bash
--parseDependency
--pd
```

## Concepto

Este ejemplo demuestra cómo nexs-swag puede parsear dependencias cuando tienes tipos definidos en paquetes separados. En este ejemplo simplificado, mostramos el concepto con un solo archivo, pero en proyectos reales tendrías:

```
myapp/
├── main.go              # Usa models.Product
└── models/
    └── product.go       # Define Product
```

## Uso

```bash
nexs-swag init --parseDependency
```

## Cómo Funciona

### SIN --parseDependency
Solo los tipos del paquete actual se incluyen en la documentación.

### CON --parseDependency
Los tipos de paquetes importados también se parsean e incluyen.

## Estructura en Proyectos Reales

```go
// main.go
package main

import "myapp/models"

// @Success 200 {object} models.Product
func GetProduct() {}

// models/product.go
package models

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

## Cómo Ejecutar

```bash
./run.sh
```

## Cuándo Usar

**Use --parseDependency cuando:**
- Models en paquetes separados
- Estructura modular
- Tipos importados
- Bibliotecas compartidas

**NO es necesario cuando:**
- Todos los tipos en el mismo paquete
- API simple
- Sin imports de models

## Niveles de Parsing

Combine con `--parseDependencyLevel` para controlar la profundidad:

```bash
# Nivel 0: Solo directorio principal
nexs-swag init --parseDependency --parseDependencyLevel 0

# Nivel 1: + 1 nivel de dependencias
nexs-swag init --parseDependency --parseDependencyLevel 1

# Nivel 2: + 2 niveles (default)
nexs-swag init --parseDependency --parseDependencyLevel 2
```

## Performance

⚠️ **Atención:** Parsear muchas dependencias puede ser lento.

Optimizaciones:
```bash
# Solo lo necesario
nexs-swag init --parseDependency --parseDependencyLevel 1

# Limitar con --exclude
nexs-swag init --parseDependency --exclude "vendor,node_modules"
```
