package empaths

import (
	"reflect"
	"strconv"
	"strings"
)

// resolvePathAgainstValue resolves a path against a reflect.Value.
// This function handles the actual resolution of a model path against a data object using reflection.
//
// Parameters:
//   - path: The path string to resolve (e.g., "User.Address.City")
//   - value: The reflect.Value to resolve the path against
//
// Returns:
//   - The resolved reflect.Value
func resolvePathAgainstValue(path string, value reflect.Value) reflect.Value {
	// Handle nil or invalid values
	if !value.IsValid() {
		return reflect.Value{}
	}

	// Remove leading dot if present
	if len(path) > 0 && path[0] == '.' {
		path = path[1:]
	}

	// If the path is empty, return the value itself
	if path == "" {
		return value
	}

	// Handle pointers and interfaces
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return reflect.Value{}
		}
		return resolvePathAgainstValue(path, value.Elem())
	}

	// Split the path into segments
	return resolvePathSegments(path, value)
}

// resolvePathSegments handles the resolution of path segments against a reflect.Value.
// It supports property access, array/slice indexing, and map access.
// Uses single-pass scanning to find delimiters efficiently.
//
// Parameters:
//   - path: The path string to resolve (e.g., "User.Address" or "Users[0]")
//   - value: The reflect.Value to resolve the path against
//
// Returns:
//   - The resolved reflect.Value
func resolvePathSegments(path string, value reflect.Value) reflect.Value {
	// Check if the path starts with an array/map index
	if len(path) > 0 && path[0] == '[' {
		return resolveArrayOrMapAccess(path, value)
	}

	// Single-pass scan to find first '.' or '['
	var currentSegment string
	var remainingPath string
	splitIdx := -1
	splitChar := byte(0)

	for i := 0; i < len(path); i++ {
		c := path[i]
		if c == '.' || c == '[' {
			splitIdx = i
			splitChar = c
			break
		}
	}

	if splitIdx == -1 {
		// No more segments, this is the final part
		currentSegment = path
		remainingPath = ""
	} else if splitChar == '.' {
		// Dot comes first
		currentSegment = path[:splitIdx]
		remainingPath = path[splitIdx+1:]
	} else {
		// Bracket comes first
		currentSegment = path[:splitIdx]
		remainingPath = path[splitIdx:]
	}

	// Resolve the current segment
	resolvedValue := resolveFieldOrMethod(currentSegment, value)

	// If we couldn't resolve the current segment or there's no remaining path, return the result
	if !resolvedValue.IsValid() || remainingPath == "" {
		return resolvedValue
	}

	// Continue resolving with the remaining path
	return resolvePathAgainstValue(remainingPath, resolvedValue)
}

// resolveArrayOrMapAccess handles array, slice, and map access with brackets.
// It processes path segments that start with '[' for accessing elements by index or key.
//
// Parameters:
//   - path: The path string to resolve (e.g., "[0]" or "[\"key\"]")
//   - value: The reflect.Value to resolve the path against
//
// Returns:
//   - The resolved reflect.Value
func resolveArrayOrMapAccess(path string, value reflect.Value) reflect.Value {
	// Find the closing bracket
	closeBracketIndex := strings.Index(path, "]")
	if closeBracketIndex == -1 {
		// Invalid path, missing closing bracket
		return reflect.Value{}
	}

	indexOrKey := path[1:closeBracketIndex]
	resolvedValue := resolveIndexOrKey(indexOrKey, value)

	// If we couldn't resolve or there's no remaining path, return the result
	if !resolvedValue.IsValid() || closeBracketIndex == len(path)-1 {
		return resolvedValue
	}

	// Continue resolving with the remaining path
	remainingPath := path[closeBracketIndex+1:]
	return resolvePathAgainstValue(remainingPath, resolvedValue)
}

// resolveIndexOrKey resolves an index or key against an array, slice, or map.
// It handles numeric indices for array/slice access and various key types for map access.
//
// Parameters:
//   - indexOrKey: The index or key string to resolve
//   - value: The reflect.Value to resolve the index/key against
//
// Returns:
//   - The resolved reflect.Value
func resolveIndexOrKey(indexOrKey string, value reflect.Value) reflect.Value {
	if !value.IsValid() {
		return reflect.Value{}
	}

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		index, err := strconv.Atoi(indexOrKey)
		if err != nil || index < 0 || index >= value.Len() {
			return reflect.Value{}
		}
		return value.Index(index)
	case reflect.Map:
		return getMapValue(indexOrKey, value)
	default:
		return reflect.Value{}
	}
}

// resolveFieldOrMethod resolves a field or method name against a value.
// It first tries to resolve the name as a method, then as a field.
//
// Parameters:
//   - name: The field or method name to resolve
//   - value: The reflect.Value to resolve the name against
//
// Returns:
//   - The resolved reflect.Value
func resolveFieldOrMethod(name string, value reflect.Value) reflect.Value {
	// Handle nil or invalid values
	if !value.IsValid() || name == "" {
		return reflect.Value{}
	}

	// Try to resolve as a method first
	methodValue := resolveMethod(name, value)
	if methodValue.IsValid() {
		return methodValue
	}

	// Then try to resolve as a field
	return resolveField(name, value)
}

// resolveMethod tries to resolve a method name against a value.
// It only resolves methods that take no arguments and returns at least one value.
//
// Parameters:
//   - name: The method name to resolve
//   - value: The reflect.Value to resolve the method against
//
// Returns:
//   - The result of calling the method, or an invalid reflect.Value if the method doesn't exist
//     or requires arguments
func resolveMethod(name string, value reflect.Value) reflect.Value {
	// Check if the value has a method with the given name
	method := value.MethodByName(name)
	if !method.IsValid() {
		return reflect.Value{}
	}

	// Check if the method requires arguments
	if method.Type().NumIn() > 0 {
		return reflect.Value{}
	}

	// Call the method
	results := method.Call(nil)
	if len(results) == 0 {
		return reflect.Value{}
	}

	// Return the first result
	return results[0]
}

// resolveField tries to resolve a field name against a value.
// It handles struct fields and map keys.
//
// Parameters:
//   - name: The field name to resolve
//   - value: The reflect.Value to resolve the field against
//
// Returns:
//   - The resolved field value, or an invalid reflect.Value if the field doesn't exist
func resolveField(name string, value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Struct:
		field := value.FieldByName(name)
		if !field.IsValid() {
			return reflect.Value{}
		}
		return field
	case reflect.Map:
		return getMapValue(name, value)
	default:
		return reflect.Value{}
	}
}
