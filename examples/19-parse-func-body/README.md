# Exemplo 19 - Parse Func Body

Demonstra o uso de `--parseFuncBody` para análise do corpo das funções.

## Flag

```bash
--parseFuncBody
```

## Comportamento

### SEM --parseFuncBody (Default)
Apenas annotations são parseadas:
```go
// @Success 200 {object} User
// @Failure 400 {object} Error
func GetUser() {
    // Este código é IGNORADO
    if err != nil {
        return Error{} // NÃO detectado
    }
}
```

### COM --parseFuncBody
Corpo das funções também é analisado:
```go
// @Success 200 {object} User
func GetUser() {
    // Este código é ANALISADO
    if err != nil {
        return Error{} // DETECTADO automaticamente
    }
}
```

## Uso

```bash
nexs-swag init --parseFuncBody
```

## O Que é Detectado

### 1. Error Responses
```go
// @Router /items [post]
func CreateItem() {
    if item.Name == "" {
        w.WriteHeader(http.StatusBadRequest)  // Detecta 400
        json.NewEncoder(w).Encode(Error{})     // Detecta type Error
        return
    }
}
```

Gera automaticamente:
```yaml
responses:
  400:
    schema:
      $ref: '#/definitions/Error'
```

### 2. Status Codes
```go
func Handler() {
    w.WriteHeader(http.StatusCreated)      // 201
    w.WriteHeader(http.StatusNoContent)    // 204
    w.WriteHeader(http.StatusNotFound)     // 404
    w.WriteHeader(http.StatusConflict)     // 409
}
```

### 3. Return Types
```go
func GetData() {
    if cached {
        return CachedData{}  // Detecta CachedData
    }
    return FreshData{}       // Detecta FreshData
}
```

## Como Executar

```bash
./run.sh
```

## Benefícios

### 1. Documentação Automática
Sem flag:
```go
// @Success 200 {object} User
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
func GetUser() { ... }
```

Com flag:
```go
// @Success 200 {object} User
// Outros responses detectados automaticamente!
func GetUser() {
    if err == sql.ErrNoRows {
        return 404, Error{}
    }
    if err != nil {
        return 500, Error{}
    }
}
```

### 2. Menos Annotations
Código fica mais limpo e DRY.

### 3. Sincronização
Documentação sempre atualizada com o código.

## Desvantagens

### 1. Performance
```bash
# Mais lento
--parseFuncBody: ~2x tempo de parsing
```

### 2. False Positives
```go
func Test() {
    // Código de teste pode gerar responses inválidos
    w.WriteHeader(http.StatusTeapot)  // 418 detectado
}
```

### 3. Complexidade
```go
// Difícil detectar lógica complexa
func Handler() {
    status := getStatus()  // Valor dinâmico
    w.WriteHeader(status)   // Não detectado
}
```

## Quando Usar

**Use --parseFuncBody quando:**
- Validações complexas no código
- Múltiplos error responses
- Quer menos annotations
- Performance não é crítica

**NÃO use quando:**
- API simples com poucos endpoints
- Performance importante
- Prefer explicit documentation
- Muitos testes junto com código

## Comparação

### Sem Flag
```go
// Tudo explícito
// @Success 200 {object} Data
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
func Handler() {
    // validações...
}
```

**Prós:** Explícito, rápido, previsível
**Contras:** Verboso, pode desatualizar

### Com Flag
```go
// Mínimo necessário
// @Success 200 {object} Data
func Handler() {
    // validações detectadas automaticamente
}
```

**Prós:** Conciso, sincronizado, DRY
**Contras:** Lento, false positives, menos controle

## Recomendação

✅ **Não use por padrão**

Prefira annotations explícitas para melhor controle e performance.

Use apenas quando tiver muitas validações e quiser reduzir verbosidade.

## Alternativa: Middleware

Em vez de parsear corpo das funções, considere usar middleware:

```go
// Middleware documenta errors comuns
// @Failure 401 {object} Error "Unauthorized"
// @Failure 403 {object} Error "Forbidden"
// @Failure 500 {object} Error "Internal Error"

// Endpoint só documenta success e específicos
// @Success 200 {object} User
// @Failure 404 {object} Error "User not found"
func GetUser() {}
```
