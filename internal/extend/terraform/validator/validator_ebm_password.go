package validator

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type validatorEbmPassword struct {
}

func EbmPassword() validator.String {
	return &validatorEbmPassword{}
}

func (v validatorEbmPassword) Description(_ context.Context) string {
	return "不满足ebm密码要求"
}

func (v validatorEbmPassword) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (v validatorEbmPassword) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	password := request.ConfigValue.ValueString()
	length := len(password)
	if length < 8 || length > 30 {
		errMessage := "ebm密码长度必须在8-30"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	hasUpperLetter := false
	hasLowerLetter := false
	hasDigit := false
	hasSpecialSymbols := false

	for _, r := range password {
		if utils.IsDigit(r) {
			hasDigit = true
		} else if utils.IsLower(r) {
			hasLowerLetter = true
		} else if utils.IsUpper(r) {
			hasUpperLetter = true
		} else if isHasSpecialSymbols(r) {
			hasSpecialSymbols = true
		} else {
			// 包含不允许的字符
			errMessage := "ebm密码只能包含大小写字母、数字和指定特殊符号"
			response.Diagnostics.AddError(errMessage, errMessage)
			return
		}
	}

	// 验证必须包含大写字母
	if !hasUpperLetter {
		errMessage := "ebm密码必须包含大写字母"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}

	// 验证必须包含小写字母
	if !hasLowerLetter {
		errMessage := "ebm密码必须包含小写字母"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}

	// 验证至少包含一个数字或者一个特殊字符
	if !hasDigit && !hasSpecialSymbols {
		errMessage := "ebm密码必须至少包含一个数字或者一个特殊字符"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
}
