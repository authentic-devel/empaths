// Package empaths provides xpath-like path resolution for Go data structures.
//
// Empaths (pronounced "em-paths" - the M stands for Model) allows reflective
// access to nested values in structs, maps, slices, and arrays using a
// simple path syntax. It is designed to be fast, memory-efficient, and easy to use.
//
// # Path Syntax
//
// Path expressions can contain multiple segments that are evaluated in sequence:
//
// Model References (start with '.'):
//
//	.Name              - Access a struct field or map key named "Name"
//	.User.Address.City - Access nested fields
//	.Users[0]          - Access array/slice element by index (zero-based)
//	.Data["key"]       - Access map element by key
//	.GetValue          - Call a zero-argument method
//
// String Literals (enclosed in quotes):
//
//	'Hello'            - Single-quoted string
//	"World"            - Double-quoted string
//	'It\'s'            - Escaped quotes within strings
//
// Negation (starts with '!'):
//
//	!.IsActive         - Negate a boolean value
//	!'true'            - Negate a string "true" -> false
//
// Comparisons (start with '?'):
//
//	?.Age=='18'        - Compare if Age equals 18
//	?.Status!='active' - Compare if Status is not "active"
//
// External References (start with ':'):
//
//	:config            - Resolve using the provided ReferenceResolver
//
// Multiple segments can be combined:
//
//	'Hello, ' .User.Name '!'  - Concatenates to "Hello, John!"
//
// # Array and Slice Access
//
// Arrays and slices are accessed using zero-based integer indices:
//
//	.Items[0]          - First element
//	.Items[1]          - Second element
//	.Matrix[0][1]      - Nested array access
//
// Out-of-bounds access returns nil rather than panicking.
// Negative indices are not supported.
//
// # Map Access
//
// Maps can be accessed either with bracket notation or dot notation:
//
//	.Data["key"]       - Bracket notation (works for any key type)
//	.Data.key          - Dot notation (string keys only)
//
// The library supports maps with various key types:
//   - string, int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - bool, float32, float64
//
// # Method Calls
//
// Zero-argument methods can be called as part of a path:
//
//	.GetFullName       - Calls GetFullName() method
//	.User.String       - Calls String() on User
//
// Methods must:
//   - Take no arguments
//   - Return at least one value (first value is used)
//
// # Error Handling
//
// The library uses graceful failure - invalid paths return nil rather than
// panicking or returning errors. This design choice simplifies usage in
// templates and other contexts where nil is an acceptable fallback.
//
// # Example Usage
//
//	type User struct {
//	    Name    string
//	    Age     int
//	    Address struct {
//	        City string
//	    }
//	}
//
//	user := User{
//	    Name: "Alice",
//	    Age:  30,
//	    Address: struct{ City string }{City: "NYC"},
//	}
//
//	// Simple field access
//	name := empaths.Resolve(".Name", user, nil)  // "Alice"
//
//	// Nested field access
//	city := empaths.Resolve(".Address.City", user, nil)  // "NYC"
//
//	// String concatenation
//	greeting := empaths.Resolve("'Hello, ' .Name '!'", user, nil)  // "Hello, Alice!"
//
//	// Comparison
//	isAdult := empaths.Resolve("?.Age=='30'", user, nil)  // true
//
// # Thread Safety
//
// All functions in this package are safe for concurrent use.
// The library does not maintain any global state.
package empaths
