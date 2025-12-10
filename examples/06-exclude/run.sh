#!/bin/bash

echo "=== Exemplo 06: Exclude Patterns ==="
echo ""

rm -rf docs-*

# Sem exclude - parseia tudo
echo "1. SEM --exclude (parseia tudo)..."
../../nexs-swag init --dir . --output ./docs-all --quiet
echo "   Endpoints encontrados:"
grep -o '"\/[^"]*"' docs-all/openapi.json 2>/dev/null | wc -l

# Com exclude
echo ""
echo "2. COM --exclude mock,testdata..."
../../nexs-swag init --dir . --output ./docs-excluded --exclude mock,testdata --quiet
echo "   Endpoints encontrados:"
grep -o '"\/[^"]*"' docs-excluded/openapi.json 2>/dev/null | wc -l

echo ""
echo "✓ Documentação gerada!"
echo ""
echo "Nota: *_test.go sempre são excluídos automaticamente"
echo ""
echo "Diretórios excluídos por padrão:"
echo "  - vendor/"
echo "  - testdata/"
echo "  - docs/"
echo "  - .git/"
echo "  - *_test.go files"
