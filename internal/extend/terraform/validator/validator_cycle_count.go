package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-ctyun/internal/business"
)

type cycleCountRange struct {
	minInclude int
	maxInclude int
}

type validatorCycleCount struct {
	yearRange  cycleCountRange
	monthRange cycleCountRange
}

func CycleCount(monthMinInclude, monthMaxInclude, yearMinInclude, yearMaxInclude int) validator.Int64 {
	return &validatorCycleCount{
		monthRange: cycleCountRange{
			minInclude: monthMinInclude,
			maxInclude: monthMaxInclude,
		},
		yearRange: cycleCountRange{
			minInclude: yearMinInclude,
			maxInclude: yearMaxInclude,
		},
	}
}

func (v validatorCycleCount) Description(_ context.Context) string {
	return v.getError()
}

func (v validatorCycleCount) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorCycleCount) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	// 没有填写值的情况
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	var cycleType types.String
	response.Diagnostics.Append(request.Config.GetAttribute(ctx, path.Root("cycle_type"), &cycleType)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 没有填写值的情况
	if cycleType.IsNull() || cycleType.IsUnknown() {
		return
	}
	cycleTypeContent := cycleType.ValueString()
	var cr cycleCountRange
	switch cycleTypeContent {
	case business.OrderCycleTypeMonth:
		cr = v.monthRange
	case business.OrderCycleTypeYear:
		cr = v.yearRange
	case business.OrderCycleTypeOnDemand:
		// 按需不需要校验对应的值
		return
	default:
		errMessage := "未知的cycleType取值：" + cycleTypeContent + "可选值：month：按月，year：按年"
		response.Diagnostics.AddError(errMessage, errMessage)
		return
	}
	target := int(request.ConfigValue.ValueInt64())
	if target < cr.minInclude || target > cr.maxInclude {
		errMessage := v.getError()
		response.Diagnostics.AddError(errMessage, errMessage)
	}
}

func (v validatorCycleCount) getError() string {
	return "超出周期范围：按年可选" + strconv.Itoa(v.yearRange.minInclude) + "-" + strconv.Itoa(v.yearRange.maxInclude) + "，按月可选" + strconv.Itoa(v.monthRange.minInclude) + "-" + strconv.Itoa(v.monthRange.maxInclude)
}
