package validator

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type validatorStringIsJson struct {
}

const (
	StringIsJsonError = "无效的 JSON 格式"
)

func StringIsJson() validator.String {
	return &validatorStringIsJson{}
}

func (v validatorStringIsJson) Description(_ context.Context) string {
	return "值必须是有效的 JSON 格式字符串"
}

func (v validatorStringIsJson) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorStringIsJson) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	jsonString := request.ConfigValue.ValueString()

	// 尝试解析 JSON 以验证其有效性
	var js json.RawMessage
	if err := json.Unmarshal([]byte(jsonString), &js); err != nil {
		response.Diagnostics.AddError(
			StringIsJsonError,
			"值必须是有效的 JSON 格式字符串",
		)
		return
	}
}
