package planmodifier

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplaceIfStateNotNull 是一个计划修改器，根据状态和计划值的情况决定是否触发替换
// 如果 state 状态为null，plan 为一个值，这个时候不触发 replace 而是继续触发 update
// 如果 state 状态不为null，plan 为一个值 与state不一致，这个时候触发 replace
type RequiresReplaceIfStateNotNull struct{}

// RequiresReplaceIfStateNotNullModifier 返回一个RequiresReplaceIfStateNotNull计划修改器实例
func RequiresReplaceIfStateNotNullModifier() RequiresReplaceIfStateNotNull {
	return RequiresReplaceIfStateNotNull{}
}

// Description 返回计划修改器的描述
func (m RequiresReplaceIfStateNotNull) Description(ctx context.Context) string {
	return "当状态值不为null且与计划值不同时触发替换，否则继续更新"
}

// MarkdownDescription 返回计划修改器的 Markdown 描述
func (m RequiresReplaceIfStateNotNull) MarkdownDescription(ctx context.Context) string {
	return "当状态值不为null且与计划值不同时触发替换，否则继续更新"
}

// PlanModifyString 实现了字符串类型的计划修改器逻辑
func (m RequiresReplaceIfStateNotNull) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// 如果状态值为null，不触发替换，继续更新
	if req.StateValue.IsNull() {
		return
	}

	// 如果状态值不为null，且计划值与状态值不同，则触发替换
	if !req.StateValue.IsNull() && !req.PlanValue.IsNull() && req.StateValue.ValueString() != req.PlanValue.ValueString() {
		resp.RequiresReplace = true
		return
	}

	// 如果状态值不为null，但计划值为null，则触发替换
	if !req.StateValue.IsNull() && req.PlanValue.IsNull() {
		resp.RequiresReplace = true
		return
	}
}

// PlanModifyBool 实现了布尔类型的计划修改器逻辑
func (m RequiresReplaceIfStateNotNull) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// 如果状态值为null，不触发替换，继续更新
	if req.StateValue.IsNull() {
		return
	}

	// 如果状态值不为null，且计划值与状态值不同，则触发替换
	if !req.StateValue.IsNull() && !req.PlanValue.IsNull() && req.StateValue.ValueBool() != req.PlanValue.ValueBool() {
		resp.RequiresReplace = true
		return
	}

	// 如果状态值不为null，但计划值为null，则触发替换
	if !req.StateValue.IsNull() && req.PlanValue.IsNull() {
		resp.RequiresReplace = true
		return
	}
}

// PlanModifyInt64 实现了Int64类型的计划修改器逻辑
func (m RequiresReplaceIfStateNotNull) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// 如果状态值为null，不触发替换，继续更新
	if req.StateValue.IsNull() {
		return
	}

	// 如果状态值不为null，且计划值与状态值不同，则触发替换
	if !req.StateValue.IsNull() && !req.PlanValue.IsNull() && req.StateValue.ValueInt64() != req.PlanValue.ValueInt64() {
		resp.RequiresReplace = true
		return
	}

	// 如果状态值不为null，但计划值为null，则触发替换
	if !req.StateValue.IsNull() && req.PlanValue.IsNull() {
		resp.RequiresReplace = true
		return
	}
}

// PlanModifyFloat64 实现了Float64类型的计划修改器逻辑
func (m RequiresReplaceIfStateNotNull) PlanModifyFloat64(ctx context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	// 如果状态值为null，不触发替换，继续更新
	if req.StateValue.IsNull() {
		return
	}

	// 如果状态值不为null，且计划值与状态值不同，则触发替换
	if !req.StateValue.IsNull() && !req.PlanValue.IsNull() && req.StateValue.ValueFloat64() != req.PlanValue.ValueFloat64() {
		resp.RequiresReplace = true
		return
	}

	// 如果状态值不为null，但计划值为null，则触发替换
	if !req.StateValue.IsNull() && req.PlanValue.IsNull() {
		resp.RequiresReplace = true
		return
	}
}

// PlanModifyNumber 实现了Number类型的计划修改器逻辑
func (m RequiresReplaceIfStateNotNull) PlanModifyNumber(ctx context.Context, req planmodifier.NumberRequest, resp *planmodifier.NumberResponse) {
	// 如果状态值为null，不触发替换，继续更新
	if req.StateValue.IsNull() {
		return
	}

	// 如果状态值不为null，且计划值与状态值不同，则触发替换
	if !req.StateValue.IsNull() && !req.PlanValue.IsNull() {
		// 比较数值是否相等
		if !req.StateValue.Equal(req.PlanValue) {
			resp.RequiresReplace = true
			return
		}
	}

	// 如果状态值不为null，但计划值为null，则触发替换
	if !req.StateValue.IsNull() && req.PlanValue.IsNull() {
		resp.RequiresReplace = true
		return
	}
}
