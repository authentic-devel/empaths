package empaths

import (
	"fmt"
	"reflect"
	"strconv"
)

// toString converts a value to its string representation efficiently.
// It uses type switches for common types to avoid the overhead of fmt.Sprintf.
func toString(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// parseMapKey parses a string into a reflect.Value of the specified key type.
// It handles string, int, uint, bool, and float key types.
//
// Parameters:
//   - keyStr: The string representation of the key
//   - keyType: The reflect.Type of the map key
//
// Returns:
//   - The parsed key as a reflect.Value, or an invalid Value if parsing fails
func parseMapKey(keyStr string, keyType reflect.Type) reflect.Value {
	key := reflect.New(keyType).Elem()

	switch keyType.Kind() {
	case reflect.String:
		key.SetString(keyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(keyStr, 10, 64)
		if err != nil {
			return reflect.Value{}
		}
		key.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(keyStr, 10, 64)
		if err != nil {
			return reflect.Value{}
		}
		key.SetUint(uintVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(keyStr)
		if err != nil {
			return reflect.Value{}
		}
		key.SetBool(boolVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(keyStr, 64)
		if err != nil {
			return reflect.Value{}
		}
		key.SetFloat(floatVal)
	default:
		return reflect.Value{}
	}

	return key
}

// getMapValue retrieves a value from a map using a string key.
// It parses the key according to the map's key type and returns a copy of the value.
//
// Parameters:
//   - keyStr: The string representation of the key
//   - mapValue: The map to retrieve the value from
//
// Returns:
//   - The map value as a reflect.Value, or an invalid Value if the key doesn't exist
func getMapValue(keyStr string, mapValue reflect.Value) reflect.Value {
	keyType := mapValue.Type().Key()
	key := parseMapKey(keyStr, keyType)
	if !key.IsValid() {
		return reflect.Value{}
	}

	result := mapValue.MapIndex(key)
	if !result.IsValid() {
		return reflect.Value{}
	}

	// Make a copy of the map value to ensure it's addressable
	copyValue := reflect.New(result.Type()).Elem()
	copyValue.Set(result)
	return copyValue
}

// extractValue converts a reflect.Value to its interface{} representation.
// It handles special cases like pointers, nil slices, nil maps, interfaces,
// and unexported fields (which cannot be accessed via Interface()).
//
// Parameters:
//   - value: The reflect.Value to convert
//
// Returns:
//   - The value as an interface{} (any), or nil for invalid, nil, or inaccessible values
func extractValue(value reflect.Value) any {
	// Handle nil or invalid values
	if !value.IsValid() {
		return nil
	}

	// Handle nil pointers, slices, maps, and interfaces
	switch value.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface:
		if value.IsNil() {
			return nil
		}
	}

	// Dereference pointers to get their actual value
	if value.Kind() == reflect.Ptr {
		return extractValue(value.Elem())
	}

	// Handle interface values
	if value.Kind() == reflect.Interface {
		return extractValue(value.Elem())
	}

	// Check if we can safely call Interface() to avoid panics on unexported fields
	if !value.CanInterface() {
		return nil
	}

	return value.Interface()
}
