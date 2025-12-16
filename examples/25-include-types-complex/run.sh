#!/bin/bash

echo "=== Example 25: Include Types Complex ==="
echo ""

# Build nexs-swag if not exists
if [ ! -f "../../nexs-swag" ]; then
    echo "Building nexs-swag..."
    (cd ../.. && go build ./cmd/nexs-swag)
fi

echo "1. Generating with --includeTypes='all' (default)"
echo "   This should include all referenced structs transitively"
../../nexs-swag init -g main.go -o ./docs-all --includeTypes="all"
echo ""

echo "2. Generating with --includeTypes='struct'"
echo "   This should include structs but exclude interfaces"
../../nexs-swag init -g main.go -o ./docs-struct --includeTypes="struct"
echo ""

echo "3. Analyzing generated schemas..."
echo ""

echo "   Total schemas in docs-all:"
grep -c '"type": "object"' ./docs-all/openapi.json 2>/dev/null || echo "   0"

echo "   Total schemas in docs-struct:"
grep -c '"type": "object"' ./docs-struct/openapi.json 2>/dev/null || echo "   0"
echo ""

echo "4. Checking transitive dependencies..."
echo "   Looking for Money (should be included via OrderItem):"
if grep -q '"Money":' ./docs-struct/openapi.json 2>/dev/null; then
    echo "   ✓ Money found (transitive dependency resolved)"
else
    echo "   ✗ Money not found (transitive dependency failed)"
fi

echo "   Looking for Address (should be included via OrderRequest):"
if grep -q '"Address":' ./docs-struct/openapi.json 2>/dev/null; then
    echo "   ✓ Address found (transitive dependency resolved)"
else
    echo "   ✗ Address not found (transitive dependency failed)"
fi
echo ""

echo "5. Checking swaggertype conversions..."
echo "   UUID field (should be string with format uuid):"
grep -A3 '"order_id"' ./docs-struct/openapi.json 2>/dev/null | head -4

echo "   Time field (should be string with format date-time):"
grep -A3 '"created_at"' ./docs-struct/openapi.json 2>/dev/null | head -4
echo ""

echo "6. Verifying unused types are excluded..."
if grep -q '"UnusedComplexModel"' ./docs-struct/openapi.json 2>/dev/null; then
    echo "   ✗ UnusedComplexModel found (should be excluded)"
else
    echo "   ✓ UnusedComplexModel not found (correctly excluded)"
fi

if grep -q '"UnusedService"' ./docs-struct/openapi.json 2>/dev/null; then
    echo "   ✗ UnusedService found (should be excluded)"
else
    echo "   ✓ UnusedService not found (correctly excluded)"
fi
echo ""

echo "7. List of all generated schemas:"
grep -o '"[A-Z][a-zA-Z.]*":' ./docs-struct/openapi.json 2>/dev/null | sort -u || echo "   (none)"
echo ""

echo "✓ Documentation generated successfully!"
echo "  - Check ./docs-all/ for default output"
echo "  - Check ./docs-struct/ for struct-only output"
echo ""
echo "To inspect the full output:"
echo "  cat docs-struct/openapi.json | python3 -m json.tool | less"
