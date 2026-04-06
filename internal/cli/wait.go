package cli

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Condition represents a parsed wait condition.
type Condition struct {
	Type     string // "state" or "custom"
	Name     string // exporter or command name
	Path     string // path or request
	Operator string // ==, !=, <, >, <=, >=
	Value    any    // expected value (parsed from JSON)
}

// ParseCondition parses a condition string in the format "type:name:path operator value".
func ParseCondition(s string) (*Condition, error) {
	// Find the operator
	operators := []string{"==", "!=", "<=", ">=", "<", ">"}
	var op string
	var opIndex int

	for _, operator := range operators {
		idx := strings.Index(s, " "+operator+" ")
		if idx != -1 {
			op = operator
			opIndex = idx
			break
		}
	}

	if op == "" {
		return nil, fmt.Errorf(`invalid condition format: expected "type:name:path operator value"`)
	}

	// Split into query part and value part
	queryPart := strings.TrimSpace(s[:opIndex])
	valuePart := strings.TrimSpace(s[opIndex+len(op)+2:])

	// Parse query part: type:name:path
	parts := strings.SplitN(queryPart, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf(`invalid condition format: expected "type:name:path operator value"`)
	}

	condType := parts[0]
	if condType != "state" && condType != "custom" {
		return nil, fmt.Errorf("invalid condition type: %q (expected 'state' or 'custom')", condType)
	}

	// Parse value as JSON
	var value any
	if err := json.Unmarshal([]byte(valuePart), &value); err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	return &Condition{
		Type:     condType,
		Name:     parts[1],
		Path:     parts[2],
		Operator: op,
		Value:    value,
	}, nil
}