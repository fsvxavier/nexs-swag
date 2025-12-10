#!/bin/bash

echo "=== Exemplo 09: Parse Dependency ==="
echo ""

rm -rf docs-*

# Sem parseDependency
echo "1. SEM --parseDependency..."
../../nexs-swag init --dir . --output ./docs-no-parse --quiet 2>&1 | grep -E "error|warning" || echo "   ⚠️ models.Product não encontrado"

echo ""

# Com parseDependency
echo "2. COM --parseDependency..."
../../nexs-swag init --dir . --output ./docs-with-parse --parseDependency --quiet
echo "   ✓ models.Product incluído"

echo ""
echo "Verificar definição no JSON:"
grep -A 5 '"Product"' docs-with-parse/openapi.json | head -8

echo ""
echo "Benefícios:"
echo "  - Parseia packages importados"
echo "  - Inclui models externos"
echo "  - Resolve todas as dependências"
