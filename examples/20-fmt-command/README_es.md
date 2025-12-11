# Ejemplo 20 - Fmt Command

🌍 [English](README.md) • [Português (Brasil)](README_pt.md) • **Español**

Demuestra el uso del comando `fmt` para formatear comentarios Swagger.

## Comando

```bash
nexs-swag fmt [flags]
```

## Concepto

Formatea y organiza comentarios Swagger en el código Go de forma consistente.

## Uso Básico

```bash
# Formatear archivo específico
nexs-swag fmt -f main.go

# Formatear directorio
nexs-swag fmt -d ./api

# Formatear todo el proyecto
nexs-swag fmt -d .
```

## Cómo Ejecutar

```bash
./run.sh
```

## Qué Formatea

### 1. Indentación
```go
// ANTES
// @Summary Create user
//       @Description Creates a new user
//  @Tags users
func CreateUser() {}

// DESPUÉS
// @Summary Create user
// @Description Creates a new user
// @Tags users
func CreateUser() {}
```

### 2. Orden de Tags
```go
// ANTES - Desordenado
// @Tags users
// @Summary Create user
// @Router /users [post]
// @Description Creates a new user
// @Accept json

// DESPUÉS - Orden standard
// @Summary Create user
// @Description Creates a new user
// @Tags users
// @Accept json
// @Router /users [post]
```

### 3. Espaciado
```go
// ANTES
// @Summary Create user
//
//
// @Description Creates a new user
// @Tags users

// DESPUÉS
// @Summary Create user
// @Description Creates a new user
// @Tags users
```

### 4. Alignment
```go
// ANTES
// @Param id path int true "User ID"
// @Param name query string false "Name"
// @Success 200 {object} User

// DESPUÉS - Alineado
// @Param   id    path   int    true  "User ID"
// @Param   name  query  string false "Name"
// @Success 200   {object} User
```

## Flags

```bash
-d, --dir string        Directory to format (default ".")
-f, --file string       File to format
--exclude string        Exclude directories/files (comma-separated)
-w, --write            Write changes to files (default: preview only)
--check                Check if files are formatted (exit 1 if not)
```

## Ejemplos

### 1. Preview (Default)
```bash
# Muestra cambios sin aplicar
nexs-swag fmt -d ./api

# Output:
# api/handlers.go: Would reformat 3 functions
# api/users.go: Would reformat 5 functions
```

### 2. Write Changes
```bash
# Aplica cambios
nexs-swag fmt -d ./api -w

# Output:
# api/handlers.go: Reformatted 3 functions
# api/users.go: Reformatted 5 functions
```

### 3. Check Mode (CI/CD)
```bash
# Exit 0 si formateado, exit 1 si no
nexs-swag fmt -d . --check

# Output:
# api/handlers.go: Not formatted
# exit code: 1
```

### 4. Exclude Directories
```bash
nexs-swag fmt -d . --exclude "vendor,testdata,mocks"
```

## Orden Standard de Tags

```go
// 1. Summary y Description
// @Summary
// @Description

// 2. Tags
// @Tags

// 3. Accept y Produce
// @Accept
// @Produce

// 4. Parameters
// @Param

// 5. Security
// @Security

// 6. Responses
// @Success
// @Failure

// 7. Router (siempre último)
// @Router
```

## Ejemplo Completo

### ANTES
```go
// @Router /users [post]
// @Param user body User true "User object"
// @Success 201 {object} User
// @Tags users
//       @Description Creates a new user in the system
// @Summary Create user
// @Accept json
// @Produce json
func CreateUser(c *gin.Context) {}
```

### DESPUÉS
```bash
nexs-swag fmt -f handlers.go -w
```

```go
// @Summary Create user
// @Description Creates a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User object"
// @Success 201 {object} User
// @Router /users [post]
func CreateUser(c *gin.Context) {}
```

## CI/CD Integration

### GitHub Actions
```yaml
name: Format Check

on: [push, pull_request]

jobs:
  format:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Install nexs-swag
        run: go install github.com/fsvxavier/nexs-swag/cmd/nexs-swag@latest
      
      - name: Check formatting
        run: nexs-swag fmt -d . --check
      
      - name: Format code
        if: failure()
        run: |
          nexs-swag fmt -d . -w
          git diff
```

### Pre-commit Hook
```bash
# .git/hooks/pre-commit
#!/bin/bash

echo "Checking Swagger comments format..."
nexs-swag fmt -d . --check

if [ $? -ne 0 ]; then
  echo "❌ Swagger comments not formatted!"
  echo "Run: nexs-swag fmt -d . -w"
  exit 1
fi

echo "✅ Swagger comments formatted correctly"
```

### Makefile
```makefile
.PHONY: fmt fmt-check

fmt:
nexs-swag fmt -d . -w

fmt-check:
nexs-swag fmt -d . --check

fmt-api:
nexs-swag fmt -d ./api -w
```

## Editor Integration

### VS Code
```json
// .vscode/tasks.json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Format Swagger",
      "type": "shell",
      "command": "nexs-swag fmt -d . -w",
      "problemMatcher": []
    }
  ]
}
```

### GoLand/IntelliJ
```xml
<!-- External Tools -->
<tool name="Format Swagger"
      description="Format Swagger comments"
      showInMainMenu="true"
      showInEditor="true">
  <exec>
    <option name="COMMAND" value="nexs-swag" />
    <option name="PARAMETERS" value="fmt -f $FilePath$ -w" />
  </exec>
</tool>
```

## Tips

### 1. Format Before Commit
```bash
# Antes de commit
nexs-swag fmt -d . -w
git add .
git commit -m "Format Swagger comments"
```

### 2. Format Specific Files
```bash
# Solo archivos modificados
git diff --name-only | grep '\.go$' | while read file; do
  nexs-swag fmt -f "$file" -w
done
```

### 3. Combine con go fmt
```bash
# Format Go y Swagger juntos
go fmt ./...
nexs-swag fmt -d . -w
```

### 4. Team Convention
```bash
# Documente en README
## Code Formatting

\`\`\`bash
# Format Go code
go fmt ./...

# Format Swagger comments
nexs-swag fmt -d . -w
\`\`\`
```

## Performance

```bash
# Proyecto pequeño (< 50 archivos)
Time: ~100ms

# Proyecto mediano (100-500 archivos)
Time: ~500ms

# Proyecto grande (> 1000 archivos)
Time: ~2s
```

## Cuándo Usar

**Use nexs-swag fmt cuando:**
- Team grande con diferentes estilos
- Quiere consistencia
- Code reviews difíciles
- CI/CD enforcement

**NO use cuando:**
- Proyecto muy pequeño
- Solo un desarrollador
- Preferencia por manual formatting
- Performance crítico en CI/CD

## Recomendaciones

**Para la mayoría de proyectos:**
```bash
# 1. Format antes de commit
nexs-swag fmt -d . -w

# 2. Check en CI/CD
nexs-swag fmt -d . --check
```

**Para proyectos grandes:**
```bash
# Format solo archivos modificados
git diff --name-only origin/main | grep '\.go$' | while read file; do
  nexs-swag fmt -f "$file" -w
done
```
