# Examples - nexs-swag

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

This directory contains usage examples for each flag and functionality of nexs-swag.

## Prerequisites

Install nexs-swag before running the examples:

```bash
# From the project root directory
cd ..
go install ./cmd/nexs-swag

# Or use the installation script
./install.sh

# Verify installation
nexs-swag --version
```

## Structure

Each subdirectory contains a specific example with:
- `main.go` - Go code with Swagger annotations
- `README.md` - Detailed usage instructions (🌍 Available in 3 languages)
- `run.sh` - Script to run the example

## Examples List

### Basic (01-08)
- [01-basic](./01-basic) - Basic usage with `--dir` and `--output`
- [02-formats](./02-formats) - Multiple formats with `--format`
- [03-general-info](./03-general-info) - Specific file with `--generalInfo`
- [04-property-strategy](./04-property-strategy) - `--propertyStrategy` (snake_case, camelCase, PascalCase)
- [05-required-default](./05-required-default) - `--requiredByDefault`
- [06-exclude](./06-exclude) - `--exclude` to exclude directories
- [07-tags-filter](./07-tags-filter) - `--tags` to filter by tags
- [08-parse-internal](./08-parse-internal) - `--parseInternal`

### Dependencies (09-11)
- [09-parse-dependency](./09-parse-dependency) - `--parseDependency`
- [10-dependency-level](./10-dependency-level) - `--parseDependencyLevel` (0-3)
- [11-parse-golist](./11-parse-golist) - `--parseGoList`

### External Content (12-14)
- [12-markdown-files](./12-markdown-files) - `--markdownFiles`
- [13-code-examples](./13-code-examples) - `--codeExampleFilesDir`
- [14-overrides-file](./14-overrides-file) - `--overridesFile`

### Configuration (15-18)
- [15-generated-time](./15-generated-time) - `--generatedTime`
- [16-instance-name](./16-instance-name) - `--instanceName`
- [17-template-delims](./17-template-delims) - `--templateDelims`
- [18-collection-format](./18-collection-format) - `--collectionFormat`

### Advanced (19-22)
- [19-parse-func-body](./19-parse-func-body) - `--parseFuncBody`
- [20-fmt-command](./20-fmt-command) - `fmt` command
- [21-struct-tags](./21-struct-tags) - swaggertype, swaggerignore, extensions
- [22-openapi-v2](./22-openapi-v2) - `--openapi-version` (Swagger 2.0 / OpenAPI 3.1.0)

## How to Use

### Run a specific example

```bash
cd 01-basic
./run.sh
```

### Run manually

```bash
cd 01-basic
nexs-swag init --dir . --output ./docs
```

### Run all examples

```bash
for dir in */; do
    echo "=== Running $dir ==="
    cd "$dir"
    ./run.sh
    cd ..
    echo ""
done
```

## Example Structure

```
XX-example-name/
├── main.go          # HTTP server with annotations
├── run.sh           # Demo script
└── README.md        # Complete documentation (🌍 3 languages)
```

## Tips

### View generated documentation

```bash
# OpenAPI 3.1.0 (default)
cat docs/openapi.json | jq
cat docs/openapi.yaml

# Swagger 2.0 (if generated with --openapi-version 2.0)
cat docs/swagger.json | jq
cat docs/swagger.yaml

# Go docs
cat docs/docs.go
```

### Serve with Swagger UI

```bash
# Install swagger ui
docker run -p 8080:8080 \
  -e SWAGGER_JSON=/docs/openapi.json \
  -v $(pwd)/docs:/docs \
  swaggerapi/swagger-ui

# Access: http://localhost:8080
```

### Integrate in projects

```go
package main

import (
    "net/http"
    
    _ "myapp/docs"  // Import generated docs
    
    httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
    http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Troubleshooting

### Error "nexs-swag: command not found"

```bash
# Check installation
which nexs-swag

# If not installed
cd ..
go install ./cmd/nexs-swag

# Check if $GOPATH/bin is in PATH
echo $PATH | grep $(go env GOPATH)/bin

# Add to PATH if necessary
export PATH=$PATH:$(go env GOPATH)/bin
```

### Error generating documentation

```bash
# Check if code compiles
go build .

# Run with more details
nexs-swag init --dir . --output ./docs --debug
```

### Clean previous documentation

```bash
rm -rf docs docs-*
```

## Resources

- [Complete Documentation](../INSTALL.md)
- [swaggo/swag - Original Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
