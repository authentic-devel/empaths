// Package empaths provides xpath-like path resolution for Go data structures.
//
// Empaths (pronounced "em-paths" - the M stands for Model) allows reflective
// access to nested values in structs, maps, slices, and arrays using a
// simple path syntax.
package empaths

// ReferenceResolver is a function type that resolves external references.
// It takes a reference name and a data context, and returns the resolved value.
// This can be used to resolve references to templates, configuration values,
// or any other external data sources.
type ReferenceResolver func(name string, data any) any

// Resolve evaluates a path expression against a data model and returns the resolved value.
//
// A path can consist of multiple segments and supports various expression types:
//   - Model references: Starts with '.' followed by field/property names (e.g., ".User.Name")
//   - String literals: Enclosed in single or double quotes (e.g., "'Hello'" or "\"World\"")
//   - Negation: Starts with '!' to negate a boolean value (e.g., "!.IsActive")
//   - External references: Starts with ':' followed by reference name (e.g., ":config")
//   - Comparisons: Starts with '?' followed by operands and operator (e.g., "?.Age=='18'")
//
// Character encoding: Path syntax elements (field names, map keys, reference names) should
// use ASCII characters only. UTF-8 content within string literals is supported, but non-ASCII
// characters in path syntax have undefined behavior.
//
// Path segments can be combined to form complex expressions, and can include:
//   - Nested properties: ".User.Address.City"
//   - Array/slice indexing: ".Users[0].Name"
//   - Map access: ".Data[\"key\"]" or ".Data.key"
//   - Method calls: ".User.GetFullName"
//
// Examples:
//   - ".User.Name" - Accesses the Name property of the User object
//   - "'Hello ' .User.Name" - Concatenates the string "Hello " with User.Name
//   - "?.IsAdmin=='true'" - Compares if IsAdmin equals true
//   - "!.IsBlocked" - Negates the IsBlocked boolean value
//   - ":config" - References an external value named "config"
//
// Parameters:
//   - path: The path expression to evaluate
//   - data: The data model to evaluate the path against
//   - referenceResolver: Optional function to resolve external references (prefixed with ':')
//
// Returns:
//
//	The resolved value from the data model based on the path expression
func Resolve(path string, data any, refResolver ReferenceResolver) any {
	if path == "" {
		return data
	}
	result, _ := resolveExpressions(path, data, refResolver, 0)
	return result
}

// ResolveModel resolves a model reference in a path expression.
// Model references start with '.' followed by a path to a property or method in the data model.
// This function can be used directly to resolve a model path against a data object.
//
// A model path can include:
//   - Simple property access: ".PropertyName"
//   - Nested property access: ".User.Address.City"
//   - Array/slice indexing: ".Users[0]"
//   - Map access: ".Data[\"key\"]" or ".Data.key"
//   - Method calls: ".GetFullName"
//
// Parameters:
//   - path: The path expression string
//   - data: The data model to evaluate against
//   - index: The current index in the path (should point to the '.' character)
//
// Returns:
//   - The resolved value from the data model
//   - The new index after processing
//   - Error if the path cannot be resolved
func ResolveModel(path string, data any, index int) (any, int, error) {
	return resolveModel(path, data, index)
}
