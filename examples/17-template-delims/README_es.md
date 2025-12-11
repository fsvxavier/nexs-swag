# Ejemplo 17 - Template Delims

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso de delimitadores personalizados para templates.

## Flags

```bash
--leftTemplateDelim <delim>
--ltd <delim>

--rightTemplateDelim <delim>
--rtd <delim>
```

Default: `{{` y `}}`

## Concepto

Cambia los delimitadores de templates Go para evitar conflictos.

## Problema

### Conflicto con Swagger UI
```yaml
# Swagger UI usa {{ }} para templating
paths:
  /users/{{id}}:  # Conflicto!
    get:
      summary: Get user by ID
```

### Conflicto con Vue.js
```html
<!-- Vue.js usa {{ }} -->
<div>{{ userName }}</div>  <!-- Conflicto! -->
```

## Solución

```bash
nexs-swag init --leftTemplateDelim "[[" --rightTemplateDelim "]]"
```

Ahora nexs-swag usa `[[ ]]` en vez de `{{ }}`:

```go
// Template con delimitadores personalizados
[[ .Title ]]
[[ .Description ]]
```

## Cómo Ejecutar

```bash
./run.sh
```

## Casos de Uso

### 1. Swagger UI Templates
```yaml
# SIN custom delims - Problema
paths:
  /users/{{userId}}:     # Swagger UI lo interpreta
    parameters:
      - name: userId
        in: path

# CON custom delims - OK
paths:
  /users/{userId}:       # Correcto OpenAPI format
    parameters:
      - name: userId
        in: path
```

### 2. Frontend Templates
```html
<!-- Vue.js template -->
<template>
  <div>
    <!-- Vue usa {{ }} -->
    <h1>{{ title }}</h1>
    
    <!-- nexs-swag usa [[ ]] -->
    <pre>[[ .SwaggerJSON ]]</pre>
  </div>
</template>
```

### 3. Documentation Templates
```markdown
<!-- docs/template.md -->

# API: {{ .Title }}

**Version:** {{ .Version }}

<!-- nexs-swag template -->
[[ .GeneratedDescription ]]
```

## Delimitadores Comunes

### Padrão Go
```go
{{ .Field }}
{{- .Field -}}
{{ if .Condition }}...{{ end }}
```

### Jinja2 Style
```bash
nexs-swag init --ltd "{%" --rtd "%}"
```
```go
{% .Field %}
{% if .Condition %}...{% endif %}
```

### Angular Style
```bash
nexs-swag init --ltd "{{" --rtd "}}"  # Default
# O
nexs-swag init --ltd "<<" --rtd ">>"
```
```html
<< .Field >>
```

### Bracket Style
```bash
nexs-swag init --ltd "[[" --rtd "]]"
```
```go
[[ .Field ]]
[[ range .Items ]]...[[ end ]]
```

### Mustache Style
```bash
nexs-swag init --ltd "{{" --rtd "}}"  # Same as default
```

## Template Variables

Disponibles en templates:

```go
[[ .Title ]]              // API title
[[ .Description ]]        // API description
[[ .Version ]]            // API version
[[ .Host ]]               // API host
[[ .BasePath ]]           // API base path
[[ .Schemes ]]            // [http, https]
[[ .Swagger ]]            // Swagger version
[[ .Info ]]               // Info object
[[ .Paths ]]              // Paths object
[[ .Definitions ]]        // Definitions/Components
```

## Exemplo Completo

### Custom Template
```html
<!-- templates/swagger.html -->
<!DOCTYPE html>
<html>
<head>
  <title>[[ .Title ]]</title>
</head>
<body>
  <h1>[[ .Title ]]</h1>
  <p>[[ .Description ]]</p>
  <p>Version: [[ .Version ]]</p>
  
  <!-- Swagger UI usa {{ }} internamente -->
  <div id="swagger-ui"></div>
  <script>
    SwaggerUIBundle({
      url: '/swagger.json',
      dom_id: '#swagger-ui'
    })
  </script>
</body>
</html>
```

### Generar
```bash
nexs-swag init \
  --leftTemplateDelim "[[" \
  --rightTemplateDelim "]]" \
  --outputTypes "go,json,yaml"
```

## Configuração

### .nexs-swag.yaml
```yaml
leftTemplateDelim: "[["
rightTemplateDelim: "]]"
```

```bash
nexs-swag init  # Lee config automáticamente
```

## Tips

### 1. Consistencia
```bash
# Defina en config file
echo "leftTemplateDelim: [[" >> .nexs-swag.yaml
echo "rightTemplateDelim: ]]" >> .nexs-swag.yaml

# Team usa misma config
git add .nexs-swag.yaml
```

### 2. Documentar
```markdown
# README.md

## Generating Docs

Este proyecto usa delimitadores personalizados:
- Left: `[[`
- Right: `]]`

\`\`\`bash
nexs-swag init  # Lee .nexs-swag.yaml
\`\`\`
```

### 3. CI/CD
```yaml
# .github/workflows/docs.yml
- name: Generate docs
  run: |
    nexs-swag init \
      --leftTemplateDelim "[[" \
      --rightTemplateDelim "]]"
```

## Debugging

### Verificar Delimitadores
```bash
# Check generated files
grep -r "\[\[" docs/
grep -r "\]\]" docs/

# Should NOT find {{ }} in nexs-swag templates
grep -r "{{" docs/ | grep -v "swagger-ui"
```

### Template Errors
```bash
# Error: template syntax
template: swagger:1: unexpected "{{" in operand

# Solución: use custom delims
nexs-swag init --ltd "[[" --rtd "]]"
```

## Cuándo Usar

**Use custom delims cuando:**
- Conflicto con Swagger UI
- Conflicto con frontend framework
- Templates anidados
- Preferencia de equipo

**Use defaults cuando:**
- Sin conflictos
- Standard Go templates
- Simplicidad
