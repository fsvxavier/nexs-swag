# Análise de Suporte a OpenAPI 2.0 / Swagger 2.0

## 📋 Resumo Executivo

**Status:** ✅ **VIÁVEL** - Implementação possível com arquitetura modular

**Complexidade:** 🟡 Média (estimativa: 40-60 horas de desenvolvimento)

**Compatibilidade:** Pode coexistir com OpenAPI 3.1.0 atual sem quebrar funcionalidades existentes

---

## 🎯 Objetivo

Adicionar suporte para geração de especificações OpenAPI 2.0 (Swagger 2.0) além da versão padrão atual 3.1.0, permitindo que usuários escolham qual formato gerar via flag CLI.

## 📊 Análise da Estrutura Atual

### Arquitetura Existente

```
nexs-swag/
├── pkg/
│   ├── openapi/          # ✅ Estruturas OpenAPI 3.1.0
│   ├── parser/           # ✅ Parser de anotações (agnóstico)
│   ├── generator/        # ⚠️ Gerador específico para 3.1.0
│   └── format/           # ✅ Formatador (agnóstico)
```

**Componentes Agnósticos (não precisam mudanças):**
- ✅ `pkg/parser` - Parse de anotações Go (independente de versão)
- ✅ `pkg/format` - Formatação de comentários
- ✅ CLI flags e comandos principais

**Componentes Dependentes (precisam adaptação):**
- ⚠️ `pkg/openapi` - Estruturas específicas do OpenAPI 3.1.0
- ⚠️ `pkg/generator` - Gerador hardcoded para 3.1.0

### Análise de Código Chave

#### 1. Estrutura OpenAPI (`pkg/openapi/openapi.go`)

```go
// Estrutura atual - OpenAPI 3.1.0
type OpenAPI struct {
    OpenAPI           string                // "3.1.0"
    Info              Info                  
    JSONSchemaDialect string                // ❌ Não existe em 2.0
    Servers           []Server              // ❌ Não existe em 2.0 (usa host+basePath)
    Paths             Paths                 
    Webhooks          map[string]*PathItem  // ❌ Não existe em 2.0
    Components        *Components           // ⚠️ Em 2.0 é "definitions"
    Security          []SecurityRequirement 
    Tags              []Tag                 
    ExternalDocs      *ExternalDocs         
}
```

**Incompatibilidades principais:**
- `servers` → Substituído por `host`, `basePath`, `schemes` em 2.0
- `webhooks` → Não existe em 2.0
- `components` → Chamado de `definitions` em 2.0
- `JSONSchemaDialect` → Não existe em 2.0
- JSON Schema 2020-12 → JSON Schema Draft 4 em 2.0

#### 2. Gerador (`pkg/generator/generator.go`)

```go
func New(spec *openapi.OpenAPI, outputDir string, outputType []string) *Generator {
    // Atualmente aceita apenas openapi.OpenAPI (3.1.0)
}
```

**Necessita:**
- Interface genérica ou tipo union para aceitar ambas versões
- Lógica de serialização diferenciada por versão

## 🏗️ Proposta de Arquitetura

### Opção 1: Estruturas Separadas + Interface (RECOMENDADO)

```
pkg/
├── openapi/
│   ├── v3/
│   │   └── openapi.go      # Estruturas OpenAPI 3.1.0 (atual)
│   ├── v2/
│   │   └── swagger.go      # Estruturas Swagger 2.0 (novo)
│   └── spec.go             # Interface comum
├── generator/
│   ├── generator.go        # Gerador genérico
│   ├── v3/
│   │   └── generator.go    # Gerador específico 3.1.0
│   └── v2/
│       └── generator.go    # Gerador específico 2.0
└── converter/              # Opcional: converter entre versões
    └── converter.go
```

**Vantagens:**
- ✅ Separação clara de responsabilidades
- ✅ Não quebra código existente
- ✅ Fácil manutenção de cada versão
- ✅ Permite evolução independente

**Interface Proposta:**

```go
// pkg/openapi/spec.go
package openapi

type Specification interface {
    GetVersion() string
    Validate() error
    MarshalJSON() ([]byte, error)
    MarshalYAML() ([]byte, error)
}

// pkg/openapi/v3/openapi.go
func (o *OpenAPI) GetVersion() string { return "3.1.0" }

// pkg/openapi/v2/swagger.go  
func (s *Swagger) GetVersion() string { return "2.0" }
```

### Opção 2: Estrutura Unificada com Tags (NÃO RECOMENDADO)

Usar uma única estrutura com tags JSON condicionais.

**Desvantagens:**
- ❌ Código complexo e difícil de manter
- ❌ Validação complicada
- ❌ Conflitos de nomenclatura (components vs definitions)

## 🔄 Mapeamento OpenAPI 3.1.0 → 2.0

### Campos Diretos (Compatíveis)

| OpenAPI 3.1.0 | Swagger 2.0 | Notas |
|---------------|-------------|-------|
| `info` | `info` | ✅ Compatível |
| `paths` | `paths` | ✅ Compatível |
| `tags` | `tags` | ✅ Compatível |
| `externalDocs` | `externalDocs` | ✅ Compatível |
| `security` | `security` | ✅ Compatível (sintaxe diferente) |

### Campos que Precisam Conversão

| OpenAPI 3.1.0 | Swagger 2.0 | Conversão |
|---------------|-------------|-----------|
| `servers[0].url` | `host` + `basePath` + `schemes` | Parse URL → componentes |
| `components.schemas` | `definitions` | Rename + adapt schema |
| `components.securitySchemes` | `securityDefinitions` | Rename + adapt |
| `requestBody` | `parameters` (body) | Converter para parameter type=body |
| `content` (MediaType) | `consumes` / `produces` | Extrair MIME types |

### Campos Exclusivos do 3.1.0 (Ignorados no 2.0)

- ❌ `webhooks` - Não existe em 2.0
- ❌ `jsonSchemaDialect` - Não existe em 2.0
- ❌ `license.identifier` - Não existe em 2.0 (usar `license.url`)
- ❌ Schema fields: `prefixItems`, `unevaluatedProperties`, etc.

### Campos Exclusivos do 2.0 (Adicionados na conversão)

- ✅ `swagger: "2.0"` - Versão da especificação
- ✅ `host` - Hostname da API
- ✅ `basePath` - Base path
- ✅ `schemes` - Protocolos (http, https)
- ✅ `consumes` - Global content types aceitos
- ✅ `produces` - Global content types produzidos

## 💻 Implementação Proposta

### 1. Nova Estrutura Swagger 2.0

```go
// pkg/openapi/v2/swagger.go
package v2

type Swagger struct {
    Swagger             string                       `json:"swagger"`                       // "2.0"
    Info                Info                         `json:"info"`
    Host                string                       `json:"host,omitempty"`
    BasePath            string                       `json:"basePath,omitempty"`
    Schemes             []string                     `json:"schemes,omitempty"`            // http, https, ws, wss
    Consumes            []string                     `json:"consumes,omitempty"`
    Produces            []string                     `json:"produces,omitempty"`
    Paths               Paths                        `json:"paths"`
    Definitions         map[string]*Schema           `json:"definitions,omitempty"`
    Parameters          map[string]*Parameter        `json:"parameters,omitempty"`
    Responses           map[string]*Response         `json:"responses,omitempty"`
    SecurityDefinitions map[string]*SecurityScheme   `json:"securityDefinitions,omitempty"`
    Security            []SecurityRequirement        `json:"security,omitempty"`
    Tags                []Tag                        `json:"tags,omitempty"`
    ExternalDocs        *ExternalDocs                `json:"externalDocs,omitempty"`
}

// Adaptar Schema para JSON Schema Draft 4
type Schema struct {
    Type                 string              `json:"type,omitempty"`
    Format               string              `json:"format,omitempty"`
    Title                string              `json:"title,omitempty"`
    Description          string              `json:"description,omitempty"`
    Default              interface{}         `json:"default,omitempty"`
    Maximum              *float64            `json:"maximum,omitempty"`
    Minimum              *float64            `json:"minimum,omitempty"`
    MaxLength            *int                `json:"maxLength,omitempty"`
    MinLength            *int                `json:"minLength,omitempty"`
    Pattern              string              `json:"pattern,omitempty"`
    MaxItems             *int                `json:"maxItems,omitempty"`
    MinItems             *int                `json:"minItems,omitempty"`
    UniqueItems          bool                `json:"uniqueItems,omitempty"`
    Enum                 []interface{}       `json:"enum,omitempty"`
    MultipleOf           *float64            `json:"multipleOf,omitempty"`
    
    // Object properties
    Properties           map[string]*Schema  `json:"properties,omitempty"`
    AdditionalProperties interface{}         `json:"additionalProperties,omitempty"` // bool or *Schema
    Required             []string            `json:"required,omitempty"`
    
    // Array items
    Items                *Schema             `json:"items,omitempty"`
    
    // Composition
    AllOf                []*Schema           `json:"allOf,omitempty"`
    
    // Reference
    Ref                  string              `json:"$ref,omitempty"`
    
    // Extensions
    Extensions           map[string]interface{} `json:"-"`
}

// Parameter em Swagger 2.0
type Parameter struct {
    Name             string      `json:"name"`
    In               string      `json:"in"` // query, header, path, formData, body
    Description      string      `json:"description,omitempty"`
    Required         bool        `json:"required,omitempty"`
    
    // Para in != body
    Type             string      `json:"type,omitempty"`
    Format           string      `json:"format,omitempty"`
    AllowEmptyValue  bool        `json:"allowEmptyValue,omitempty"`
    Items            *Items      `json:"items,omitempty"`
    CollectionFormat string      `json:"collectionFormat,omitempty"` // csv, ssv, tsv, pipes, multi
    Default          interface{} `json:"default,omitempty"`
    Maximum          *float64    `json:"maximum,omitempty"`
    Minimum          *float64    `json:"minimum,omitempty"`
    MaxLength        *int        `json:"maxLength,omitempty"`
    MinLength        *int        `json:"minLength,omitempty"`
    Pattern          string      `json:"pattern,omitempty"`
    Enum             []interface{} `json:"enum,omitempty"`
    
    // Para in = body
    Schema           *Schema     `json:"schema,omitempty"`
}

// Response em Swagger 2.0
type Response struct {
    Description string             `json:"description"`
    Schema      *Schema            `json:"schema,omitempty"`
    Headers     map[string]*Header `json:"headers,omitempty"`
    Examples    map[string]interface{} `json:"examples,omitempty"`
}
```

### 2. Converter de 3.1.0 para 2.0

```go
// pkg/converter/converter.go
package converter

import (
    v2 "github.com/fsvxavier/nexs-swag/pkg/openapi/v2"
    v3 "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

type Converter struct {
    warnings []string
}

func (c *Converter) ConvertToV2(spec *v3.OpenAPI) (*v2.Swagger, error) {
    swagger := &v2.Swagger{
        Swagger: "2.0",
        Info:    c.convertInfo(spec.Info),
        Paths:   c.convertPaths(spec.Paths),
        Tags:    spec.Tags,
    }
    
    // Converter servers[0] para host/basePath/schemes
    if len(spec.Servers) > 0 {
        host, basePath, schemes := c.parseServerURL(spec.Servers[0].URL)
        swagger.Host = host
        swagger.BasePath = basePath
        swagger.Schemes = schemes
    }
    
    // Converter components para definitions
    if spec.Components != nil {
        swagger.Definitions = c.convertSchemas(spec.Components.Schemas)
        swagger.SecurityDefinitions = c.convertSecuritySchemes(spec.Components.SecuritySchemes)
    }
    
    // Ignorar webhooks (não existe em 2.0)
    if len(spec.Webhooks) > 0 {
        c.warnings = append(c.warnings, "webhooks are not supported in OpenAPI 2.0 and were ignored")
    }
    
    return swagger, nil
}

func (c *Converter) parseServerURL(url string) (host, basePath string, schemes []string) {
    // Parse URL: https://api.example.com/v1 
    // → host: api.example.com, basePath: /v1, schemes: [https]
    // Implementação...
    return
}

func (c *Converter) convertRequestBody(rb *v3.RequestBody, op *v2.Operation) {
    // Converter requestBody para parameter type=body
    for mediaType, content := range rb.Content {
        if mediaType == "application/json" {
            param := &v2.Parameter{
                Name:     "body",
                In:       "body",
                Required: rb.Required,
                Schema:   c.convertSchema(content.Schema),
            }
            op.Parameters = append(op.Parameters, param)
        }
    }
}

func (c *Converter) convertSchema(schema *v3.Schema) *v2.Schema {
    // Converter JSON Schema 2020-12 → Draft 4
    v2Schema := &v2.Schema{
        Type:        c.convertType(schema.Type),
        Format:      schema.Format,
        Description: schema.Description,
        // ... outros campos
    }
    
    // Remover campos não suportados em Draft 4
    // - prefixItems → ignorar
    // - unevaluatedProperties → ignorar
    // - dependentSchemas → ignorar
    
    return v2Schema
}

func (c *Converter) convertType(t interface{}) string {
    // Em 3.1.0, type pode ser array: ["string", "null"]
    // Em 2.0, type é sempre string
    if arr, ok := t.([]interface{}); ok {
        // Pegar primeiro tipo não-null
        for _, typ := range arr {
            if s, ok := typ.(string); ok && s != "null" {
                return s
            }
        }
    }
    if s, ok := t.(string); ok {
        return s
    }
    return ""
}
```

### 3. CLI Flag para Versão

```go
// cmd/nexs-swag/main.go

&cli.StringFlag{
    Name:    "openapi-version",
    Aliases: []string{"ov"},
    Value:   "3.1.0",
    Usage:   "OpenAPI version: 2.0, 3.0.0, 3.1.0",
},
```

### 4. Lógica no Parser

```go
// pkg/parser/parser.go

type Parser struct {
    // ...
    openapiVersion string  // "2.0" ou "3.1.0"
}

func (p *Parser) SetOpenAPIVersion(version string) {
    p.openapiVersion = version
}

func (p *Parser) GetSpecification() openapi.Specification {
    if p.openapiVersion == "2.0" {
        // Converter ou construir diretamente em 2.0
        return p.buildSwagger2()
    }
    return p.openapi // 3.1.0 (default)
}
```

### 5. Gerador Adaptado

```go
// pkg/generator/generator.go

func New(spec openapi.Specification, outputDir string, outputType []string) *Generator {
    return &Generator{
        spec:       spec,
        outputDir:  outputDir,
        outputType: outputType,
    }
}

func (g *Generator) generateJSON() error {
    var filename string
    switch g.spec.GetVersion() {
    case "2.0":
        filename = "swagger.json"
    case "3.0.0", "3.1.0":
        filename = "openapi.json"
    default:
        filename = "openapi.json"
    }
    
    data, err := json.MarshalIndent(g.spec, "", "  ")
    // ...
}
```

## 📝 Plano de Implementação

### Fase 1: Estruturas Base (8-12h)
1. ✅ Criar `pkg/openapi/v2/swagger.go` com estruturas Swagger 2.0
2. ✅ Criar `pkg/openapi/spec.go` com interface comum
3. ✅ Mover estruturas atuais para `pkg/openapi/v3/`
4. ✅ Atualizar imports no projeto

### Fase 2: Conversor (12-16h)
1. ✅ Implementar `pkg/converter/converter.go`
2. ✅ Converter Info, Paths básicos
3. ✅ Converter Servers → host/basePath/schemes
4. ✅ Converter Components → Definitions
5. ✅ Converter RequestBody → Parameter
6. ✅ Converter Responses
7. ✅ Converter Security Schemes
8. ✅ Testes unitários do conversor

### Fase 3: Parser e Gerador (8-12h)
1. ✅ Adicionar flag `--openapi-version` no CLI
2. ✅ Adaptar Parser para versão configurável
3. ✅ Adaptar Generator para aceitar interface
4. ✅ Atualizar lógica de geração (nomes de arquivo, etc)

### Fase 4: Testes e Documentação (12-16h)
1. ✅ Testes unitários para estruturas v2
2. ✅ Testes de integração (gerar 2.0 e 3.1.0)
3. ✅ Adicionar exemplos em `examples/22-openapi-v2/`
4. ✅ Atualizar documentação (README, etc)
5. ✅ Validar specs geradas com validadores externos

### Fase 5: Refinamento (4-8h)
1. ✅ Tratamento de edge cases
2. ✅ Mensagens de warning para recursos não suportados
3. ✅ Otimizações de performance
4. ✅ Code review e refatoração

**Total Estimado:** 44-64 horas

## ⚠️ Limitações e Advertências

### Recursos do 3.1.0 que Não Podem Ser Convertidos

1. **Webhooks** - Não existe em 2.0
2. **JSON Schema 2020-12** - Precisa downgrade para Draft 4
   - `prefixItems` → Ignorado
   - `unevaluatedProperties` → Ignorado
   - `dependentSchemas` → Ignorado
3. **Nullable como array de tipos** - `type: ["string", "null"]` → `type: "string"` + `x-nullable: true`
4. **Multiple servers** - Apenas primeiro server é convertido
5. **License.identifier** - Não suportado em 2.0

### Estratégia para Nullable

```go
// OpenAPI 3.1.0
type: ["string", "null"]

// Swagger 2.0 (com extensão)
type: "string"
x-nullable: true
```

## 🎯 Benefícios da Implementação

### Para os Usuários

1. ✅ **Compatibilidade retroativa** - Suportar ferramentas que só entendem 2.0
2. ✅ **Escolha flexível** - Gerar ambas versões simultaneamente
3. ✅ **Migração gradual** - Migrar de 2.0 para 3.1.0 aos poucos
4. ✅ **Ferramentas legadas** - Integração com Swagger UI antigo, codegen, etc

### Para o Projeto

1. ✅ **Diferencial competitivo** - swaggo/swag só gera 2.0, nexs-swag geraria ambos
2. ✅ **Adoção ampliada** - Atinge usuários presos em 2.0
3. ✅ **Showcase técnico** - Demonstra arquitetura sólida e extensível

## 🚀 Exemplo de Uso Proposto

```bash
# Gerar OpenAPI 3.1.0 (padrão atual)
nexs-swag init

# Gerar Swagger 2.0
nexs-swag init --openapi-version 2.0

# Gerar ambos
nexs-swag init --openapi-version 2.0,3.1.0

# Especificar nome de arquivo
nexs-swag init --openapi-version 2.0 -o docs --format json
# Gera: docs/swagger.json

nexs-swag init --openapi-version 3.1.0 -o docs --format json
# Gera: docs/openapi.json
```

## 📊 Comparação com swaggo/swag

| Recurso | swaggo/swag | nexs-swag (atual) | nexs-swag (após implementação) |
|---------|-------------|-------------------|--------------------------------|
| Swagger 2.0 | ✅ | ❌ | ✅ |
| OpenAPI 3.0 | ❌ | ❌ | ⚠️ (possível) |
| OpenAPI 3.1 | ❌ | ✅ | ✅ |
| Escolha de versão | ❌ | ❌ | ✅ |
| Conversão entre versões | ❌ | ❌ | ✅ |

## 🎬 Recomendação Final

**RECOMENDO A IMPLEMENTAÇÃO** pelos seguintes motivos:

1. ✅ **Arquitetura permite** - Estrutura modular facilita adição
2. ✅ **Não quebra existente** - Retrocompatível com código atual
3. ✅ **Diferencial forte** - swaggo/swag não oferece escolha de versão
4. ✅ **Demanda real** - Muitas empresas ainda usam Swagger 2.0
5. ✅ **Complexidade gerenciável** - 44-64h é razoável para o valor entregue

### Priorização Sugerida

**Alta Prioridade:**
- ✅ Suporte a Swagger 2.0 básico
- ✅ Conversor 3.1.0 → 2.0
- ✅ Flag CLI `--openapi-version`

**Média Prioridade:**
- ⚠️ Suporte a OpenAPI 3.0 (intermediário)
- ⚠️ Geração simultânea de múltiplas versões

**Baixa Prioridade:**
- 🔵 Conversão reversa 2.0 → 3.1.0
- 🔵 Detecção automática de versão desejada

---

## 📚 Referências

- [OpenAPI 2.0 Specification](https://swagger.io/specification/v2/)
- [OpenAPI 3.1.0 Specification](https://spec.openapis.org/oas/v3.1.0)
- [JSON Schema Draft 4](https://json-schema.org/specification-links.html#draft-4)
- [JSON Schema 2020-12](https://json-schema.org/specification-links.html#2020-12)
- [Swagger Converter (by OpenAPI)](https://converter.swagger.io/)

---

**Autor:** Análise técnica gerada para nexs-swag
**Data:** 10 de dezembro de 2025
**Versão:** 1.0



