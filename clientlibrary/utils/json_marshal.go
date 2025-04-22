package utils

import "encoding/json"

// Marshal any type into a JSON-encoding byte slice
func GenericMarshal[T any](input T) ([]byte, error) {
	return json.Marshal(input)
}

// Unmarshal a JSON-encoded byte slice into any type
func GenericUnmarshal[T any](body []byte) (*T, error) {
	var output *T
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}
	return output, nil
}
