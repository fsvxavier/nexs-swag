#!/bin/bash

echo "=== Exemplo 19: Parse Func Body ==="
echo ""

rm -rf docs-*

# Sem parseFuncBody
echo "1. SEM --parseFuncBody..."
../../nexs-swag init --dir . --output ./docs-no-parse --quiet
echo "   ✓ Apenas annotations são parseadas"

echo ""

# Com parseFuncBody
echo "2. COM --parseFuncBody..."
../../nexs-swag init --dir . --output ./docs-with-parse --parseFuncBody --quiet
echo "   ✓ Corpo das funções também é analisado"

echo ""
echo "Benefícios:"
echo "  • Detecta validações no código"
echo "  • Infere responses adicionais"
echo "  • Análise mais profunda"
echo ""
echo "⚠️ Atenção:"
echo "  • Mais lento"
echo "  • Pode gerar false positives"
echo "  • Use apenas se necessário"
