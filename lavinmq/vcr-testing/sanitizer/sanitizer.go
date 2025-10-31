package sanitizer

import (
	"strings"
)

func FilterSensitiveData(jsonBody, value, placeholder string) string {
	if len(value) == 0 {
		return jsonBody
	}
	return strings.ReplaceAll(jsonBody, value, placeholder)
}
