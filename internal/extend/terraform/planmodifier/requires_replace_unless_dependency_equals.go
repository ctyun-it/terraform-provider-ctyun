package planmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func RequiresReplaceUnlessDependencyEqualsBool(dependencyPath path.Expression, ignoreValues ...attr.Value) planmodifier.Bool {
	return &requiresReplaceUnlessDependencyEqualsModifier{
		dependencyPath: dependencyPath,
		ignoreValues:   ignoreValues,
	}
}

func RequiresReplaceUnlessDependencyEqualsInt64(dependencyPath path.Expression, ignoreValues ...attr.Value) planmodifier.Int64 {
	return &requiresReplaceUnlessDependencyEqualsModifier{
		dependencyPath: dependencyPath,
		ignoreValues:   ignoreValues,
	}
}

// RequiresReplaceUnlessDependencyEqualsString returns a plan modifier for String attributes.
func RequiresReplaceUnlessDependencyEqualsString(dependencyPath path.Expression, ignoreValues ...attr.Value) planmodifier.String {
	return &requiresReplaceUnlessDependencyEqualsModifier{
		dependencyPath: dependencyPath,
		ignoreValues:   ignoreValues,
	}
}

// RequiresReplaceUnlessDependencyEqualsInt32 returns a plan modifier for Int32 attributes.
func RequiresReplaceUnlessDependencyEqualsInt32(dependencyPath path.Expression, ignoreValues ...attr.Value) planmodifier.Int32 {
	return &requiresReplaceUnlessDependencyEqualsModifier{
		dependencyPath: dependencyPath,
		ignoreValues:   ignoreValues,
	}
}

type requiresReplaceUnlessDependencyEqualsModifier struct {
	dependencyPath path.Expression
	ignoreValues   []attr.Value
}

func (m *requiresReplaceUnlessDependencyEqualsModifier) Description(_ context.Context) string {
	return fmt.Sprintf("当依赖属性 %q 的值为指定值时，本属性变更不触发重建", m.dependencyPath.String())
}

func (m *requiresReplaceUnlessDependencyEqualsModifier) MarkdownDescription(_ context.Context) string {
	return m.Description(context.Background())
}

// checkDependencyEquals reads the dependency attribute from config and returns
// true if the dependency value matches any of the ignoreValues (should NOT replace).
// If the dependency is null/unknown or does not match, it returns false (should replace).
func (m *requiresReplaceUnlessDependencyEqualsModifier) checkDependencyEquals(ctx context.Context, config tfsdk.Config, currentPath path.Path, currentPathExpr path.Expression, diags *diag.Diagnostics) bool {
	expressions := currentPathExpr.MergeExpressions(m.dependencyPath)
	for _, expression := range expressions {
		matchedPaths, d := config.PathMatches(ctx, expression)
		diags.Append(d...)
		if d.HasError() {
			continue
		}
		for _, mp := range matchedPaths {
			if mp.Equal(currentPath) {
				continue
			}
			var depVal attr.Value
			d := config.GetAttribute(ctx, mp, &depVal)
			diags.Append(d...)
			if d.HasError() {
				continue
			}
			if depVal.IsNull() || depVal.IsUnknown() {
				return false
			}
			for _, ignoreVal := range m.ignoreValues {
				if depVal.Equal(ignoreVal) {
					return true
				}
			}
		}
	}
	return false
}

func (m *requiresReplaceUnlessDependencyEqualsModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || req.PlanValue.Equal(req.StateValue) {
		return
	}
	if req.StateValue.IsNull() {
		return
	}
	if m.checkDependencyEquals(ctx, req.Config, req.Path, req.PathExpression, &resp.Diagnostics) {
		return
	}
	resp.RequiresReplace = true
}

func (m *requiresReplaceUnlessDependencyEqualsModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || req.PlanValue.Equal(req.StateValue) {
		return
	}
	if req.StateValue.IsNull() {
		return
	}
	if m.checkDependencyEquals(ctx, req.Config, req.Path, req.PathExpression, &resp.Diagnostics) {
		return
	}
	resp.RequiresReplace = true
}

func (m *requiresReplaceUnlessDependencyEqualsModifier) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || req.PlanValue.Equal(req.StateValue) {
		return
	}
	if req.StateValue.IsNull() {
		return
	}
	if m.checkDependencyEquals(ctx, req.Config, req.Path, req.PathExpression, &resp.Diagnostics) {
		return
	}
	resp.RequiresReplace = true
}

func (m *requiresReplaceUnlessDependencyEqualsModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || req.PlanValue.Equal(req.StateValue) {
		return
	}
	if req.StateValue.IsNull() {
		return
	}
	if m.checkDependencyEquals(ctx, req.Config, req.Path, req.PathExpression, &resp.Diagnostics) {
		return
	}
	resp.RequiresReplace = true
}