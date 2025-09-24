package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

var uuidRegex = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")

const (
	UUIDError = "不满足UUID格式"
)

type validatorUUID struct {
}

func UUID() validator.String {
	return &validatorUUID{}
}

func (v validatorUUID) Description(_ context.Context) string {
	return UUIDError
}

func (v validatorUUID) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorUUID) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if !uuidRegex.MatchString(request.ConfigValue.ValueString()) {
		response.Diagnostics.AddError(UUIDError, UUIDError)
	}
}
