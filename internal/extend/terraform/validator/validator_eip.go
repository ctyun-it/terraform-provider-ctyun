package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

const (
	EipError = "不满足EipID格式"
)

type validatorEip struct {
}

func (v validatorEip) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// 正则表达式：Eip-开头 + 10个小写字母或数字
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	if value == "" {
		return
	}
	pattern := `^eip-[a-z0-9]{10}$`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		response.Diagnostics.AddError(IpError, IpError)
		return
	}
}

func EipValidate() validator.String {
	return &validatorEip{}
}

func (v validatorEip) Description(ctx context.Context) string {
	return EipError
}

func (v validatorEip) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
