package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
	"unicode"
)

const (
	PgsqlDatabaseNameError = "不满足postgresql数据库实例数据库名称规范"
)

type validatorPgsqlDatabaseName struct {
}

func (v validatorPgsqlDatabaseName) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	name := request.ConfigValue.ValueString()
	// 规则1: 长度在2~63个字符之间
	if len(name) < 2 || len(name) > 63 {
		errMessage := fmt.Sprintf("数据库名称长度必须在2~63个字符之间（当前长度：%d）", len(name))
		response.Diagnostics.AddError(PgsqlDatabaseNameError, errMessage)
		return
	}

	// 规则2: 以字母开头
	firstChar := rune(name[0])
	if !unicode.IsLetter(firstChar) {
		errMessage := fmt.Sprintf("数据库名称必须以字母开头（当前开头字符：%c）", firstChar)
		response.Diagnostics.AddError(PgsqlDatabaseNameError, errMessage)
		return
	}

	// 规则3: 以字母或数字结尾
	lastChar := rune(name[len(name)-1])
	if !unicode.IsLetter(lastChar) && !unicode.IsNumber(lastChar) {
		errMessage := fmt.Sprintf("数据库名称必须以字母或数字结尾（当前结尾字符：%c）", lastChar)
		response.Diagnostics.AddError(PgsqlDatabaseNameError, errMessage)
		return
	}

	// 规则4: 只能包含小写字母、数字、下划线或中划线
	validPattern := regexp.MustCompile(`^[a-z0-9_-]+$`)
	if !validPattern.MatchString(name) {
		// 找出非法字符
		var invalidChars []rune
		for _, char := range name {
			if !(char >= 'a' && char <= 'z') &&
				!(char >= '0' && char <= '9') &&
				char != '_' && char != '-' {
				invalidChars = append(invalidChars, char)
			}
		}
		errMessage := fmt.Sprintf("数据库名称只能包含小写字母、数字、下划线(_)或中划线(-)（非法字符：%v）", invalidChars)
		response.Diagnostics.AddError(PgsqlDatabaseNameError, errMessage)
		return
	}
}

func PgsqlDatabaseName() validator.String {
	return &validatorPgsqlDatabaseName{}
}

func (v validatorPgsqlDatabaseName) Description(ctx context.Context) string {
	return PgsqlDatabaseNameError
}

func (v validatorPgsqlDatabaseName) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
