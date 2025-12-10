#!/bin/bash

echo "=== Exemplo 10: Dependency Level ==="
echo ""

rm -rf docs-*

# Level 0
echo "1. Level 0 (apenas types referenciados diretamente)..."
../../nexs-swag init --dir . --output ./docs-level0 --parseDependency --parseDependencyLevel 0 --quiet
echo "   Definitions encontradas:"
grep -o '"[^"]*\(Order\|Item\|Meta\)"' docs-level0/openapi.json 2>/dev/null | sort -u || echo "   Nenhum"

echo ""

# Level 1
echo "2. Level 1 (+ 1 nível de referências)..."
../../nexs-swag init --dir . --output ./docs-level1 --parseDependency --parseDependencyLevel 1 --quiet
echo "   Definitions encontradas:"
grep -o '"[^"]*\(Order\|Item\|Meta\)"' docs-level1/openapi.json 2>/dev/null | sort -u

echo ""

# Level 2
echo "3. Level 2 (+ 2 níveis de referências)..."
../../nexs-swag init --dir . --output ./docs-level2 --parseDependency --parseDependencyLevel 2 --quiet
echo "   Definitions encontradas:"
grep -o '"[^"]*\(Order\|Item\|Meta\)"' docs-level2/openapi.json 2>/dev/null | sort -u

echo ""
echo "Estrutura de dependências:"
echo "  Order → Item → Meta"
echo ""
echo "Cada nível adiciona mais types referenciados na cadeia de dependências."
