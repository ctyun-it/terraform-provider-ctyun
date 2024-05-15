package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"
	"strings"
)

const (
	RangeError = "不满足范围格式要求"
)

type validatorRange struct {
	splitter string
	min      int
	max      int
}

func Range(splitter string, min int, max int) validator.String {
	return &validatorRange{
		splitter: splitter,
		min:      min,
		max:      max,
	}
}

func (v validatorRange) Description(_ context.Context) string {
	return RangeError
}

func (v validatorRange) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorRange) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	getError := v.getError()
	// 如果有分割的字符，那么就直接拆分
	if strings.Contains(value, v.splitter) {
		strs := strings.Split(value, v.splitter)
		if len(strs) != 2 {
			response.Diagnostics.AddError(getError, getError)
			return
		}
		start, err := strconv.Atoi(strs[0])
		if err != nil {
			response.Diagnostics.AddError(getError, getError)
			return
		}
		end, err := strconv.Atoi(strs[1])
		if err != nil {
			response.Diagnostics.AddError(getError, getError)
			return
		}
		if v.min > start || v.max < end {
			response.Diagnostics.AddError(getError, getError)
			return
		}
	} else {
		// 如果没有找到分割字符，内容一定为一个数字
		target, err := strconv.Atoi(value)
		if err != nil {
			errMessage := "无法转换值：" + value
			response.Diagnostics.AddError(errMessage, errMessage)
			return
		}
		if v.min > target || v.max < target {
			response.Diagnostics.AddError(getError, getError)
			return
		}
	}
}

func (v validatorRange) getError() string {
	return RangeError + "：开始值" + v.splitter + "结束值"
}
