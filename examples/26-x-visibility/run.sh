#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Example 26: X-Visibility${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Step 1: Generate documentation
echo -e "${YELLOW}Step 1: Generating OpenAPI documentation with x-visibility...${NC}"
../../nexs-swag init --output ./docs --ov 3.1 --pd=true --pdl=3

if [ $? -ne 0 ]; then
    echo -e "\n${RED}✗ Documentation generation failed${NC}"
    exit 1
fi

# Step 2: Display public endpoints
echo -e "\n${YELLOW}Step 2: Public Endpoints${NC}"
echo "Generated in: docs/openapi_public.json"
jq '.paths | keys' docs/openapi_public.json

# Step 3: Display private endpoints
echo -e "\n${YELLOW}Step 3: Private Endpoints${NC}"
echo "Generated in: docs/openapi_private.json"
jq '.paths | keys' docs/openapi_private.json

# Step 4: Display public schemas
echo -e "\n${YELLOW}Step 4: Public Schemas${NC}"
jq '.components.schemas | keys' docs/openapi_public.json

# Step 5: Display private schemas
echo -e "\n${YELLOW}Step 5: Private Schemas${NC}"
jq '.components.schemas | keys' docs/openapi_private.json

# Summary
echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Example completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}\n"

echo "Generated files:"
echo "  - docs/openapi_public.json  (public API endpoints)"
echo "  - docs/openapi_private.json (private/admin endpoints)"
echo "  - docs/openapi_public.yaml"
echo "  - docs/openapi_private.yaml"
echo "  - docs/docs_public.go"
echo "  - docs/docs_private.go"

echo -e "\n${BLUE}Key Features Demonstrated:${NC}"
echo "  • @x-visibility public - Endpoint only in public spec"
echo "  • @x-visibility private - Endpoint only in private spec"
echo "  • No annotation - Endpoint in both specs"
echo "  • Automatic schema filtering based on usage"
echo "  • Recursive schema dependency collection"

echo -e "\n${BLUE}What to notice:${NC}"
echo "  1. Public spec has UserPublic schema (no sensitive fields)"
echo "  2. Private spec has UserPrivate schema (includes email, password, role)"
echo "  3. Private spec also has UserPublic (used by shared POST /users endpoint)"
echo "  4. ErrorResponse appears in both (shared schema)"
echo "  5. Admin endpoints only in private spec"
echo "  6. Public GET /users/{id} only in public spec"
echo "  7. CreateUser POST /users (no annotation) appears in both specs"
