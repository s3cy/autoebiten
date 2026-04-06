package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/internal/rpc"
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

// CheckCondition compares a queried value against an expected value using the given operator.
func CheckCondition(queried any, operator string, expected any) (bool, error) {
	// Determine the type of the queried value
	switch q := queried.(type) {
	case float64:
		e, ok := expected.(float64)
		if !ok {
			return false, fmt.Errorf("type mismatch: queried is number, expected is %T", expected)
		}
		return checkNumber(q, operator, e)

	case string:
		e, ok := expected.(string)
		if !ok {
			return false, fmt.Errorf("type mismatch: queried is string, expected is %T", expected)
		}
		return checkString(q, operator, e)

	case bool:
		e, ok := expected.(bool)
		if !ok {
			return false, fmt.Errorf("type mismatch: queried is bool, expected is %T", expected)
		}
		return checkBool(q, operator, e)

	default:
		return false, fmt.Errorf("cannot compare objects/arrays, use specific path to a primitive value")
	}
}

func checkNumber(queried float64, operator string, expected float64) (bool, error) {
	switch operator {
	case "==":
		return queried == expected, nil
	case "!=":
		return queried != expected, nil
	case "<":
		return queried < expected, nil
	case ">":
		return queried > expected, nil
	case "<=":
		return queried <= expected, nil
	case ">=":
		return queried >= expected, nil
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}

func checkString(queried, operator, expected string) (bool, error) {
	switch operator {
	case "==":
		return queried == expected, nil
	case "!=":
		return queried != expected, nil
	default:
		return false, fmt.Errorf("operator %s not supported for string values", operator)
	}
}

func checkBool(queried bool, operator string, expected bool) (bool, error) {
	switch operator {
	case "==":
		return queried == expected, nil
	case "!=":
		return queried != expected, nil
	default:
		return false, fmt.Errorf("operator %s not supported for boolean values", operator)
	}
}

// RunWaitForCommand polls until a condition is met or timeout expires.
func (e *CommandExecutor) RunWaitForCommand(conditionStr, timeoutStr, intervalStr string) error {
	// Parse condition
	cond, err := ParseCondition(conditionStr)
	if err != nil {
		return err
	}

	// Parse timeout
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return fmt.Errorf("invalid timeout: %w", err)
	}

	// Parse interval (default to 100ms)
	interval := 100 * time.Millisecond
	if intervalStr != "" {
		interval, err = time.ParseDuration(intervalStr)
		if err != nil {
			return fmt.Errorf("invalid interval: %w", err)
		}
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Determine the custom command name based on condition type
	var customName string
	if cond.Type == "state" {
		customName = autoebiten.StateExporterPathPrefix + cond.Name
	} else {
		customName = cond.Name
	}

	start := time.Now()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout exceeded after %s waiting for condition", timeoutStr)

		case <-ticker.C:
			// Query the value
			params := &rpc.CustomParams{
				Name:    customName,
				Request: cond.Path,
			}

			req, err := rpc.BuildRequest("custom", params)
			if err != nil {
				continue // Retry on transient errors
			}

			resp, err := rpc.SendRequestSocket(req)
			if err != nil {
				continue // Retry on transient errors
			}

			if resp.Error != nil {
				continue // Retry on game-side errors (e.g., path not found yet)
			}

			var result rpc.CustomResult
			if err := json.Unmarshal(resp.Result, &result); err != nil {
				continue
			}

			// Parse the response as JSON
			var queried any
			if err := json.Unmarshal([]byte(result.Response), &queried); err != nil {
				continue
			}

			// Check condition
			met, err := CheckCondition(queried, cond.Operator, cond.Value)
			if err != nil {
				return err
			}

			if met {
				elapsed := time.Since(start).Round(100 * time.Millisecond)
				e.writer.Success(fmt.Sprintf("condition met after %s", elapsed))
				return nil
			}
		}
	}
}
