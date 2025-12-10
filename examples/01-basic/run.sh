#!/bin/bash

# Exemplo 01 - Básico
# Demonstra uso básico do nexs-swag com flags --dir e --output

echo "=== Exemplo 01: Uso Básico ==="
echo ""

# Limpar documentação anterior
rm -rf docs

# Gerar documentação
echo "Gerando documentação..."
../../nexs-swag init --dir . --output ./docs

echo ""
echo "✓ Documentação gerada em ./docs"
echo ""
echo "Arquivos criados:"
ls -lh docs/

echo ""
echo "Para visualizar a documentação:"
echo "  cat docs/openapi.json"
echo "  cat docs/openapi.yaml"
