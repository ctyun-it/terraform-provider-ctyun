package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

type validatorAclID struct {
}

const (
	AclIDError = "不满足acl id要求"
)

func AclID() validator.String {
	return &validatorAclID{}
}

func (v validatorAclID) Description(_ context.Context) string {
	return "不满足acl id 格式"
}

func (v validatorAclID) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorAclID) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	aclID := request.ConfigValue.ValueString()
	pattern := `^acl-[a-zA-Z0-9]{10}$`
	matched, err := regexp.MatchString(pattern, aclID)
	if err != nil {
		errMessage := "不满足acl id 格式"
		response.Diagnostics.AddError(AclIDError, errMessage)
		return
	}
	if !matched {
		errMessage := "不满足acl id 格式"
		response.Diagnostics.AddError(AclIDError, errMessage)
		return
	}

}
