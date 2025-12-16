#!/bin/bash

echo "=== Example 24: Include Types Simple ==="
echo ""

# Build nexs-swag if not exists
if [ ! -f "../../nexs-swag" ]; then
    echo "Building nexs-swag..."
    (cd ../.. && go build ./cmd/nexs-swag)
fi

echo "1. Generating with --includeTypes='all' (default)"
../../nexs-swag init -g main.go -o ./docs-all --includeTypes="all"
echo ""

echo "2. Generating with --includeTypes='struct'"
../../nexs-swag init -g main.go -o ./docs-struct --includeTypes="struct"
echo ""

echo "3. Comparing outputs..."
echo "   Schemas in docs-all:"
grep -o '"[A-Z][a-zA-Z]*":' ./docs-all/openapi.json | sort -u || echo "   (none)"
echo "   Schemas in docs-struct:"
grep -o '"[A-Z][a-zA-Z]*":' ./docs-struct/openapi.json | sort -u || echo "   (none)"
echo ""

echo "4. Checking swaggertype conversion for time.Time:"
grep -A2 '"created_at"' ./docs-struct/openapi.json || echo "   (not found)"
echo ""

echo "✓ Documentation generated successfully!"
echo "  - Check ./docs-all/ for default output"
echo "  - Check ./docs-struct/ for struct-only output"
