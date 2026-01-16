package empaths

import (
	"errors"
	"reflect"
	"strings"
)

// resolveComparison evaluates a comparison expression in a path.
// Comparison expressions start with '?' and compare two operands with either '==' or '!=' operators.
//
// Parameters:
//   - path: The path expression as a string
//   - data: The data model to evaluate against
//   - index: The current index in the path
//   - refResolver: Function to resolve external references
//
// Returns:
//   - The boolean result of the comparison
//   - The new index after processing
func resolveComparison(path string, data any, index int, refResolver ReferenceResolver) (bool, int) {
	// skip over the ? prefix
	index++
	leftOperand, index := resolveOperand(path, data, refResolver, index)
	equalsOperator, index, err := parseOperator(path, index)
	if err != nil {
		// Invalid operator - return false as comparison result
		return false, index
	}

	leftStr := toString(leftOperand)

	rightOperand, index := resolveOperand(path, data, refResolver, index)
	rightStr := toString(rightOperand)

	if equalsOperator {
		return leftStr == rightStr, index
	}
	return leftStr != rightStr, index
}

// parseOperator determines the comparison operator (== or !=) in a comparison expression.
// Returns true for equals (==) and false for not equals (!=).
//
// Parameters:
//   - path: The path expression as a string
//   - index: The current index in the path
//
// Returns:
//   - true for equals operator (==), false for not equals operator (!=)
//   - The new index after processing
//   - Error if an invalid operator is found
func parseOperator(path string, index int) (bool, int, error) {
	if index >= len(path)-1 {
		return false, index + 1, errors.New("no operator found for comparison")
	}
	if path[index] == '!' && path[index+1] == '=' {
		return false, index + 2, nil
	}
	if path[index] == '=' && path[index+1] == '=' {
		return true, index + 2, nil
	}
	return false, index + 1, errors.New("invalid operator")
}

// resolveReference processes an external reference.
// External references start with ':' followed by the reference name.
//
// Parameters:
//   - path: The path expression as a string
//   - data: The data model to evaluate against
//   - index: The current index in the path
//   - refResolver: Function to resolve external references
//
// Returns:
//   - The resolved value from the external reference
//   - The new index after processing
func resolveReference(path string, data any, index int, refResolver ReferenceResolver) (any, int) {
	// Skip over the ':' prefix
	index++
	referenceName, index := readUntilTerminatorASCII(path, index)

	if refResolver == nil {
		return nil, index
	}
	referenceValue := refResolver(referenceName, data)
	return referenceValue, index
}

// resolveNegation processes a negation expression in a path.
// Negation expressions start with '!' and negate a boolean value or convert a value to its boolean opposite.
//
// Parameters:
//   - path: The path expression as a string
//   - data: The data model to evaluate against
//   - index: The current index in the path
//   - refResolver: Function to resolve external references
//
// Returns:
//   - The negated boolean value
//   - The new index after processing
func resolveNegation(path string, data any, index int, refResolver ReferenceResolver) (any, int) {
	// skip over the ! prefix
	index++

	value, newIndex := resolveOperand(path, data, refResolver, index)
	// If it's already a boolean, just negate it
	if boolValue, ok := value.(bool); ok {
		return !boolValue, newIndex
	}

	// Try to convert to boolean
	strValue := toString(value)
	lowerStr := strings.ToLower(strValue)

	if lowerStr == "true" {
		return false, newIndex
	}
	if lowerStr == "false" {
		return true, newIndex
	}
	return false, newIndex
}

// resolveModel resolves a model reference in a path expression.
// Model references start with '.' followed by a path to a property or method in the data model.
//
// Parameters:
//   - path: The path expression as a string
//   - data: The data model to evaluate against
//   - index: The current index in the path (should point to the '.' character)
//
// Returns:
//   - The resolved value from the data model
//   - The new index after processing
//   - Error if the path cannot be resolved
func resolveModel(path string, data any, index int) (any, int, error) {
	// skip over the '.'
	index++
	modelPath, index := readUntilTerminatorASCII(path, index)
	if data == nil {
		return nil, index, nil
	}
	value := reflect.ValueOf(data)
	result := resolvePathAgainstValue(modelPath, value)

	return extractValue(result), index, nil
}
