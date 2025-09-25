package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ScalingPolicyDayError      = "day取值有误，当cycle=monthly时，取值范围为[1,31]，当cycle=weekly时，取值范围为[1,7]"
	ScalingPolicyDayDailyError = "day取值有误，当cycle=daily时，无需填写"
	ScalingPolicyDayNullError  = "day取值有误，当cycle=monthly和weekly时，day不能为空"
)

type validatorScalingPolicyDay struct {
}

func (v validatorScalingPolicyDay) ValidateSet(ctx context.Context, request validator.SetRequest, response *validator.SetResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	config := request.Config
	var cycle types.String
	var day types.Set
	// 先获取cycle，确认是否有值
	diags := config.GetAttribute(context.Background(), path.Root("cycle"), &cycle)
	if diags.HasError() {
		response.Diagnostics.AddError(ScalingPolicyDayError, diags[0].Detail())
		return
	}
	// 获取day，并解析成[]int32
	diags = config.GetAttribute(context.Background(), path.Root("day"), &day)
	var days []int32
	diags = day.ElementsAs(ctx, &days, true)
	if diags.HasError() {
		response.Diagnostics.AddError(ScalingPolicyDayError, diags[0].Detail())
		return
	}
	// 如果cycle为空，不做任何操作直接返回。因为伸缩策略可能是其他类型
	if cycle.IsNull() || cycle.IsUnknown() {
		return
	}
	// 当cycle = monthly，验证是否为【1，31】
	if cycle.ValueString() == "monthly" {
		if days == nil || len(days) == 0 {
			response.Diagnostics.AddError(ScalingPolicyDayNullError, ScalingPolicyDayNullError)
			return
		}
		for _, date := range days {
			if date < 1 || date > 31 {
				response.Diagnostics.AddError(ScalingPolicyDayError, ScalingPolicyDayError)
				return
			}
		}
	} else if cycle.ValueString() == "weekly" {
		if days == nil || len(days) == 0 {
			response.Diagnostics.AddError(ScalingPolicyDayNullError, ScalingPolicyDayNullError)
			return
		}
		for _, date := range days {
			if date < 1 || date > 7 {
				response.Diagnostics.AddError(ScalingPolicyDayError, ScalingPolicyDayError)
				return
			}
		}
	} else {
		// 当cycle = daily，无需填写
		if !day.IsNull() && !day.IsUnknown() {
			response.Diagnostics.AddError(ScalingPolicyDayDailyError, ScalingPolicyDayDailyError)
			return
		}
	}

}

func (v validatorScalingPolicyDay) Description(ctx context.Context) string {
	return ScalingPolicyDayError
}

func (v validatorScalingPolicyDay) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func ScalingPolicyDayValidate() validator.Set {
	return validatorScalingPolicyDay{}
}
