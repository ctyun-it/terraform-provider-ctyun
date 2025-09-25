package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

var descRegex = regexp.MustCompile("^[a-zA-Z0-9\\x{4e00}-\\x{9fa5}~!@#$%^&*()_\\-+= <>?:'{},./;'[\\]·~！@#￥%……&*（） ——+={}]*$")

const (
	DescError = "不满足描述格式 支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·~！@#￥%……&*（） —— -+={}"
)

type validatorDesc struct {
}

func Desc() validator.String {
	return &validatorDesc{}
}

func (v validatorDesc) Description(_ context.Context) string {
	return DescError
}

func (v validatorDesc) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorDesc) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if !descRegex.MatchString(request.ConfigValue.ValueString()) {
		response.Diagnostics.AddError(DescError, DescError)
	}
}
