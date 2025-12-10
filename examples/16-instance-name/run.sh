#!/bin/bash

echo "=== Exemplo 16: Instance Name ==="
echo ""

rm -rf docs-*

# Instance padrão
echo "1. SEM --instanceName (padrão: swagger)..."
../../nexs-swag init --dir . --output ./docs-default --quiet
echo "   Package:"
head -3 docs-default/docs.go | tail -1

echo ""

# Instance custom
echo "2. COM --instanceName customapi..."
../../nexs-swag init --dir . --output ./docs-custom --instanceName customapi --quiet
echo "   Package:"
head -3 docs-custom/docs.go | tail -1

echo ""
echo "✓ Observe a diferença no package name!"
echo ""
echo "Uso no código:"
echo "  Default:  import _ \"./docs-default\""
echo "  Custom:   import _ \"./docs-custom\""
echo ""
echo "Múltiplas instâncias:"
echo "  ../../nexs-swag init --instanceName apiv1 --output ./docs/v1"
echo "  ../../nexs-swag init --instanceName apiv2 --output ./docs/v2"
