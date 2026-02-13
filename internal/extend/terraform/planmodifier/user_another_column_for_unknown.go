package planmodifier

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type useAnother struct {
	another path.Path
}

func UseAnotherColumnForUnknown(another path.Path) planmodifier.String {
	return &useAnother{another: another}
}

func (c useAnother) Description(ctx context.Context) string {
	return fmt.Sprintf("使用%s的值", c.another.String())
}

func (c useAnother) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("使用%s的值", c.another.String())
}

func (c useAnother) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}

	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.ConfigValue.IsUnknown() {
		return
	}

	var anotherValue types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, c.another, &anotherValue)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = anotherValue
}
