#!/bin/bash

echo "=== Exemplo 07: Tag Filtering ==="
echo ""

rm -rf docs-*

# Sem filtro - todos os endpoints
echo "1. SEM filtro (todos os endpoints)..."
../../nexs-swag init --dir . --output ./docs-all --quiet
echo "   Endpoints:"
grep '"tags"' docs-all/openapi.json | head -4

echo ""

# Apenas users
echo "2. Apenas tag 'users'..."
../../nexs-swag init --dir . --output ./docs-users --tags users --quiet
echo "   Endpoints com tag users:"
grep -o '"\/[^"]*"' docs-users/openapi.json

echo ""

# Apenas admin
echo "3. Apenas tag 'admin'..."
../../nexs-swag init --dir . --output ./docs-admin --tags admin --quiet
echo "   Endpoints com tag admin:"
grep -o '"\/[^"]*"' docs-admin/openapi.json

echo ""

# Excluir internal
echo "4. Excluir tag 'internal'..."
../../nexs-swag init --dir . --output ./docs-no-internal --tags '!internal' --quiet
echo "   Endpoints (sem internal):"
grep -o '"\/[^"]*"' docs-no-internal/openapi.json

echo ""

# Múltiplas tags
echo "5. Tags 'users' e 'admin'..."
../../nexs-swag init --dir . --output ./docs-multi --tags users,admin --quiet
echo "   Endpoints:"
grep -o '"\/[^"]*"' docs-multi/openapi.json

echo ""
echo "✓ Documentação gerada com diferentes filtros!"
