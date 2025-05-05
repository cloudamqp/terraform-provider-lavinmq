package converters

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name  string
	Age   int
	Admin bool
}

func TestStructToMap(t *testing.T) {
	// Input struct
	input := TestStruct{
		Name:  "John Doe",
		Age:   30,
		Admin: true,
	}

	// Expected output
	expected := map[string]any{
		"Name":  "John Doe",
		"Age":   30,
		"Admin": true,
	}

	// Call the function
	result := StructToMap(input)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StructToMap() = %v, want %v", result, expected)
	}
}
