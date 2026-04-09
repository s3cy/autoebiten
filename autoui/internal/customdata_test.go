package internal_test

import (
	"testing"

	"github.com/s3cy/autoebiten/autoui/internal"
)

// TestExtractCustomData_Nil tests nil input.
func TestExtractCustomData_Nil(t *testing.T) {
	result := internal.ExtractCustomData(nil)
	if result != nil {
		t.Errorf("Expected nil for nil input, got %v", result)
	}
}

// TestExtractCustomData_Map tests map[string]string extraction.
func TestExtractCustomData_Map(t *testing.T) {
	input := map[string]string{
		"id":   "test-1",
		"role": "button",
	}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["id"] != "test-1" {
		t.Errorf("Expected id='test-1', got '%s'", result["id"])
	}
	if result["role"] != "button" {
		t.Errorf("Expected role='button', got '%s'", result["role"])
	}
}

// TestExtractCustomData_String tests string extraction.
func TestExtractCustomData_String(t *testing.T) {
	result := internal.ExtractCustomData("my-custom-data")
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["custom_data"] != "my-custom-data" {
		t.Errorf("Expected custom_data='my-custom-data', got '%s'", result["custom_data"])
	}
}

// TestExtractCustomData_Int tests integer extraction.
func TestExtractCustomData_Int(t *testing.T) {
	result := internal.ExtractCustomData(42)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["custom_data"] != "42" {
		t.Errorf("Expected custom_data='42', got '%s'", result["custom_data"])
	}
}

// TestExtractCustomData_Float tests float extraction.
func TestExtractCustomData_Float(t *testing.T) {
	result := internal.ExtractCustomData(3.14)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Float formatting may vary, so just check it's not empty
	if result["custom_data"] == "" {
		t.Error("Expected custom_data to have a value")
	}
}

// TestExtractCustomData_Bool tests boolean extraction.
func TestExtractCustomData_Bool(t *testing.T) {
	result := internal.ExtractCustomData(true)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["custom_data"] != "true" {
		t.Errorf("Expected custom_data='true', got '%s'", result["custom_data"])
	}
}

// TestExtractCustomData_StructWithTags tests struct with ae tags.
func TestExtractCustomData_StructWithTags(t *testing.T) {
	type WidgetMeta struct {
		ID      string `ae:"id"`
		Name    string `ae:"name"`
		Section string `ae:"section"`
	}

	input := WidgetMeta{
		ID:      "btn-1",
		Name:    "Start Button",
		Section: "main",
	}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["id"] != "btn-1" {
		t.Errorf("Expected id='btn-1', got '%s'", result["id"])
	}
	if result["name"] != "Start Button" {
		t.Errorf("Expected name='Start Button', got '%s'", result["name"])
	}
	if result["section"] != "main" {
		t.Errorf("Expected section='main', got '%s'", result["section"])
	}
}

// TestExtractCustomData_StructWithoutTags tests struct without xml tags.
func TestExtractCustomData_StructWithoutTags(t *testing.T) {
	type SimpleMeta struct {
		ID   string
		Name string
	}

	input := SimpleMeta{
		ID:   "btn-2",
		Name: "Quit Button",
	}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Field names should be preserved as-is (uppercase)
	if result["ID"] != "btn-2" {
		t.Errorf("Expected ID='btn-2', got '%s'", result["ID"])
	}
	if result["Name"] != "Quit Button" {
		t.Errorf("Expected Name='Quit Button', got '%s'", result["Name"])
	}
}

// TestExtractCustomData_NestedStruct tests nested struct flattening.
func TestExtractCustomData_NestedStruct(t *testing.T) {
	type Inner struct {
		Value string `ae:"value"`
	}

	type Outer struct {
		ID    string `ae:"id"`
		Inner Inner  `ae:"inner"`
	}

	input := Outer{
		ID:    "widget-1",
		Inner: Inner{Value: "nested-value"},
	}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["id"] != "widget-1" {
		t.Errorf("Expected id='widget-1', got '%s'", result["id"])
	}
	if result["inner.value"] != "nested-value" {
		t.Errorf("Expected inner.value='nested-value', got '%s'", result["inner.value"])
	}
}

// TestExtractCustomData_Pointer tests pointer input.
func TestExtractCustomData_Pointer(t *testing.T) {
	type Meta struct {
		ID string `ae:"id"`
	}

	input := &Meta{ID: "ptr-test"}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["id"] != "ptr-test" {
		t.Errorf("Expected id='ptr-test', got '%s'", result["id"])
	}
}

// TestExtractCustomData_NilPointer tests nil pointer input.
func TestExtractCustomData_NilPointer(t *testing.T) {
	var input *struct{ ID string }

	result := internal.ExtractCustomData(input)
	if result != nil {
		t.Errorf("Expected nil for nil pointer, got %v", result)
	}
}

// TestExtractCustomData_AETag tests struct with ae tags for custom attribute names.
func TestExtractCustomData_AETag(t *testing.T) {
	type WidgetMeta struct {
		ID      string `ae:"widget_id"`
		Name    string `ae:"display_name"`
		Section string // No tag, uses field name
	}

	input := WidgetMeta{
		ID:      "btn-1",
		Name:    "Start Button",
		Section: "main",
	}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["widget_id"] != "btn-1" {
		t.Errorf("Expected widget_id='btn-1', got '%s'", result["widget_id"])
	}
	if result["display_name"] != "Start Button" {
		t.Errorf("Expected display_name='Start Button', got '%s'", result["display_name"])
	}
	if result["Section"] != "main" {
		t.Errorf("Expected Section='main', got '%s'", result["Section"])
	}
}

// TestExtractCustomData_AETagIgnore tests ae:"-" to skip fields.
func TestExtractCustomData_AETagIgnore(t *testing.T) {
	type WidgetMeta struct {
		ID       string `ae:"widget_id"`
		Internal string `ae:"-"` // Should be ignored
		Name     string `ae:"display_name"`
	}

	input := WidgetMeta{
		ID:       "btn-1",
		Internal: "secret-data",
		Name:     "Start Button",
	}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["widget_id"] != "btn-1" {
		t.Errorf("Expected widget_id='btn-1', got '%s'", result["widget_id"])
	}
	if result["display_name"] != "Start Button" {
		t.Errorf("Expected display_name='Start Button', got '%s'", result["display_name"])
	}

	// Internal field should be ignored (not in result)
	if _, exists := result["Internal"]; exists {
		t.Errorf("Expected 'Internal' to be ignored, but found '%s'", result["Internal"])
	}
	if _, exists := result["-"]; exists {
		t.Errorf("Expected '-' to be ignored, but found '%s'", result["-"])
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 fields, got %d: %v", len(result), result)
	}
}


// TestExtractCustomData_NestedAETag tests nested struct with ae tags.
func TestExtractCustomData_NestedAETag(t *testing.T) {
	type Inner struct {
		Value string `ae:"inner_value"`
	}

	type Outer struct {
		ID    string `ae:"outer_id"`
		Inner Inner  `ae:"nested"`
	}

	input := Outer{
		ID:    "widget-1",
		Inner: Inner{Value: "nested-value"},
	}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["outer_id"] != "widget-1" {
		t.Errorf("Expected outer_id='widget-1', got '%s'", result["outer_id"])
	}
	if result["nested.inner_value"] != "nested-value" {
		t.Errorf("Expected nested.inner_value='nested-value', got '%s'", result["nested.inner_value"])
	}
}