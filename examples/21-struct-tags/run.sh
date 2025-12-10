#!/bin/bash

echo "=== Exemplo 21: Struct Tags ==="
echo ""

rm -rf docs

echo "Gerando documentação..."
../../nexs-swag init --dir . --output ./docs --quiet

echo ""
echo "✓ Documentação gerada!"
echo ""
echo "Tags de validação no User:"
grep -A 20 '"User"' docs/openapi.json | grep -E '(minLength|maxLength|minimum|maximum|format|example)' | head -10

echo ""
echo "Campo ignorado (Password):"
grep -c '"password"' docs/openapi.json || echo "0 ocorrências (ignorado com sucesso!)"

echo ""
echo "Tags suportadas:"
echo "  • example - Valor de exemplo"
echo "  • format - Formato (email, date, etc)"
echo "  • minLength/maxLength - Validação de string"
echo "  • minimum/maximum - Validação numérica"
echo "  • enums - Valores permitidos"
echo "  • default - Valor padrão"
echo "  • swaggertype - Override de tipo"
echo "  • swaggerignore - Ignorar campo"
echo "  • readonly - Campo read-only"
echo "  • x-* - Extensões customizadas"
