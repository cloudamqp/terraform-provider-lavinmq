package converters

import "reflect"

func StructToMap(input any) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	for i := range val.NumField() {
		field := typ.Field(i)
		// Use Interface() to preserve the original type
		result[field.Name] = val.Field(i).Interface()
	}

	return result
}
