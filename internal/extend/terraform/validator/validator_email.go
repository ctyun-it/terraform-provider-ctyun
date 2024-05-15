package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")

const (
	EmailError = "不满足邮箱格式"
)

type validatorEmail struct {
}

func Email() validator.String {
	return &validatorEmail{}
}

func (v validatorEmail) Description(_ context.Context) string {
	return EmailError
}

func (v validatorEmail) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorEmail) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if !emailRegex.MatchString(request.ConfigValue.ValueString()) {
		response.Diagnostics.AddError(EmailError, EmailError)
	}
}
