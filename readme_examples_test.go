package empaths

import (
	"testing"
)

// =============================================================================
// Test Data Structures matching README examples
// =============================================================================

// QuickStartUser matches the User type in Quick Start section
type QuickStartUser struct {
	Name    string
	Age     int
	Address struct {
		City    string
		Country string
	}
	Tags []string
}

// ConcatUser has FirstName/LastName for concatenation examples
type ConcatUser struct {
	FirstName string
	LastName  string
	Name      string
	Age       int
	Address   struct {
		City string
	}
}

// MethodUser demonstrates method calls
type MethodUser struct {
	FirstName string
	LastName  string
}

func (u MethodUser) FullName() string {
	return u.FirstName + " " + u.LastName
}

// NegationData for negation examples
type NegationData struct {
	IsActive bool
	Status   string
}

// ComparisonData for comparison examples
type ComparisonData struct {
	Age          int
	Status       string
	Name         string
	ExpectedName string
}

// =============================================================================
// Quick Start Examples
// =============================================================================

func TestReadmeQuickStart(t *testing.T) {
	user := QuickStartUser{
		Name: "Alice",
		Age:  30,
		Tags: []string{"developer", "gopher"},
	}
	user.Address.City = "New York"
	user.Address.Country = "USA"

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"Access Name field", ".Name", "Alice"},
		{"Access nested Address.City", ".Address.City", "New York"},
		{"Access first tag", ".Tags[0]", "developer"},
		{"Concatenate greeting", "'Hello, ' .Name '! You are ' .Age ' years old.'", "Hello, Alice! You are 30 years old."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, user, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v (%T), want %v (%T)", tt.path, result, result, tt.expected, tt.expected)
			}
		})
	}
}

// =============================================================================
// Expression Concatenation Examples
// =============================================================================

func TestReadmeExpressionConcatenation(t *testing.T) {
	user := ConcatUser{
		FirstName: "John",
		LastName:  "Doe",
		Name:      "Alice",
		Age:       30,
	}
	user.Address.City = "New York"

	// Test that single expression preserves type
	t.Run("Single expression preserves int type", func(t *testing.T) {
		result := Resolve(".Age", user, nil)
		if _, ok := result.(int); !ok {
			t.Errorf("Expected int type, got %T", result)
		}
		if result != 30 {
			t.Errorf("Expected 30, got %v", result)
		}
	})

	// Test concatenation examples
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"Name is Age years old", ".Name ' is ' .Age ' years old'", "Alice is 30 years old"},
		{"User: Name", "'User: ' .Name", "User: Alice"},
		{"FirstName LastName", ".FirstName ' ' .LastName", "John Doe"},
		{"Welcome greeting", "'Hello, ' .Name '! Welcome to ' .Address.City '.'", "Hello, Alice! Welcome to New York."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, user, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}

	// Test with external reference
	t.Run("Concatenation with external reference", func(t *testing.T) {
		resolver := func(name string, data any) any {
			if name == "greeting" {
				return "Hello"
			}
			return nil
		}
		result := Resolve(":greeting ', ' .Name '!'", user, resolver)
		expected := "Hello, Alice!"
		if result != expected {
			t.Errorf("Resolve with resolver = %v, want %v", result, expected)
		}
	})
}

// =============================================================================
// Field Access Examples
// =============================================================================

func TestReadmeFieldAccess(t *testing.T) {
	type Address struct {
		City string
	}
	type User struct {
		Address Address
	}
	type Data struct {
		User User
		Key  string
	}

	data := Data{
		User: User{Address: Address{City: "Boston"}},
		Key:  "value",
	}

	mapData := map[string]string{"key": "mapvalue"}

	tests := []struct {
		name     string
		path     string
		data     any
		expected any
	}{
		{"Simple field .Name style", ".Key", data, "value"},
		{"Nested fields .User.Address.City", ".User.Address.City", data, "Boston"},
		{"Map with string key .Data.key style", ".key", mapData, "mapvalue"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, tt.data, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Array/Slice Indexing Examples
// =============================================================================

func TestReadmeArraySliceIndexing(t *testing.T) {
	type User struct {
		Name string
	}

	items := []string{"first", "second", "third"}
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}}
	users := []User{{Name: "UserOne"}, {Name: "UserTwo"}}

	tests := []struct {
		name     string
		path     string
		data     any
		expected any
	}{
		{"First element .Items[0]", ".[0]", items, "first"},
		{"Third element .Items[2]", ".[2]", items, "third"},
		{"Nested array .Matrix[0][1]", ".[0][1]", matrix, 2},
		{"Field of array element .Users[0].Name", ".[0].Name", users, "UserOne"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, tt.data, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Map Access Examples
// =============================================================================

func TestReadmeMapAccess(t *testing.T) {
	config := map[string]any{
		"timeout": 30,
	}
	scores := map[int]string{
		42: "answer",
	}
	flags := map[bool]string{
		true: "enabled",
	}

	tests := []struct {
		name     string
		path     string
		data     any
		expected any
	}{
		{"Dot notation .Config.timeout", ".timeout", config, 30},
		{"Bracket notation .Config[timeout]", ".[timeout]", config, 30},
		{"Integer key .Scores[42]", ".[42]", scores, "answer"},
		{"Boolean key .Flags[true]", ".[true]", flags, "enabled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, tt.data, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// String Literals Examples
// =============================================================================

func TestReadmeStringLiterals(t *testing.T) {
	user := struct{ Name string }{Name: "Alice"}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"Single quotes", "'Hello'", "Hello"},
		{"Double quotes", "\"World\"", "World"},
		{"Concatenation with field", "'Hello, ' .Name '!'", "Hello, Alice!"},
		{"Escaped single quote", "'It\\'s working'", "It's working"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, user, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Comparisons Examples
// =============================================================================

func TestReadmeComparisons(t *testing.T) {
	data := ComparisonData{
		Age:          30,
		Status:       "active",
		Name:         "Alice",
		ExpectedName: "Alice",
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Age equals 30", "?.Age=='30'", true},
		{"Age equals 25 (false)", "?.Age=='25'", false},
		{"Status not inactive", "?.Status!='inactive'", true},
		{"Compare two fields equal", "?.Name==.ExpectedName", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, data, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Negation Examples
// =============================================================================

func TestReadmeNegation(t *testing.T) {
	data := NegationData{
		IsActive: true,
		Status:   "active",
	}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"Negate boolean field", "!.IsActive", false},
		{"Negate literal true", "!'true'", false},
		{"Negate literal false", "!'false'", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, data, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// External References Examples
// =============================================================================

func TestReadmeExternalReferences(t *testing.T) {
	user := struct{ Name string }{Name: "Alice"}

	resolver := func(name string, data any) any {
		switch name {
		case "greeting":
			return "Hello"
		case "config":
			return map[string]string{"key": "value"}
		}
		return nil
	}

	t.Run("External reference with concatenation", func(t *testing.T) {
		result := Resolve(":greeting ', ' .Name", user, resolver)
		expected := "Hello, Alice"
		if result != expected {
			t.Errorf("Resolve = %v, want %v", result, expected)
		}
	})
}

// =============================================================================
// Method Calls Examples
// =============================================================================

func TestReadmeMethodCalls(t *testing.T) {
	user := MethodUser{
		FirstName: "John",
		LastName:  "Doe",
	}

	result := Resolve(".FullName", user, nil)
	expected := "John Doe"
	if result != expected {
		t.Errorf("Resolve(.FullName) = %v, want %v", result, expected)
	}
}

// =============================================================================
// Working with Different Types - Structs
// =============================================================================

func TestReadmeStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "Bob", Age: 25}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"Struct string field", ".Name", "Bob"},
		{"Struct int field", ".Age", 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, person, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Working with Different Types - Maps
// =============================================================================

func TestReadmeMaps(t *testing.T) {
	data := map[string]any{
		"name": "Charlie",
		"scores": map[string]int{
			"math":    95,
			"science": 88,
		},
	}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"Map string value", ".name", "Charlie"},
		{"Nested map dot notation", ".scores.math", 95},
		{"Nested map bracket notation", ".scores[science]", 88},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, data, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Working with Different Types - Slices and Arrays
// =============================================================================

func TestReadmeSlicesAndArrays(t *testing.T) {
	items := []string{"apple", "banana", "cherry"}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"First element", ".[0]", "apple"},
		{"Third element", ".[2]", "cherry"},
		{"Out of bounds returns nil", ".[99]", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, items, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Working with Different Types - Pointers
// =============================================================================

func TestReadmePointers(t *testing.T) {
	type User struct {
		Name string
	}

	user := &User{Name: "Diana"}

	result := Resolve(".Name", user, nil)
	expected := "Diana"
	if result != expected {
		t.Errorf("Resolve(.Name) on pointer = %v, want %v", result, expected)
	}
}

// =============================================================================
// Nil Safety Examples
// =============================================================================

func TestReadmeNilSafety(t *testing.T) {
	type User struct {
		Name  string
		Items []string
	}

	var nilUser *User = nil
	data := User{Name: "Test", Items: []string{"a", "b"}}

	tests := []struct {
		name     string
		path     string
		data     any
		expected any
	}{
		{"Nil pointer access", ".Name", nilUser, nil},
		{"Non-existent field", ".NonExistent", data, nil},
		{"Out of bounds index", ".Items[999]", data, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, tt.data, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Error Handling Examples
// =============================================================================

func TestReadmeErrorHandling(t *testing.T) {
	type Data struct {
		Field string
		Items []string
	}
	data := Data{Field: "value", Items: []string{"a"}}

	tests := []struct {
		name     string
		path     string
		data     any
		expected any
	}{
		{"Non-existent field returns nil", ".NonExistent", data, nil},
		{"Out of bounds returns nil", ".Items[999]", data, nil},
		{"Field on nil returns nil", ".Field", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, tt.data, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Additional edge cases from README that should work
// =============================================================================

func TestReadmeAdditionalExamples(t *testing.T) {
	// Test escaped double quotes in string literal
	t.Run("Escaped double quotes", func(t *testing.T) {
		result := Resolve("\"Say \\\"Hi\\\"\"", nil, nil)
		expected := "Say \"Hi\""
		if result != expected {
			t.Errorf("Escaped double quotes = %v, want %v", result, expected)
		}
	})

	// Test comparison with two field references
	t.Run("Compare two different fields", func(t *testing.T) {
		data := struct {
			Name         string
			ExpectedName string
		}{Name: "Test", ExpectedName: "Different"}
		result := Resolve("?.Name==.ExpectedName", data, nil)
		if result != false {
			t.Errorf("Comparing different values should be false, got %v", result)
		}
	})
}
