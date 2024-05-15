package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"slices"
	"strings"
	"terraform-provider-ctyun/internal/utils"
)

type validatorEcsPassword struct {
}

func EcsPassword() validator.String {
	return &validatorEcsPassword{}
}

func (v validatorEcsPassword) Description(_ context.Context) string {
	return "不满足ecs密码要求"
}

func (v validatorEcsPassword) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorEcsPassword) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	password := request.ConfigValue.ValueString()
	length := len(password)
	if length < 8 || length > 30 {
		errMessage := "ecs密码长度必须在8-30"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	if strings.HasPrefix(password, "/") {
		errMessage := "ecs密码不能以/开头"
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
			errMessage := "ecs密码只能为数字，大小写字母以及特殊符号"
			response.Diagnostics.AddError(errMessage, errMessage)
			return
		}
	}

	count := 0
	if hasUpperLetter {
		count++
	}
	if hasLowerLetter {
		count++
	}
	if hasDigit {
		count++
	}
	if hasSpecialSymbols {
		count++
	}

	if count < 3 {
		errMessage := "ecs密码必须包含大小写字母、数字、特殊符号中的至少三种"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
}

// isHasSpecialSymbols 是否包含特殊字符
func isHasSpecialSymbols(target int32) bool {
	return slices.Contains([]int32{'(', ')', '`', '~', '!', '@', '#', '$', '%', '^', '&', '*', '_', '-', '+', '=', '|', '{', '}', '[', ']', ':', ';', '\'', '<', '>', ',', '.', '?', '/', '\\'}, target)
}
