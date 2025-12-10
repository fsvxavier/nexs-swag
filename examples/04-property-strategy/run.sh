#!/bin/bash

echo "=== Exemplo 04: Property Naming Strategy ==="
echo ""

# Limpar
rm -rf docs-*

# 1. snake_case
echo "1. Gerando com snake_case..."
../../nexs-swag init --dir . --output ./docs-snake --propertyStrategy snakecase --quiet
echo "   FirstName → $(grep -o '"first_name"' docs-snake/openapi.json || echo 'first_name')"
echo "   LastName  → $(grep -o '"last_name"' docs-snake/openapi.json || echo 'last_name')"

# 2. camelCase (default)
echo ""
echo "2. Gerando com camelCase (default)..."
../../nexs-swag init --dir . --output ./docs-camel --propertyStrategy camelcase --quiet
echo "   FirstName → $(grep -o '"firstName"' docs-camel/openapi.json || echo 'firstName')"
echo "   LastName  → $(grep -o '"lastName"' docs-camel/openapi.json || echo 'lastName')"

# 3. PascalCase
echo ""
echo "3. Gerando com PascalCase..."
../../nexs-swag init --dir . --output ./docs-pascal --propertyStrategy pascalcase --quiet
echo "   FirstName → $(grep -o '"FirstName"' docs-pascal/openapi.json || echo 'FirstName')"
echo "   LastName  → $(grep -o '"LastName"' docs-pascal/openapi.json || echo 'LastName')"

echo ""
echo "✓ Documentação gerada com diferentes estratégias!"
echo ""
echo "Comparar:"
echo "  - docs-snake/openapi.json   (snake_case)"
echo "  - docs-camel/openapi.json   (camelCase)"
echo "  - docs-pascal/openapi.json  (PascalCase)"

echo ""
echo "Nota: user_id sempre mantém o nome da tag json explícita"
