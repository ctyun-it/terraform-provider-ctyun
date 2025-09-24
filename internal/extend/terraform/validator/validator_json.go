package validator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
)

const (
	JsonError = "不满足JSON格式"
)

type validatorJson struct {
	mustExistKey []string
}

func Json(key ...string) validator.String {
	return &validatorJson{mustExistKey: key}
}

func (v validatorJson) Description(_ context.Context) string {
	return JsonError
}

func (v validatorJson) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorJson) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	// 检查空字符串
	if strings.TrimSpace(value) == "" {
		response.Diagnostics.AddError(JsonError, "不能为空")
		return
	}

	// 尝试解析 JSON
	var data map[string]interface{}
	err := json.Unmarshal([]byte(value), &data)
	if err != nil {
		response.Diagnostics.AddError(JsonError, fmt.Sprintf("解析错误: %v", err))
		return
	}

	for _, key := range v.mustExistKey {
		if _, ok := data[key]; !ok {
			response.Diagnostics.AddError(JsonError, fmt.Sprintf("必须存在key: %s", key))
			return
		}
	}

	return
}
