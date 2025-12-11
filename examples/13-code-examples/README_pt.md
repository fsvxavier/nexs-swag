# Exemplo 13 - Code Examples

🌍 [English](README.md) • **Português (Brasil)** • [Español](README_es.md)

Demonstra como adicionar exemplos de código em múltiplas linguagens.

## Flag

```bash
--codeExampleFilesDir <directory>
--cef <directory>
```

## Uso

```bash
nexs-swag init --codeExampleFilesDir ./code_samples
```

## Estrutura

```
13-code-examples/
├── main.go
└── code_samples/
    ├── create_user.go    # Exemplo em Go
    ├── create_user.js    # Exemplo em JavaScript
    ├── create_user.py    # Exemplo em Python
    └── create_user.sh    # Exemplo em Bash
```

## Como Funciona

### 1. No código, use a annotation:
```go
// @x-codeSamples file(create_user)
// @Router /users [post]
func CreateUser() {}
```

### 2. Crie arquivos com o prefixo:
- `create_user.go`
- `create_user.js`
- `create_user.py`
- etc.

### 3. O nexs-swag detecta automaticamente a linguagem pela extensão

## Linguagens Suportadas

| Extensão | Linguagem |
|----------|-----------|
| .go | Go |
| .js | JavaScript |
| .ts | TypeScript |
| .py | Python |
| .java | Java |
| .rb | Ruby |
| .php | PHP |
| .cs | C# |
| .cpp | C++ |
| .c | C |
| .sh | Bash |
| .rs | Rust |
| .swift | Swift |
| .kt | Kotlin |
| .dart | Dart |
| .scala | Scala |
| ...e mais |

## OpenAPI Gerado

```json
{
  "paths": {
    "/users": {
      "post": {
        "x-codeSamples": [
          {
            "lang": "Go",
            "source": "// código do create_user.go..."
          },
          {
            "lang": "JavaScript",
            "source": "// código do create_user.js..."
          },
          {
            "lang": "Python",
            "source": "# código do create_user.py..."
          },
          {
            "lang": "Bash",
            "source": "#!/bin/bash\n# código do create_user.sh..."
          }
        ]
      }
    }
  }
}
```

## Como Executar

```bash
./run.sh
```

## Benefícios

1. **Múltiplas Linguagens:** Exemplos para diferentes clientes
2. **Detecção Automática:** Extensão → linguagem
3. **Reutilização:** Mesmo arquivo para múltiplos endpoints
4. **Testáveis:** Code samples são código real que pode ser executado

## Casos de Uso

- **SDKs:** Exemplos para cada linguagem suportada
- **Documentação Rica:** Facilita integração
- **Onboarding:** Desenvolvedores copiam e colam
- **Testing:** Exemplos validados e funcionais
