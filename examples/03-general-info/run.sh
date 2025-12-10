#!/bin/bash

echo "=== Exemplo 03: General Info File ==="
echo ""

rm -rf docs

# SEM --generalInfo: Parseia todos os arquivos (pode gerar erro)
echo "1. SEM --generalInfo (pode gerar conflito)..."
../../nexs-swag init --dir . --output ./docs 2>&1 | head -5

rm -rf docs

# COM --generalInfo: Apenas main.go tem info geral
echo ""
echo "2. COM --generalInfo main.go (correto)..."
../../nexs-swag init --dir . --output ./docs --generalInfo main.go

echo ""
echo "✓ Documentação gerada corretamente!"
echo ""
echo "Verificar info no JSON:"
grep -A 3 '"info"' docs/openapi.json

echo ""
echo "Operações encontradas:"
grep -o '"\/[^"]*"' docs/openapi.json | head -5
