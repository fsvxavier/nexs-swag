# Examples - nexs-swag

Este diretório contém exemplos de uso para cada flag e funcionalidade do nexs-swag.

## Pré-requisitos

Instale o nexs-swag antes de executar os exemplos:

```bash
# Do diretório raiz do projeto
cd ..
go install ./cmd/nexs-swag

# Ou use o script de instalação
./install.sh

# Verificar instalação
nexs-swag --version
```

## Estrutura

Cada subdiretório contém um exemplo específico com:
- `main.go` - Código Go com annotations Swagger
- `README.md` - Instruções detalhadas de uso
- `run.sh` - Script para executar o exemplo

## Lista de Exemplos

### Básicos (01-08)
- [01-basic](./01-basic) - Uso básico com `--dir` e `--output`
- [02-formats](./02-formats) - Múltiplos formatos com `--format`
- [03-general-info](./03-general-info) - Arquivo específico com `--generalInfo`
- [04-property-strategy](./04-property-strategy) - `--propertyStrategy` (snake_case, camelCase, PascalCase)
- [05-required-default](./05-required-default) - `--requiredByDefault`
- [06-exclude](./06-exclude) - `--exclude` para excluir diretórios
- [07-tags-filter](./07-tags-filter) - `--tags` para filtrar por tags
- [08-parse-internal](./08-parse-internal) - `--parseInternal`

### Dependências (09-11)
- [09-parse-dependency](./09-parse-dependency) - `--parseDependency`
- [10-dependency-level](./10-dependency-level) - `--parseDependencyLevel` (0-3)
- [11-parse-golist](./11-parse-golist) - `--parseGoList`

### Conteúdo Externo (12-14)
- [12-markdown-files](./12-markdown-files) - `--markdownFiles`
- [13-code-examples](./13-code-examples) - `--codeExampleFilesDir`
- [14-overrides-file](./14-overrides-file) - `--overridesFile`

### Configurações (15-18)
- [15-generated-time](./15-generated-time) - `--generatedTime`
- [16-instance-name](./16-instance-name) - `--instanceName`
- [17-template-delims](./17-template-delims) - `--templateDelims`
- [18-collection-format](./18-collection-format) - `--collectionFormat`

### Avançados (19-21)
- [19-parse-func-body](./19-parse-func-body) - `--parseFuncBody`
- [20-fmt-command](./20-fmt-command) - Comando `fmt`
- [21-struct-tags](./21-struct-tags) - swaggertype, swaggerignore, extensions

## Como Usar

### Executar um exemplo específico

```bash
cd 01-basic
./run.sh
```

### Executar manualmente

```bash
cd 01-basic
nexs-swag init --dir . --output ./docs
```

### Executar todos os exemplos

```bash
for dir in */; do
    echo "=== Executando $dir ==="
    cd "$dir"
    ./run.sh
    cd ..
    echo ""
done
```

## Estrutura de Cada Exemplo

```
XX-nome-exemplo/
├── main.go          # Servidor HTTP com annotations
├── run.sh           # Script de demonstração
└── README.md        # Documentação completa
```

## Dicas

### Visualizar documentação gerada

```bash
# JSON
cat docs/openapi.json | jq

# YAML
cat docs/openapi.yaml

# Docs Go
cat docs/docs.go
```

### Servir com Swagger UI

```bash
# Instalar swagger ui
docker run -p 8080:8080 \
  -e SWAGGER_JSON=/docs/openapi.json \
  -v $(pwd)/docs:/docs \
  swaggerapi/swagger-ui

# Acessar: http://localhost:8080
```

### Integrar em projetos

```go
package main

import (
    "net/http"
    
    _ "myapp/docs"  // Importar docs gerados
    
    httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
    http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Troubleshooting

### Erro "nexs-swag: command not found"

```bash
# Verificar instalação
which nexs-swag

# Se não estiver instalado
cd ..
go install ./cmd/nexs-swag

# Verificar se $GOPATH/bin está no PATH
echo $PATH | grep $(go env GOPATH)/bin

# Adicionar ao PATH se necessário
export PATH=$PATH:$(go env GOPATH)/bin
```

### Erro ao gerar documentação

```bash
# Verificar se o código compila
go build .

# Executar com mais detalhes
nexs-swag init --dir . --output ./docs --debug
```

### Limpar documentação anterior

```bash
rm -rf docs docs-*
```

## Recursos

- [Documentação Completa](../INSTALL.md)
- [swaggo/swag - Documentação Original](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
