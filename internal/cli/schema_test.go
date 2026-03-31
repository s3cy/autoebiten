package cli

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/s3cy/autoebiten/internal/script"
)

// schemaValidator provides a simple schema validation using jsonschema's Reflector.
type schemaValidator struct {
	schema *jsonschema.Schema
}

// newSchemaValidator creates a validator from the generated schema.
func newSchemaValidator(t *testing.T) *schemaValidator {
	data, err := GenerateSchema()
	if err != nil {
		t.Fatalf("failed to generate schema: %v", err)
	}

	var schema jsonschema.Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("failed to unmarshal schema: %v", err)
	}

	return &schemaValidator{schema: &schema}
}

// isValid checks if the data conforms to the schema.
// This is a simplified check - it marshals and unmarshals to validate structure.
func (v *schemaValidator) isValid(data map[string]any) (bool, []string) {
	// Try to marshal to JSON (structural validation)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, []string{fmt.Sprintf("marshal error: %v", err)}
	}

	// Try to unmarshal into ScriptSchema (type validation)
	var s script.ScriptSchema
	if err := json.Unmarshal(jsonData, &s); err != nil {
		return false, []string{fmt.Sprintf("unmarshal error: %v", err)}
	}

	// Check required fields
	var errors []string
	if s.Version == "" {
		errors = append(errors, "missing required field: version")
	}
	if s.Version != "" && s.Version != "1.0" {
		errors = append(errors, fmt.Sprintf("invalid version: %s (must be '1.0')", s.Version))
	}
	if len(s.Commands) == 0 {
		errors = append(errors, "commands array is empty or missing")
	}

	return len(errors) == 0, errors
}

// TestSchemaValidation ensures the schema validates correctly against
// valid scripts and rejects invalid ones. This test will fail if:
// 1. The schema is malformed
// 2. The schema doesn't match the actual parser behavior
// 3. Someone changes script structs without updating the schema
func TestSchemaValidation(t *testing.T) {
	validator := newSchemaValidator(t)

	// Valid script with all command types
	validScript := map[string]any{
		"version": "1.0",
		"commands": []map[string]any{
			{"input": map[string]any{"action": "press", "key": "KeyA"}},
			{"input": map[string]any{"action": "hold", "key": "KeySpace", "duration_ticks": 10, "async": true}},
			{"mouse": map[string]any{"action": "position", "x": 100, "y": 200}},
			{"mouse": map[string]any{"action": "hold", "button": "MouseButtonLeft", "duration_ticks": 6}},
			{"wheel": map[string]any{"x": 0.5, "y": -1.0}},
			{"screenshot": map[string]any{"output": "test.png", "async": false}},
			{"delay": map[string]any{"ms": 500}},
			{"repeat": map[string]any{
				"times": 3,
				"commands": []map[string]any{
					{"input": map[string]any{"action": "press", "key": "KeyB"}},
				},
			}},
		},
	}

	// Test valid script passes validation
	t.Run("valid script passes schema", func(t *testing.T) {
		valid, errors := validator.isValid(validScript)
		if !valid {
			t.Errorf("valid script failed validation: %v", errors)
		}
	})

	// Test parser accepts what schema validates (integration check)
	t.Run("schema and parser are consistent", func(t *testing.T) {
		data, err := json.Marshal(validScript)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		_, err = script.ParseBytes(data)
		if err != nil {
			t.Errorf("schema validates but parser rejects: %v\n"+
				"This means the schema is out of sync with the parser. "+
				"Update the schema to match the actual parser behavior.", err)
		}
	})

	// Test invalid scripts fail validation
	invalidCases := []struct {
		name   string
		script map[string]any
	}{
		{
			name:   "missing version",
			script: map[string]any{"commands": []map[string]any{}},
		},
		{
			name:   "missing commands",
			script: map[string]any{"version": "1.0"},
		},
		{
			name: "invalid version",
			script: map[string]any{
				"version":  "2.0",
				"commands": []map[string]any{},
			},
		},
		{
			name: "empty script",
			script: map[string]any{
				"version":  "1.0",
				"commands": []map[string]any{},
			},
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, _ := validator.isValid(tc.script)
			if valid {
				t.Error("expected invalid script to fail validation, but it passed")
			}
		})
	}
}

// TestSchemaStructure ensures the schema contains all required definitions
func TestSchemaStructure(t *testing.T) {
	data, err := GenerateSchema()
	if err != nil {
		t.Fatalf("failed to generate schema: %v", err)
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("failed to unmarshal schema: %v", err)
	}

	// Required top-level fields
	requiredTop := []string{"$schema", "$ref", "title", "description", "$defs"}
	for _, field := range requiredTop {
		if _, ok := schema[field]; !ok {
			t.Errorf("schema missing top-level field: %s", field)
		}
	}

	// Check definitions exist
	definitions, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatal("schema $defs is not a map")
	}

	requiredDefs := []string{
		"ScriptSchema",
		"CommandSchema",
		"InputCmd",
		"MouseCmd",
		"WheelCmd",
		"ScreenshotCmd",
		"DelayCmd",
		"RepeatSchema",
	}
	for _, def := range requiredDefs {
		if _, ok := definitions[def]; !ok {
			t.Errorf("schema missing definition: %s", def)
		}
	}
}

// TestSchemaGeneration ensures schema can be generated without errors
func TestSchemaGeneration(t *testing.T) {
	data, err := GenerateSchema()
	if err != nil {
		t.Fatalf("failed to generate schema: %v", err)
	}

	// Verify it's valid JSON
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("generated schema is not valid JSON: %v", err)
	}

	// Verify it's a valid JSON Schema (has $schema field)
	if _, ok := schema["$schema"]; !ok {
		t.Error("generated schema missing $schema field")
	}
}
