package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

const (
	DatabaseSecurityGroupError = "不满足DatabaseSecurityGroup格式，数据库安全组格式要求：sg-xxxx,sg-xxxx"
)

type validatorDatabaseSecurityGroup struct {
}

func (v validatorDatabaseSecurityGroup) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	securityGroupID := request.ConfigValue.ValueString()
	// 先判断是否以,开头或结尾，
	flag := isSecurityGroupID(securityGroupID)
	if !flag {
		errMessage := fmt.Sprintf("输入的securityGroupID不符合要求：%s", securityGroupID)
		response.Diagnostics.AddError(DnsNameError, errMessage)
		return
	}
}

func isSecurityGroupID(securityGroupID string) bool {
	// 正则表达式：sg-开头 + 10个小写字母或数字

	// 匹配uuid，说明是3.0资源池
	if uuidRegex.MatchString(securityGroupID) {
		return true
	}
	pattern := `^sg-[a-z0-9]{10}$`
	matched, _ := regexp.MatchString(pattern, securityGroupID)
	if !matched {
		return false
	}
	return true
}

func DatabaseSecurityGroupValidate() validator.String {
	return &validatorDatabaseSecurityGroup{}
}

func (v validatorDatabaseSecurityGroup) Description(ctx context.Context) string {
	return DatabaseSecurityGroupError
}

func (v validatorDatabaseSecurityGroup) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
