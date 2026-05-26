package validator

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &validatorSearchPassword{}

type validatorSearchPassword struct {
}

func SearchPassword() validator.String {
	return &validatorSearchPassword{}
}

func (v validatorSearchPassword) Description(ctx context.Context) string {
	return "密码应为数字、大写字母、小写字母、特殊符号 (@$!%*#_~?) 的组合，长度在 12－26 位"
}

func (v validatorSearchPassword) MarkdownDescription(ctx context.Context) string {
	return "密码应为数字、大写字母、小写字母、特殊符号 (@$!%*#_~?) 的组合，长度在 12－26 位"
}

func (v validatorSearchPassword) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.String() == "" {
		return
	}

	value := req.ConfigValue.ValueString()

	// 检查长度 12-26 位
	if len(value) < 12 || len(value) > 26 {
		resp.Diagnostics.AddError(
			"密码长度无效",
			"密码长度必须在 12-26 位之间",
		)
		return
	}

	// 检查是否包含数字
	if !regexp.MustCompile(`[0-9]`).MatchString(value) {
		resp.Diagnostics.AddError(
			"密码复杂度不足",
			"密码必须包含至少一个数字",
		)
		return
	}

	// 检查是否包含小写字母
	if !regexp.MustCompile(`[a-z]`).MatchString(value) {
		resp.Diagnostics.AddError(
			"密码复杂度不足",
			"密码必须包含至少一个小写字母",
		)
		return
	}

	// 检查是否包含大写字母
	if !regexp.MustCompile(`[A-Z]`).MatchString(value) {
		resp.Diagnostics.AddError(
			"密码复杂度不足",
			"密码必须包含至少一个大写字母",
		)
		return
	}

	// 检查是否包含特殊符号
	if !regexp.MustCompile(`[@$!%*#_~?]`).MatchString(value) {
		resp.Diagnostics.AddError(
			"密码复杂度不足",
			"密码必须包含至少一个特殊符号 (@$!%*#_~?)",
		)
		return
	}
}
