# Example 12 - Markdown Files

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates how to use markdown files for detailed descriptions.

## Flag

```bash
--markdownFiles <directory>
--md <directory>
```

## Usage

```bash
nexs-swag init --markdownFiles ./docs
```

## Syntax

In Go code, use `file(name.md)` in the description:

```go
// @Description file(create-user.md)
// @Router /users [post]
func CreateUser() {}
```

nexs-swag will:
1. Read `docs/create-user.md`
2. Replace `file(create-user.md)` with file content
3. Add content to OpenAPI description

## Structure

```
12-markdown-files/
├── main.go                   # API with file() references
└── docs/
    ├── create-user.md        # Detailed description
    └── get-user.md           # Another description
```

## Benefits

### 1. Detailed Descriptions
```markdown
# Create User

Creates a new user with validation.

## Request Body
- name: string (required)
- email: string (required, valid email)

## Validation Rules
...
```

### 2. Separation of Concerns
- Go code: application logic
- Markdown: detailed documentation
- Easier to maintain

### 3. Reusability
```go
// Multiple endpoints can use same markdown
// @Description file(auth-required.md)
func Endpoint1() {}

// @Description file(auth-required.md)
func Endpoint2() {}
```

### 4. Rich Formatting
- Headers
- Lists
- Code blocks
- Tables
- Links

## Real Example

```go
// @Description file(user-endpoints.md)
// @Router /users [post]
func CreateUser() {}
```

**user-endpoints.md:**
```markdown
# User Management

## Authentication Required
All user endpoints require Bearer token.

## Rate Limiting
- 100 requests/minute

## Response Codes
| Code | Description |
|------|-------------|
| 200  | Success |
| 401  | Unauthorized |
| 429  | Rate limit exceeded |
```

## How to Run

```bash
./run.sh
```

## Comparison

### Without Markdown
```json
{
  "description": "file(create-user.md)"
}
```

### With Markdown
```json
{
  "description": "# Create User Endpoint\n\nCreates a new user in the system...\n\n## Request Body\n..."
}
```

## Use Cases

- Complex APIs with lots of documentation
- Collaborative documentation (tech writers)
- Documentation versioning separate from code
- Documentation generation from Wiki
