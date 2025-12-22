package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
)

type validatorAclName struct {
}

const (
	AclNameError = "不满足acl名称要求,支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32"
)

func AclName() validator.String {
	return &validatorAclName{}
}

func (v validatorAclName) Description(_ context.Context) string {
	return "不满足acl名称要求,支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32"
}

func (v validatorAclName) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorAclName) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	name := request.ConfigValue.ValueString()
	length := len(name)
	if length < 2 || length > 32 {
		errMessage := "长度必须在2-32"
		response.Diagnostics.AddError(AclNameError, errMessage)
		return
	}

	// 2. 检查不能以 http: 或 https: 开头（不区分大小写）
	lowerStr := strings.ToLower(name)
	if strings.HasPrefix(lowerStr, "http:") || strings.HasPrefix(lowerStr, "https:") {
		errMessage := "不得以http:/https:开头"
		response.Diagnostics.AddError(AclNameError, errMessage)
		return
	}

	// 3. 检查第一个字符必须是中文或英文字母
	firstChar := []rune(name)[0]
	if !isChinese(firstChar) && !isEnglishLetter(firstChar) {
		errMessage := "只能中文 / 英文字母开头"
		response.Diagnostics.AddError(AclNameError, errMessage)
		return
	}

	// 4. 检查所有字符是否合法
	// 允许的字符：中文、英文字母、数字、下划线、连字符
	// 注意：Go的正则不支持\p{Han}，所以需要自己实现验证逻辑
	for _, r := range name {
		if !isValidRune(r) {
			errMessage := "只能支持拉丁字母、中文、数字，下划线，连字符"
			response.Diagnostics.AddError(AclNameError, errMessage)
			return
		}
	}
}

// 辅助函数
func isChinese(r rune) bool {
	// 中文字符的Unicode范围
	// 基本汉字: 4E00-9FFF
	// 扩展A: 3400-4DBF
	// 扩展B: 20000-2A6DF
	// 扩展C: 2A700-2B73F
	// 扩展D: 2B740-2B81F
	// 扩展E: 2B820-2CEAF
	// 扩展F: 2CEB0-2EBEF
	return (r >= 0x4E00 && r <= 0x9FFF) || // 基本汉字
		(r >= 0x3400 && r <= 0x4DBF) || // 扩展A
		(r >= 0x20000 && r <= 0x2A6DF) || // 扩展B
		(r >= 0x2A700 && r <= 0x2B73F) || // 扩展C
		(r >= 0x2B740 && r <= 0x2B81F) || // 扩展D
		(r >= 0x2B820 && r <= 0x2CEAF) || // 扩展E
		(r >= 0x2CEB0 && r <= 0x2EBEF) // 扩展F
}

func isEnglishLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isValidRune(r rune) bool {
	return isChinese(r) || isEnglishLetter(r) || isDigit(r) || r == '_' || r == '-'
}
