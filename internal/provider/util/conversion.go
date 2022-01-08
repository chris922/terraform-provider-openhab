package util

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringToType(v *string) types.String {
	if v == nil {
		return types.String{
			Null: true,
		}
	}

	return types.String{
		Value: *v,
	}
}

func TypeToString(v types.String) *string {
	if v.Unknown || v.Null {
		return nil
	}

	return &v.Value
}

func StringArrayToType(v *[]string) types.List {
	if v == nil {
		return types.List{
			ElemType: types.StringType,
			Null:     true,
		}
	}

	var vTyped []attr.Value
	for _, v2 := range *v {
		vTyped = append(vTyped, StringToType(&v2))
	}

	return types.List{
		ElemType: types.StringType,
		Elems:    vTyped,
	}
}

func TypeToStringArray(v types.List) *[]string {
	if v.Unknown || v.Null {
		return nil
	}

	var va []string
	for _, v2 := range v.Elems {
		va = append(va, *TypeToString(v2.(types.String)))
	}

	return &va
}

func StringMapToType(v *map[string]string) types.Map {
	if v == nil {
		return types.Map{
			Null: true,
		}
	}

	vTyped := make(map[string]attr.Value, len(*v))
	for k, v2 := range *v {
		vTyped[k] = StringToType(&v2)
	}

	return types.Map{
		ElemType: types.StringType,
		Elems:    vTyped,
	}
}

func TypeToStringMap(v types.Map) *map[string]string {
	if v.Unknown || v.Null {
		return nil
	}

	va := make(map[string]string, len(v.Elems))
	for k, v2 := range v.Elems {
		va[k] = *TypeToString(v2.(types.String))
	}

	return &va
}
