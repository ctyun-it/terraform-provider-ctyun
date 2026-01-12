package explanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func isIgnoredProjectID(v string) bool {
	return v == "" || v == "0"
}

func Project() planmodifier.String {
	return projectIDModifier{}
}

type projectIDModifier struct{}

func (m projectIDModifier) Description(_ context.Context) string {
	return "仅当企业项目ID为非空/非0值且新旧值不一致时触发资源替换；空字符串/0的任何变化均不触发替换"
}

func (m projectIDModifier) MarkdownDescription(_ context.Context) string {
	return "仅当企业项目ID为非空/非0值且新旧值不一致时触发资源替换；空字符串/0的任何变化均不触发替换"
}

func (m projectIDModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() {
		return
	}
	if req.Plan.Raw.IsNull() {
		return
	}

	p, s := req.PlanValue.ValueString(), req.StateValue.ValueString()
	if isIgnoredProjectID(p) && isIgnoredProjectID(s) {
		resp.PlanValue = req.StateValue
		resp.RequiresReplace = false
		return
	}
	resp.RequiresReplace = true
}
