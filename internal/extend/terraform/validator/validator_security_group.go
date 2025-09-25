package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

const (
	SecurityGroupError = "不满足SecurityGroup格式"
)

type validatorSecurityGroup struct {
}

func (v validatorSecurityGroup) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// 正则表达式：sg-开头 + 10个小写字母或数字
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	if value == "" {
		return
	}
	pattern := `^sg-[a-z0-9]{10}$`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		response.Diagnostics.AddError(SecurityGroupError, SecurityGroupError)
		return
	}
}

func SecurityGroupValidate() validator.String {
	return &validatorSecurityGroup{}
}

func (v validatorSecurityGroup) Description(ctx context.Context) string {
	return SecurityGroupError
}

func (v validatorSecurityGroup) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
