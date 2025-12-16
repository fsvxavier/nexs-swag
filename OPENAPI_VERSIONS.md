# Suporte de Versões OpenAPI

O nexs-swag suporta a geração de documentação para todas as versões oficiais do OpenAPI Specification, conforme documentado em https://spec.openapis.org/oas/.

## Versões Suportadas

### OpenAPI 2.0 (Swagger 2.0)
- **Versão**: `2.0.0`
- **Uso**: `--openapi-version 2.0` ou `--ov 2.0`
- **Características**:
  - Formato legado (Swagger)
  - Estrutura `swagger`, `info`, `paths`, `definitions`
  - Não suporta webhooks ou múltiplos servers
  - Conversão automática de OpenAPI 3.x para 2.0

### OpenAPI 3.0.x
Versões suportadas: **3.0.0**, **3.0.1**, **3.0.2**, **3.0.3**, **3.0.4**

- **Uso**: `--openapi-version 3.0.x` (ex: `--ov 3.0.4`)
- **Características principais**:
  - Introdução de `servers` (múltiplos servidores)
  - Componentes reutilizáveis em `components`
  - RequestBody separado dos parâmetros
  - Callbacks para webhooks assíncronos
  - Links entre operações
  - Schema baseado em JSON Schema Draft 5

**Diferenças entre patch versions:**
- **3.0.0**: Versão inicial
- **3.0.1**: Correções de documentação
- **3.0.2**: Clarificações sobre serialização
- **3.0.3**: Melhorias em examples e discriminator
- **3.0.4**: Patch mais recente com correções menores

### OpenAPI 3.1.x
Versões suportadas: **3.1.0**, **3.1.1**, **3.1.2**

- **Uso**: `--openapi-version 3.1.x` (ex: `--ov 3.1.2`)
- **Características principais**:
  - **Compatível com JSON Schema Draft 2020-12**
  - Campo `webhooks` para webhooks incoming
  - Campo `jsonSchemaDialect` para especificar dialeto
  - Suporte a `$schema` em Schema Objects
  - Propriedade `summary` no Info Object
  - Campo `identifier` no License Object (SPDX)
  - Exemplos usando `examples` (array) ao invés de `example`

**Diferenças entre patch versions:**
- **3.1.0**: Versão inicial com compatibilidade JSON Schema
- **3.1.1**: Correções de documentação e clarificações
- **3.1.2**: Patch mais recente com melhorias menores

### OpenAPI 3.2.0 (Mais Recente)
- **Versão**: `3.2.0`
- **Uso**: `--openapi-version 3.2.0` ou `--ov 3.2`
- **Características principais** (adições sobre 3.1.x):
  - **Método HTTP QUERY** no Path Item Object
  - **OAuth2 Device Authorization Flow** (`deviceAuthorization`)
  - **oauth2MetadataUrl** no Security Scheme
  - **itemSchema** para streaming em Media Types
  - **itemEncoding** para streaming multipart
  - Suporte melhorado para Server-Sent Events
  - Campo `deprecated` no Security Scheme Object
  - Melhorias no suporte a sequential media types

## Como Usar

### Via Linha de Comando

```bash
# Gerar OpenAPI 2.0 (Swagger)
nexs-swag init --ov 2.0

# Gerar OpenAPI 3.0.4
nexs-swag init --ov 3.0.4

# Gerar OpenAPI 3.1.2
nexs-swag init --ov 3.1.2

# Gerar OpenAPI 3.2.0 (mais recente)
nexs-swag init --ov 3.2.0

# Forma curta também funciona
nexs-swag init --openapi-version 3.1
```

### Atalhos Aceitos

O sistema aceita versões curtas e as normaliza automaticamente:

| Entrada | Versão Gerada | Observação |
|---------|---------------|------------|
| `2`, `2.0`, `2.0.0` | `2.0.0` | Swagger 2.0 |
| `3` | `3.1.0` | Default para versão 3.1.0 |
| `3.0`, `3.0.0` | `3.0.0` | OpenAPI 3.0.0 |
| `3.0.4` | `3.0.4` | OpenAPI 3.0.4 |
| `3.1`, `3.1.0` | `3.1.0` | OpenAPI 3.1.0 |
| `3.1.2` | `3.1.2` | OpenAPI 3.1.2 |
| `3.2`, `3.2.0` | `3.2.0` | OpenAPI 3.2.0 (latest) |

## Diferenças Importantes por Família de Versões

### 2.0 vs 3.x
- **2.0**: Usa `swagger`, `host`, `basePath`, `schemes`
- **3.x**: Usa `openapi`, `servers` (array com URLs completas)

### 3.0 vs 3.1
- **3.0**: JSON Schema Draft 5
- **3.1**: JSON Schema Draft 2020-12 (totalmente compatível)
- **3.1+**: Adiciona `webhooks` e `jsonSchemaDialect`

### 3.1 vs 3.2
- **3.2**: Adiciona método QUERY, streaming avançado, OAuth2 device flow
- **3.2**: Melhora suporte para sequential media types (NDJSON, JSON Lines)

## Conversão Automática

Quando você especifica `--ov 2.0`, o sistema:
1. Gera internamente usando OpenAPI 3.1.0
2. Converte automaticamente para Swagger 2.0
3. Emite warnings sobre recursos não suportados

## Validação

O sistema valida automaticamente a especificação gerada de acordo com a versão escolhida usando:
- Schemas oficiais da OpenAPI Initiative
- Regras específicas de cada versão

## Recomendações

- **Novos projetos**: Use `3.2.0` para ter acesso aos recursos mais recentes
- **Compatibilidade máxima**: Use `3.0.4` para ferramentas que ainda não suportam 3.1+
- **JSON Schema compliance**: Use `3.1.x` para total compatibilidade com JSON Schema 2020-12
- **Legado**: Use `2.0` apenas se precisar dar suporte a ferramentas antigas

## Referências Oficiais

- [OpenAPI Specification](https://spec.openapis.org/oas/)
- [v3.2.0](https://spec.openapis.org/oas/v3.2.0.html)
- [v3.1.2](https://spec.openapis.org/oas/v3.1.2.html)
- [v3.0.4](https://spec.openapis.org/oas/v3.0.4.html)
- [v2.0](https://spec.openapis.org/oas/v2.0.html)

## Changelog de Implementação

### Versão 1.0.6
- ✅ Suporte completo para OpenAPI v2.0.0
- ✅ Suporte completo para OpenAPI v3.0.0-3.0.4
- ✅ Suporte completo para OpenAPI v3.1.0-3.1.2
- ✅ Suporte completo para OpenAPI v3.2.0
- ✅ Normalização automática de versões
- ✅ Validação de versão na linha de comando
- ✅ Conversão automática 3.x → 2.0
