# Exemplo 16 - Instance Name

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra o uso de `--instanceName` para customizar o nome da instância do Swagger.

## Flag

```bash
--instanceName <name>
--in <name>
```

Default: `swagger`

## Uso

```bash
nexs-swag init --instanceName customapi
```

## Comportamento

### Package Name

**Default:**
```go
package swagger
```

**Custom:**
```go
package customapi
```

### Registração

**Default:**
```go
func init() {
    swagger.SwaggerInfo.Title = "My API"
}
```

**Custom:**
```go
func init() {
    customapi.SwaggerInfo.Title = "My API"
}
```

## Múltiplas Instâncias

Gerar documentação para múltiplas APIs no mesmo projeto:

```bash
# API v1
nexs-swag init \
  --dir ./api/v1 \
  --output ./docs/v1 \
  --instanceName apiv1

# API v2
nexs-swag init \
  --dir ./api/v2 \
  --output ./docs/v2 \
  --instanceName apiv2
```

### Código

```go
package main

import (
    _ "myapp/docs/v1"  // Registra apiv1
    _ "myapp/docs/v2"  // Registra apiv2
    
    "github.com/swaggo/http-swagger/v2"
)

func main() {
    // API v1
    http.HandleFunc("/v1/swagger/", httpSwagger.Handler(
        httpSwagger.InstanceName("apiv1"),
    ))
    
    // API v2
    http.HandleFunc("/v2/swagger/", httpSwagger.Handler(
        httpSwagger.InstanceName("apiv2"),
    ))
    
    http.ListenAndServe(":8080", nil)
}
```

## Estrutura

```
myproject/
├── api/
│   ├── v1/
│   │   └── main.go      # @title API v1
│   └── v2/
│       └── main.go      # @title API v2
└── docs/
    ├── v1/
    │   └── docs.go      # package apiv1
    └── v2/
        └── docs.go      # package apiv2
```

## Como Executar

```bash
./run.sh
```

## Casos de Uso

### 1. Versionamento de API
```bash
--instanceName apiv1  # Versão 1
--instanceName apiv2  # Versão 2
```

### 2. Múltiplos Serviços
```bash
--instanceName users      # Serviço de usuários
--instanceName products   # Serviço de produtos
```

### 3. Ambientes
```bash
--instanceName prod   # Produção
--instanceName dev    # Desenvolvimento
```

### 4. Modularização
```bash
--instanceName public   # API pública
--instanceName admin    # API administrativa
```

## Acesso no Browser

```bash
# Default instance
http://localhost:8080/swagger/

# Custom instance v1
http://localhost:8080/v1/swagger/

# Custom instance v2
http://localhost:8080/v2/swagger/
```

## Benefícios

- **Isolamento:** Cada API com sua documentação
- **Versionamento:** Múltiplas versões simultâneas
- **Organização:** Código e docs organizados
- **Flexibilidade:** Várias APIs no mesmo processo
