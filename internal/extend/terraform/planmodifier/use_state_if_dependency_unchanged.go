package planmodifier

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// 适用字段A需要满足以下条件：
// A字段为Optional+Computed或仅Computed
// A字段可以通过修改A字段本身或修改B字段而产生变化

type useState struct {
	dependency path.Path
}

func UseStringStateIfDependencyUnchanged(dependency path.Path) planmodifier.String {
	return &useState{dependency: dependency}
}

func (c useState) Description(ctx context.Context) string {
	return "依赖字段没有变化时，使用state中的值"
}

func (c useState) MarkdownDescription(ctx context.Context) string {
	return "依赖字段没有变化时，使用state中的值"
}

func (c useState) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}

	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.ConfigValue.IsUnknown() {
		return
	}

	var dependencyPlan attr.Value
	var dependencyState attr.Value
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, c.dependency, &dependencyPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, c.dependency, &dependencyState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if dependencyPlan != nil && dependencyState != nil && dependencyPlan.Equal(dependencyState) {
		resp.PlanValue = req.StateValue
	}
}

func UseListStateIfDependencyUnchanged(dependency path.Path) planmodifier.List {
	return &useState{dependency: dependency}
}
func (c useState) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.StateValue.IsNull() {
		return
	}

	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.ConfigValue.IsUnknown() {
		return
	}

	var dependencyPlan attr.Value
	var dependencyState attr.Value
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, c.dependency, &dependencyPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, c.dependency, &dependencyState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if dependencyPlan != nil && dependencyState != nil && dependencyPlan.Equal(dependencyState) {
		resp.PlanValue = req.StateValue
	}
}
