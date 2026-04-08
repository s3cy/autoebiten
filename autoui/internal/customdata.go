package internal

import (
	"fmt"
	"reflect"
)

// ExtractCustomData flattens custom data into a string map for XML attributes.
// It handles multiple input types:
//   - map[string]string: keys become attributes directly
//   - string/int/float/bool: becomes "custom_data" attribute
//   - struct with xml tags: uses tag values as attribute names
//   - struct without tags: uses field names as attribute names (preserving case)
func ExtractCustomData(data any) map[string]string {
	if data == nil {
		return nil
	}

	result := make(map[string]string)

	v := reflect.ValueOf(data)
	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		// Handle map[string]string
		if v.Type().Key().Kind() == reflect.String {
			iter := v.MapRange()
			for iter.Next() {
				key := iter.Key().String()
				val := iter.Value()
				if val.Kind() == reflect.String {
					result[key] = val.String()
				} else {
					result[key] = fmt.Sprintf("%v", val.Interface())
				}
			}
		}

	case reflect.String:
		result["custom_data"] = v.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result["custom_data"] = fmt.Sprintf("%d", v.Int())

	case reflect.Float32, reflect.Float64:
		result["custom_data"] = fmt.Sprintf("%v", v.Float())

	case reflect.Bool:
		result["custom_data"] = fmt.Sprintf("%v", v.Bool())

	case reflect.Struct:
		// Extract struct fields, checking for xml tags
		extractStructFields(v, result, "")

	default:
		// For other types, convert to string as fallback
		result["custom_data"] = fmt.Sprintf("%v", data)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

// extractStructFields recursively extracts fields from a struct.
// Nested structs are flattened with dot notation (e.g., "parent.child").
func extractStructFields(v reflect.Value, result map[string]string, prefix string) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Determine attribute name
		name := getXMLAttributeName(field)

		// Handle nested structs
		if fieldValue.Kind() == reflect.Struct {
			nestedPrefix := name
			if prefix != "" {
				nestedPrefix = prefix + "." + name
			}
			extractStructFields(fieldValue, result, nestedPrefix)
			continue
		}

		// Handle pointers to structs
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() && fieldValue.Elem().Kind() == reflect.Struct {
			nestedPrefix := name
			if prefix != "" {
				nestedPrefix = prefix + "." + name
			}
			extractStructFields(fieldValue.Elem(), result, nestedPrefix)
			continue
		}

		// Build full name with prefix
		fullName := name
		if prefix != "" {
			fullName = prefix + "." + name
		}

		// Convert value to string
		result[fullName] = valueToString(fieldValue)
	}
}

// getXMLAttributeName extracts the attribute name from struct field.
// If an xml tag is present with an "attr" option, it uses the tag name.
// Otherwise, it uses the field name as-is (preserving case).
func getXMLAttributeName(field reflect.StructField) string {
	tag := field.Tag.Get("xml")
	if tag != "" {
		// Parse xml tag: "name,attr" or just "name"
		// We only care about the name part
		name := parseXMLTagName(tag)
		if name != "" {
			return name
		}
	}

	// Use field name as-is (preserving case per spec)
	return field.Name
}

// parseXMLTagName extracts the name from an xml tag.
// Handles formats like "id,attr", "name", ",attr", etc.
func parseXMLTagName(tag string) string {
	for i := 0; i < len(tag); i++ {
		if tag[i] == ',' {
			return tag[:i]
		}
	}
	return tag
}

// valueToString converts a reflect.Value to its string representation.
func valueToString(v reflect.Value) string {
	// Handle nil pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", v.Float())
	case reflect.Bool:
		return fmt.Sprintf("%v", v.Bool())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}