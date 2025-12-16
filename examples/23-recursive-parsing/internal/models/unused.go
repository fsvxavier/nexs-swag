package models

// UnusedModel is a model that is NOT referenced in any API operation
// This should NOT appear in the generated OpenAPI schema
type UnusedModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// AnotherUnusedStruct should also not appear
type AnotherUnusedStruct struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}
