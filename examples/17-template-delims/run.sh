#!/bin/bash

echo "=== Exemplo 17: Template Delimiters ==="
echo ""

rm -rf docs-*

# Delimiters padrão
echo "1. SEM --templateDelims (padrão: {{ }})..."
../../nexs-swag init --dir . --output ./docs-default --quiet
echo "   ✓ Gerado com delimiters padrão"

echo ""

# Delimiters customizados
echo "2. COM --templateDelims \"[[ ]]\"..."
../../nexs-swag init --dir . --output ./docs-custom --templateDelims "[[ ]]" --quiet
echo "   ✓ Gerado com delimiters customizados"

echo ""
echo "Quando usar:"
echo "  - Conflito com template engines"
echo "  - Mustache, Jinja2, etc"
echo "  - Frontend frameworks que usam {{ }}"
echo ""
echo "Delimiters suportados:"
echo "  [[ ]]   - Recomendado"
echo "  {{{ }}} - Triple mustache"
echo "  <% %>   - ERB style"
echo "  {% %}   - Jinja2 style"
