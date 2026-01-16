package empaths

// NOTE: Path Expression Character Encoding
//
// This parser is optimized for ASCII path expressions and processes paths byte-by-byte
// rather than as Unicode code points. This is a deliberate performance optimization
// since the vast majority of path expressions use ASCII-only syntax.
//
// Supported: ASCII field names, operators, brackets, quotes, and UTF-8 string literal content.
// Undefined behavior: Non-ASCII characters in field names, map keys, or reference names.
//
// If full Unicode support is needed in the future, the parser would need to be rewritten
// to use []rune instead of direct byte indexing, which would incur a performance cost.

import (
	"strings"
)

// resolveExpressions processes a path expression and evaluates it against the provided data.
// It handles multiple expression types and concatenates the results if multiple expressions are found.
//
// This implementation works directly with string bytes for ASCII paths (the common case),
// avoiding the overhead of []rune conversion. It also uses a stack-allocated approach
// for the common single-value result case.
//
// Parameters:
//   - path: The path expression as a string
//   - data: The data model to evaluate against
//   - refResolver: Function to resolve external references
//   - startIndex: The starting index in the path string
//
// Returns:
//   - The resolved value
//   - The new index after processing
func resolveExpressions(
	path string,
	data any,
	refResolver ReferenceResolver,
	startIndex int,
) (any, int) {
	if len(path) == 0 {
		return data, startIndex
	}

	index := startIndex

	// Optimization: most paths resolve to a single value.
	// Use stack-allocated first value to avoid slice allocation in the common case.
	var first any
	var hasFirst bool
	var rest []any // only allocated if we have multiple values

	for index < len(path) {
		c := path[index]
		switch c {
		case '.':
			modelResult, newIndex, err := resolveModel(path, data, index)
			if err != nil {
				return nil, index
			}
			index = newIndex
			if !hasFirst {
				first = modelResult
				hasFirst = true
			} else {
				rest = append(rest, modelResult)
			}
		case '\'':
			stringResult, newIndex := resolveStringLiteralASCII(path, index, '\'')
			index = newIndex
			if !hasFirst {
				first = stringResult
				hasFirst = true
			} else {
				rest = append(rest, stringResult)
			}
		case '"':
			stringResult, newIndex := resolveStringLiteralASCII(path, index, '"')
			index = newIndex
			if !hasFirst {
				first = stringResult
				hasFirst = true
			} else {
				rest = append(rest, stringResult)
			}
		case '!':
			negResult, newIndex := resolveNegation(path, data, index, refResolver)
			index = newIndex
			if !hasFirst {
				first = negResult
				hasFirst = true
			} else {
				rest = append(rest, negResult)
			}
		case ':':
			referenceResult, newIndex := resolveReference(path, data, index, refResolver)
			index = newIndex
			if !hasFirst {
				first = referenceResult
				hasFirst = true
			} else {
				rest = append(rest, referenceResult)
			}
		case '?':
			comparisonResult, newIndex := resolveComparison(path, data, index, refResolver)
			index = newIndex
			if !hasFirst {
				first = comparisonResult
				hasFirst = true
			} else {
				rest = append(rest, comparisonResult)
			}
		case ' ':
			index++
		default:
			index++
		}
	}

	// Return the result. If there's only one element, return it directly (no allocation).
	// If there are multiple elements, concatenate them as strings.
	if len(rest) > 0 {
		var sb strings.Builder
		sb.WriteString(toString(first))
		for _, v := range rest {
			sb.WriteString(toString(v))
		}
		return sb.String(), index
	}
	if hasFirst {
		return first, index
	}
	return data, index
}

// resolveOperand evaluates a single operand in a path expression.
// An operand can be a model reference, string literal, negation, or external reference.
//
// Parameters:
//   - path: The path expression as a string
//   - data: The data model to evaluate against
//   - refResolver: Function to resolve external references
//   - startIndex: The starting index in the path string
//
// Returns:
//   - The resolved value of the operand
//   - The new index after processing
func resolveOperand(
	path string,
	data any,
	refResolver ReferenceResolver,
	startIndex int,
) (any, int) {
	if len(path) == 0 {
		return data, startIndex
	}
	index := startIndex
	for index < len(path) {
		c := path[index]
		switch c {
		case '.':
			modelResult, newIndex, err := resolveModel(path, data, index)
			if err != nil {
				return nil, index
			}
			return modelResult, newIndex
		case '\'':
			stringResult, newIndex := resolveStringLiteralASCII(path, index, '\'')
			return stringResult, newIndex
		case '"':
			stringResult, newIndex := resolveStringLiteralASCII(path, index, '"')
			return stringResult, newIndex
		case '!':
			negResult, newIndex := resolveNegation(path, data, index, refResolver)
			return negResult, newIndex
		case ':':
			referenceResult, newIndex := resolveReference(path, data, index, refResolver)
			return referenceResult, newIndex
		case ' ':
			index++
		default:
			index++
		}
	}
	return data, index
}

// resolveStringLiteralASCII processes a string literal working directly with bytes.
// This is optimized for ASCII-only paths which is the common case.
// String literals are enclosed in single (') or double (") quotes and can include escaped characters.
//
// Parameters:
//   - path: The path expression as a string
//   - index: The current index in the path
//   - quoteChar: The quote character used (single or double quote)
//
// Returns:
//   - The string literal value
//   - The new index after processing
func resolveStringLiteralASCII(path string, index int, quoteChar byte) (string, int) {
	// skip over the opening quote
	index++
	start := index
	escaping := false
	hasEscapes := false

	// First pass: find the end and check for escapes
	for index < len(path) {
		c := path[index]
		if escaping {
			escaping = false
			hasEscapes = true
			index++
			continue
		}
		if c == quoteChar {
			break
		}
		if c == '\\' {
			escaping = true
			index++
			continue
		}
		index++
	}

	// If no escapes, we can return a substring directly (no allocation for the content)
	if !hasEscapes {
		return path[start:index], index + 1
	}

	// With escapes, we need to build the string
	var sb strings.Builder
	sb.Grow(index - start)
	escaping = false
	for i := start; i < index; i++ {
		c := path[i]
		if escaping {
			escaping = false
			sb.WriteByte(c)
			continue
		}
		if c == '\\' {
			escaping = true
			continue
		}
		sb.WriteByte(c)
	}
	return sb.String(), index + 1
}

// readUntilTerminatorASCII reads characters from a path until a terminator character is found.
// This works directly with string bytes for efficiency.
// Terminator characters include space, exclamation mark, and equals sign.
//
// Parameters:
//   - path: The path expression as a string
//   - index: The starting index in the path
//
// Returns:
//   - The segment read from the path as a string
//   - The new index after processing
func readUntilTerminatorASCII(path string, index int) (string, int) {
	start := index
	for index < len(path) {
		c := path[index]
		if c == ' ' || c == '!' || c == '=' {
			break
		}
		index++
	}
	return path[start:index], index
}
