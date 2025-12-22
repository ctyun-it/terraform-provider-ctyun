package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

var descNotStartWithHttpRegex = regexp.MustCompile("^([^h]|h[^t]|ht[^t]|htt[^p]|http[^s]|https.).*$")

const (
	DescNotStartWithHttpError = "不能以http:或https:开头"
)

type validatorDescNotStartWithHttp struct {
}

func DescNotStartWithHttp() validator.String {
	return &validatorDescNotStartWithHttp{}
}

func (v validatorDescNotStartWithHttp) Description(_ context.Context) string {
	return DescNotStartWithHttpError
}

func (v validatorDescNotStartWithHttp) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorDescNotStartWithHttp) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if !descNotStartWithHttpRegex.MatchString(request.ConfigValue.ValueString()) {
		response.Diagnostics.AddError(DescNotStartWithHttpError, DescNotStartWithHttpError)
	}
}
