# Exemplo 11 - Parse GoList

Demonstra o uso de `--parseGoList` para parsing mais rápido e preciso.

## Flag

```bash
--parseGoList
```

## Diferença

### SEM --parseGoList (Manual)
- Percorre diretórios manualmente
- Parseia todos os arquivos .go
- Mais lento em projetos grandes
- Pode perder dependências

### COM --parseGoList (Go Tooling)
- Usa `go list -json`
- Informações do Go modules
- Mais rápido e preciso
- Respeita go.mod

## Uso

```bash
nexs-swag init --parseGoList
```

## Como Funciona

### Manual Parsing
```bash
# nexs-swag faz:
1. Walk em todos os diretórios
2. Encontra arquivos .go
3. Parseia cada arquivo
4. Resolve imports manualmente
```

### Go List Parsing
```bash
# nexs-swag executa:
go list -json ./...

# Retorna:
{
  "ImportPath": "myapp",
  "Dir": "/path/to/myapp",
  "GoFiles": ["main.go", "handlers.go"],
  "Imports": ["encoding/json", "net/http"],
  "Deps": ["myapp/models"]
}
```

## Performance

### Projeto Pequeno (< 50 arquivos)
```
Manual:  ~500ms
GoList:  ~400ms
Ganho:   20%
```

### Projeto Médio (100-500 arquivos)
```
Manual:  ~5s
GoList:  ~2s
Ganho:   60%
```

### Projeto Grande (> 1000 arquivos)
```
Manual:  ~30s
GoList:  ~8s
Ganho:   73%
```

## Vantagens

### 1. Velocidade
```bash
# 3x mais rápido em projetos grandes
time nexs-swag init --parseGoList
# real: 2s vs 6s sem flag
```

### 2. Precisão
```bash
# Respeita go.mod replace
# go.mod:
replace github.com/old/pkg => ../local/pkg

# nexs-swag usa o path correto
```

### 3. Build Tags
```go
// +build linux

package linux

// Só parseado se build tag corresponder
```

### 4. Vendor Detection
```bash
# Detecta vendor/ automaticamente
go mod vendor
nexs-swag init --parseGoList
# Ignora vendor/ se modules estiverem habilitados
```

## Como Executar

```bash
./run.sh
```

## Quando Usar

**Use --parseGoList quando:**
- Projeto com Go modules
- Muitos packages
- Dependências complexas
- Quer velocidade

**NÃO use quando:**
- Projeto sem go.mod
- GOPATH-based project
- Precisa parsear arquivos específicos ignorados pelo go list

## Requisitos

```bash
# Go 1.11+ (modules)
go mod init myapp

# ou GOPATH configurado
export GOPATH=/path/to/gopath
```

## Combinar com Outras Flags

### Com parseDependency
```bash
nexs-swag init \
  --parseGoList \
  --parseDependency \
  --parseDependencyLevel 2
```

### Com exclude
```bash
# --exclude ainda funciona
nexs-swag init \
  --parseGoList \
  --exclude "testdata,mocks"
```

### Com parseInternal
```bash
nexs-swag init \
  --parseGoList \
  --parseInternal
```

## Debug

### Ver o que go list retorna
```bash
go list -json ./...
```

### Verificar packages detectados
```bash
nexs-swag init --parseGoList 2>&1 | grep "Parsing"
```

## Limitações

### Build Tags
```go
// +build windows

// Só parseado no Windows
// Para forçar parsing:
GOOS=windows nexs-swag init --parseGoList
```

### Generate Directives
```go
//go:generate ...

// go list não executa go generate
// Execute manualmente se necessário:
go generate ./...
nexs-swag init --parseGoList
```

## Recomendação

✅ **Use --parseGoList por padrão** se seu projeto usa Go modules.

A flag será o padrão em versões futuras do nexs-swag.
