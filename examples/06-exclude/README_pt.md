# Exemplo 06 - Exclude Patterns

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra como excluir diretórios e arquivos do parsing.

## Flag

```bash
--exclude pattern1,pattern2,pattern3
```

## Uso

```bash
# Excluir um diretório
nexs-swag init --exclude mock

# Excluir múltiplos
nexs-swag init --exclude mock,testdata,vendor

# Excluir com wildcards
nexs-swag init --exclude "*.test.go,*_mock.go"
```

## Exclusões Automáticas

Sempre excluídos (não precisa especificar):
- `vendor/` - Dependências
- `testdata/` - Dados de teste
- `docs/` - Documentação gerada
- `.git/` - Repositório Git
- `*_test.go` - Arquivos de teste

## Estrutura do Exemplo

```
06-exclude/
├── main.go           # ✅ Será parseado
├── main_test.go      # ❌ Excluído (test)
├── mock/
│   └── mock.go       # ❌ Excluído (com flag)
└── testdata/
    └── data.go       # ❌ Excluído (automático)
```

## Como Executar

```bash
./run.sh
```

## Casos de Uso

- **mock:** Código de mocking para testes
- **testdata:** Fixtures e dados de teste
- **vendor:** Dependências (se usar vendor)
- **examples:** Código de exemplo
- **internal:** Pacotes internos (use --parseInternal para incluir)
