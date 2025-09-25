package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

const (
	SubnetError = "不满足Subnet格式"
)

type validatorSubnet struct {
}

func (v validatorSubnet) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// 正则表达式：Subnet-开头 + 10个小写字母或数字
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	if value == "" {
		return
	}
	pattern := `^subnet-[a-z0-9]{10}$`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		response.Diagnostics.AddError(SubnetError, SubnetError)
		return
	}
}

func SubnetValidate() validator.String {
	return &validatorSubnet{}
}

func (v validatorSubnet) Description(ctx context.Context) string {
	return SubnetError
}

func (v validatorSubnet) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
