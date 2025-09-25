package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"unicode"
)

const (
	ScalingConfigNameError  = "不满足scaling config name格式"
	ScalingConfigNameLength = "名称长度必须为2～15个字符"
	ScalingConfigNameStart  = "不能以点号(.)或连字符(-)开头"
	ScalingConfigNameEnd    = "不能以点号(.)或连字符(-)结尾"
	ScalingConfigNameNum    = "不能仅使用数字"
)

type validatorScalingConfigName struct {
}

func ScalingConfigNameValidate() validator.String {
	return &validatorScalingConfigName{}
}

func (v validatorScalingConfigName) Description(ctx context.Context) string {
	return ScalingConfigNameError
}

func (v validatorScalingConfigName) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorScalingConfigName) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	name := request.ConfigValue.ValueString()
	if name == "" {
		return
	}
	// 规则 1: 长度验证 (2-15个字符)
	n := len(name)
	if n < 2 || n > 15 {
		response.Diagnostics.AddError(ScalingConfigNameLength, ScalingConfigNameLength)
		return
	}

	// 规则 2: 首尾字符验证
	first, last := rune(name[0]), rune(name[n-1])
	if first == '-' || first == '.' {
		response.Diagnostics.AddError(ScalingConfigNameStart, ScalingConfigNameStart)
		return
	}
	if last == '-' || last == '.' {
		response.Diagnostics.AddError(ScalingConfigNameEnd, ScalingConfigNameEnd)
		return
	}

	// 规则 3: 仅数字验证
	allDigits := true
	for _, r := range name {
		if !unicode.IsDigit(r) {
			allDigits = false
			break
		}
	}
	if allDigits {
		response.Diagnostics.AddError(ScalingConfigNameNum, ScalingConfigNameNum)
		return
	}

	prev := rune(0) // 跟踪上一个字符
	for i, r := range name {
		// 规则 4: 有效字符验证
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-') {
			response.Diagnostics.AddError(fmt.Sprintf("位置[%d]: 仅允许大小写字母、数字或连字符(-)", i), ScalingConfigNameError)
			return
		}

		// 规则 5: 连续符号验证
		if (r == '-' || r == '.') && (prev == '-' || prev == '.') {
			response.Diagnostics.AddError(fmt.Sprintf("位置[%d-%d]: 不能连续使用点号(.)或连字符(-)", i-1, i), ScalingConfigNameError)
			return
		}
		prev = r
	}

}
