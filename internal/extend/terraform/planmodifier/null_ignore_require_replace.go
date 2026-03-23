package planmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func NullIgnoreString() planmodifier.String {
	return nullModifierString{}
}

type nullModifierString struct{}

func (m nullModifierString) Description(_ context.Context) string {
	return "旧值为null，忽略更新"
}

func (m nullModifierString) MarkdownDescription(_ context.Context) string {
	return "旧值为null，忽略更新"
}

func (m nullModifierString) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() {
		return
	}
	if req.Plan.Raw.IsNull() {
		return
	}
	if req.ConfigValue.IsNull() {
		return
	}
	if req.StateValue.IsNull() {
		resp.RequiresReplace = false
		return
	}
	resp.RequiresReplace = true
}

func NullIgnoreBool() planmodifier.Bool {
	return nullModifierBool{}
}

type nullModifierBool struct{}

func (m nullModifierBool) Description(_ context.Context) string {
	return "旧值为null，忽略更新"
}

func (m nullModifierBool) MarkdownDescription(_ context.Context) string {
	return "旧值为null，忽略更新"
}

func (m nullModifierBool) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.State.Raw.IsNull() {
		return
	}
	if req.Plan.Raw.IsNull() {
		return
	}
	if req.ConfigValue.IsNull() {
		return
	}
	if req.StateValue.IsNull() {
		resp.RequiresReplace = false
		return
	}
	resp.RequiresReplace = true
}

func NullIgnoreInt32() planmodifier.Int32 {
	return nullModifierInt32{}
}

type nullModifierInt32 struct{}

func (m nullModifierInt32) Description(_ context.Context) string {
	return "旧值为null，忽略更新"
}

func (m nullModifierInt32) MarkdownDescription(_ context.Context) string {
	return "旧值为null，忽略更新"
}

func (m nullModifierInt32) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
	if req.State.Raw.IsNull() {
		return
	}
	if req.Plan.Raw.IsNull() {
		return
	}
	if req.ConfigValue.IsNull() {
		return
	}
	if req.StateValue.IsNull() {
		resp.RequiresReplace = false
		return
	}
	resp.RequiresReplace = true
}

func NullIgnoreSet() planmodifier.Set {
	return nullModifierSet{}
}

type nullModifierSet struct{}

func (m nullModifierSet) Description(_ context.Context) string {
	return "旧值为null，忽略更新"
}

func (m nullModifierSet) MarkdownDescription(_ context.Context) string {
	return "旧值为null，忽略更新"
}

func (m nullModifierSet) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.State.Raw.IsNull() {
		return
	}
	if req.Plan.Raw.IsNull() {
		return
	}
	if req.ConfigValue.IsNull() {
		return
	}
	if req.StateValue.IsNull() {
		resp.RequiresReplace = false
		return
	}
	resp.RequiresReplace = true
}
