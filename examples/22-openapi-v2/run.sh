#!/bin/bash

# Script to generate both OpenAPI 3.1.0 and Swagger 2.0 documentation

set -e

echo "🚀 nexs-swag OpenAPI 2.0 Example"
echo "================================="
echo ""

# Build nexs-swag if not installed
if ! command -v nexs-swag &> /dev/null; then
    echo "📦 Building nexs-swag..."
    (cd ../.. && go build -o nexs-swag cmd/nexs-swag/main.go)
    NEXS_SWAG="../../nexs-swag"
else
    NEXS_SWAG="nexs-swag"
fi

# Clean old docs
echo "🧹 Cleaning old documentation..."
rm -rf docs

# Generate OpenAPI 3.1.0
echo ""
echo "📄 Generating OpenAPI 3.1.0..."
$NEXS_SWAG init -g main.go -o ./docs/v3 --openapi-version 3.1 --format json,yaml

# Generate Swagger 2.0
echo ""
echo "📄 Generating Swagger 2.0..."
$NEXS_SWAG init -g main.go -o ./docs/v2 --openapi-version 2.0 --format json,yaml

echo ""
echo "✅ Documentation generated successfully!"
echo ""
echo "📂 Generated files:"
echo "   OpenAPI 3.1.0:"
echo "     - docs/v3/openapi.json"
echo "     - docs/v3/openapi.yaml"
echo ""
echo "   Swagger 2.0:"
echo "     - docs/v2/swagger.json"
echo "     - docs/v2/swagger.yaml"
echo ""
echo "🔍 View differences:"
echo "   diff docs/v2/swagger.json docs/v3/openapi.json"
echo ""
echo "🌐 Run the API:"
echo "   go run main.go"
