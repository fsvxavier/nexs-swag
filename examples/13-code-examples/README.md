# Example 13 - Code Examples

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates how to add code examples in multiple languages.

## Flag

```bash
--codeExampleFilesDir <directory>
--cef <directory>
```

## Usage

```bash
nexs-swag init --codeExampleFilesDir ./code_samples
```

## Structure

```
13-code-examples/
├── main.go
└── code_samples/
    ├── create_user.go    # Example in Go
    ├── create_user.js    # Example in JavaScript
    ├── create_user.py    # Example in Python
    └── create_user.sh    # Example in Bash
```

## How It Works

### 1. In code, use the annotation:
```go
// @x-codeSamples file(create_user)
// @Router /users [post]
func CreateUser() {}
```

### 2. Create files with the prefix:
- `create_user.go`
- `create_user.js`
- `create_user.py`
- etc.

### 3. nexs-swag automatically detects language by extension

## Supported Languages

| Extension | Language |
|-----------|----------|
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
| ...and more |

## Generated OpenAPI

```json
{
  "paths": {
    "/users": {
      "post": {
        "x-codeSamples": [
          {
            "lang": "Go",
            "source": "// code from create_user.go..."
          },
          {
            "lang": "JavaScript",
            "source": "// code from create_user.js..."
          },
          {
            "lang": "Python",
            "source": "# code from create_user.py..."
          },
          {
            "lang": "Bash",
            "source": "#!/bin/bash\n# code from create_user.sh..."
          }
        ]
      }
    }
  }
}
```

## How to Run

```bash
./run.sh
```

## Benefits

1. **Multiple Languages:** Examples for different clients
2. **Auto-Detection:** Extension → language
3. **Reusability:** Same file for multiple endpoints
4. **Testable:** Code samples are real code that can be executed

## Use Cases

- **SDKs:** Examples for each supported language
- **Rich Documentation:** Facilitates integration
- **Onboarding:** Developers copy and paste
- **Testing:** Validated and functional examples
