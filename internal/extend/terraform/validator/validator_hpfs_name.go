package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type validatorHpfsName struct {
}

const (
	HpfsNameError = "不满足hpfs名称要求，并行文件名，仅允许英文字母数字及-，开头必须为字母，结尾不允许为-，且长度为2-255字符"
)

func HpfsName() validator.String {
	return &validatorHpfsName{}
}

func (v validatorHpfsName) Description(_ context.Context) string {
	return HpfsNameError
}

func (v validatorHpfsName) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorHpfsName) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	name := request.ConfigValue.ValueString()
	length := len(name)
	// 检查总长度
	if len(name) < 2 || len(name) > 255 {
		errMessage := "长度必须为2-255个字符"
		response.Diagnostics.AddError(HpfsNameError, errMessage)
		return
	}
	// 检查第一个字符必须是字母
	firstChar := name[0]
	if !(('a' <= firstChar && firstChar <= 'z') || ('A' <= firstChar && firstChar <= 'Z')) {
		errMessage := "name开头必须为字母"
		response.Diagnostics.AddError(HpfsNameError, errMessage)
		return
	}

	// 检查最后一个字符不能是连字符
	lastChar := name[length-1]
	if lastChar == '-' {
		errMessage := "name最后一个字符不能是连字符"
		response.Diagnostics.AddError(HpfsNameError, errMessage)
		return
	}

	// 检查所有字符
	for i := 0; i < length; i++ {
		c := name[i]
		// 允许: 字母、数字、连字符
		if !(('a' <= c && c <= 'z') ||
			('A' <= c && c <= 'Z') ||
			('0' <= c && c <= '9') ||
			c == '-') {
			errMessage := "name包含非法字符："
			response.Diagnostics.AddError(HpfsNameError, errMessage+string(c))
			return
		}
	}
}
