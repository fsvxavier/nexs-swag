# Exemplo 01 - Básico

Demonstra o uso básico do nexs-swag com as flags essenciais.

## Flags Utilizadas

- `--dir .` - Diretório com código Go
- `--output ./docs` - Diretório de saída

## Estrutura

```
01-basic/
├── main.go      # API simples com 2 endpoints
├── run.sh       # Script de execução
└── README.md    # Este arquivo
```

**Nota:** Este exemplo usa o go.mod da raiz do projeto.

## Como Executar

```bash
./run.sh
```

## O que é gerado

1. **docs/openapi.json** - Especificação OpenAPI em JSON
2. **docs/openapi.yaml** - Especificação OpenAPI em YAML
3. **docs/docs.go** - Código Go com a especificação

## API Endpoints

- `GET /api/v1/users/{id}` - Obter usuário
- `POST /api/v1/users` - Criar usuário

## Testar a API

```bash
# Executar o servidor
go run main.go

# Em outro terminal
curl http://localhost:8080/api/v1/users/1
```
