#!/bin/bash

echo "=== Exemplo 08: Parse Internal Packages ==="
echo ""

rm -rf docs-*

# Sem parseInternal - ignora internal/
echo "1. SEM --parseInternal (ignora internal/)..."
../../nexs-swag init --dir . --output ./docs-no-internal --quiet
echo "   Schemas encontrados:"
grep -o '"[A-Z][a-zA-Z]*"' docs-no-internal/openapi.json | sort -u
echo "   Endpoints:"
grep -o '"\/[^"]*"' docs-no-internal/openapi.json

echo ""

# Com parseInternal - inclui internal/
echo "2. COM --parseInternal (inclui internal/)..."
../../nexs-swag init --dir . --output ./docs-with-internal --parseInternal --quiet
echo "   Schemas encontrados:"
grep -o '"[A-Z][a-zA-Z]*"' docs-with-internal/openapi.json | sort -u
echo "   Endpoints:"
grep -o '"\/[^"]*"' docs-with-internal/openapi.json

echo ""
echo "✓ Compare os resultados!"
echo ""
echo "SEM flag: apenas /users"
echo "COM flag: /users + /internal/config"
