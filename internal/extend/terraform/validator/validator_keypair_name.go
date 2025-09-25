package validator

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
)

type validatorKeypairName struct {
}

func KeypairName() validator.String {
	return &validatorKeypairName{}
}

func (v validatorKeypairName) Description(_ context.Context) string {
	return "不满足密钥对名称要求"
}

func (v validatorKeypairName) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorKeypairName) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	name := request.ConfigValue.ValueString()
	length := len(name)
	if length < 2 || length > 63 {
		errMessage := "长度必须在2-63"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	strs := strings.Split(name, "-")
	for _, str := range strs {
		if str == "" {
			errMessage := "不能以-开头和结尾，不能使用连续的-"
			response.Diagnostics.AddError(errMessage, errMessage)
			return
		}
		for _, c := range str {
			if utils.IsLetter(c) {
				continue
			} else if !utils.IsDigit(c) {
				errMessage := "只能包含大小写字母，数字以及-"
				response.Diagnostics.AddError(errMessage, errMessage)
				return
			}
		}
	}
}
