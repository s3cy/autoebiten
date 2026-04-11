package autoui

import (
	"fmt"
	"reflect"

	"github.com/ebitenui/ebitenui/widget"
)

// InvokeMethod invokes a method on a widget using reflection with whitelist safety.
// Only whitelisted method signatures are allowed for security.
// Supported signatures: func(), func(bool), func(int), func(float64), func(string)
func InvokeMethod(w widget.PreferredSizeLocateableWidget, methodName string, args []any) error {
	if w == nil {
		return fmt.Errorf("widget is nil")
	}

	// Get the method by name
	method := reflect.ValueOf(w).MethodByName(methodName)
	if !method.IsValid() {
		return fmt.Errorf("method '%s' not found on widget type %T", methodName, w)
	}

	// Get method type for signature checking
	methodType := method.Type()

	// Check if signature is whitelisted
	if !isWhitelistedSignature(methodType) {
		return fmt.Errorf("method '%s' has non-whitelisted signature %s", methodName, methodType)
	}

	// Check argument count
	expectedArgs := methodType.NumIn()
	if len(args) != expectedArgs {
		return fmt.Errorf("method '%s' expects %d arguments, got %d", methodName, expectedArgs, len(args))
	}

	// Convert arguments to reflect.Value
	convertedArgs, err := convertArgs(args, methodType)
	if err != nil {
		return fmt.Errorf("argument conversion failed: %w", err)
	}

	// Call the method
	results := method.Call(convertedArgs)

	// Check for error return value (if method returns error)
	if len(results) > 0 {
		if errVal, ok := results[0].Interface().(error); ok && errVal != nil {
			return fmt.Errorf("method '%s' returned error: %w", methodName, errVal)
		}
	}

	return nil
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
// Handles numeric conversions:
// - float64 -> int (truncation)
// - int -> float64 (conversion)
// - bool, string direct matches only
func convertArg(arg any, targetType reflect.Type) (reflect.Value, error) {
	if arg == nil {
		return reflect.Value{}, fmt.Errorf("nil argument not supported")
	}

	argType := reflect.TypeOf(arg)
	argValue := reflect.ValueOf(arg)

	// Direct type match
	if argType == targetType {
		return argValue, nil
	}

	// Handle conversions based on target type
	switch targetType.Kind() {
	case reflect.Bool:
		// Only accept bool for bool parameters
		if argType.Kind() == reflect.Bool {
			return argValue, nil
		}
		return reflect.Value{}, fmt.Errorf("cannot convert %s to bool", argType)

	case reflect.Int, reflect.Int32, reflect.Int64:
		// Accept int types and float64 (with truncation)
		switch argType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return argValue.Convert(targetType), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return argValue.Convert(targetType), nil
		case reflect.Float32, reflect.Float64:
			// Truncate float to int
			return argValue.Convert(targetType), nil
		default:
			return reflect.Value{}, fmt.Errorf("cannot convert %s to %s", argType, targetType)
		}

	case reflect.Float32, reflect.Float64:
		// Accept float types and int types (conversion)
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
		// Only accept string for string parameters
		if argType.Kind() == reflect.String {
			return argValue, nil
		}
		return reflect.Value{}, fmt.Errorf("cannot convert %s to string", argType)

	default:
		return reflect.Value{}, fmt.Errorf("unsupported target type %s", targetType)
	}
}
