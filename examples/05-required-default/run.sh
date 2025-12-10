#!/bin/bash

echo "=== Exemplo 05: Required By Default ==="
echo ""

rm -rf docs-*

# Sem flag (default: campos opcionais)
echo "1. SEM --requiredByDefault (campos opcionais)..."
../../nexs-swag init --dir . --output ./docs-optional --quiet
echo "   Required fields:"
grep -A 20 '"Product"' docs-optional/openapi.json | grep -A 2 '"required"' || echo "   (nenhum)"

echo ""

# Com flag (campos required por padrão)
echo "2. COM --requiredByDefault (campos required)..."
../../nexs-swag init --dir . --output ./docs-required --requiredByDefault --quiet
echo "   Required fields:"
grep -A 20 '"Product"' docs-required/openapi.json | grep -A 5 '"required"'

echo ""
echo "✓ Comparar os schemas gerados:"
echo ""
echo "Campos REQUIRED com a flag:"
echo "  - id, name, price (campos normais)"
echo ""
echo "Campos OPTIONAL mesmo com a flag:"
echo "  - description (tem omitempty)"
echo "  - discount (é pointer)"
echo "  - category (tem binding:omitempty)"

echo ""
echo "Arquivos gerados:"
echo "  - docs-optional/openapi.json  (sem flag)"
echo "  - docs-required/openapi.json  (com flag)"
