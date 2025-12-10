# Exemplo 09 - Parse Dependency

Demonstra o uso de `--parseDependency` para incluir types de packages importados.

## Flag

```bash
--parseDependency
--pd
```

## Conceito

Este exemplo demonstra como o nexs-swag pode parsear dependências quando você tem types definidos em packages separados. Neste exemplo simplificado, mostramos o conceito com um único arquivo, mas em projetos reais você teria:

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

## Como Funciona

### SEM --parseDependency
Apenas types do package atual são incluídos na documentação.

### COM --parseDependency
Types de packages importados também são parseados e incluídos.

## Estrutura em Projetos Reais

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

## Como Executar

```bash
./run.sh
```

## Quando Usar

**Use --parseDependency quando:**
- Models em packages separados
- Estrutura modular
- Types importados
- Libraries compartilhadas

**NÃO precisa quando:**
- Todos os types no mesmo package
- API simples
- Sem imports de models

## Níveis de Parsing

Combine com `--parseDependencyLevel` para controlar profundidade:

```bash
# Nível 0: Apenas diretório principal
nexs-swag init --parseDependency --parseDependencyLevel 0

# Nível 1: + 1 nível de dependências
nexs-swag init --parseDependency --parseDependencyLevel 1

# Nível 2: + 2 níveis (default)
nexs-swag init --parseDependency --parseDependencyLevel 2
```

## Performance

⚠️ **Atenção:** Parsear muitas dependências pode ser lento.

Otimizações:
```bash
# Apenas o necessário
nexs-swag init --parseDependency --parseDependencyLevel 1

# Limitar com --exclude
nexs-swag init --parseDependency --exclude "vendor,node_modules"
```
