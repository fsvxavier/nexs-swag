#!/bin/bash

echo "=== Exemplo 02: Múltiplos Formatos ==="
echo ""

rm -rf docs

# Gerar apenas JSON
echo "1. Gerando apenas JSON..."
../../nexs-swag init --dir . --output ./docs --format json
echo "   Criado: $(ls docs/)"

rm -rf docs

# Gerar apenas YAML
echo ""
echo "2. Gerando apenas YAML..."
../../nexs-swag init --dir . --output ./docs --format yaml
echo "   Criado: $(ls docs/)"

rm -rf docs

# Gerar JSON e YAML
echo ""
echo "3. Gerando JSON e YAML..."
../../nexs-swag init --dir . --output ./docs --format json,yaml
echo "   Criado: $(ls docs/)"

rm -rf docs

# Gerar todos os formatos (default)
echo ""
echo "4. Gerando todos os formatos (json,yaml,go)..."
../../nexs-swag init --dir . --output ./docs --format json,yaml,go
echo "   Criado: $(ls docs/)"

echo ""
echo "✓ Exemplos de múltiplos formatos gerados!"
