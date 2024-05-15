package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"net"
)

const (
	IpError = "不满足ip格式"
)

type validatorIp struct {
}

func Ip() validator.String {
	return &validatorIp{}
}

func (v validatorIp) Description(_ context.Context) string {
	return IpError
}

func (v validatorIp) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorIp) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	if value == "" {
		return
	}
	ip := net.ParseIP(value)
	if ip == nil {
		response.Diagnostics.AddError(IpError, IpError)
		return
	}
}
