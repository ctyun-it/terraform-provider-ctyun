package planmodifier

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type UseStateIfUnspecified struct{}

// Description 返回计划修改器的描述
func (m UseStateIfUnspecified) Description(ctx context.Context) string {
	return "如果用户未指定值，则使用状态值"
}

// MarkdownDescription 返回计划修改器的 Markdown 描述
func (m UseStateIfUnspecified) MarkdownDescription(ctx context.Context) string {
	return "如果用户未指定值，则使用状态值"
}

// PlanModifyString 实现了计划修改器的逻辑
func (m UseStateIfUnspecified) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// 如果用户未指定值
	if req.PlanValue.IsNull() {
		// 使用状态值
		resp.PlanValue = req.StateValue
	}
}
