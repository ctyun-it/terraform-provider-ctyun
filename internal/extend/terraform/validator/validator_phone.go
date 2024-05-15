package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

var phoneRegex = regexp.MustCompile(`^1([3456789])\d{9}$`)

const (
	PhoneError = "不满足手机号格式"
)

type validatorPhone struct {
}

func Phone() validator.String {
	return &validatorPhone{}
}

func (v validatorPhone) Description(_ context.Context) string {
	return PhoneError
}

func (v validatorPhone) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorPhone) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if !phoneRegex.MatchString(request.ConfigValue.ValueString()) {
		response.Diagnostics.AddError(PhoneError, PhoneError)
	}
}
