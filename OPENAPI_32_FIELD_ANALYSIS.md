# Análise Técnica: Campos OpenAPI 3.2.0 - Status de Implementação

**Data da Análise:** 15 de dezembro de 2025  
**Versão do nexs-swag:** Atual (main branch)  
**Campos Analisados:** 6 campos específicos introduzidos no OpenAPI 3.2.0

---

## Resumo Executivo

| Campo | Status | pkg/parser | pkg/converter | pkg/format |
|-------|--------|------------|---------------|------------|
| **PathItem.Query** | ❌ Não implementado | ❌ | ❌ | ⚠️ |
| **SecurityScheme.Deprecated** | ✅ Implementado | ✅ | ⚠️ | ✅ |
| **SecurityScheme.OAuth2MetadataURL** | ✅ Implementado | ✅ | ⚠️ | ✅ |
| **OAuthFlows.DeviceAuthorization** | ✅ Implementado | ✅ | ⚠️ | ✅ |
| **MediaType.ItemSchema** | ✅ Implementado | ❌ | ✅ | ✅ |
| **MediaType.ItemEncoding** | ✅ Implementado | ❌ | ✅ | ✅ |

**Legenda:**
- ✅ Completamente implementado
- ⚠️ Parcialmente implementado
- ❌ Não implementado

---

## Análise Detalhada por Campo

### 1. PathItem.Query (método HTTP QUERY)

**Status Geral:** ❌ **NÃO IMPLEMENTADO**

#### Definição OpenAPI v3
```go
// Em pkg/openapi/v3/openapi.go linha 113
type PathItem struct {
    // ... outros métodos HTTP ...
    Query *Operation `json:"query,omitempty" yaml:"query,omitempty"` // QUERY operation (new in 3.2.0)
}
```

#### pkg/parser (operation.go)
**Status:** ❌ **Não implementado**

**Arquivo:** `pkg/parser/parser.go` linhas 340-360

**Problema:** O parser **não reconhece** o método QUERY no switch statement que processa `@Router`:

```go
// Métodos HTTP suportados atualmente:
switch strings.ToLower(routeInfo.Method) {
    case "get":
        pathItem.Get = op
    case "post":
        pathItem.Post = op
    case "put":
        pathItem.Put = op
    case "delete":
        pathItem.Delete = op
    case "patch":
        pathItem.Patch = op
    case "options":
        pathItem.Options = op
    case "head":
        pathItem.Head = op
    case "trace":
        pathItem.Trace = op
    // FALTA: case "query"
}
```

**Impacto:** 
- Anotações `@Router /path [QUERY]` são **ignoradas silenciosamente**
- Nenhum erro ou warning é gerado
- Operação QUERY não é adicionada ao PathItem

#### pkg/converter (converter.go)
**Status:** ❌ **Não implementado**

**Arquivo:** `pkg/converter/converter.go` linhas 160-198

**Problema:** O método `convertPathItem` **não converte** PathItem.Query:

```go
func (c *Converter) convertPathItem(pathItem *openapi.PathItem) *swagger.PathItem {
    // ... conversão de outros métodos ...
    if pathItem.Patch != nil {
        v2PathItem.Patch = c.convertOperation(pathItem.Patch)
    }
    // FALTA: Tratamento de pathItem.Query
    
    return v2PathItem
}
```

**Impacto:**
- Conversão V3→V2: QUERY é perdido (esperado, pois v2.0 não suporta)
- Conversão V2→V3: Não aplicável
- **Falta warning** informando que QUERY não é suportado em v2.0

#### pkg/format (format.go)
**Status:** ⚠️ **Parcialmente implementado**

**Arquivo:** `pkg/format/format.go`

O formatter não valida métodos HTTP específicos, então **preservaria** anotações `@Router ... [QUERY]` se existissem, mas não há validação ou normalização específica.

#### Solução Recomendada

**1. Atualizar pkg/parser/parser.go:**

```go
// Adicionar case "query" ao switch statement (linha ~358)
switch strings.ToLower(routeInfo.Method) {
    case "get":
        pathItem.Get = op
    case "post":
        pathItem.Post = op
    case "put":
        pathItem.Put = op
    case "delete":
        pathItem.Delete = op
    case "patch":
        pathItem.Patch = op
    case "options":
        pathItem.Options = op
    case "head":
        pathItem.Head = op
    case "trace":
        pathItem.Trace = op
    case "query":  // NOVO
        pathItem.Query = op
    default:
        // Log warning para métodos desconhecidos
        fmt.Printf("Warning: unknown HTTP method '%s' for path '%s'\n", routeInfo.Method, routeInfo.Path)
}
```

**2. Atualizar pkg/converter/converter.go:**

```go
// Adicionar no método convertPathItem (linha ~193)
func (c *Converter) convertPathItem(pathItem *openapi.PathItem) *swagger.PathItem {
    // ... código existente ...
    
    if pathItem.Patch != nil {
        v2PathItem.Patch = c.convertOperation(pathItem.Patch)
    }
    
    // NOVO: Tratar Query method
    if pathItem.Query != nil {
        c.warnings = append(c.warnings, 
            "QUERY method is not supported in Swagger 2.0 (OpenAPI 3.2.0+ only)")
        // Opcionalmente, converter para GET com extension x-method-query
        // ou simplesmente ignorar
    }
    
    return v2PathItem
}
```

**3. Atualizar validação em pkg/parser/parser.go:**

```go
// No método validateOperation, adicionar Query à lista (linha ~433)
operations := []*openapi.Operation{
    pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Delete,
    pathItem.Patch, pathItem.Options, pathItem.Head, pathItem.Trace,
    pathItem.Query, // NOVO
}
```

---

### 2. SecurityScheme.Deprecated (boolean)

**Status Geral:** ✅ **IMPLEMENTADO** (estrutura) / ⚠️ **Parcialmente implementado** (parser/converter)

#### Definição OpenAPI v3
```go
// Em pkg/openapi/v3/openapi.go linha 314
type SecurityScheme struct {
    // ... outros campos ...
    Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"` // Deprecated (new in 3.2.0)
}
```

#### pkg/parser
**Status:** ❌ **Não implementado no parser de anotações**

**Arquivo:** `pkg/parser/general_info.go` linhas 237-250

**Problema:** O parser de `@securityDefinitions` **não reconhece** o atributo `deprecated`:

```go
// Regex atual:
securityBasicRegex  = regexp.MustCompile(`^@securityDefinitions\.basic\s+(\S+)\s*(.*)$`)
securityAPIKeyRegex = regexp.MustCompile(`^@securityDefinitions\.apikey\s+(\S+)\s+(\w+)\s+(\w+)\s*(.*)$`)
securityOAuth2Regex = regexp.MustCompile(`^@securityDefinitions\.oauth2\.(\w+)\s+(\S+)\s*(.*)$`)

// Criação do SecurityScheme:
g.openapi.Components.SecuritySchemes[matches[1]] = &openapi.SecurityScheme{
    Type:        "http",
    Scheme:      "basic",
    Description: strings.TrimSpace(matches[2]),
    // FALTA: Deprecated não é parseado
}
```

**Uso esperado:**
```go
// @securityDefinitions.basic BasicAuth Basic Authentication (deprecated)
// @securityDefinitions.deprecated BasicAuth true
```

Mas não há suporte para isso atualmente.

#### pkg/converter
**Status:** ⚠️ **Parcialmente implementado**

**V3→V2 (Swagger 2.0):**
```go
// pkg/converter/converter.go linha 823
func (c *Converter) convertSecurityScheme(scheme *openapi.SecurityScheme) *swagger.SecurityScheme {
    // ...
    // Swagger 2.0 não tem campo Deprecated nativo
    // FALTA: Converter para extension x-deprecated
}
```

**V2→V3:**
```go
// pkg/converter/converter.go linha 1566
func (c *Converter) convertSecuritySchemeToV3(scheme *swagger.SecurityScheme) *openapi.SecurityScheme {
    // ...
    // FALTA: Ler extension x-deprecated e popular Deprecated
}
```

#### Solução Recomendada

**1. Estender parser de SecurityDefinitions:**

```go
// Em pkg/parser/general_info.go

// Adicionar novo regex para atributos de SecurityScheme
var securityAttrRegex = regexp.MustCompile(`^@securityDefinitions\.(\w+)\.(\w+)\s+(\S+)$`)

// No método Process, adicionar:
case securityAttrRegex.MatchString(text):
    matches := securityAttrRegex.FindStringSubmatch(text)
    schemeName := matches[1]
    attrName := matches[2]
    attrValue := matches[3]
    
    if scheme, exists := g.openapi.Components.SecuritySchemes[schemeName]; exists {
        switch strings.ToLower(attrName) {
        case "deprecated":
            if attrValue == "true" {
                scheme.Deprecated = true
            }
        case "oauth2metadataurl":
            scheme.OAuth2MetadataURL = attrValue
        }
    }
```

**Uso:**
```go
// @securityDefinitions.basic BasicAuth
// @securityDefinitions.BasicAuth.deprecated true
```

**2. Atualizar converter V3→V2:**

```go
// Em pkg/converter/converter.go, método convertSecurityScheme
func (c *Converter) convertSecurityScheme(scheme *openapi.SecurityScheme) *swagger.SecurityScheme {
    // ... código existente ...
    
    // Adicionar suporte para Deprecated
    if scheme.Deprecated {
        if v2Scheme.Extensions == nil {
            v2Scheme.Extensions = make(map[string]interface{})
        }
        v2Scheme.Extensions["x-deprecated"] = true
        c.warnings = append(c.warnings, 
            "SecurityScheme.deprecated is not natively supported in Swagger 2.0, converted to x-deprecated extension")
    }
    
    return v2Scheme
}
```

**3. Atualizar converter V2→V3:**

```go
// Método convertSecuritySchemeToV3
func (c *Converter) convertSecuritySchemeToV3(scheme *swagger.SecurityScheme) *openapi.SecurityScheme {
    // ... código existente ...
    
    // Ler x-deprecated extension
    if scheme.Extensions != nil {
        if deprecated, ok := scheme.Extensions["x-deprecated"].(bool); ok {
            v3Scheme.Deprecated = deprecated
        }
    }
    
    return v3Scheme
}
```

---

### 3. SecurityScheme.OAuth2MetadataURL (string)

**Status Geral:** ✅ **IMPLEMENTADO** (estrutura) / ❌ **Não implementado** (parser/converter)

#### Definição OpenAPI v3
```go
// Em pkg/openapi/v3/openapi.go linha 315
type SecurityScheme struct {
    // ... outros campos ...
    OAuth2MetadataURL string `json:"oauth2MetadataUrl,omitempty" yaml:"oauth2MetadataUrl,omitempty"` // OAuth2 metadata URL (new in 3.2.0)
}
```

#### Análise
A estrutura está definida, mas:
- **Parser:** Não há suporte para ler este campo via anotações
- **Converter V3→V2:** Não converte (esperado, v2.0 não suporta)
- **Converter V2→V3:** Não aplicável

#### Solução Recomendada

**Usar a mesma abordagem de SecurityScheme.Deprecated acima:**

```go
// @securityDefinitions.oauth2.password OAuth2Password
// @securityDefinitions.OAuth2Password.oauth2metadataurl https://auth.example.com/.well-known/oauth-authorization-server
```

**Converter V3→V2:**
```go
if scheme.OAuth2MetadataURL != "" {
    if v2Scheme.Extensions == nil {
        v2Scheme.Extensions = make(map[string]interface{})
    }
    v2Scheme.Extensions["x-oauth2-metadata-url"] = scheme.OAuth2MetadataURL
    c.warnings = append(c.warnings, 
        "OAuth2MetadataURL is not supported in Swagger 2.0, converted to x-oauth2-metadata-url extension")
}
```

---

### 4. OAuthFlows.DeviceAuthorization (OAuthFlow)

**Status Geral:** ✅ **IMPLEMENTADO** (estrutura) / ❌ **Não implementado** (parser/converter)

#### Definição OpenAPI v3
```go
// Em pkg/openapi/v3/openapi.go linha 324
type OAuthFlows struct {
    Implicit            *OAuthFlow `json:"implicit,omitempty"`
    Password            *OAuthFlow `json:"password,omitempty"`
    ClientCredentials   *OAuthFlow `json:"clientCredentials,omitempty"`
    AuthorizationCode   *OAuthFlow `json:"authorizationCode,omitempty"`
    DeviceAuthorization *OAuthFlow `json:"deviceAuthorization,omitempty"` // Device authorization flow (new in 3.2.0)
}
```

#### pkg/parser
**Status:** ❌ **Não implementado**

Não há regex ou processamento para `@securityDefinitions.oauth2.deviceAuthorization`.

#### pkg/converter
**Status:** ⚠️ **Parcialmente implementado**

**V3→V2:**
```go
// pkg/converter/converter.go linhas 862-895
func (c *Converter) convertOAuth2Flows(scheme *openapi.SecurityScheme, v2Scheme *swagger.SecurityScheme) {
    // Prioridade: implicit > password > application > authorizationCode
    // FALTA: Tratar DeviceAuthorization
}
```

Atualmente, DeviceAuthorization é **silenciosamente ignorado**.

**V2→V3:**
Não aplicável (Swagger 2.0 não tem device authorization).

#### Solução Recomendada

**1. Parser:**

```go
// Em pkg/parser/general_info.go

// Adicionar novo regex:
securityOAuth2DeviceRegex = regexp.MustCompile(`^@securityDefinitions\.oauth2\.deviceAuthorization\s+(\S+)\s+(\S+)\s*(.*)$`)

// No Process:
case securityOAuth2DeviceRegex.MatchString(text):
    matches := securityOAuth2DeviceRegex.FindStringSubmatch(text)
    schemeName := matches[1]
    tokenURL := matches[2]
    description := strings.TrimSpace(matches[3])
    
    if _, exists := g.openapi.Components.SecuritySchemes[schemeName]; !exists {
        g.openapi.Components.SecuritySchemes[schemeName] = &openapi.SecurityScheme{
            Type:  "oauth2",
            Flows: &openapi.OAuthFlows{},
        }
    }
    
    scheme := g.openapi.Components.SecuritySchemes[schemeName]
    if scheme.Flows == nil {
        scheme.Flows = &openapi.OAuthFlows{}
    }
    
    scheme.Flows.DeviceAuthorization = &openapi.OAuthFlow{
        TokenURL: tokenURL,
        Scopes:   make(map[string]string),
    }
    scheme.Description = description
```

**2. Converter V3→V2:**

```go
func (c *Converter) convertOAuth2Flows(scheme *openapi.SecurityScheme, v2Scheme *swagger.SecurityScheme) {
    if scheme.Flows == nil {
        return
    }
    
    // ... código existente ...
    
    // Adicionar após authorizationCode:
    if scheme.Flows.DeviceAuthorization != nil {
        c.warnings = append(c.warnings, 
            "OAuth2 Device Authorization flow is not supported in Swagger 2.0 (RFC 8628)")
        // Opcionalmente, converter para extension
        if v2Scheme.Extensions == nil {
            v2Scheme.Extensions = make(map[string]interface{})
        }
        v2Scheme.Extensions["x-device-authorization-token-url"] = scheme.Flows.DeviceAuthorization.TokenURL
    }
    
    // Atualizar contador de flows
    if scheme.Flows.DeviceAuthorization != nil {
        flowCount++
    }
}
```

---

### 5. MediaType.ItemSchema (Schema para streaming)

**Status Geral:** ✅ **IMPLEMENTADO** (estrutura) / ❌ **Não usado** (parser) / ✅ **Mantido** (converter)

#### Definição OpenAPI v3
```go
// Em pkg/openapi/v3/openapi.go linha 162
type MediaType struct {
    Schema       *Schema              `json:"schema,omitempty"`
    Example      interface{}          `json:"example,omitempty"`
    Examples     map[string]*Example  `json:"examples,omitempty"`
    Encoding     map[string]*Encoding `json:"encoding,omitempty"`
    ItemSchema   *Schema              `json:"itemSchema,omitempty"`   // Schema for streaming items (new in 3.2.0)
    ItemEncoding map[string]*Encoding `json:"itemEncoding,omitempty"` // Encoding for streaming items (new in 3.2.0)
}
```

#### pkg/parser
**Status:** ❌ **Não suportado**

O parser não tem anotações para especificar `ItemSchema` ou `ItemEncoding`. Atualmente, apenas `Schema` é suportado via:

```go
// @Success 200 {object} ResponseType
// @Param body body RequestType true "description"
```

**Não há suporte para:**
```go
// Exemplo hipotético (não implementado):
// @Success 200 {stream} EventType "Server-sent events stream"
```

#### pkg/converter
**Status:** ✅ **Preservado na conversão**

As estruturas `MediaType` são copiadas diretamente ou convertidas field-by-field, então `ItemSchema` e `ItemEncoding` são **preservados** se estiverem presentes no JSON de entrada.

**V3→V2:**
```go
// ItemSchema e ItemEncoding são ignorados (não existem em v2.0)
// Deveria haver warning
```

**V2→V3:**
Não aplicável.

#### Solução Recomendada

**1. Parser: Adicionar suporte para streaming types**

```go
// Em pkg/parser/operation.go

// Adicionar novo tipo de resposta "stream":
// @Success 200 {stream} EventType "SSE stream of events"

func (o *OperationProcessor) processResponse(text string, regex *regexp.Regexp, op *openapi.Operation) {
    matches := regex.FindStringSubmatch(text)
    // ...
    
    responseType := matches[2]
    schemaRef := matches[3]
    
    // Adicionar case para stream:
    if responseType == "stream" {
        mediaType := &openapi.MediaType{
            ItemSchema: o.parseSchemaType(schemaRef),
        }
        response.Content["text/event-stream"] = mediaType
    } else if responseType == typeObject || responseType == typeArray {
        // ... código existente ...
    }
}
```

**2. Converter: Adicionar warnings**

```go
// Em pkg/converter/converter.go

func (c *Converter) convertMediaType(mt *openapi.MediaType) interface{} {
    if mt.ItemSchema != nil {
        c.warnings = append(c.warnings, 
            "MediaType.itemSchema (streaming) is not supported in Swagger 2.0")
    }
    if mt.ItemEncoding != nil {
        c.warnings = append(c.warnings, 
            "MediaType.itemEncoding (streaming) is not supported in Swagger 2.0")
    }
    // ... conversão normal ...
}
```

---

### 6. MediaType.ItemEncoding (map[string]*Encoding para streaming)

**Status:** Igual ao `ItemSchema` (item 5 acima)

Mesmo status e soluções aplicam-se a `ItemEncoding`.

---

## Plano de Implementação Prioritário

### Prioridade ALTA (Impacto em funcionalidade básica)

1. **PathItem.Query**
   - Arquivos: `pkg/parser/parser.go`, `pkg/converter/converter.go`
   - Esforço: 2-3 horas
   - Impacto: Alto (método HTTP não funciona)

### Prioridade MÉDIA (Campos menos utilizados)

2. **SecurityScheme.Deprecated**
   - Arquivos: `pkg/parser/general_info.go`, `pkg/converter/converter.go`
   - Esforço: 3-4 horas
   - Impacto: Médio (recurso de documentação)

3. **OAuthFlows.DeviceAuthorization**
   - Arquivos: `pkg/parser/general_info.go`, `pkg/converter/converter.go`
   - Esforço: 4-5 horas
   - Impacto: Médio (flow OAuth2 moderno)

### Prioridade BAIXA (Features avançadas)

4. **MediaType.ItemSchema / ItemEncoding**
   - Arquivos: `pkg/parser/operation.go`, `pkg/converter/converter.go`
   - Esforço: 6-8 horas
   - Impacto: Baixo (streaming é caso de uso avançado)

5. **SecurityScheme.OAuth2MetadataURL**
   - Arquivos: `pkg/parser/general_info.go`, `pkg/converter/converter.go`
   - Esforço: 2-3 horas
   - Impacto: Baixo (discovery automático de OAuth2)

---

## Arquivos que Precisam de Modificação

### 1. pkg/parser/parser.go
- Adicionar suporte para `case "query"` no switch de métodos HTTP
- Adicionar `pathItem.Query` na validação de operações

### 2. pkg/parser/general_info.go
- Adicionar regex para atributos de SecurityScheme (deprecated, oauth2MetadataURL)
- Adicionar processamento de `@securityDefinitions.oauth2.deviceAuthorization`

### 3. pkg/parser/operation.go
- Adicionar suporte para tipo de resposta `{stream}` (ItemSchema)

### 4. pkg/converter/converter.go
- Método `convertPathItem`: tratar Query method
- Método `convertSecurityScheme`: converter Deprecated e OAuth2MetadataURL
- Método `convertOAuth2Flows`: tratar DeviceAuthorization
- Adicionar warnings para ItemSchema/ItemEncoding

### 5. pkg/format/format.go
- Adicionar `@query` às anotações reconhecidas (se aplicável)

---

## Testes Recomendados

Para cada implementação, criar testes em:

1. **pkg/parser/parser_test.go**
   ```go
   func TestQueryMethod(t *testing.T) {
       // Testar @Router /search [QUERY]
   }
   ```

2. **pkg/parser/general_info_test.go**
   ```go
   func TestSecuritySchemeDeprecated(t *testing.T) {
       // Testar @securityDefinitions.BasicAuth.deprecated true
   }
   ```

3. **pkg/converter/converter_test.go**
   ```go
   func TestConvertQueryMethodToV2(t *testing.T) {
       // Verificar warning quando Query existe
   }
   ```

---

## Conclusão

A maioria dos campos OpenAPI 3.2.0 está **estruturalmente implementada** nas definições de tipos, mas **não é utilizada** pelo parser de anotações nem pelo converter.

**Próximos passos:**
1. Implementar suporte ao método QUERY (prioritário)
2. Adicionar parsing de atributos de SecurityScheme
3. Adicionar warnings apropriados no converter
4. Criar testes de integração

**Compatibilidade:**
- Todas as implementações devem gerar **warnings apropriados** ao converter para Swagger 2.0
- Campos não suportados devem ser convertidos para extensions `x-*` quando possível
