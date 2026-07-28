package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
	"strings"
)

type validatorDnsName struct {
}

const (
	DnsNameError = "不满足dns名称要求，由多个以点分隔的字符串组成，可包含字母、数字中划线、中划线不能在开头或末尾，单个字符串不超过63个字符，域名总长度不超过254个字符"
)

func DnsName() validator.String {
	return &validatorDnsName{}
}

func (v validatorDnsName) Description(_ context.Context) string {
	return DnsNameError
}

func (v validatorDnsName) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorDnsName) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	domain := request.ConfigValue.ValueString()
	// 检查总长度
	if len(domain) > 254 {
		errMessage := "长度必须小于254"
		response.Diagnostics.AddError(DnsNameError, errMessage)
		return
	}

	// 空域名无效
	if len(domain) == 0 {
		errMessage := "name不得为空"
		response.Diagnostics.AddError(DnsNameError, errMessage)
		return
	}

	// 检查开头和结尾不能是.和-
	if domain[0] == '.' || domain[len(domain)-1] == '.' || domain[0] == '-' || domain[len(domain)-1] == '-' {
		errMessage := "不得以.和-开头或结尾"
		response.Diagnostics.AddError(DnsNameError, errMessage)
		return
	}

	// 检查连续点
	if strings.Contains(domain, "..") {
		errMessage := "不能出现连续."
		response.Diagnostics.AddError(DnsNameError, errMessage)
		return
	}
	// 至少包含一个点
	if !strings.Contains(domain, ".") {
		errMessage := "必须由多个以点分隔的字符串组成"
		response.Diagnostics.AddError(DnsNameError, errMessage)
		return
	}

	// 正则表达式: 验证整个域名
	// 1. 整个字符串由点分隔的标签组成
	// 2. 每个标签: 字母数字开头，字母数字或中划线中间，字母数字结尾
	// 3. 每个标签长度1-63
	pattern := `^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])(\.[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])*$`
	matched, err := regexp.MatchString(pattern, domain)
	if err != nil {
		errMessage := "内网DNS的name为非法输入"
		response.Diagnostics.AddError(DnsNameError, errMessage)
		return
	}

	if !matched {
		errMessage := "内网DNS的name为非法输入"
		response.Diagnostics.AddError(DnsNameError, errMessage)
		return
	}
}
