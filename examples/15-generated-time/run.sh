#!/bin/bash

echo "=== Exemplo 15: Generated Time ==="
echo ""

rm -rf docs-*

# Sem timestamp
echo "1. SEM --generatedTime..."
../../nexs-swag init --dir . --output ./docs-no-time --quiet
echo "   Header do docs.go:"
head -3 docs-no-time/docs.go

echo ""

# Com timestamp
echo "2. COM --generatedTime..."
../../nexs-swag init --dir . --output ./docs-with-time --generatedTime --quiet
echo "   Header do docs.go:"
head -3 docs-with-time/docs.go

echo ""
echo "✓ Compare os headers!"
echo ""
echo "Benefícios do timestamp:"
echo "  - Rastreamento de quando foi gerado"
echo "  - Útil para debug e versionamento"
echo "  - Pode ser usado em CI/CD"
