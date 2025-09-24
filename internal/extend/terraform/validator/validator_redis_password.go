package validator

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"slices"
	"strings"
)

type validatorRedisPassword struct {
}

func RedisPassword() validator.String {
	return &validatorRedisPassword{}
}

func (v validatorRedisPassword) Description(_ context.Context) string {
	return "不满足Redis密码要求"
}

func (v validatorRedisPassword) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorRedisPassword) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	password := request.ConfigValue.ValueString()
	length := len(password)
	if length < 8 || length > 26 {
		errMessage := "redis密码长度必须在8-26"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	if strings.Contains(password, " ") {
		errMessage := "redis密码不能包含空格"
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
		} else if redisSpecialSymbols(r) {
			hasSpecialSymbols = true
		} else {
			errMessage := "redis密码只能为数字，大小写字母以及特殊符号(@%^*_+!$-=.)"
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
		errMessage := "redis密码必须同时包含大写字母、小写字母、数字、英文格式特殊符号(@%^*_+!$-=.)中的三种类型"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
}

func redisSpecialSymbols(target int32) bool {
	return slices.Contains([]int32{'(', '@', '%', '^', '*', '_', '+', '!', '$', '-', '=', '.', ')'}, target)
}
