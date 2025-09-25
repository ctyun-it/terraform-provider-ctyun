package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

const (
	VpcError = "不满足vpc格式"
)

type validatorVpc struct {
}

func (v validatorVpc) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// 正则表达式：vpc-开头 + 10个小写字母或数字
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	if value == "" {
		return
	}
	pattern := `^vpc-[a-z0-9]{10}$`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		response.Diagnostics.AddError(IpError, IpError)
		return
	}
}

func VpcValidate() validator.String {
	return &validatorVpc{}
}

func (v validatorVpc) Description(ctx context.Context) string {
	return VpcError
}

func (v validatorVpc) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
