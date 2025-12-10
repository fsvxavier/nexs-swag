#!/bin/bash

echo "=== Exemplo 13: Code Examples ==="
echo ""

rm -rf docs

echo "Gerando documentação com code examples..."
../../nexs-swag init --dir . --output ./docs --codeExampleFilesDir ./code_samples --quiet

echo ""
echo "✓ Documentação gerada!"
echo ""
echo "Verificando x-codeSamples no JSON:"
grep -A 20 "x-codeSamples" docs/openapi.json | head -25

echo ""
echo "Linguagens detectadas:"
echo "  - create_user.go  → Go"
echo "  - create_user.js  → JavaScript"
echo "  - create_user.py  → Python"
echo "  - create_user.sh  → Bash"
