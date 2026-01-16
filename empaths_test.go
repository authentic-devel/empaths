package empaths

import (
	"testing"
)

// Test data structures
type Address struct {
	Street string
	City   string
	Zip    int
}

type Person struct {
	Name    string
	Age     int
	Active  bool
	Address Address
	Tags    []string
	Scores  map[string]int
}

func (p Person) GetFullName() string {
	return "Mr/Ms " + p.Name
}

func (p Person) IsAdult() bool {
	return p.Age >= 18
}

// createTestPerson returns a Person for testing
func createTestPerson() Person {
	return Person{
		Name:   "Alice",
		Age:    30,
		Active: true,
		Address: Address{
			Street: "123 Main St",
			City:   "NYC",
			Zip:    10001,
		},
		Tags:   []string{"developer", "gopher", "tester"},
		Scores: map[string]int{"math": 95, "science": 88},
	}
}

func TestResolve_SimpleField(t *testing.T) {
	person := createTestPerson()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"string field", ".Name", "Alice"},
		{"int field", ".Age", 30},
		{"bool field", ".Active", true},
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

func TestResolve_NestedField(t *testing.T) {
	person := createTestPerson()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"nested string", ".Address.City", "NYC"},
		{"nested string 2", ".Address.Street", "123 Main St"},
		{"nested int", ".Address.Zip", 10001},
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

func TestResolve_SliceAccess(t *testing.T) {
	person := createTestPerson()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"first element", ".Tags[0]", "developer"},
		{"second element", ".Tags[1]", "gopher"},
		{"third element", ".Tags[2]", "tester"},
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

func TestResolve_SliceOutOfBounds(t *testing.T) {
	person := createTestPerson()

	result := Resolve(".Tags[99]", person, nil)
	if result != nil {
		t.Errorf("Resolve with out of bounds index should return nil, got %v", result)
	}
}

func TestResolve_MapAccess(t *testing.T) {
	person := createTestPerson()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"bracket notation", ".Scores[math]", 95},
		{"dot notation", ".Scores.science", 88},
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

func TestResolve_MapKeyNotFound(t *testing.T) {
	person := createTestPerson()

	result := Resolve(".Scores[nonexistent]", person, nil)
	if result != nil {
		t.Errorf("Resolve with nonexistent map key should return nil, got %v", result)
	}
}

func TestResolve_MethodCall(t *testing.T) {
	person := createTestPerson()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"string method", ".GetFullName", "Mr/Ms Alice"},
		{"bool method", ".IsAdult", true},
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

func TestResolve_StringLiteral(t *testing.T) {
	person := createTestPerson()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"single quotes", "'Hello'", "Hello"},
		{"double quotes", "\"World\"", "World"},
		{"escaped single", "'It\\'s'", "It's"},
		{"escaped double", "\"Say \\\"Hi\\\"\"", "Say \"Hi\""},
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

func TestResolve_Concatenation(t *testing.T) {
	person := createTestPerson()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"string and field", "'Hello, ' .Name", "Hello, Alice"},
		{"field and string", ".Name ' is here'", "Alice is here"},
		{"multiple parts", "'Name: ' .Name ', Age: ' .Age", "Name: Alice, Age: 30"},
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

func TestResolve_Negation(t *testing.T) {
	person := createTestPerson()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"negate true", "!.Active", false},
		{"negate string true", "!'true'", false},
		{"negate string false", "!'false'", true},
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

func TestResolve_Comparison(t *testing.T) {
	person := createTestPerson()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"int equals string literal true", "?.Age=='30'", true},
		{"int equals string literal false", "?.Age=='25'", false},
		{"int not equals string literal true", "?.Age!='25'", true},
		{"int not equals string literal false", "?.Age!='30'", false},
		{"string comparison", "?.Name=='Alice'", true},
		{"bool comparison", "?.Active=='true'", true},
		{"nested int comparison", "?.Address.Zip=='10001'", true},
		// Comparisons convert both sides to strings, so int 30 becomes "30"
		// Note: bare integer literals (e.g., ?.Age==30) are NOT supported - use string literals
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

func TestResolve_ComparisonFieldToField(t *testing.T) {
	data := map[string]any{
		"value":    30,
		"expected": 30,
		"other":    25,
	}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"field equals field true", "?.value==.expected", true},
		{"field equals field false", "?.value==.other", false},
		{"field not equals field true", "?.value!=.other", true},
		{"field not equals field false", "?.value!=.expected", false},
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

func TestResolve_ExternalReference(t *testing.T) {
	person := createTestPerson()

	resolver := func(name string, data any) any {
		switch name {
		case "greeting":
			return "Hello"
		case "suffix":
			return "!"
		default:
			return nil
		}
	}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"simple reference", ":greeting", "Hello"},
		{"reference with field", ":greeting ', ' .Name", "Hello, Alice"},
		{"multiple references", ":greeting .Name :suffix", "HelloAlice!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, person, resolver)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestResolve_NilResolver(t *testing.T) {
	person := createTestPerson()

	result := Resolve(":something", person, nil)
	if result != nil {
		t.Errorf("Resolve with nil resolver should return nil, got %v", result)
	}
}

func TestResolve_EmptyPath(t *testing.T) {
	data := "test data"

	result := Resolve("", data, nil)
	if result != data {
		t.Errorf("Resolve with empty path should return data, got %v", result)
	}
}

func TestResolve_NilData(t *testing.T) {
	result := Resolve(".Name", nil, nil)
	if result != nil {
		t.Errorf("Resolve with nil data should return nil, got %v", result)
	}
}

func TestResolve_InvalidField(t *testing.T) {
	person := createTestPerson()

	result := Resolve(".NonExistent", person, nil)
	if result != nil {
		t.Errorf("Resolve with invalid field should return nil, got %v", result)
	}
}

func TestResolve_Pointer(t *testing.T) {
	person := createTestPerson()
	ptr := &person

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"pointer to struct field", ".Name", "Alice"},
		{"pointer to nested field", ".Address.City", "NYC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, ptr, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestResolve_NestedPointer(t *testing.T) {
	type Inner struct {
		Value string
	}
	type Outer struct {
		Inner *Inner
	}

	outer := Outer{Inner: &Inner{Value: "nested"}}

	result := Resolve(".Inner.Value", outer, nil)
	if result != "nested" {
		t.Errorf("Resolve with nested pointer = %v, want %v", result, "nested")
	}
}

func TestResolve_NilPointer(t *testing.T) {
	type Outer struct {
		Inner *Address
	}

	outer := Outer{Inner: nil}

	result := Resolve(".Inner.City", outer, nil)
	if result != nil {
		t.Errorf("Resolve with nil pointer should return nil, got %v", result)
	}
}

func TestResolve_SliceOfPointers(t *testing.T) {
	type Item struct {
		Name  string
		Value int
	}

	items := []*Item{
		{Name: "first", Value: 1},
		{Name: "second", Value: 2},
		{Name: "third", Value: 3},
	}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"first pointer element name", ".[0].Name", "first"},
		{"second pointer element value", ".[1].Value", 2},
		{"third pointer element name", ".[2].Name", "third"},
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

func TestResolve_SliceOfPointersWithNil(t *testing.T) {
	type Item struct {
		Name string
	}

	items := []*Item{
		{Name: "first"},
		nil,
		{Name: "third"},
	}

	// Accessing a field on a nil element should return nil
	result := Resolve(".[1].Name", items, nil)
	if result != nil {
		t.Errorf("Resolve on nil slice element should return nil, got %v", result)
	}

	// Non-nil elements should still work
	result = Resolve(".[2].Name", items, nil)
	if result != "third" {
		t.Errorf("Resolve on non-nil element = %v, want %v", result, "third")
	}
}

func TestResolve_PointerToPointer(t *testing.T) {
	value := "deep value"
	ptr := &value
	ptrptr := &ptr

	result := Resolve("", ptrptr, nil)
	// With empty path, should return the data itself
	if result != ptrptr {
		t.Errorf("Resolve with empty path = %v, want %v", result, ptrptr)
	}
}

func TestResolve_PointerToPointerStruct(t *testing.T) {
	inner := Address{
		Street: "456 Deep St",
		City:   "Boston",
		Zip:    02101,
	}
	ptr := &inner
	ptrptr := &ptr

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"pointer to pointer field", ".City", "Boston"},
		{"pointer to pointer street", ".Street", "456 Deep St"},
		{"pointer to pointer zip", ".Zip", 02101},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.path, ptrptr, nil)
			if result != tt.expected {
				t.Errorf("Resolve(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestResolve_PointerToPointerNil(t *testing.T) {
	var ptr *Address = nil
	ptrptr := &ptr

	result := Resolve(".City", ptrptr, nil)
	if result != nil {
		t.Errorf("Resolve on pointer to nil pointer should return nil, got %v", result)
	}
}

func TestResolve_MapWithIntKey(t *testing.T) {
	data := map[int]string{
		1: "one",
		2: "two",
	}

	result := Resolve(".[1]", data, nil)
	if result != "one" {
		t.Errorf("Resolve with int map key = %v, want %v", result, "one")
	}
}

func TestResolve_MapWithBoolKey(t *testing.T) {
	data := map[bool]string{
		true:  "yes",
		false: "no",
	}

	result := Resolve(".[true]", data, nil)
	if result != "yes" {
		t.Errorf("Resolve with bool map key = %v, want %v", result, "yes")
	}
}

func TestResolve_Array(t *testing.T) {
	data := [3]string{"a", "b", "c"}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"first", ".[0]", "a"},
		{"second", ".[1]", "b"},
		{"third", ".[2]", "c"},
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

func TestResolve_Interface(t *testing.T) {
	var data any = map[string]any{
		"name": "test",
		"nested": map[string]any{
			"value": 42,
		},
	}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{"interface map access", ".name", "test"},
		{"nested interface", ".nested.value", 42},
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

// Test the toString helper function
func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, ""},
		{"string", "hello", "hello"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"int", 42, "42"},
		{"int64", int64(123), "123"},
		{"float64", 3.14, "3.14"},
		{"struct", struct{ X int }{X: 1}, "{1}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toString(tt.input)
			if result != tt.expected {
				t.Errorf("toString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkResolve_SimpleField(b *testing.B) {
	person := createTestPerson()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve(".Name", person, nil)
	}
}

func BenchmarkResolve_NestedField(b *testing.B) {
	person := createTestPerson()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve(".Address.City", person, nil)
	}
}

func BenchmarkResolve_SliceAccess(b *testing.B) {
	person := createTestPerson()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve(".Tags[0]", person, nil)
	}
}

func BenchmarkResolve_MapAccess(b *testing.B) {
	person := createTestPerson()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve(".Scores[math]", person, nil)
	}
}

func BenchmarkResolve_Concatenation(b *testing.B) {
	person := createTestPerson()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve("'Hello, ' .Name '!'", person, nil)
	}
}

func BenchmarkResolve_Comparison(b *testing.B) {
	person := createTestPerson()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve("?.Age=='30'", person, nil)
	}
}

// Allocation benchmarks
func BenchmarkResolve_SimpleField_Allocs(b *testing.B) {
	person := createTestPerson()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve(".Name", person, nil)
	}
}

func BenchmarkResolve_NestedField_Allocs(b *testing.B) {
	person := createTestPerson()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve(".Address.City", person, nil)
	}
}
