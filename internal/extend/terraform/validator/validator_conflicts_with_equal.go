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

type validatorConflictsWithEqual struct {
	expression path.Expression
	objs       []attr.Value
}

type conflictsWithEqualValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type conflictsWithEqualValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func ConflictsWithEqualStrings(expression path.Expression, objs ...attr.Value) validator.String {
	return &validatorConflictsWithEqual{
		expression: expression,
		objs:       objs,
	}
}

func ConflictsWithEqualInt64(expression path.Expression, objs ...attr.Value) validator.Int64 {
	return &validatorConflictsWithEqual{
		expression: expression,
		objs:       objs,
	}
}

func ConflictsWithEqualBool(expression path.Expression, objs ...attr.Value) validator.Bool {
	return &validatorConflictsWithEqual{
		expression: expression,
		objs:       objs,
	}
}

func (v validatorConflictsWithEqual) Validate(ctx context.Context, req conflictsWithEqualValidatorRequest, res *conflictsWithEqualValidatorResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	expressions := req.PathExpression.MergeExpressions(v.expression)
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
			if !mpVal.IsNull() {
				for _, obj := range v.objs {
					if mpVal.Equal(obj) {
						res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
							req.Path,
							fmt.Sprintf("当属性 %q 值为 %s 时， 属性 %q 不能被指定，必须保持为空", mp, obj.String(), req.Path),
						))
					}
				}
			}
		}
	}
}

func (v validatorConflictsWithEqual) Description(_ context.Context) string {
	return "矛盾值出现"
}

func (v validatorConflictsWithEqual) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorConflictsWithEqual) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	validateReq := conflictsWithEqualValidatorRequest{
		Config:         request.Config,
		ConfigValue:    request.ConfigValue,
		Path:           request.Path,
		PathExpression: request.PathExpression,
	}
	validateResp := &conflictsWithEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	response.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorConflictsWithEqual) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := conflictsWithEqualValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &conflictsWithEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorConflictsWithEqual) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := conflictsWithEqualValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &conflictsWithEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
