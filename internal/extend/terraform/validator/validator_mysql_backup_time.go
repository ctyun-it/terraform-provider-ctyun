package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"
)

const (
	BackupTimeError = "不满足（小时：分钟）备份时间格式"
)

type validatorBackupTime struct {
}

func (v validatorBackupTime) Description(ctx context.Context) string {
	return BackupTimeError
}

func (v validatorBackupTime) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorBackupTime) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	timeStr := request.ConfigValue.ValueString()
	// 1. 检查基本格式
	if len(timeStr) != 5 {
		errMessage := "长度有误，长度必须为5"
		response.Diagnostics.AddError(BackupTimeError, errMessage)
		return
	}
	if timeStr[2] != ':' {
		errMessage := "输入时间有误，必须为00:00格式"
		response.Diagnostics.AddError(BackupTimeError, errMessage)
		return
	}

	// 2. 提取小时和分钟部分
	hourStr := timeStr[0:2]
	minuteStr := timeStr[3:5]

	// 3. 验证是否为数字
	if !isDigits(hourStr) || !isDigits(minuteStr) {
		response.Diagnostics.AddError(BackupTimeError, BackupTimeError)
		return
	}

	// 4. 转换为数字并验证范围
	hour, _ := strconv.Atoi(hourStr)
	minute, _ := strconv.Atoi(minuteStr)
	flag := hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59
	if !flag {
		response.Diagnostics.AddError(BackupTimeError, BackupTimeError)
		return
	}
}

func BackupTimeValidator() validator.String {
	return &validatorBackupTime{}
}

// isDigits 检查字符串是否全为数字
func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
