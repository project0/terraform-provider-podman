package utils

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MapStringEmpty() types.Map {
	m, _ := types.MapValue(types.StringType, make(map[string]attr.Value))
	return m
}

// MapStringToMapType maps a native golang map to a terraform map type
func MapStringToMapType(m map[string]string, diags *diag.Diagnostics) types.Map {
	elems := make(map[string]attr.Value)
	for k, v := range m {
		elems[k] = types.StringValue(v)
	}
	v, d := types.MapValue(types.StringType, elems)
	diags.Append(d...)
	return v
}

// MapStringValueToStringType extracts a terraform string value from a map
func MapStringValueToStringType(m map[string]string, key string) types.String {
	if val, exist := m[key]; exist {
		return types.StringValue(val)
	}
	return types.StringNull()
}

// MapStringValueToIntType extracts a terraform int value from a map with string
func MapStringValueToIntType(m map[string]string, key string, diags *diag.Diagnostics) types.Int64 {
	val, exist := m[key]
	i := 0
	if exist {
		var err error
		i, err = strconv.Atoi(val)
		if err != nil {
			diags.AddError("Cannot convert string to integer", fmt.Sprintf("Received value %s for key %s is not convertable: %s ", val, key, err.Error()))
		}
		return types.Int64Value(int64(i))
	}
	return types.Int64Null()
}
