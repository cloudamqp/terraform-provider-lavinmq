package utils

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Admin bool   `json:"admin"`
}

func TestGenericMarshal(t *testing.T) {
	// Input data
	inputData := TestStruct{
		Name:  "John Doe",
		Age:   30,
		Admin: true,
	}

	// Expected JSON data
	expectedJSON := `{"name":"John Doe","age":30,"admin":true}`

	// Call the function
	result, err := GenericMarshal(inputData)
	if err != nil {
		t.Fatalf("GenericMarshal() returned an error: %v", err)
	}

	resultString := string(result)

	// Compare the result with the expected output
	if resultString != expectedJSON {
		t.Errorf("GenericMarshal() = %s, want %s", resultString, expectedJSON)
	}
}

func TestGenericUnmarshal(t *testing.T) {
	// Input JSON
	jsonData := `{"name":"John Doe","age":30,"admin":true}`

	// Expected output
	expected := TestStruct{
		Name:  "John Doe",
		Age:   30,
		Admin: true,
	}

	// Call the function
	result, err := GenericUnmarshal[TestStruct]([]byte(jsonData))
	if err != nil {
		t.Fatalf("GenericUnmarshal() returned an error: %v", err)
	}

	// Compare the result with the expected output
	if !reflect.DeepEqual(*result, expected) {
		t.Errorf("GenericUnmarshal() = %v, want %v", *result, expected)
	}
}
