package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
	"unicode"
)

type validatorDBPassword struct {
	minLength    int
	maxLength    int
	typeNum      int
	name         string
	specialChars string
	errorMessage string
}

func (v validatorDBPassword) Description(ctx context.Context) string {
	return v.errorMessage
}

func (v validatorDBPassword) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorDBPassword) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	password := request.ConfigValue.ValueString()
	// 检查长度
	if len(password) < v.minLength || len(password) > v.maxLength {
		errMessage := fmt.Sprintf("%s实例密码长度需要保持%d~%d位", v.name, v.minLength, v.maxLength)
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	// 定义允许的特殊字符
	specialChars := v.specialChars
	// 初始化字符类型标记
	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
		validChars = true
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		default:
			validChars = false // 发现非法字符
		}
	}

	// 统计满足的字符类型数量
	typeCount := 0
	if hasUpper {
		typeCount++
	}
	if hasLower {
		typeCount++
	}
	if hasDigit {
		typeCount++
	}
	if hasSpecial {
		typeCount++
	}

	// 验证结果
	if !validChars {
		errMessage := "存在非法字符，" + fmt.Sprintf("密码仅支持大写字母、小写字母、数字和特殊字符%s", v.specialChars)
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	if typeCount < v.typeNum {
		errMessage := fmt.Sprintf("密码组合必须包括大写字母、小写字母、数字和特殊字符中的任意%d种及以上", v.typeNum)
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}

}

func DBPassword(minLength, maxLength, typeNum int, name, specialChars string) validator.String {
	return &validatorDBPassword{
		minLength:    minLength,
		maxLength:    maxLength,
		typeNum:      typeNum,
		errorMessage: fmt.Sprintf("%s密码长度为%d~%d个字符，至少包含大写字母、小写字母、数字和特殊字符中的%d种，特殊字符支持%s", name, minLength, maxLength, typeNum, specialChars),
		specialChars: specialChars,
		name:         name,
	}
}
