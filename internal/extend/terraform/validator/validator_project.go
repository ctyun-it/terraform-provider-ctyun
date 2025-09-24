package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

var ProjectRegex = regexp.MustCompile("^(|0|[a-z0-9]{32})$")

const (
	ProjectError = "不满足 允许 \"\",\"0\"和长度为32位的数字+小写"
)

type validatorProject struct {
}

func Project() validator.String {
	return &validatorProject{}
}

func (v validatorProject) Description(_ context.Context) string {
	return ProjectError
}

func (v validatorProject) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorProject) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if !ProjectRegex.MatchString(request.ConfigValue.ValueString()) {
		response.Diagnostics.AddError(ProjectError, ProjectError)
	}
}
