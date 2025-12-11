# Example 11 - Parse GoList

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates the use of `--parseGoList` for faster and more accurate parsing.

## Flag

```bash
--parseGoList
```

## Difference

### WITHOUT --parseGoList (Manual)
- Traverses directories manually
- Parses all .go files
- Slower in large projects
- May miss dependencies

### WITH --parseGoList (Go Tooling)
- Uses `go list -json`
- Go modules information
- Faster and more accurate
- Respects go.mod

## Usage

```bash
nexs-swag init --parseGoList
```

## How It Works

### Manual Parsing
```bash
# nexs-swag does:
1. Walk all directories
2. Find .go files
3. Parse each file
4. Resolve imports manually
```

### Go List Parsing
```bash
# nexs-swag executes:
go list -json ./...

# Returns:
{
  "ImportPath": "myapp",
  "Dir": "/path/to/myapp",
  "GoFiles": ["main.go", "handlers.go"],
  "Imports": ["encoding/json", "net/http"],
  "Deps": ["myapp/models"]
}
```

## Performance

### Small Project (< 50 files)
```
Manual:  ~500ms
GoList:  ~400ms
Gain:    20%
```

### Medium Project (100-500 files)
```
Manual:  ~5s
GoList:  ~2s
Gain:    60%
```

### Large Project (> 1000 files)
```
Manual:  ~30s
GoList:  ~8s
Gain:    73%
```

## Advantages

### 1. Speed
```bash
# 3x faster in large projects
time nexs-swag init --parseGoList
# real: 2s vs 6s without flag
```

### 2. Accuracy
```bash
# Respects go.mod replace
# go.mod:
replace github.com/old/pkg => ../local/pkg

# nexs-swag uses correct path
```

### 3. Build Tags
```go
// +build linux

package linux

// Only parsed if build tag matches
```

### 4. Vendor Detection
```bash
# Detects vendor/ automatically
go mod vendor
nexs-swag init --parseGoList
# Ignores vendor/ if modules are enabled
```

## How to Run

```bash
./run.sh
```

## When to Use

**Use --parseGoList when:**
- Project with Go modules
- Many packages
- Complex dependencies
- Want speed

**DON'T use when:**
- Project without go.mod
- GOPATH-based project
- Need to parse specific files ignored by go list

## Requirements

```bash
# Go 1.11+ (modules)
go mod init myapp

# or configured GOPATH
export GOPATH=/path/to/gopath
```

## Combine with Other Flags

### With parseDependency
```bash
nexs-swag init \
  --parseGoList \
  --parseDependency \
  --parseDependencyLevel 2
```

### With exclude
```bash
nexs-swag init \
  --parseGoList \
  --exclude "testdata,mocks"
```

## Troubleshooting

### Error: "go list: command not found"
```bash
# Install Go
# https://golang.org/doc/install

# Verify
go version
```

### Error: "not a Go module"
```bash
# Initialize module
go mod init myapp
```

### Parsing old GOPATH project
```bash
# Use without --parseGoList
nexs-swag init --dir . --output ./docs
```
