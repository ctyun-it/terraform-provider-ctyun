package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type validatorCrossField struct {
	sourceExpression path.Expression // 源字段
	sourceValues     []attr.Value    // 源字段的触发值列表
	targetValues     []attr.Value    // 目标字段允许的值列表
}

type validatorCrossFieldRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type validatorCrossFieldResponse struct {
	Diagnostics diag.Diagnostics
}

func CrossFieldBool(expression path.Expression, sourceValues, targetValues []attr.Value) validator.Bool {
	return &validatorCrossField{sourceExpression: expression, sourceValues: sourceValues, targetValues: targetValues}
}

func (v validatorCrossField) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := validatorCrossFieldRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &validatorCrossFieldResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorCrossField) Validate(ctx context.Context, req validatorCrossFieldRequest, res *validatorCrossFieldResponse) {
	expressions := req.PathExpression.MergeExpressions(v.sourceExpression)
	for _, expression := range expressions {
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)
		res.Diagnostics.Append(diags...)
		if diags.HasError() {
			continue
		}
		for _, mp := range matchedPaths {
			if mp.Equal(req.Path) {
				continue
			}
			var mpVal attr.Value
			diags := req.Config.GetAttribute(ctx, mp, &mpVal)
			res.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}
			if mpVal.IsUnknown() {
				return
			}
			var targetValid bool
			for _, t := range v.targetValues {
				if req.ConfigValue.Equal(t) {
					targetValid = true
				}
			}

			if !mpVal.IsNull() {
				for _, obj := range v.sourceValues {
					if mpVal.Equal(obj) && !targetValid {
						res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
							req.Path,
							fmt.Sprintf("当属性 %q 值为 %s 时， 属性 %q 取值范围为 %v", mp, obj.String(), req.Path, v.targetValues),
						))
					}
				}
			}
		}
	}
}

func (v validatorCrossField) Description(_ context.Context) string {
	return "矛盾值出现"
}

func (v validatorCrossField) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
