package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
	"unicode"
)

const (
	ScalingConfigNamePwError        = "不满足scaling config password格式"
	ScalingConfigNamePwLengthError  = "密码长度必须为8～30个字符"
	ScalingConfigNamePwStartError   = "不能以斜线号(/)开头"
	ScalingConfigNamePwWindowsError = "Windows密码不能包含管理员用户名"
	ScalingConfigNamePwTypeError    = "密码必须包含大写字母、小写字母、数字和特殊符号中的至少3类"
	ScalingConfigNamePwSeqCharError = "密码不能包含3个及以上连续字符"
	ScalingConfigNamePwSeqNumError  = "密码不能包含3个及以上连续数字"
)

type validatorScalingConfigPassword struct {
}

func ScalingConfigPasswordValidate() validator.String {
	return &validatorScalingConfigPassword{}
}

func (v validatorScalingConfigPassword) Description(ctx context.Context) string {
	return ScalingConfigNamePwError
}

func (v validatorScalingConfigPassword) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorScalingConfigPassword) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	password := request.ConfigValue.ValueString()
	if password == "" {
		return
	}
	// 规则1: 长度验证 (8～30个字符)
	n := len(password)
	if n < 8 || n > 30 {
		response.Diagnostics.AddError(ScalingConfigNamePwLengthError, ScalingConfigNamePwLengthError)
		return
	}

	// 规则2: 不能以斜线号(/)开头
	if strings.HasPrefix(password, "/") {
		response.Diagnostics.AddError(ScalingConfigNamePwStartError, ScalingConfigNamePwStartError)
		return
	}

	// 规则3: 检测Windows用户名
	//if isWindows {
	//	lowerPass := strings.ToLower(password)
	//	if strings.Contains(lowerPass, "administrator") {
	//		response.Diagnostics.AddError(ScalingConfigNamePwWindowsError, ScalingConfigNamePwWindowsError)
	//		return
	//	}
	//}

	// 规则4: 检查必须包含的字符类型
	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	// 定义允许的特殊字符
	specialChars := "()`~!@#$%^&*_-+=|{}[]:;'<>,.?/"
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case strings.ContainsRune(specialChars, r):
			hasSpecial = true
		}
	}
	// 检查是否包含至少三类字符
	typesCount := 0
	if hasUpper {
		typesCount++
	}
	if hasLower {
		typesCount++
	}
	if hasDigit {
		typesCount++
	}
	if hasSpecial {
		typesCount++
	}

	if typesCount < 3 {
		response.Diagnostics.AddError(ScalingConfigNamePwTypeError, ScalingConfigNamePwTypeError)
		return
	}

	// 规则5: 检查连续字符
	for i := 0; i <= len(password)-3; i++ {
		seq := password[i : i+3]
		// 检查字母连续 (如abc, xyz)
		if isAlphabeticSequence(seq) {
			response.Diagnostics.AddError(fmt.Sprintf("密码不能包含3个及以上连续字符 (%q)", seq), ScalingConfigNamePwSeqCharError)
			return
		}

		// 检查数字连续 (如123, 789)
		if isNumericSequence(seq) {
			response.Diagnostics.AddError(fmt.Sprintf("密码不能包含3个及以上连续数字 (%q)", seq), ScalingConfigNamePwSeqNumError)
			return
		}
	}

}

//func getStringAttribute(config tfsdk.Config, attrPath path.Path) (types.String, error) {
//	var value types.String
//	diags := config.GetAttribute(context.Background(), attrPath, &value)
//	if diags.HasError() {
//		return types.StringValue(""), fmt.Errorf("failed to get attribute: %v", diags.Errors())
//	}
//	return value, nil
//}

// 检查是否是连续的字母序列
func isAlphabeticSequence(s string) bool {
	if len(s) != 3 {
		return false
	}

	// 如果包含非字母字符则跳过
	if !unicode.IsLetter(rune(s[0])) || !unicode.IsLetter(rune(s[1])) || !unicode.IsLetter(rune(s[2])) {
		return false
	}

	// 转换为小写进行判断
	s = strings.ToLower(s)
	return (s[1]-s[0] == 1 && s[2]-s[1] == 1) || // 升序 (abc)
		(s[0]-s[1] == 1 && s[1]-s[2] == 1) // 降序 (cba)
}

// 检查是否是连续的数字序列
func isNumericSequence(s string) bool {
	if len(s) != 3 {
		return false
	}

	// 如果包含非数字字符则跳过
	if !unicode.IsDigit(rune(s[0])) || !unicode.IsDigit(rune(s[1])) || !unicode.IsDigit(rune(s[2])) {
		return false
	}

	return (s[1]-s[0] == 1 && s[2]-s[1] == 1) || // 升序 (123)
		(s[0]-s[1] == 1 && s[1]-s[2] == 1) // 降序 (321)
}
