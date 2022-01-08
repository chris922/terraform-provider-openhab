package provider

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
