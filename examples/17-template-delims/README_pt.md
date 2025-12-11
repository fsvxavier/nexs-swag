# Exemplo 17 - Template Delimiters

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra o uso de `--templateDelims` para customizar delimitadores de template.

## Flag

```bash
--templateDelims "<left> <right>"
--td "<left> <right>"
```

Default: `"{{ }}"`

## Problema

Templates Go usam `{{ }}`, mas isso pode conflitar com:
- Frontend frameworks (Vue, Angular, Svelte)
- Template engines (Mustache, Handlebars)
- Documentation systems

## Solução

Customizar os delimitadores:

```bash
nexs-swag init --templateDelims "[[ ]]"
```

## Delimitadores Suportados

### Recomendados

```bash
--templateDelims "[[ ]]"     # Melhor para evitar conflitos
--templateDelims "{{{ }}}"   # Mustache style
--templateDelims "<< >>"     # Shell style
```

### Outros

```bash
--templateDelims "<% %>"     # ERB style
--templateDelims "{% %}"     # Jinja2/Twig style
--templateDelims "${ }"      # ES6 style
```

## Caso de Uso: Vue.js Conflict

### Problema

```html
<!-- index.html -->
<div id="app">
  <!-- Vue usa {{ }} -->
  <p>{{ message }}</p>
  
  <!-- Swagger também usa {{ }} -->
  <script>
    const spec = {{ .SwaggerJSON }}; // CONFLITO!
  </script>
</div>
```

### Solução

```bash
nexs-swag init --templateDelims "[[ ]]"
```

```html
<!-- index.html -->
<div id="app">
  <!-- Vue continua com {{ }} -->
  <p>{{ message }}</p>
  
  <!-- Swagger usa [[ ]] -->
  <script>
    const spec = [[ .SwaggerJSON ]]; // SEM conflito!
  </script>
</div>
```

## No Código Go

```go
package main

import (
    "text/template"
)

func main() {
    // Template padrão
    tmpl1 := template.Must(template.New("t1").Parse("Hello {{.Name}}"))
    
    // Template com delimiters customizados
    tmpl2 := template.New("t2")
    tmpl2.Delims("[[", "]]")
    tmpl2 = template.Must(tmpl2.Parse("Hello [[.Name]]"))
}
```

## Como Executar

```bash
./run.sh
```

## Exemplo Completo: Swagger UI Custom

### 1. Gerar com delimiters customizados

```bash
nexs-swag init --templateDelims "[[ ]]"
```

### 2. Template HTML

```html
<!DOCTYPE html>
<html>
<head>
    <title>API Docs</title>
</head>
<body>
    <!-- Frontend framework pode usar {{ }} -->
    <div id="vue-app">{{ vueMessage }}</div>
    
    <!-- Swagger usa [[ ]] -->
    <script>
        const swaggerSpec = [[ .SwaggerJSON ]];
        SwaggerUIBundle({
            spec: swaggerSpec,
            dom_id: '#swagger-ui'
        });
    </script>
</body>
</html>
```

### 3. Servir o HTML

```go
package main

import (
    "html/template"
    "net/http"
    
    _ "myapp/docs"
)

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.New("swagger")
    tmpl.Delims("[[", "]]")  // Mesmos delimiters!
    
    tmpl, _ = tmpl.ParseFiles("swagger.html")
    tmpl.Execute(w, nil)
}

func main() {
    http.HandleFunc("/swagger/", swaggerHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Benefícios

- **Sem conflitos:** Frontend frameworks funcionam normalmente
- **Flexibilidade:** Escolha os delimiters que preferir
- **Compatibilidade:** Funciona com qualquer template engine
- **Clareza:** Código mais legível

## Quando Usar

**Use --templateDelims quando:**
- Usar Vue.js, Angular, Svelte
- Integrar com Mustache, Handlebars
- Ter conflitos de sintaxe
- Preferir outro estilo

**NÃO precisa quando:**
- API pura sem frontend
- Servir JSON direto
- Usar Swagger UI standalone
