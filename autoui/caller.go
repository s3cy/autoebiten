package autoui

import (
	"fmt"
	"reflect"

	"github.com/ebitenui/ebitenui/widget"
)

// InvokeMethod invokes a method on a widget using reflection with whitelist safety.
// Also handles proxy methods that bypass reflection (e.g., SelectEntryByIndex).
//
// Only whitelisted method signatures are allowed for security:
// - func()
// - func(any/interface{})
// - func(bool), func(int), func(float64), func(string)
// - func(enum types) - types with underlying basic kind
//
// Return values are converted to Go's "largest" representation for JSON compatibility:
// - Integers (int, int32, int64, enums) → int64
// - Floats (float32, float64) → float64
// - bool, string, slices → unchanged
func InvokeMethod(w widget.PreferredSizeLocateableWidget, methodName string, args []any) (any, error) {
	if w == nil {
		return nil, fmt.Errorf("widget is nil")
	}

	// Check proxy registry first
	if handler := GetProxyHandler(methodName); handler != nil {
		return handler(w, args)
	}

	method := reflect.ValueOf(w).MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method '%s' not found on widget type %T", methodName, w)
	}

	methodType := method.Type()

	if !isWhitelistedSignature(methodType) {
		return nil, fmt.Errorf("method '%s' has non-whitelisted signature %s", methodName, methodType)
	}

	expectedArgs := methodType.NumIn()
	if len(args) != expectedArgs {
		return nil, fmt.Errorf("method '%s' expects %d arguments, got %d", methodName, expectedArgs, len(args))
	}

	convertedArgs, err := convertArgs(args, methodType)
	if err != nil {
		return nil, fmt.Errorf("argument conversion failed: %w", err)
	}

	results := method.Call(convertedArgs)

	// Capture return value
	if len(results) > 0 {
		ret := results[0]

		// Check for error return
		if errVal, ok := ret.Interface().(error); ok && errVal != nil {
			return nil, fmt.Errorf("method '%s' returned error: %w", methodName, errVal)
		}

		// Convert enum types to underlying type for JSON serialization
		switch ret.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			return ret.Int(), nil
		case reflect.Float32, reflect.Float64:
			return ret.Float(), nil
		case reflect.Bool:
			return ret.Bool(), nil
		case reflect.String:
			return ret.String(), nil
		default:
			// Slice, struct, interface - return as-is
			return ret.Interface(), nil
		}
	}

	return nil, nil
}

// isWhitelistedSignature checks if a method signature is allowed for invocation.
// Whitelisted signatures:
// - func()
// - func(any/interface{})
// - func(bool), func(int), func(float64), func(string)
// - func(enum types) - types with underlying basic kind (int, float, string, bool)
// - func() any/interface{} (return value capture)
// - func() []any, func() []string, func() []int, etc. (slice returns)
// - func() error (no args, returns error)
// - func(error) error
// - func(any) error, func(bool) error, func(int) error, etc.
//
// Note: Custom types based on basic types (enums like widget.WidgetState) ARE now whitelisted.
func isWhitelistedSignature(t reflect.Type) bool {
	numIn := t.NumIn()
	numOut := t.NumOut()

	// Check return values - must be either 0 or 1 (single return)
	if numOut > 1 {
		return false
	}

	// If there's a return value, check if it's allowed
	if numOut == 1 {
		returnType := t.Out(0)

		// Allow error
		if returnType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return true
		}

		// Allow any/interface{} (empty PkgPath, interface Kind)
		if returnType.Kind() == reflect.Interface && returnType.PkgPath() == "" {
			return true
		}

		// Allow slices of basic types or any
		if returnType.Kind() == reflect.Slice {
			elemType := returnType.Elem()
			// []any
			if elemType.Kind() == reflect.Interface && elemType.PkgPath() == "" {
				return true
			}
			// []string, []int, []bool, etc. (element has empty PkgPath = built-in)
			if elemType.PkgPath() == "" {
				switch elemType.Kind() {
				case reflect.Bool, reflect.Int, reflect.Int32, reflect.Int64,
					reflect.Float32, reflect.Float64, reflect.String:
					return true
				}
			}
		}

		// Allow types with underlying basic Kind (enums and built-in types)
		switch returnType.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.String:
			return true
		}

		return false
	}

	// Check input parameters against whitelist
	switch numIn {
	case 0:
		// func() - always whitelisted
		return true
	case 1:
		// Check single parameter type
		paramType := t.In(0)

		// Allow any/interface{} (empty PkgPath with interface Kind)
		if paramType.Kind() == reflect.Interface && paramType.PkgPath() == "" {
			return true
		}

		// Allow built-in basic types (PkgPath is empty)
		if paramType.PkgPath() == "" {
			switch paramType.Kind() {
			case reflect.Bool, reflect.Int, reflect.Int32, reflect.Int64,
				reflect.Float32, reflect.Float64, reflect.String:
				return true
			default:
				return false
			}
		}

		// Allow enum types (custom types with underlying basic kind)
		switch paramType.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.String:
			return true
		default:
			return false
		}
	default:
		// Multiple parameters not whitelisted
		return false
	}
}

// convertArgs converts []any arguments to []reflect.Value matching the method's parameter types.
func convertArgs(args []any, methodType reflect.Type) ([]reflect.Value, error) {
	if len(args) == 0 {
		return nil, nil
	}

	converted := make([]reflect.Value, len(args))
	for i, arg := range args {
		targetType := methodType.In(i)
		convVal, convErr := convertArg(arg, targetType)
		if convErr != nil {
			return nil, fmt.Errorf("argument %d: %w", i, convErr)
		}
		converted[i] = convVal
	}

	return converted, nil
}

// convertArg converts a single argument to the target reflect.Type.
// Handles:
// - Direct type match
// - interface{} targets - accept any non-nil value (must check BEFORE ConvertibleTo)
// - Convertible types (int → enum, etc.) via reflect.ConvertibleTo
// - Numeric conversions: float64 → int (truncation), int → float64
// - bool, string direct matches only
func convertArg(arg any, targetType reflect.Type) (reflect.Value, error) {
	if arg == nil {
		return reflect.Value{}, fmt.Errorf("nil argument not supported")
	}

	argValue := reflect.ValueOf(arg)
	argType := argValue.Type()

	// Direct type match
	if argType == targetType {
		return argValue, nil
	}

	// Handle interface{} targets FIRST - accept any non-nil value
	// Must check before ConvertibleTo since all types are convertible to interface{}
	if targetType.Kind() == reflect.Interface && targetType.PkgPath() == "" {
		return argValue, nil
	}

	// If arg is convertible to target type (int → enum, etc.)
	if argType.ConvertibleTo(targetType) {
		return argValue.Convert(targetType), nil
	}

	// Handle conversions based on target type's underlying kind
	switch targetType.Kind() {
	case reflect.Bool:
		if argType.Kind() == reflect.Bool {
			return argValue, nil
		}
		return reflect.Value{}, fmt.Errorf("cannot convert %s to bool", argType)

	case reflect.Int, reflect.Int32, reflect.Int64:
		switch argType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return argValue.Convert(targetType), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return argValue.Convert(targetType), nil
		case reflect.Float32, reflect.Float64:
			return argValue.Convert(targetType), nil
		default:
			return reflect.Value{}, fmt.Errorf("cannot convert %s to %s", argType, targetType)
		}

	case reflect.Float32, reflect.Float64:
		switch argType.Kind() {
		case reflect.Float32, reflect.Float64:
			return argValue.Convert(targetType), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return argValue.Convert(targetType), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return argValue.Convert(targetType), nil
		default:
			return reflect.Value{}, fmt.Errorf("cannot convert %s to %s", argType, targetType)
		}

	case reflect.String:
		if argType.Kind() == reflect.String {
			return argValue, nil
		}
		return reflect.Value{}, fmt.Errorf("cannot convert %s to string", argType)

	case reflect.Interface:
		// For non-empty interface targets that weren't caught by ConvertibleTo check,
		// the argument does not implement the interface - return error
		return reflect.Value{}, fmt.Errorf("cannot convert %s to %s (interface not implemented)", argType, targetType)

	default:
		return reflect.Value{}, fmt.Errorf("unsupported target type %s", targetType)
	}
}
