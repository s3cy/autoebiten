package autoui_test

import (
	"encoding/json"
	"testing"

	"github.com/s3cy/autoebiten/autoui"
)

func TestExistsResponse_JSON(t *testing.T) {
	// Test found=true case
	resp := autoui.ExistsResponse{Found: true, Count: 2}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal ExistsResponse: %v", err)
	}
	expected := `{"found":true,"count":2}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}

	// Test found=false case
	resp = autoui.ExistsResponse{Found: false, Count: 0}
	data, err = json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal ExistsResponse: %v", err)
	}
	expected = `{"found":false,"count":0}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}
