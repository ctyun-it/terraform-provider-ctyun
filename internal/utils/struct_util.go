package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
)

// StructToTFObjectTypes将结构体转换为types.ObjectType类型
func StructToTFObjectTypes(s interface{}) types.ObjectType {
	result := make(map[string]attr.Type)
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tfsdk")
		if tag == "" {
			panic(fmt.Sprintf("StructToTFObjectTypes %s must have tfsdk tag", field.Name))
		}
		var fieldType attr.Type
		switch field.Type {
		case reflect.TypeOf(types.String{}):
			fieldType = types.StringType
		case reflect.TypeOf(types.Bool{}):
			fieldType = types.BoolType
		case reflect.TypeOf(types.Int64{}):
			fieldType = types.Int64Type
		case reflect.TypeOf(types.Int32{}):
			fieldType = types.Int32Type
		case reflect.TypeOf(types.Float64{}):
			fieldType = types.Float64Type
		case reflect.TypeOf(types.Float32{}):
			fieldType = types.Float32Type
		case reflect.TypeOf(types.List{}):
			fieldType = types.ListType{}
		//case reflect.TypeOf(types.List{}):
		//  fieldType = types.ListType{ElemType: types.StringType}
		//case reflect.TypeOf(types.Map{}):
		//  fieldType = types.MapType{ElemType: types.StringType}
		default:
			panic(fmt.Sprintf("StructToTFObjectTypes not support %s", field.Type.String()))
		}
		result[tag] = fieldType
	}
	return types.ObjectType{AttrTypes: result}
}

// DifferenceStructArray 获取两个结构体切片的差集
func DifferenceStructArray[T comparable](a, b []T) ([]T, []T) {
	setB := make(map[T]bool)
	for _, item := range b {
		setB[item] = true
	}

	var onlyInA []T
	for _, item := range a {
		if !setB[item] {
			onlyInA = append(onlyInA, item)
		}
	}

	setA := make(map[T]bool)
	for _, item := range a {
		setA[item] = true
	}

	var onlyInB []T
	for _, item := range b {
		if !setA[item] {
			onlyInB = append(onlyInB, item)
		}
	}

	return onlyInA, onlyInB
}
