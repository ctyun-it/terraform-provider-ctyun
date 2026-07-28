package planmodifier

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// 适用字段A需要满足以下条件：
// A字段为Optional+Computed
// A字段在B字段=1时必填，在B字段=2时不可填
// A字段在B字段=2时，state写入值为null

type useStateNull struct {
	dependency path.Path
}

func (c useStateNull) Description(ctx context.Context) string {
	return "本字段未指定情况下，当依赖字段变化时，将本字段设置为null。当依赖资源未变化时，则本字段使用state中的值"
}

func (c useStateNull) MarkdownDescription(ctx context.Context) string {
	return "本字段未指定情况下，当依赖字段变化时，将本字段设置为null。当依赖资源未变化时，则本字段使用state中的值"
}

func SetStateNullIfDependencyChangeString(dependency path.Path) planmodifier.String {
	return &useStateNull{dependency: dependency}
}

func (c useStateNull) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
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
	// 依赖值没变
	if dependencyPlan.Equal(dependencyState) {
		resp.PlanValue = req.StateValue
	} else {
		resp.PlanValue = types.StringNull()
	}
}

func SetStateNullIfDependencyChangeInt64(dependency path.Path) planmodifier.Int64 {
	return &useStateNull{dependency: dependency}
}

func (c useStateNull) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// Do nothing if there is a known planned value.
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
	// 依赖值没变
	if dependencyPlan.Equal(dependencyState) {
		resp.PlanValue = req.StateValue
	} else {
		resp.PlanValue = types.Int64Null()
	}
}

func SetStateNullIfDependencyChangeInt32(dependency path.Path) planmodifier.Int32 {
	return &useStateNull{dependency: dependency}
}

func (c useStateNull) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
	// Do nothing if there is a known planned value.
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
	// 依赖值没变
	if dependencyPlan.Equal(dependencyState) {
		resp.PlanValue = req.StateValue
	} else {
		resp.PlanValue = types.Int32Null()
	}
}
