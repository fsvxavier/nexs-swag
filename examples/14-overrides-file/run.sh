#!/bin/bash

# Exemplo 14 - overrides-file

echo "=== Exemplo 14: overrides-file ==="
echo ""

# Limpar documentação anterior
rm -rf docs

# Gerar OpenAPI 3.1
echo "Gerando OpenAPI 3.1..."
../../nexs-swag init --dir . --output ./docs/v3 --openapi-version 3.1 --overridesFile .swaggo

echo ""

# Gerar Swagger 2.0
echo "Gerando Swagger 2.0..."
../../nexs-swag init --dir . --output ./docs/v2 --openapi-version 2.0 --overridesFile .swaggo

echo ""
echo "✓ OpenAPI 3.1 gerada em ./docs/v3"
echo "✓ Swagger 2.0 gerada em ./docs/v2"
echo ""
echo "Arquivos OpenAPI 3.1:"
ls -lh docs/v3/
echo ""
echo "Arquivos Swagger 2.0:"
ls -lh docs/v2/
