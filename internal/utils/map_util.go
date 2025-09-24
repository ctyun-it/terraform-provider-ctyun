package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/url"
)

// TypesMapToStringMap 函数用于将 types.Map 转换为 map[string]string
func TypesMapToStringMap(ctx context.Context, m types.Map) (map[string]string, error) {
	// 检查是否为 null 或未知状态
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}

	result := make(map[string]string)
	// 遍历 types.Map 中的元素
	diags := m.ElementsAs(ctx, &result, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert types.Map to map[string]string: %v", diags)
	}

	return result, nil
}

// MapStringToTypesMap 将 map[string]string 转换为 types.Map
func MapStringToTypesMap(ctx context.Context, input map[string]string) (types.Map, error) {
	elements := make(map[string]types.String, len(input))
	for key, value := range input {
		elements[key] = types.StringValue(value)
	}
	result, diags := types.MapValueFrom(ctx, types.StringType, elements)
	if diags.HasError() {
		return types.MapNull(types.StringType), fmt.Errorf("failed to convert map[string]string to types.Map: %v", diags)
	}
	return result, nil
}

// MapToQueryString 函数将 map[string]string 转换为 k=v&k=v 格式的字符串
func MapToQueryString(m map[string]string) string {
	params := url.Values{}
	for k, v := range m {
		params.Add(k, v)
	}
	return params.Encode()
}
