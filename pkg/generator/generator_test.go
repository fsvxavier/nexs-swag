package generator

import (
	"testing"

	"github.com/fsvxavier/nexs-swag/pkg/openapi"
)

func TestNew(t *testing.T) {
	spec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
	}
	outputDir := "/tmp/test"
	outputTypes := []string{"json", "yaml"}

	gen := New(spec, outputDir, outputTypes)
	if gen == nil {
		t.Fatal("New() returned nil")
	}

	if gen.spec != spec {
		t.Error("Generator spec not set correctly")
	}
	if gen.outputDir != outputDir {
		t.Error("Generator outputDir not set correctly")
	}
	if len(gen.outputType) != len(outputTypes) {
		t.Errorf("Generator outputType length = %d, want %d", len(gen.outputType), len(outputTypes))
	}
}

func TestSetInstanceName(t *testing.T) {
	gen := &Generator{
		spec: &openapi.OpenAPI{},
	}

	tests := []struct {
		name         string
		instanceName string
	}{
		{"default swagger", "swagger"},
		{"custom name", "api"},
		{"empty name", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen.SetInstanceName(tt.instanceName)
			if gen.instanceName != tt.instanceName {
				t.Errorf("SetInstanceName() = %q, want %q", gen.instanceName, tt.instanceName)
			}
		})
	}
}

func TestSetGeneratedTime(t *testing.T) {
	gen := &Generator{
		spec: &openapi.OpenAPI{},
	}

	gen.SetGeneratedTime(true)
	if !gen.generatedTime {
		t.Error("SetGeneratedTime(true) failed")
	}

	gen.SetGeneratedTime(false)
	if gen.generatedTime {
		t.Error("SetGeneratedTime(false) failed")
	}
}

func TestSetTemplateDelims(t *testing.T) {
	gen := &Generator{
		spec: &openapi.OpenAPI{},
	}

	tests := []struct {
		name      string
		delims    string
		wantLeft  string
		wantRight string
	}{
		{"default", "{{,}}", "{{", "}}"},
		{"custom", "<%,%>", "<%", "%>"},
		{"brackets", "[[,]]", "[[", "]]"},
		{"empty uses default", "", "{{", "}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen.SetTemplateDelims(tt.delims)
			if len(gen.templateDelims) != 2 {
				t.Fatalf("templateDelims length = %d, want 2", len(gen.templateDelims))
			}
			if gen.templateDelims[0] != tt.wantLeft {
				t.Errorf("leftDelim = %q, want %q", gen.templateDelims[0], tt.wantLeft)
			}
			if gen.templateDelims[1] != tt.wantRight {
				t.Errorf("rightDelim = %q, want %q", gen.templateDelims[1], tt.wantRight)
			}
		})
	}
}
