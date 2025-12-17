#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Example 27: X-Visibility with Swagger 2.0${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Step 1: Generate documentation
echo -e "${YELLOW}Step 1: Generating Swagger 2.0 documentation with x-visibility...${NC}"
../../nexs-swag init --output ./docs --ov 2.0

if [ $? -ne 0 ]; then
    echo -e "\n${RED}✗ Documentation generation failed${NC}"
    exit 1
fi

# Step 2: Display public endpoints
echo -e "\n${YELLOW}Step 2: Public Endpoints${NC}"
echo "Generated in: docs/swagger_public.json"
jq '.paths | keys' docs/swagger_public.json

# Step 3: Display private endpoints
echo -e "\n${YELLOW}Step 3: Private Endpoints${NC}"
echo "Generated in: docs/swagger_private.json"
jq '.paths | keys' docs/swagger_private.json

# Step 4: Display public definitions
echo -e "\n${YELLOW}Step 4: Public Definitions (Schemas)${NC}"
jq '.definitions | keys' docs/swagger_public.json

# Step 5: Display private definitions
echo -e "\n${YELLOW}Step 5: Private Definitions (Schemas)${NC}"
jq '.definitions | keys' docs/swagger_private.json

# Summary
echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Example completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}\n"

echo "Generated files:"
echo "  - docs/swagger_public.json  (public API endpoints)"
echo "  - docs/swagger_private.json (private/admin endpoints)"
echo "  - docs/swagger_public.yaml"
echo "  - docs/swagger_private.yaml"
echo "  - docs/docs_public.go"
echo "  - docs/docs_private.go"

echo -e "\n${BLUE}Key Features Demonstrated:${NC}"
echo "  • @x-visibility public - Endpoint only in public spec"
echo "  • @x-visibility private - Endpoint only in private spec"
echo "  • No annotation - Endpoint in both specs"
echo "  • Automatic schema filtering based on usage"
echo "  • Swagger 2.0 compatibility with x-visibility"

echo -e "\n${BLUE}What to notice:${NC}"
echo "  1. Public spec has UserPublic definition (no sensitive fields)"
echo "  2. Private spec has UserPrivate definition (includes email, password, role)"
echo "  3. Private spec also has UserPublic (used by shared POST /users endpoint)"
echo "  4. ErrorResponse appears in both (shared schema)"
echo "  5. Admin endpoints only in private spec"
echo "  6. Public GET /users/{id} only in public spec"
echo "  7. CreateUser POST /users (no annotation) appears in both specs"
echo "  8. Works seamlessly with Swagger 2.0 (definitions instead of components.schemas)"
