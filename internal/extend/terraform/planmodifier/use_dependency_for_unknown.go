package planmodifier

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type useDependency struct {
	dependency path.Path
}

func UseDependencyForUnknown(dependency path.Path) planmodifier.String {
	return &useDependency{dependency: dependency}
}

func (c useDependency) Description(ctx context.Context) string {
	return "使用依赖字段的值"
}

func (c useDependency) MarkdownDescription(ctx context.Context) string {
	return "使用依赖字段的值"
}

func (c useDependency) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}

	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.ConfigValue.IsUnknown() {
		return
	}

	var dependencyValue types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, c.dependency, &dependencyValue)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = dependencyValue
}
