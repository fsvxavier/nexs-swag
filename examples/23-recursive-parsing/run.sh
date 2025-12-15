#!/bin/bash

echo "=========================================="
echo "Example 23: Recursive Parsing"
echo "=========================================="
echo ""

# Test 1: Parse only main package (skip internal/)
echo "Test 1: WITHOUT --parseInternal (should skip internal/)"
echo "Command: nexs-swag init --output ./docs-no-internal --ov 3.1"
echo ""
../../nexs-swag init --output ./docs-no-internal --ov 3.1 --quiet
echo "✓ Generated in: ./docs-no-internal"
echo ""

# Test 2: Parse with internal packages
echo "Test 2: WITH --parseInternal (should include internal/)"
echo "Command: nexs-swag init --output ./docs-with-internal --ov 3.1 --parseInternal"
echo ""
../../nexs-swag init --output ./docs-with-internal --ov 3.1 --parseInternal --quiet
echo "✓ Generated in: ./docs-with-internal"
echo ""

# Test 3: Parse with internal but exclude config/
echo "Test 3: WITH --parseInternal and --exclude config"
echo "Command: nexs-swag init --output ./docs-exclude-config --ov 3.1 --parseInternal --exclude config"
echo ""
../../nexs-swag init --output ./docs-exclude-config --ov 3.1 --parseInternal --exclude config --quiet
echo "✓ Generated in: ./docs-exclude-config"
echo ""

# Test 4: Full command with parseDependency (CORRECT SYNTAX)
echo "Test 4: FULL command (parseInternal + parseDependency + exclude)"
echo "Command: nexs-swag init --output ./docs --ov 3.1 --pd --pdl 3 --parseInternal --validate --exclude config"
echo ""
../../nexs-swag init --output ./docs --ov 3.1 --pd --pdl 3 --parseInternal --validate --exclude config
echo ""

# Show results comparison
echo "=========================================="
echo "Results Comparison:"
echo "=========================================="
echo ""
echo "Endpoints in docs-no-internal:"
jq -r '.paths | keys[]' ./docs-no-internal/openapi.json 2>/dev/null | sort
echo ""
echo "Endpoints in docs-with-internal:"
jq -r '.paths | keys[]' ./docs-with-internal/openapi.json 2>/dev/null | sort
echo ""
echo "Endpoints in docs-exclude-config:"
jq -r '.paths | keys[]' ./docs-exclude-config/openapi.json 2>/dev/null | sort
echo ""
echo "Endpoints in docs (full command):"
jq -r '.paths | keys[]' ./docs/openapi.json 2>/dev/null | sort
echo ""
