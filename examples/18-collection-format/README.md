# Exemplo 18 - Collection Format

Demonstra diferentes formatos para arrays em query parameters.

## Flag

```bash
--collectionFormat <format>
--cf <format>
```

Default: `csv`

## Formatos Suportados

### 1. CSV (Comma Separated Values)
```
?ids=1,2,3
```

```go
// @Param ids query []int true "IDs" collectionFormat(csv)
```

OpenAPI:
```yaml
collectionFormat: csv
```

### 2. Multi (Multiple Parameters)
```
?tags=go&tags=api&tags=web
```

```go
// @Param tags query []string true "Tags" collectionFormat(multi)
```

OpenAPI:
```yaml
collectionFormat: multi
explode: true
```

### 3. Pipes
```
?statuses=active|pending|done
```

```go
// @Param statuses query []string true "Statuses" collectionFormat(pipes)
```

OpenAPI:
```yaml
collectionFormat: pipes
```

### 4. TSV (Tab Separated Values)
```
?values=a\tb\tc
```

```go
// @Param values query []string true "Values" collectionFormat(tsv)
```

OpenAPI:
```yaml
collectionFormat: tsv
```

### 5. SSV (Space Separated Values)
```
?items=one two three
# URL encoded: ?items=one%20two%20three
```

```go
// @Param items query []string true "Items" collectionFormat(ssv)
```

OpenAPI:
```yaml
collectionFormat: ssv
```

## Uso

### Prioridade

1. **Annotation** (maior prioridade)
```go
// @Param ids query []int true "IDs" collectionFormat(multi)
```

2. **CLI Flag**
```bash
nexs-swag init --collectionFormat multi
```

3. **Default:** csv

### Global + Override

```bash
# Global: csv
nexs-swag init --collectionFormat csv
```

```go
// Mas este endpoint usa multi
// @Param tags query []string true "Tags" collectionFormat(multi)
// @Router /search [get]
func Search() {}

// Este usa o global (csv)
// @Param ids query []int true "IDs"
// @Router /filter [get]
func Filter() {}
```

## Comparação

| Formato | Exemplo | URL Encoding | Uso Comum |
|---------|---------|--------------|-----------|
| CSV | `1,2,3` | Não | REST APIs |
| Multi | `?id=1&id=2&id=3` | Não | HTML forms |
| Pipes | `a\|b\|c` | Sim (\|) | Filters |
| TSV | `a\tb\tc` | Sim (\t) | Data export |
| SSV | `a b c` | Sim (space) | Raro |

## Cliente HTTP

### CSV
```bash
curl "http://api.com/items?ids=1,2,3"
```

```javascript
fetch('http://api.com/items?ids=1,2,3')
```

### Multi
```bash
curl "http://api.com/items?tag=go&tag=api&tag=web"
```

```javascript
const params = new URLSearchParams();
params.append('tag', 'go');
params.append('tag', 'api');
params.append('tag', 'web');
fetch('http://api.com/items?' + params)
```

### Pipes
```bash
curl "http://api.com/items?status=active|pending"
```

```javascript
fetch('http://api.com/items?status=active|pending')
```

## Server-side Parsing

```go
func SearchHandler(w http.ResponseWriter, r *http.Request) {
    // CSV: "1,2,3"
    idsCSV := r.URL.Query().Get("ids")
    ids := strings.Split(idsCSV, ",")
    
    // Multi: ["go", "api", "web"]
    tags := r.URL.Query()["tags"]
    
    // Pipes: "active|pending"
    statusPipes := r.URL.Query().Get("status")
    statuses := strings.Split(statusPipes, "|")
}
```

## Como Executar

```bash
./run.sh
```

## Recomendações

**Use CSV quando:**
- REST API padrão
- Compatibilidade máxima
- Arrays pequenos

**Use Multi quando:**
- HTML forms
- Muitos valores
- Melhor legibilidade

**Use Pipes quando:**
- Valores podem ter vírgulas
- Filtros complexos

**Evite TSV/SSV:**
- Precisam URL encoding
- Pouco suportados
- Confusos para usuários

## OpenAPI 3.0 Style

```yaml
# Multi (OpenAPI 3.0)
parameters:
  - name: tags
    in: query
    schema:
      type: array
      items:
        type: string
    style: form
    explode: true

# CSV (OpenAPI 3.0)
parameters:
  - name: ids
    in: query
    schema:
      type: array
      items:
        type: integer
    style: form
    explode: false
```
