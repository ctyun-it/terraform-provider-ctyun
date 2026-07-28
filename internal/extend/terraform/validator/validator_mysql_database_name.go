package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"unicode"
)

const (
	MysqlDatabaseNameError = "不满足mysql数据库实例数据库名称规范"
)

type validatorMysqlDatabaseName struct {
}

func (v validatorMysqlDatabaseName) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	name := request.ConfigValue.ValueString()

	// 检查第一个字符是否为小写字母
	firstChar := rune(name[0])
	if !unicode.IsLower(firstChar) {
		errMessage := fmt.Sprintf("必须以小写字母开头（当前开头字符：%c）", firstChar)
		response.Diagnostics.AddError(MysqlDatabaseNameError, errMessage)
		return
	}

	// 检查最后一个字符是否为小写字母或数字
	lastChar := rune(name[len(name)-1])
	if !unicode.IsLower(lastChar) && !unicode.IsNumber(lastChar) {
		errMessage := fmt.Sprintf("必须以小写字母或数字结尾（当前结尾字符：%c）", lastChar)
		response.Diagnostics.AddError(MysqlDatabaseNameError, errMessage)
		return
	}

	// 检查所有字符是否合法
	for i, char := range name {
		if !isValidCharacter(char) {
			errMessage := fmt.Sprintf("位置 %d 包含非法字符：%c（只允许小写字母、数字或下划线）", i+1, char)
			response.Diagnostics.AddError(MysqlDatabaseNameError, errMessage)
			return
		}
	}

}

// isValidCharacter 检查字符是否合法
func isValidCharacter(char rune) bool {
	// 允许小写字母
	if unicode.IsLower(char) {
		return true
	}

	// 允许数字
	if unicode.IsNumber(char) {
		return true
	}

	// 允许下划线
	if char == '_' {
		return true
	}

	return false
}

func MysqlDatabaseName() validator.String {
	return &validatorMysqlDatabaseName{}
}

func (v validatorMysqlDatabaseName) Description(ctx context.Context) string {
	return MysqlDatabaseNameError
}

func (v validatorMysqlDatabaseName) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
