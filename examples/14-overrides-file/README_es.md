# Ejemplo 14 - Overrides File

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de un archivo de overrides para modificar la spec generada.

## Flag

```bash
--overridesFile <path>
--of <path>
```

## Concepto

Permite sobrescribir partes de la spec OpenAPI sin modificar el código Go.

## Uso

```bash
nexs-swag init --overridesFile overrides.yaml
```

## Formato

```yaml
# overrides.yaml
info:
  title: "Custom Title"
  description: "Custom Description"
  version: "2.0.0"

paths:
  /users:
    get:
      summary: "Custom Summary"
      security:
        - bearerAuth: []

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

## Cómo Ejecutar

```bash
./run.sh
```

## Casos de Uso

### 1. Environment-Specific
```yaml
# dev.overrides.yaml
servers:
  - url: http://localhost:8080
    description: Development

# prod.overrides.yaml  
servers:
  - url: https://api.example.com
    description: Production
```

```bash
# Development
nexs-swag init --of dev.overrides.yaml

# Production
nexs-swag init --of prod.overrides.yaml
```

### 2. Security Schemes
```yaml
components:
  securitySchemes:
    apiKey:
      type: apiKey
      in: header
      name: X-API-Key
    
    oauth2:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://auth.example.com/oauth/authorize
          tokenUrl: https://auth.example.com/oauth/token
          scopes:
            read: Read access
            write: Write access

security:
  - apiKey: []
  - oauth2: [read, write]
```

### 3. Adicionar Endpoints Externos
```yaml
paths:
  /health:
    get:
      summary: Health check
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
```

### 4. Modificar Responses
```yaml
paths:
  /users:
    get:
      responses:
        '200':
          description: Success
          headers:
            X-RateLimit-Limit:
              schema:
                type: integer
              description: Request limit per hour
            X-RateLimit-Remaining:
              schema:
                type: integer
              description: Remaining requests
```

### 5. Tags y Metadata
```yaml
tags:
  - name: users
    description: User management operations
    externalDocs:
      description: User API Guide
      url: https://docs.example.com/users
  
  - name: orders
    description: Order processing
    externalDocs:
      url: https://docs.example.com/orders

externalDocs:
  description: Complete API Documentation
  url: https://docs.example.com
```

## Merge Strategy

nexs-swag hace deep merge:

### Generated
```yaml
info:
  title: "My API"
  version: "1.0.0"
paths:
  /users:
    get:
      summary: "Get users"
```

### Overrides
```yaml
info:
  version: "2.0.0"
paths:
  /users:
    get:
      security:
        - bearerAuth: []
```

### Result
```yaml
info:
  title: "My API"        # kept from generated
  version: "2.0.0"       # overridden
paths:
  /users:
    get:
      summary: "Get users"  # kept from generated
      security:             # added from overrides
        - bearerAuth: []
```

## Tips

### 1. Versionamiento
```
overrides/
├── base.yaml           # Common config
├── v1.yaml            # v1 specific
└── v2.yaml            # v2 specific
```

```bash
nexs-swag init --of overrides/v2.yaml
```

### 2. Combinar Múltiples
```bash
# Unix/Linux
cat base.yaml v2.yaml > combined.yaml
nexs-swag init --of combined.yaml

# O use yq
yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' \
  base.yaml v2.yaml > combined.yaml
```

### 3. Validación
```bash
# Valide el YAML antes
yamllint overrides.yaml

# O use yq
yq eval overrides.yaml
```

## Cuándo Usar

**Use overrides cuando:**
- Info varía por ambiente
- Security schemes complejos
- Metadata adicional
- Endpoints externos no en código
- CI/CD con múltiples configs

**NO use cuando:**
- Info puede estar en código
- Solo un ambiente
- Spec simple
- Preferencia por single source of truth en código
