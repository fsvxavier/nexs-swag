package openapi

import (
	"testing"

	swagger "github.com/fsvxavier/nexs-swag/pkg/openapi/v2"
	openapi "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

func TestSpecificationInterfaceV2(t *testing.T) {
	swagger := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:   "Test API V2",
			Version: "1.0.0",
		},
	}

	var spec Specification = swagger

	version := spec.GetVersion()
	if version != "2.0" {
		t.Errorf("GetVersion() = %v, want 2.0", version)
	}

	title := spec.GetTitle()
	if title != "Test API V2" {
		t.Errorf("GetTitle() = %v, want Test API V2", title)
	}

	info := spec.GetInfo()
	if info == nil {
		t.Error("GetInfo() returned nil")
	}
}

func TestSpecificationInterfaceV3(t *testing.T) {
	openapi := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:   "Test API V3",
			Version: "2.0.0",
		},
	}

	var spec Specification = openapi

	version := spec.GetVersion()
	if version != "3.1.0" {
		t.Errorf("GetVersion() = %v, want 3.1.0", version)
	}

	title := spec.GetTitle()
	if title != "Test API V3" {
		t.Errorf("GetTitle() = %v, want Test API V3", title)
	}

	info := spec.GetInfo()
	if info == nil {
		t.Error("GetInfo() returned nil")
	}
}

func TestSpecificationInterfacePolymorphism(t *testing.T) {
	specs := []Specification{
		&swagger.Swagger{
			Swagger: "2.0",
			Info: swagger.Info{
				Title:   "Swagger API",
				Version: "1.0.0",
			},
		},
		&openapi.OpenAPI{
			OpenAPI: "3.1.0",
			Info: openapi.Info{
				Title:   "OpenAPI",
				Version: "1.0.0",
			},
		},
	}

	expectedVersions := []string{"2.0", "3.1.0"}
	expectedTitles := []string{"Swagger API", "OpenAPI"}

	for i, spec := range specs {
		version := spec.GetVersion()
		if version != expectedVersions[i] {
			t.Errorf("spec[%d].GetVersion() = %v, want %v", i, version, expectedVersions[i])
		}

		title := spec.GetTitle()
		if title != expectedTitles[i] {
			t.Errorf("spec[%d].GetTitle() = %v, want %v", i, title, expectedTitles[i])
		}
	}
}

func TestSpecificationInterfaceInfoReturnType(t *testing.T) {
	swaggerSpec := &swagger.Swagger{
		Swagger: "2.0",
		Info: swagger.Info{
			Title:       "Test",
			Description: "Test Swagger",
			Version:     "1.0.0",
		},
	}

	openapiSpec := &openapi.OpenAPI{
		OpenAPI: "3.1.0",
		Info: openapi.Info{
			Title:       "Test",
			Description: "Test OpenAPI",
			Version:     "1.0.0",
		},
	}

	// Test V2 Info
	v2Info := swaggerSpec.GetInfo()
	if v2InfoTyped, ok := v2Info.(swagger.Info); ok {
		if v2InfoTyped.Description != "Test Swagger" {
			t.Errorf("V2 Info.Description = %v, want Test Swagger", v2InfoTyped.Description)
		}
	} else {
		t.Error("V2 GetInfo() did not return swagger.Info type")
	}

	// Test V3 Info
	v3Info := openapiSpec.GetInfo()
	if v3InfoTyped, ok := v3Info.(openapi.Info); ok {
		if v3InfoTyped.Description != "Test OpenAPI" {
			t.Errorf("V3 Info.Description = %v, want Test OpenAPI", v3InfoTyped.Description)
		}
	} else {
		t.Error("V3 GetInfo() did not return openapi.Info type")
	}
}
