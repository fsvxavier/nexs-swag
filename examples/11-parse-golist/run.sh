#!/bin/bash

echo "=== Exemplo 11: Parse GoList ==="
echo ""

rm -rf docs-*

# Sem parseGoList (parsing manual)
echo "1. SEM --parseGoList (parsing manual)..."
time ../../nexs-swag init --dir . --output ./docs-manual --quiet 2>&1 | tail -1

echo ""

# Com parseGoList (usa go list)
echo "2. COM --parseGoList (usa 'go list')..."
time ../../nexs-swag init --dir . --output ./docs-golist --parseGoList --quiet 2>&1 | tail -1

echo ""
echo "✓ Documentação gerada!"
echo ""
echo "Benefícios do --parseGoList:"
echo "  • Mais rápido em projetos grandes"
echo "  • Usa informações do Go modules"
echo "  • Detecta dependências automaticamente"
echo "  • Respeita go.mod e replace directives"
echo ""
echo "Quando usar:"
echo "  • Projetos com muitos packages"
echo "  • Go modules configurado"
echo "  • Dependências complexas"
