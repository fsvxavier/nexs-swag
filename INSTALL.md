# nexs-swag

Ferramenta para gerar documentação OpenAPI/Swagger a partir de annotations em código Go.

## Instalação

### Opção 1: Instalar localmente (recomendado para desenvolvimento)

```bash
# Clonar o repositório
git clone https://github.com/fsvxavier/nexs-swag.git
cd nexs-swag

# Instalar no $GOPATH/bin
go install ./cmd/nexs-swag

# Verificar instalação
nexs-swag --version
```

### Opção 2: Compilar binário

```bash
# Compilar
go build -o nexs-swag ./cmd/nexs-swag

# Mover para um local no PATH (opcional)
sudo mv nexs-swag /usr/local/bin/
```

### Opção 3: Usar sem instalar

```bash
# Executar direto do código
go run ./cmd/nexs-swag --help
```

## Requisitos

- Go 1.24 ou superior
- $GOPATH/bin no PATH (para usar após go install)

### Configurar PATH (se necessário)

```bash
# Adicionar ao ~/.bashrc ou ~/.zshrc
export PATH=$PATH:$GOPATH/bin

# Ou se GOPATH não estiver definido
export PATH=$PATH:$(go env GOPATH)/bin

# Recarregar o shell
source ~/.bashrc  # ou source ~/.zshrc
```

## Uso Básico

```bash
# Gerar documentação
nexs-swag init --dir . --output ./docs

# Ver todas as opções
nexs-swag --help

# Ver opções do comando init
nexs-swag init --help

# Formatar annotations
nexs-swag fmt
```

## Comandos

### init
Gera documentação OpenAPI a partir do código Go.

```bash
nexs-swag init [options]
```

**Flags principais:**
- `--dir <path>` - Diretório com código Go (default: .)
- `--output <path>` - Diretório de saída (default: ./docs)
- `--format <formats>` - Formatos de saída: json,yaml,go (default: json,yaml,go)
- `--generalInfo <file>` - Arquivo com informações gerais
- `--exclude <patterns>` - Padrões para excluir
- `--parseDependency` - Parsear dependências
- `--parseInternal` - Parsear packages internos

### fmt
Formata annotations Swagger no código.

```bash
nexs-swag fmt [options]
```

**Flags:**
- `--dir <path>` - Diretório a formatar (default: .)
- `--exclude <patterns>` - Padrões para excluir
- `--ext <extensions>` - Extensões de arquivo (default: .go)

## Exemplos

O projeto inclui 21 exemplos completos em [`examples/`](./examples/):

```bash
# Ver todos os exemplos
ls examples/

# Executar um exemplo
cd examples/01-basic
./run.sh

# Ou executar manualmente
nexs-swag init --dir . --output ./docs
```

### Lista de Exemplos

1. **01-basic** - Uso básico
2. **02-formats** - Múltiplos formatos
3. **03-general-info** - Arquivo de info geral
4. **04-property-strategy** - Estratégias de nomes
5. **05-required-default** - Required by default
6. **06-exclude** - Excluir diretórios
7. **07-tags-filter** - Filtrar por tags
8. **08-parse-internal** - Parse packages internos
9. **09-parse-dependency** - Parse dependências
10. **10-dependency-level** - Níveis de dependência
11. **11-parse-golist** - Usar go list
12. **12-markdown-files** - Descrições em markdown
13. **13-code-examples** - Exemplos multi-linguagem
14. **14-overrides-file** - Arquivo de overrides
15. **15-generated-time** - Timestamp na documentação
16. **16-instance-name** - Nome customizado
17. **17-template-delims** - Delimitadores customizados
18. **18-collection-format** - Formatos de array
19. **19-parse-func-body** - Parse corpo das funções
20. **20-fmt-command** - Comando de formatação
21. **21-struct-tags** - Tags de struct

Ver [examples/README.md](./examples/README.md) para detalhes.

## Annotations Suportadas

```go
// Informações gerais da API
// @title           API Title
// @version         1.0
// @description     API Description
// @host            localhost:8080
// @BasePath        /api

// Endpoint
// @Summary      Get user
// @Description  Get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  User
// @Failure      404  {object}  Error
// @Router       /users/{id} [get]
func GetUser() {}
```

## Struct Tags

```go
type User struct {
    ID       int       `json:"id" example:"1"`
    Name     string    `json:"name" example:"John" minLength:"2" maxLength:"50"`
    Email    string    `json:"email" format:"email"`
    Password string    `json:"-" swaggerignore:"true"`
    Age      int       `json:"age" minimum:"0" maximum:"150"`
    Status   string    `json:"status" enums:"active,inactive"`
}
```

## Desenvolvimento

```bash
# Clonar repositório
git clone https://github.com/fsvxavier/nexs-swag.git
cd nexs-swag

# Instalar dependências
go mod download

# Executar testes
go test ./...

# Compilar
go build -o nexs-swag ./cmd/nexs-swag

# Instalar localmente
go install ./cmd/nexs-swag
```

## Estrutura do Projeto

```
nexs-swag/
├── cmd/
│   └── nexs-swag/          # CLI
│       └── main.go
├── pkg/
│   ├── format/             # Formatação de código
│   ├── generator/          # Geração de documentação
│   ├── openapi/            # Estruturas OpenAPI
│   └── parser/             # Parser de annotations
├── examples/               # Exemplos de uso
│   ├── 01-basic/
│   ├── 02-formats/
│   └── ...
├── go.mod
├── go.sum
└── README.md
```

## Troubleshooting

### Comando não encontrado após go install

```bash
# Verificar se $GOPATH/bin está no PATH
echo $PATH | grep $(go env GOPATH)/bin

# Se não estiver, adicionar
export PATH=$PATH:$(go env GOPATH)/bin

# Adicionar ao ~/.bashrc para permanente
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### Erro "module not found"

```bash
# Limpar cache e reinstalar
go clean -modcache
cd nexs-swag
go mod download
go install ./cmd/nexs-swag
```

### Erro ao parsear código

```bash
# Verificar se o código compila
go build ./...

# Usar --debug para mais informações
nexs-swag init --dir . --output ./docs --debug
```

## Licença

MIT License

## Contribuindo

Contribuições são bem-vindas! Por favor:

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Baseado em

Este projeto é baseado no [swaggo/swag](https://github.com/swaggo/swag) com melhorias e customizações.
