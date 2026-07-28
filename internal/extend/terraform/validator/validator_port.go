package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

const (
	PortError = "不满足port格式"
)

type validatorPort struct {
}

func (v validatorPort) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// 正则表达式：port-开头 + 10个小写字母或数字
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	if value == "" {
		return
	}
	// 匹配uuid，说明是3.0资源池
	if uuidRegex.MatchString(request.ConfigValue.ValueString()) {
		return
	}
	pattern := `^port-[a-z0-9]{10}$`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		response.Diagnostics.AddError(PortError, PortError)
		return
	}
}

func PortValidate() validator.String {
	return &validatorPort{}
}

func (v validatorPort) Description(ctx context.Context) string {
	return PortError
}

func (v validatorPort) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
