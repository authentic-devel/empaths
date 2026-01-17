# empaths

[![Go Reference](https://pkg.go.dev/badge/github.com/authentic-devel/empaths.svg)](https://pkg.go.dev/github.com/authentic-devel/empaths)
[![CI](https://github.com/authentic-devel/empaths/actions/workflows/ci.yml/badge.svg)](https://github.com/authentic-devel/empaths/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/authentic-devel/empaths)](https://goreportcard.com/report/github.com/authentic-devel/empaths)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**empaths** (em-paths) is a lightweight, performant Go library for accessing nested values in data structures using 
string path expressions. The "em" stands for "Model" — think of it as "model paths."

## Features

- **Simple path syntax** — Access deeply nested fields with intuitive dot notation
- **Expression concatenation** — Combine multiple values and literals into formatted strings
- **Universal data access** — Works with structs, maps, slices, arrays, and pointers
- **Zero dependencies** — Pure Go standard library, no external dependencies
- **Performant** — Optimized for minimal allocations and fast execution
- **Type-safe** — Graceful handling of nil values and type mismatches
- **Extensible** — Support for custom reference resolvers

## Installation

```bash
go get github.com/authentic-devel/empaths
```

Requires Go 1.21 or later.

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/authentic-devel/empaths"
)

type User struct {
    Name    string
    Age     int
    Address struct {
        City    string
        Country string
    }
    Tags []string
}

func main() {
    user := User{
        Name: "Alice",
        Age:  30,
        Address: struct {
            City    string
            Country string
        }{City: "New York", Country: "USA"},
        Tags: []string{"developer", "gopher"},
    }

    // Access nested fields
    name := empaths.Resolve(".Name", user, nil)
    fmt.Println(name) // "Alice"

    city := empaths.Resolve(".Address.City", user, nil)
    fmt.Println(city) // "New York"

    // Access slice elements
    firstTag := empaths.Resolve(".Tags[0]", user, nil)
    fmt.Println(firstTag) // "developer"

    // Concatenate multiple expressions into a string
    greeting := empaths.Resolve("'Hello, ' .Name '! You are ' .Age ' years old.'", user, nil)
    fmt.Println(greeting) // "Hello, Alice! You are 30 years old."
}
```

## Path Syntax

### Expression Concatenation

A path can contain **multiple expressions** separated by spaces. When multiple expressions are present, each is evaluated and their **string representations are concatenated** into a single result:

```go
// Single expression → returns the value directly (preserves type)
empaths.Resolve(".Age", user, nil)  // returns int: 30

// Multiple expressions → concatenated as strings
empaths.Resolve(".Name ' is ' .Age ' years old'", user, nil)
// returns string: "Alice is 30 years old"

empaths.Resolve("'User: ' .Name", user, nil)
// returns string: "User: Alice"

empaths.Resolve(".FirstName ' ' .LastName", user, nil)
// returns string: "John Doe"
```

This is powerful for building dynamic strings from multiple data sources:

```go
// Build a greeting
empaths.Resolve("'Hello, ' .Name '! Welcome to ' .Address.City '.'", user, nil)
// → "Hello, Alice! Welcome to New York."

// Combine with external references
empaths.Resolve(":greeting ', ' .Name '!'", user, resolver)
// → "Hello, Alice!"
```

> **Note:** When a path contains only a single expression, the original type is preserved. When multiple expressions are present, the result is always a string.

### Field Access

Use dot notation to access struct fields or map keys:

```go
".Name"              // Access field "Name"
".User.Address.City" // Access nested fields
".Data.key"          // Access map with string key
```

### Array/Slice Indexing

Use bracket notation with zero-based indices:

```go
".Items[0]"          // First element
".Items[2]"          // Third element
".Matrix[0][1]"      // Nested array access
".Users[0].Name"     // Field of array element
```

### Map Access

Maps support both dot and bracket notation:

```go
".Config.timeout"    // Dot notation (string keys)
".Config[timeout]"   // Bracket notation
".Scores[42]"        // Integer key
".Flags[true]"       // Boolean key
```

### String Literals

Embed literal strings in expressions:

```go
"'Hello'"                    // Single quotes
"\"World\""                  // Double quotes
"'Hello, ' .Name '!'"        // Concatenation → "Hello, Alice!"
"'It\\'s working'"           // Escaped quotes
```

### Comparisons

Compare values using `==` or `!=`:

```go
"?.Age=='30'"                // Equals comparison → true/false
"?.Status!='inactive'"       // Not equals comparison
"?.Name==.ExpectedName"      // Compare two fields
```

### Negation

Negate boolean values with `!`:

```go
"!.IsActive"                 // Negate boolean field
"!'true'"                    // Negate literal → false
```

### External References

Resolve custom references with a resolver function:

```go
resolver := func(name string, data any) any {
    switch name {
    case "greeting":
        return "Hello"
    case "config":
        return someConfig
    }
    return nil
}

result := empaths.Resolve(":greeting ', ' .Name", user, resolver)
// → "Hello, Alice"
```

## Method Calls

Zero-argument methods can be called as part of a path:

```go
type User struct {
    FirstName string
    LastName  string
}

func (u User) FullName() string {
    return u.FirstName + " " + u.LastName
}

// Usage
result := empaths.Resolve(".FullName", user, nil)
// → "John Doe"
```

## Working with Different Types

### Structs

```go
type Person struct {
    Name string
    Age  int
}

person := Person{Name: "Bob", Age: 25}
empaths.Resolve(".Name", person, nil) // "Bob"
empaths.Resolve(".Age", person, nil)  // 25
```

### Maps

```go
data := map[string]any{
    "name": "Charlie",
    "scores": map[string]int{
        "math":    95,
        "science": 88,
    },
}

empaths.Resolve(".name", data, nil)          // "Charlie"
empaths.Resolve(".scores.math", data, nil)   // 95
empaths.Resolve(".scores[science]", data, nil) // 88
```

### Slices and Arrays

```go
items := []string{"apple", "banana", "cherry"}

empaths.Resolve(".[0]", items, nil) // "apple"
empaths.Resolve(".[2]", items, nil) // "cherry"
empaths.Resolve(".[99]", items, nil) // nil (out of bounds)
```

### Pointers

Pointers are automatically dereferenced:

```go
user := &User{Name: "Diana"}
empaths.Resolve(".Name", user, nil) // "Diana"
```

### Nil Safety

The library handles nil values gracefully:

```go
var user *User = nil
empaths.Resolve(".Name", user, nil)      // nil (no panic)
empaths.Resolve(".NonExistent", data, nil) // nil
empaths.Resolve(".Items[999]", data, nil)  // nil
```

## API Reference

### Resolve

```go
func Resolve(path string, data any, refResolver ReferenceResolver) any
```

Evaluates a path expression against a data model and returns the resolved value.

**Parameters:**
- `path` — The path expression to evaluate
- `data` — The data model to evaluate against
- `refResolver` — Optional function to resolve external references (can be nil)

**Returns:** The resolved value, or nil if the path cannot be resolved.

### ReferenceResolver

```go
type ReferenceResolver func(name string, data any) any
```

Function type for resolving external references (paths starting with `:`).

## Error Handling

empaths uses **graceful failure** — invalid paths return `nil` rather than panicking or returning errors. 
This design simplifies usage in templates and other contexts where nil is an acceptable fallback.


```go
// All of these return nil without panicking:
empaths.Resolve(".NonExistent", data, nil)
empaths.Resolve(".Items[999]", data, nil)
empaths.Resolve("invalid path", data, nil)
empaths.Resolve(".Field", nil, nil)
```

## Character Encoding

Path expressions should use **ASCII characters only** for path syntax elements (field names, operators, brackets, quotes). The parser is optimized for ASCII and processes paths byte-by-byte rather than as Unicode code points.

**What works:**
- ASCII field/method names: `.User.Name`, `.GetValue`
- ASCII map keys: `.Config[timeout]`, `.Data["key"]`
- UTF-8 content in string literals: `'こんにちは'`, `"日本語"`
- UTF-8 values in your data structures (these are not affected)

**What to avoid:**
- Non-ASCII field names in paths: `.用户.名前` (undefined behavior)
- Non-ASCII map keys in bracket notation: `.Data[キー]` (undefined behavior)
- Non-ASCII reference names: `:配置` (undefined behavior)

```go
// Safe: ASCII path syntax, UTF-8 data values
data := map[string]string{"greeting": "こんにちは"}
empaths.Resolve(".greeting", data, nil)  // Returns "こんにちは"

// Safe: UTF-8 in string literals
empaths.Resolve("'Hello ' .Name '!'", user, nil)

// Undefined: Non-ASCII in path syntax
type User struct {
    名前 string
}
empaths.Resolve(".名前", user, nil)  // May work but not guaranteed
```

## Use Cases

- **Template engines** — Dynamic value resolution in templates
- **Configuration** — Access nested config values by path
- **API responses** — Extract values from JSON-decoded maps
- **Testing** — Assert on deeply nested values
- **Data transformation** — Map values between structures


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## AI Disclaimer

This project was written manually by me, but AI was used to help with improvements and documentation.  

