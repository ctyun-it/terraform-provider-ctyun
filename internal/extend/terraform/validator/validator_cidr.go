package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"net"
)

const (
	CidrError = "不满足cidr格式"
)

type validatorCidr struct {
}

func Cidr() validator.String {
	return &validatorCidr{}
}

func (v validatorCidr) Description(_ context.Context) string {
	return CidrError
}

func (v validatorCidr) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorCidr) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	if value == "" {
		return
	}
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		response.Diagnostics.AddError(CidrError, CidrError)
		return
	}
	if ipnet == nil || value != ipnet.String() {
		response.Diagnostics.AddError(CidrError, CidrError)
		return
	}
}
