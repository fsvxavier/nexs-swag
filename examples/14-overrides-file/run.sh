#!/bin/bash

echo "=== Exemplo 14: Overrides File ==="
echo ""

rm -rf docs-*

# Sem overrides - tipos complexos
echo "1. SEM --overridesFile..."
../../nexs-swag init --dir . --output ./docs-no-override --quiet
echo "   Schema de Account (ID field):"
grep -A 5 '"ID"' docs-no-override/openapi.json | head -6

echo ""

# Com overrides - tipos simples
echo "2. COM --overridesFile .swaggo..."
../../nexs-swag init --dir . --output ./docs-with-override --overridesFile .swaggo --quiet
echo "   Schema de Account (ID field):"
grep -A 5 '"ID"' docs-with-override/openapi.json | head -6

echo ""
echo "✓ Comparar os schemas!"
echo ""
echo "SEM override:"
echo "  ID: object (sql.NullInt64 completo)"
echo ""
echo "COM override:"
echo "  ID: integer (tipo simples)"
echo ""
echo "Ver arquivo .swaggo para configurar outros tipos"
