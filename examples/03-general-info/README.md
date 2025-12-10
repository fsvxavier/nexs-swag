# Exemplo 03 - General Info File

Demonstra o uso de `--generalInfo` para especificar qual arquivo contém as annotations gerais da API.

## Problema

Quando você tem múltiplos arquivos Go, o parser pode encontrar annotations de info geral (@title, @version) em vários lugares, causando conflitos.

## Solução

Use `--generalInfo` para especificar exatamente qual arquivo contém a info geral:

```bash
nexs-swag init --generalInfo main.go
```

## Estrutura

```
03-general-info/
├── main.go       # ✅ TEM @title, @version, @host, etc
├── products.go   # ❌ Apenas endpoints de produtos
├── orders.go     # ❌ Apenas endpoints de orders
└── run.sh
```

## Regra

- **Arquivo de Info Geral:** Deve ter @title, @version, @host, @BasePath
- **Outros Arquivos:** Devem ter APENAS endpoints (@Router, @Summary, etc)

## Como Executar

```bash
chmod +x run.sh
./run.sh
```

## Comparação

### Sem --generalInfo
```bash
nexs-swag init --dir .
# Pode gerar erro se encontrar @title em múltiplos arquivos
```

### Com --generalInfo
```bash
nexs-swag init --dir . --generalInfo main.go
# ✅ Correto: apenas main.go é parseado para info geral
# ✅ products.go e orders.go fornecem apenas endpoints
```

## Benefícios

1. **Evita conflitos:** Um único local para info da API
2. **Mais rápido:** Parser não precisa verificar todos os arquivos para info geral
3. **Organização:** Separa concerns (info geral vs endpoints)
