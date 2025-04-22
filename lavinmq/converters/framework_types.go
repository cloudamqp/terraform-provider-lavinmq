package converters

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AttrValuesToStrings(values []attr.Value) ([]string, error) {
	var result []string
	for _, v := range values {
		str, ok := v.(types.String)
		if !ok {
			return nil, fmt.Errorf("expected types.String, got %T", v)
		}
		if !str.IsNull() && !str.IsUnknown() {
			result = append(result, str.ValueString())
		}
	}
	return result, nil
}

func StringsToAttrValues(strings []string) []attr.Value {
	attrValues := make([]attr.Value, len(strings))
	for i, s := range strings {
		attrValues[i] = types.StringValue(s)
	}
	return attrValues
}
