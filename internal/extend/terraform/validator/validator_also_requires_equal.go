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

type validatorAlsoRequiresEqual struct {
	expression path.Expression
	objs       []attr.Value
}

type alsoRequiresEqualValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type alsoRequiresEqualValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func AlsoRequiresEqualString(expression path.Expression, objs ...attr.Value) validator.String {
	return &validatorAlsoRequiresEqual{
		expression: expression,
		objs:       objs,
	}
}

func AlsoRequiresEqualInt64(expression path.Expression, objs ...attr.Value) validator.Int64 {
	return &validatorAlsoRequiresEqual{
		expression: expression,
		objs:       objs,
	}
}

func AlsoRequiresEqualInt32(expression path.Expression, objs ...attr.Value) validator.Int32 {
	return &validatorAlsoRequiresEqual{
		expression: expression,
		objs:       objs,
	}
}

func AlsoRequiresEqualBool(expression path.Expression, objs ...attr.Value) validator.Bool {
	return &validatorAlsoRequiresEqual{
		expression: expression,
		objs:       objs,
	}
}

func AlsoRequiresEqualSet(expression path.Expression, objs ...attr.Value) validator.Set {
	return &validatorAlsoRequiresEqual{
		expression: expression,
		objs:       objs,
	}
}

func AlsoRequiresEqualObject(expression path.Expression, objs ...attr.Value) validator.Object {
	return &validatorAlsoRequiresEqual{
		expression: expression,
		objs:       objs,
	}
}

func AlsoRequiresEqualList(expression path.Expression, objs ...attr.Value) validator.List {
	return &validatorAlsoRequiresEqual{
		expression: expression,
		objs:       objs,
	}
}

func (v validatorAlsoRequiresEqual) Validate(ctx context.Context, req alsoRequiresEqualValidatorRequest, res *alsoRequiresEqualValidatorResponse) {
	// if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
	// 	return
	// }
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
					if mpVal.Equal(obj) && req.ConfigValue.IsNull() {
						res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
							req.Path,
							fmt.Sprintf("当属性 %q 值为 %s 时， 属性 %q 必须被指定，不能为空", mp, obj.String(), req.Path),
						))
					}
				}
			}
		}
	}
}

func (v validatorAlsoRequiresEqual) Description(_ context.Context) string {
	return "矛盾值出现"
}

func (v validatorAlsoRequiresEqual) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorAlsoRequiresEqual) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	validateReq := alsoRequiresEqualValidatorRequest{
		Config:         request.Config,
		ConfigValue:    request.ConfigValue,
		Path:           request.Path,
		PathExpression: request.PathExpression,
	}
	validateResp := &alsoRequiresEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	response.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAlsoRequiresEqual) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := alsoRequiresEqualValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &alsoRequiresEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAlsoRequiresEqual) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	validateReq := alsoRequiresEqualValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &alsoRequiresEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAlsoRequiresEqual) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := alsoRequiresEqualValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &alsoRequiresEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAlsoRequiresEqual) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := alsoRequiresEqualValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &alsoRequiresEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAlsoRequiresEqual) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := alsoRequiresEqualValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &alsoRequiresEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAlsoRequiresEqual) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := alsoRequiresEqualValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &alsoRequiresEqualValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
