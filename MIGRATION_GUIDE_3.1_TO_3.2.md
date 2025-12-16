# Migration Guide: OpenAPI 3.1.x → 3.2.0

🌍 **English** • [Português (Brasil)](MIGRATION_GUIDE_3.1_TO_3.2_pt.md) • [Español](MIGRATION_GUIDE_3.1_TO_3.2_es.md)

**Document Version:** 1.0  
**Last Updated:** December 15, 2025  
**Target Audience:** API developers using nexs-swag

---

## Table of Contents

1. [Overview](#overview)
2. [What's New in OpenAPI 3.2.0](#whats-new-in-openapi-32)
3. [Breaking Changes](#breaking-changes)
4. [Migration Strategy](#migration-strategy)
5. [Feature Adoption Guide](#feature-adoption-guide)
6. [Conversion Tools](#conversion-tools)
7. [Common Issues & Solutions](#common-issues--solutions)
8. [Best Practices](#best-practices)
9. [FAQ](#faq)

---

## Overview

OpenAPI 3.2.0 introduces several enhancements while maintaining **full backward compatibility** with OpenAPI 3.1.x specifications. This guide helps you migrate your API documentation from OpenAPI 3.1.x to 3.2.0 using nexs-swag.

### Key Benefits of Upgrading

- ✅ **QUERY HTTP Method** - Secure complex queries with request bodies
- ✅ **Streaming Support** - First-class SSE and streaming response documentation
- ✅ **Enhanced Security** - Deprecation markers and OAuth2 metadata URL
- ✅ **Device Authorization** - OAuth 2.0 Device Grant (RFC 8628) support
- ✅ **Backward Compatible** - No breaking changes to existing annotations

### Should You Migrate?

| Use Case | Recommendation |
|----------|----------------|
| New projects | ✅ **Start with 3.2.0** |
| Existing 3.1.x projects | ⚠️ **Migrate only if needed** |
| Using streaming/SSE | ✅ **Migrate to use itemSchema** |
| OAuth2 device flow | ✅ **Migrate for deviceAuthorization** |
| Legacy security schemes | ✅ **Migrate to mark as deprecated** |
| Stable projects | ⏸️ **No rush, migrate when convenient** |

---

## What's New in OpenAPI 3.2.0

### 1. QUERY HTTP Method

**Problem Solved:** Complex searches requiring structured data in the request body while maintaining idempotency and safety.

**Before (OpenAPI 3.1.x):**
```go
// @Router /products/search [post]  // ❌ Not idempotent
func SearchProducts(c *gin.Context) {}
```

**After (OpenAPI 3.2.0):**
```go
// @Router /products/search [query]  // ✅ Idempotent, safe, with body
func SearchProducts(c *gin.Context) {}
```

**Impact:** Semantic clarity, better caching, improved API design.

---

### 2. SecurityScheme.Deprecated

**Problem Solved:** Marking legacy authentication methods for future removal without breaking existing clients.

**Before (OpenAPI 3.1.x):**
```go
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description ⚠️ DEPRECATED: Use OAuth2 instead
```

**After (OpenAPI 3.2.0):**
```go
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @deprecated true
// @description Use OAuth2 instead
```

**Generated OpenAPI:**
```yaml
securitySchemes:
  ApiKeyAuth:
    type: apiKey
    deprecated: true  # ✅ Machine-readable deprecation
```

**Impact:** Automated tooling can detect deprecated auth, warn developers.

---

### 3. OAuth2MetadataURL

**Problem Solved:** OAuth2 server discovery via RFC 8414 (Authorization Server Metadata).

**Before (OpenAPI 3.1.x):**
```go
// @securityDefinitions.oauth2.application OAuth2
// @tokenUrl https://auth.example.com/token
// Manual configuration required in client
```

**After (OpenAPI 3.2.0):**
```go
// @securityDefinitions.oauth2.application OAuth2
// @tokenUrl https://auth.example.com/token
// @oauth2metadataurl https://auth.example.com/.well-known/oauth-authorization-server
```

**Generated OpenAPI:**
```yaml
securitySchemes:
  OAuth2:
    oauth2MetadataUrl: https://auth.example.com/.well-known/oauth-authorization-server
```

**Impact:** Client libraries auto-discover endpoints, scopes, grant types.

---

### 4. DeviceAuthorization Flow

**Problem Solved:** Documenting OAuth 2.0 for input-constrained devices (TVs, IoT, CLI tools).

**Before (OpenAPI 3.1.x):**
```go
// No standard way to document device flow
// Custom documentation or extensions required
```

**After (OpenAPI 3.2.0):**
```go
// @securityDefinitions.oauth2.deviceAuth OAuth2Device
// @deviceAuthorization https://auth.example.com/device https://auth.example.com/token device-code
// @scopes.device:access Access device resources
```

**Generated OpenAPI:**
```yaml
securitySchemes:
  OAuth2Device:
    flows:
      urn:ietf:params:oauth:grant-type:device_code:
        deviceAuthorizationUrl: https://auth.example.com/device
        tokenUrl: https://auth.example.com/token
        scopes:
          device:access: Access device resources
```

**Impact:** Proper documentation for smart TVs, CLI tools, IoT devices.

---

### 5. MediaType.ItemSchema (Streaming)

**Problem Solved:** Document Server-Sent Events (SSE) and streaming responses with schema validation.

**Before (OpenAPI 3.1.x):**
```go
// @Success 200 {object} Event  // ❌ Implies single object
// @Produce text/event-stream
func StreamEvents(c *gin.Context) {}
```

**After (OpenAPI 3.2.0):**
```go
// @Success 200 {stream} Event  // ✅ Each event is an Event object
// @Produce text/event-stream
func StreamEvents(c *gin.Context) {}
```

**Generated OpenAPI:**
```yaml
responses:
  '200':
    content:
      text/event-stream:
        itemSchema:  # ✅ Schema for each streamed item
          $ref: '#/components/schemas/Event'
```

**Impact:** Code generators create proper streaming clients, validation per event.

---

## Breaking Changes

### ✅ None! 

OpenAPI 3.2.0 is **fully backward compatible** with 3.1.x:

- ✅ All 3.1.x annotations work unchanged
- ✅ New fields are **optional** (`omitempty` in JSON/YAML)
- ✅ Parsers ignore unknown fields (per spec)
- ✅ No schema version conflicts

**Conversion to Swagger 2.0:**
- ⚠️ New features generate **warnings** (not errors)
- ⚠️ Extensions used where possible (e.g., `x-deprecated`)
- ✅ Existing features unaffected

---

## Migration Strategy

### Phase 1: Assessment (1-2 hours)

1. **Inventory your APIs:**
   ```bash
   # Count endpoints by HTTP method
   grep -r "@Router" . | grep -oP '\[(get|post|put|delete|patch|head|options)\]' | sort | uniq -c
   
   # Find streaming endpoints
   grep -r "text/event-stream\|application/stream" .
   
   # Find security definitions
   grep -r "@securityDefinitions" .
   ```

2. **Identify candidates for new features:**
   - POST endpoints that are idempotent → QUERY
   - SSE/streaming endpoints → `{stream}`
   - Deprecated auth methods → `@deprecated`
   - OAuth2 servers → `@oauth2metadataurl`

### Phase 2: Update nexs-swag (5 minutes)

```bash
# Update to latest version
go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest

# Verify version supports 3.2.0
nexs-swag --version
```

### Phase 3: Incremental Migration (1-3 days)

**Option A: One Feature at a Time**
```bash
# Day 1: Migrate QUERY methods
# Update annotations, regenerate, test

# Day 2: Add deprecation markers
# Mark old auth, regenerate, test

# Day 3: Document streaming
# Add {stream}, regenerate, test
```

**Option B: Per-Service Migration**
```bash
# Week 1: Migrate user-service
# Week 2: Migrate product-service
# Week 3: Migrate order-service
```

### Phase 4: Validation (2-4 hours)

```bash
# Regenerate OpenAPI specs
nexs-swag init -g main.go

# Validate against spec
docker run --rm -v $(pwd):/local openapitools/openapi-generator-cli validate \
  -i /local/docs/openapi.json

# Test with client generators
docker run --rm -v $(pwd):/local openapitools/openapi-generator-cli generate \
  -i /local/docs/openapi.json \
  -g typescript-fetch \
  -o /local/client
```

---

## Feature Adoption Guide

### QUERY Method Migration

**Identify candidates:**
```go
// ❌ Idempotent POST (should be QUERY)
// @Router /analytics/report [post]
func GenerateReport(c *gin.Context) {
    // Always returns same result for same input
    // No side effects, no data modification
}
```

**Migration:**
```diff
- // @Router /analytics/report [post]
+ // @Router /analytics/report [query]
```

**Testing:**
```bash
# Test idempotency
curl -X QUERY http://localhost:8080/analytics/report \
  -H "Content-Type: application/json" \
  -d '{"start":"2024-01-01","end":"2024-12-31"}' > report1.json

curl -X QUERY http://localhost:8080/analytics/report \
  -H "Content-Type: application/json" \
  -d '{"start":"2024-01-01","end":"2024-12-31"}' > report2.json

diff report1.json report2.json  # Should be identical
```

---

### Security Deprecation Migration

**Step 1: Mark as deprecated**
```diff
  // @securityDefinitions.apikey ApiKeyAuth
  // @in header
  // @name X-API-Key
+ // @deprecated true
+ // @description This authentication method will be removed on 2025-12-31. Migrate to OAuth2.
```

**Step 2: Add replacement**
```go
// @securityDefinitions.oauth2.application OAuth2
// @tokenUrl https://auth.example.com/token
// @scopes.read Read access
// @scopes.write Write access
```

**Step 3: Dual support period**
```go
// @Security ApiKeyAuth || OAuth2
// @Router /users [get]
func GetUsers(c *gin.Context) {
    // Accept both auth methods during transition
}
```

**Step 4: Monitor usage**
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if apiKey := c.GetHeader("X-API-Key"); apiKey != "" {
            log.Warn("Deprecated API key used", "endpoint", c.Request.URL.Path)
            metrics.DeprecatedAuthUsage.Inc()
        }
    }
}
```

**Step 5: Remove after sunset date**
```diff
- // @securityDefinitions.apikey ApiKeyAuth
- // @Security ApiKeyAuth || OAuth2
+ // @Security OAuth2
```

---

### Streaming Migration

**Before:**
```go
// @Success 200 {array} Event "Stream of events"
// @Produce text/event-stream
func StreamEvents(c *gin.Context) {
    // OpenAPI spec doesn't capture streaming nature
}
```

**After:**
```go
// @Success 200 {stream} Event "Stream of events"
// @Produce text/event-stream
func StreamEvents(c *gin.Context) {
    c.Stream(func(w io.Writer) bool {
        // Each message validated against Event schema
        event := Event{ID: uuid.New(), Data: "..."}
        c.SSEvent("message", event)
        return true
    })
}
```

**Generated Client (TypeScript):**
```typescript
// Before: unclear type
const response = await fetch('/events');
const events = await response.json(); // any[]?

// After: proper streaming client
const eventSource = new EventSource('/events');
eventSource.onmessage = (event: MessageEvent<Event>) => {
  const data: Event = JSON.parse(event.data); // ✅ Typed!
};
```

---

### OAuth2 Metadata URL Migration

**Benefits:**
- ✅ Auto-discovery of token endpoint, scopes, grant types
- ✅ Dynamic configuration updates
- ✅ Standardized across OAuth2 providers (Google, Okta, Auth0)

**Implementation:**
```diff
  // @securityDefinitions.oauth2.application OAuth2App
  // @tokenUrl https://auth.example.com/oauth/token
+ // @oauth2metadataurl https://auth.example.com/.well-known/oauth-authorization-server
```

**Server Setup (Example with Okta):**
```bash
# Discovery endpoint
curl https://{yourOktaDomain}/.well-known/oauth-authorization-server

# Returns:
{
  "issuer": "https://{yourOktaDomain}",
  "authorization_endpoint": "https://{yourOktaDomain}/oauth2/v1/authorize",
  "token_endpoint": "https://{yourOktaDomain}/oauth2/v1/token",
  "device_authorization_endpoint": "https://{yourOktaDomain}/oauth2/v1/device/authorize",
  "introspection_endpoint": "https://{yourOktaDomain}/oauth2/v1/introspect",
  "scopes_supported": ["openid", "profile", "email", "offline_access"]
}
```

---

## Conversion Tools

### nexs-swag CLI

```bash
# Generate OpenAPI 3.2.0 (default)
nexs-swag init -g main.go -o ./docs

# Generate Swagger 2.0 (with conversion warnings)
nexs-swag init -g main.go -o ./docs/v2 --openapi-version 2.0
```

### OpenAPI Generator

```bash
# Validate 3.2.0 spec
docker run --rm -v $(pwd):/local openapitools/openapi-generator-cli validate \
  -i /local/docs/openapi.json

# Convert to Swagger 2.0 (external tool)
npm install -g swagger2openapi
swagger2openapi --outfile swagger2.json openapi32.json --reverse
```

### Redocly CLI

```bash
# Lint and validate
npx @redocly/cli lint docs/openapi.yaml

# Bundle multi-file specs
npx @redocly/cli bundle docs/openapi.yaml -o openapi-bundled.yaml
```

### Spectral (Linting)

```bash
# Install
npm install -g @stoplight/spectral-cli

# Create ruleset (.spectral.yaml)
extends: [[spectral:oas, all]]
rules:
  operation-query-method:
    description: Ensure QUERY methods are idempotent
    severity: warn
    given: $.paths[*].query
    then:
      field: description
      function: pattern
      functionOptions:
        match: "(?i)(idempotent|safe|cacheable)"

# Run
spectral lint docs/openapi.yaml
```

---

## Common Issues & Solutions

### Issue 1: QUERY Method Not Supported in Client

**Problem:**
```
Error: Unknown HTTP method "QUERY" in generated client
```

**Cause:** Client generator doesn't support QUERY (OpenAPI 3.2.0 feature).

**Solutions:**

**Option A: Use OpenAPI Generator 7.0+**
```bash
docker run --rm -v $(pwd):/local \
  openapitools/openapi-generator-cli:v7.1.0 generate \
  -i /local/docs/openapi.json \
  -g typescript-fetch
```

**Option B: Fallback to POST**
```diff
- // @Router /search [query]
+ // @Router /search [post]
+ // @Description This endpoint is idempotent and safe (should be QUERY in OpenAPI 3.2.0)
```

**Option C: Custom generator template**
```java
// Extend generator to map QUERY → POST for HTTP clients
public class CustomGenerator extends TypeScriptFetchClientCodegen {
    @Override
    public String toApiName(String name) {
        if ("QUERY".equals(name)) return "POST";
        return super.toApiName(name);
    }
}
```

---

### Issue 2: Streaming Schema Validation Fails

**Problem:**
```
Client expects array, receives stream of objects
```

**Cause:** Confusion between `{array}` and `{stream}`.

**Solution:**
```diff
- // @Success 200 {array} Event  // ❌ Implies single response with array
+ // @Success 200 {stream} Event  // ✅ Continuous stream of Event objects
```

**Client Implementation:**
```typescript
// Before (wrong)
const response = await fetch('/events');
const events: Event[] = await response.json(); // ❌ Expects complete array

// After (correct)
const eventSource = new EventSource('/events');
eventSource.onmessage = (msg) => {
  const event: Event = JSON.parse(msg.data); // ✅ Each message is one Event
};
```

---

### Issue 3: Deprecated Security Not Visible

**Problem:**
Clients still use deprecated API keys without warnings.

**Solution:**

**1. Ensure proper annotation:**
```go
// @securityDefinitions.apikey ApiKeyAuth
// @deprecated true  // ✅ Must be present
```

**2. Add sunset header:**
```go
func DeprecatedAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if apiKey := c.GetHeader("X-API-Key"); apiKey != "" {
            // RFC 8594: Sunset header
            c.Header("Sunset", "Sat, 31 Dec 2025 23:59:59 GMT")
            c.Header("Warning", `299 - "API Key authentication is deprecated. Use OAuth2."`)
        }
    }
}
```

**3. Update client to check warnings:**
```typescript
const response = await fetch('/api/users');
const sunset = response.headers.get('Sunset');
if (sunset) {
  console.warn(`This auth method will be removed on ${new Date(sunset)}`);
}
```

---

### Issue 4: OAuth2 Metadata URL Not Used

**Problem:**
Generated clients ignore `oauth2MetadataUrl`.

**Cause:** Client generator version doesn't support RFC 8414 discovery.

**Solution:**

**Verify client support:**
```bash
# Check generator version
openapi-generator-cli version-manager list

# Use version 7.2.0+ for OAuth2 metadata support
openapi-generator-cli version-manager set 7.2.0
```

**Manual discovery fallback:**
```typescript
// Fetch metadata before creating OAuth2 client
const metadata = await fetch('https://auth.example.com/.well-known/oauth-authorization-server');
const config = await metadata.json();

const oauth = new OAuth2Client({
  tokenEndpoint: config.token_endpoint,
  authorizationEndpoint: config.authorization_endpoint,
  scopes: config.scopes_supported,
});
```

---

## Best Practices

### 1. Version Your API Specs

```bash
# Tag releases
git tag v1.0.0-openapi-3.1.0
git tag v1.1.0-openapi-3.2.0

# Keep multiple versions
docs/
  v1.0/openapi.yaml  # OpenAPI 3.1.0
  v1.1/openapi.yaml  # OpenAPI 3.2.0
  v2.0/swagger.json  # Swagger 2.0 (converted)
```

### 2. Document Migration in Changelog

```markdown
## [1.1.0] - 2025-12-15

### Changed
- Migrated to OpenAPI 3.2.0
- `/search` endpoint now uses QUERY method (was POST)
- Marked API Key auth as deprecated (remove 2025-12-31)

### Added
- OAuth2 metadata URL for automatic client configuration
- Streaming schema for `/events` endpoint

### Deprecated
- API Key authentication (use OAuth2 instead)
```

### 3. Gradual Rollout

```go
// Feature flag for QUERY method
if config.EnableQueryMethod {
    // @Router /search [query]
} else {
    // @Router /search [post]
}
```

### 4. Monitor Client Versions

```go
func AnalyticsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userAgent := c.GetHeader("User-Agent")
        
        // Track client versions
        metrics.ClientVersion.WithLabelValues(
            parseClientVersion(userAgent),
        ).Inc()
        
        // Warn old clients about QUERY method
        if strings.Contains(c.Request.Method, "POST") &&
           shouldBeQuery(c.Request.URL.Path) {
            log.Info("Client using POST instead of QUERY",
                "client", userAgent,
                "path", c.Request.URL.Path)
        }
    }
}
```

### 5. Test with Multiple Clients

```bash
# Test matrix
clients=(
  "typescript-fetch"
  "python"
  "go"
  "java"
  "csharp-netcore"
)

for client in "${clients[@]}"; do
  echo "Testing $client client..."
  openapi-generator-cli generate \
    -i docs/openapi.json \
    -g $client \
    -o clients/$client
  
  cd clients/$client && npm test && cd ../..
done
```

---

## FAQ

### Q1: Do I need to migrate immediately?

**A:** No. OpenAPI 3.2.0 is backward compatible. Migrate when:
- Starting new projects
- Adding streaming/SSE endpoints
- Deprecating authentication methods
- Using OAuth2 device flow

---

### Q2: Can I use 3.2.0 features with Swagger 2.0 output?

**A:** Yes, nexs-swag converts gracefully:
- QUERY → warning, operation omitted
- `deprecated` → `x-deprecated` extension
- `itemSchema` → warning
- Existing features work unchanged

---

### Q3: Will old clients break?

**A:** No, if you:
- ✅ Keep supporting old HTTP methods during transition
- ✅ Use dual security (`ApiKeyAuth || OAuth2`)
- ✅ Set sunset dates for deprecated features
- ✅ Version your API specs

---

### Q4: How do I test QUERY method locally?

```bash
# curl (version 7.85+)
curl -X QUERY http://localhost:8080/search \
  -H "Content-Type: application/json" \
  -d '{"term":"laptop"}'

# httpie
http QUERY localhost:8080/search term=laptop

# Postman: Use custom HTTP method
```

---

### Q5: Can I mix 3.1.x and 3.2.0 features?

**A:** Yes! All 3.1.x features work in 3.2.0 specs. Use new features where beneficial:
```go
// Mix old and new
// @Router /legacy [get]          // ✅ 3.1.x style
// @Router /search [query]        // ✅ 3.2.0 feature
// @Success 200 {object} Result   // ✅ 3.1.x style
// @Success 200 {stream} Event    // ✅ 3.2.0 feature
```

---

## Additional Resources

- [OpenAPI 3.2.0 Specification](https://spec.openapis.org/oas/v3.2.0.html)
- [RFC 8414: OAuth 2.0 Authorization Server Metadata](https://datatracker.ietf.org/doc/html/rfc8414)
- [RFC 8628: OAuth 2.0 Device Authorization Grant](https://datatracker.ietf.org/doc/html/rfc8628)
- [nexs-swag Examples](./examples/)
- [OpenAPI Generator Support Matrix](https://openapi-generator.tech/docs/generators)

---

**Need Help?** Open an issue at [github.com/fsvxavier/nexs-swag/issues](https://github.com/fsvxavier/nexs-swag/issues)
