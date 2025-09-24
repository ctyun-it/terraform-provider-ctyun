package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	AuthAddrError     = "不满足auth_addr ip CIDR格式"
	AuthAddrIpError   = "auth_addr的ip部分有误"
	AuthAddrMaskError = "auth_addr的掩码部分有误"
)

type validatorAuthAddr struct {
}

func (v validatorAuthAddr) Description(ctx context.Context) string {
	return AuthAddrError
}

func (v validatorAuthAddr) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorAuthAddr) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	cidr := request.ConfigValue.ValueString()
	if cidr == "" {
		return
	}
	// 检查基本格式：必须包含斜杠
	if !strings.Contains(cidr, "/") {
		response.Diagnostics.AddError(AuthAddrError, AuthAddrError)
		return
	}

	// 分割 IP 和掩码部分
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		response.Diagnostics.AddError(AuthAddrError, AuthAddrError)
		return
	}

	ipStr := parts[0]
	maskStr := parts[1]

	// 验证 IP 部分
	if !ValidateIP(ipStr) {
		response.Diagnostics.AddError(AuthAddrIpError, AuthAddrIpError)
		return
	}

	// 验证掩码部分
	mask, err := strconv.Atoi(maskStr)
	if err != nil {
		response.Diagnostics.AddError(AuthAddrMaskError, AuthAddrMaskError)
		return
	}

	// 根据 IP 类型验证掩码范围
	if ValidateIPv4(ipStr) {
		if mask < 0 || mask > 32 {
			response.Diagnostics.AddError(AuthAddrMaskError, AuthAddrMaskError)
			return
		}
	} else if ValidateIPv6(ipStr) {
		if mask < 0 || mask > 128 {
			response.Diagnostics.AddError(AuthAddrMaskError, AuthAddrMaskError)
			return
		}
	} else {
		response.Diagnostics.AddError(AuthAddrError, AuthAddrError)
		return
	}

	// 使用 net.ParseCIDR 进行最终验证
	_, _, err = net.ParseCIDR(cidr)
	return
}

// ValidateIP 校验 IP 地址格式是否正确（IPv4 或 IPv6）
func ValidateIP(ip string) bool {
	return ValidateIPv4(ip) || ValidateIPv6(ip)
}

func AuthAddr() validator.String {
	return &validatorAuthAddr{}
}

// ValidateIPv4 校验 IPv4 地址格式
func ValidateIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || parsedIP.To4() == nil {
		return false
	}

	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}

		if len(part) > 1 && part[0] == '0' {
			return false
		}

		num, err := strconv.Atoi(part)
		if err != nil {
			return false
		}

		if num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// ValidateIPv6 校验 IPv6 地址格式
func ValidateIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || parsedIP.To16() == nil || parsedIP.To4() != nil {
		return false
	}

	// 处理特殊格式
	if strings.Contains(ip, "%") {
		ip = strings.Split(ip, "%")[0]
	}

	if strings.Contains(ip, "]") {
		ip = strings.Split(ip, "]")[0]
		if strings.HasPrefix(ip, "[") {
			ip = ip[1:]
		}
	}

	// IPv6 正则表达式
	ipv6Pattern := `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|` +
		`^(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?$|` +
		`^([0-9a-fA-F]{1,4}:){6}(\d{1,3}\.){3}\d{1,3}$|` +
		`^::([fF]{4}:)?(\d{1,3}\.){3}\d{1,3}$`

	matched, _ := regexp.MatchString(ipv6Pattern, ip)
	return matched
}
