package converters

import "encoding/json"

func StructToMap(obj any) (newMap map[string]any) {
	data, _ := json.Marshal(obj)
	_ = json.Unmarshal(data, &newMap)
	return newMap
}
