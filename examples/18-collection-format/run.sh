#!/bin/bash

echo "=== Exemplo 18: Collection Format ==="
echo ""

rm -rf docs

echo "Gerando documentação..."
../../nexs-swag init --dir . --output ./docs --collectionFormat multi --quiet

echo ""
echo "✓ Documentação gerada!"
echo ""
echo "Formatos suportados:"
echo ""
echo "1. CSV (comma):"
echo "   ?ids=1,2,3"
echo ""
echo "2. Multi (multiple params):"
echo "   ?tags=go&tags=api&tags=web"
echo ""
echo "3. Pipes:"
echo "   ?statuses=active|pending|done"
echo ""
echo "4. TSV (tab):"
echo "   ?values=a\tb\tc"
echo ""
echo "5. SSV (space):"
echo "   ?items=one%20two%20three"
echo ""
echo "Verificando formato no JSON:"
grep -A 3 "collectionFormat" docs/openapi.json
