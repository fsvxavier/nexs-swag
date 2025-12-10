package parser

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	f := NewFormatter()
	if f == nil {
		t.Fatal("NewFormatter() returned nil")
	}
}

func TestFormatSimpleFile(t *testing.T) {
	content := `package main

import "fmt"

// @title Test API
// @version 1.0
func main() {
	fmt.Println("Hello")
}
`
	f := NewFormatter()
	result, err := f.Format("test.go", []byte(content))
	if err != nil {
		t.Fatalf("Format() returned error: %v", err)
	}

	if len(result) == 0 {
		t.Error("Format() returned empty result")
	}

	resultStr := string(result)
	if !strings.Contains(resultStr, "package main") {
		t.Error("Formatted code should contain package declaration")
	}
}

func TestFormatInvalidGo(t *testing.T) {
	content := `this is not valid go code {{{`

	f := NewFormatter()
	_, err := f.Format("invalid.go", []byte(content))
	if err == nil {
		t.Error("Format() should return error for invalid Go code")
	}
}

func TestFormatEmptyFile(t *testing.T) {
	content := ``

	f := NewFormatter()
	_, err := f.Format("empty.go", []byte(content))
	if err == nil {
		t.Error("Format() should return error for empty file")
	}
}

func TestFormatSwaggerAnnotations(t *testing.T) {
	content := `package main

// @title My API
// @version 1.0.0
// @description This is my API
// @host localhost:8080
func main() {}
`
	f := NewFormatter()
	result, err := f.Format("test.go", []byte(content))
	if err != nil {
		t.Fatalf("Format() returned error: %v", err)
	}

	resultStr := string(result)
	if !strings.Contains(resultStr, "@title") {
		t.Error("Formatted code should contain @title annotation")
	}
	if !strings.Contains(resultStr, "@version") {
		t.Error("Formatted code should contain @version annotation")
	}
}

func TestFormatWithMultipleImports(t *testing.T) {
	content := `package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {}
`
	f := NewFormatter()
	result, err := f.Format("test.go", []byte(content))
	if err != nil {
		t.Fatalf("Format() returned error: %v", err)
	}

	resultStr := string(result)
	if !strings.Contains(resultStr, "import") {
		t.Error("Formatted code should contain import statement")
	}
}

func TestFormatWithComments(t *testing.T) {
	content := `package main

// Regular comment
// Another comment
func hello() {}

/* Block comment */
func world() {}
`
	f := NewFormatter()
	result, err := f.Format("test.go", []byte(content))
	if err != nil {
		t.Fatalf("Format() returned error: %v", err)
	}

	if len(result) == 0 {
		t.Error("Format() returned empty result")
	}
}

func TestSwaggerAnnotationsMap(t *testing.T) {
	expectedAnnotations := []string{
		"@summary",
		"@description",
		"@tags",
		"@accept",
		"@produce",
		"@param",
		"@success",
		"@failure",
		"@router",
		"@security",
	}

	for _, ann := range expectedAnnotations {
		if !swaggerAnnotations[ann] {
			t.Errorf("Swagger annotation %q should be in map", ann)
		}
	}
}

func TestFormatPreservesPackageName(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		packageName string
	}{
		{
			name:        "main package",
			content:     "package main\n\nfunc main() {}",
			packageName: "main",
		},
		{
			name:        "custom package",
			content:     "package mypackage\n\nfunc Hello() {}",
			packageName: "mypackage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter()
			result, err := f.Format("test.go", []byte(tt.content))
			if err != nil {
				t.Fatalf("Format() returned error: %v", err)
			}

			resultStr := string(result)
			if !strings.Contains(resultStr, "package "+tt.packageName) {
				t.Errorf("Formatted code should contain 'package %s'", tt.packageName)
			}
		})
	}
}

func TestFormatHandlesNoComments(t *testing.T) {
	content := `package main

func main() {
	println("no comments")
}
`
	f := NewFormatter()
	result, err := f.Format("test.go", []byte(content))
	if err != nil {
		t.Fatalf("Format() returned error: %v", err)
	}

	if len(result) == 0 {
		t.Error("Format() returned empty result")
	}
}

func TestFormatBuffer(t *testing.T) {
	f := NewFormatter()
	var buf bytes.Buffer

	content := `package test

func Test() {}`

	result, err := f.Format("test.go", []byte(content))
	if err != nil {
		t.Fatalf("Format() returned error: %v", err)
	}

	buf.Write(result)
	if buf.Len() == 0 {
		t.Error("Buffer should not be empty after formatting")
	}
}
