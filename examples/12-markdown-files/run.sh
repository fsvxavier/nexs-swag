#!/bin/bash

echo "=== Exemplo 12: Markdown Files ==="
echo ""

rm -rf openapi-docs

# Sem markdown - description curta
echo "1. SEM --markdownFiles..."
../../nexs-swag init --dir . --output ./openapi-docs --quiet
echo "   Description de POST /users:"
grep -A 5 '"post"' openapi-docs/openapi.json | grep '"description"' | head -1

rm -rf openapi-docs

echo ""

# Com markdown - description completa
echo "2. COM --markdownFiles..."
../../nexs-swag init --dir . --output ./openapi-docs --markdownFiles ./docs --quiet
echo "   Description de POST /users (primeiras 200 chars):"
grep -A 5 '"post"' openapi-docs/openapi.json | grep '"description"' | head -1 | cut -c1-200

echo ""
echo "✓ Documentação gerada!"
echo ""
echo "Compare:"
echo "  - Sem flag: usa 'file(create-user.md)' literal"
echo "  - Com flag: substitui pelo conteúdo do arquivo markdown"
echo ""
echo "Ver conteúdo completo:"
echo "  cat openapi-docs/openapi.json | jq '.paths.\"/users\".post.description'"
