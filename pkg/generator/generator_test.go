package generator

import (
	"testing"

	v2 "github.com/fsvxavier/nexs-swag/pkg/generator/v2"
	v3 "github.com/fsvxavier/nexs-swag/pkg/generator/v3"
	swaggerv2 "github.com/fsvxavier/nexs-swag/pkg/openapi/v2"
	openapiV3 "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

func TestGeneratorInterfaceV2(t *testing.T) {
	spec := &swaggerv2.Swagger{
		Swagger: "2.0",
		Info: swaggerv2.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	tmpDir := t.TempDir()
	gen := v2.New(spec, tmpDir, []string{"json"})

	var g Generator = gen

	// Test interface methods
	g.SetInstanceName("myapi")
	g.SetGeneratedTime(true)

	err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}
}

func TestGeneratorInterfaceV3(t *testing.T) {
	spec := &openapiV3.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapiV3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	tmpDir := t.TempDir()
	gen := v3.New(spec, tmpDir, []string{"json"})

	var g Generator = gen

	// Test interface methods
	g.SetInstanceName("myapi")
	g.SetGeneratedTime(true)

	err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}
}

func TestGeneratorInterfacePolymorphism(t *testing.T) {
	tmpDir := t.TempDir()

	generators := []Generator{
		v2.New(&swaggerv2.Swagger{
			Swagger: "2.0",
			Info: swaggerv2.Info{
				Title:   "API V2",
				Version: "1.0.0",
			},
		}, tmpDir+"/v2", []string{"json"}),
		v3.New(&openapiV3.OpenAPI{
			OpenAPI: "3.1.0",
			Info: openapiV3.Info{
				Title:   "API V3",
				Version: "1.0.0",
			},
		}, tmpDir+"/v3", []string{"json"}),
	}

	for i, gen := range generators {
		gen.SetInstanceName("testapi")
		gen.SetGeneratedTime(false)

		err := gen.Generate()
		if err != nil {
			t.Errorf("generators[%d].Generate() returned error: %v", i, err)
		}
	}
}

func TestGeneratorInterfaceMethods(t *testing.T) {
	spec := &swaggerv2.Swagger{
		Swagger: "2.0",
		Info: swaggerv2.Info{
			Title:   "Test",
			Version: "1.0.0",
		},
	}

	tmpDir := t.TempDir()
	gen := v2.New(spec, tmpDir, []string{"json"})

	// Test SetInstanceName
	gen.SetInstanceName("custom")

	// Test SetGeneratedTime
	gen.SetGeneratedTime(true)
	gen.SetGeneratedTime(false)

	// Verify Generate works
	err := gen.Generate()
	if err != nil {
		t.Errorf("Generate() returned error: %v", err)
	}
}
