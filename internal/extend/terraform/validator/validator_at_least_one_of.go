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

type validatorAtLeastOneOf struct {
	expressions []path.Expression
}

type atLeastOneOfValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type atLeastOneOfValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func AtLeastOneOfString(expressions ...path.Expression) validator.String {
	return &validatorAtLeastOneOf{
		expressions: expressions,
	}
}

func AtLeastOneOfInt64(expressions ...path.Expression) validator.Int64 {
	return &validatorAtLeastOneOf{
		expressions: expressions,
	}
}

func AtLeastOneOfInt32(expressions ...path.Expression) validator.Int32 {
	return &validatorAtLeastOneOf{
		expressions: expressions,
	}
}

func AtLeastOneOfBool(expressions ...path.Expression) validator.Bool {
	return &validatorAtLeastOneOf{
		expressions: expressions,
	}
}

func AtLeastOneOfSet(expressions ...path.Expression) validator.Set {
	return &validatorAtLeastOneOf{
		expressions: expressions,
	}
}

func AtLeastOneOfObject(expressions ...path.Expression) validator.Object {
	return &validatorAtLeastOneOf{
		expressions: expressions,
	}
}

func AtLeastOneOfList(expressions ...path.Expression) validator.List {
	return &validatorAtLeastOneOf{
		expressions: expressions,
	}
}

func (v validatorAtLeastOneOf) Validate(ctx context.Context, req atLeastOneOfValidatorRequest, res *atLeastOneOfValidatorResponse) {
	// 获取当前路径的值是否为空
	currentIsNull := req.ConfigValue.IsNull()

	// 检查所有指定路径的值是否都为空
	allOthersNull := true
	for _, expr := range v.expressions {
		expressions := req.PathExpression.MergeExpressions(expr)
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
					// 如果有任何一个字段是未知的，跳过验证
					return
				}

				if !mpVal.IsNull() {
					allOthersNull = false
				}
			}
		}
	}

	// 如果当前值为空，且所有其他指定字段也都为空，则违反了"至少一个"的要求
	if currentIsNull && allOthersNull {
		// 构建错误消息中的字段列表
		fieldNames := make([]string, 0)
		for _, expr := range v.expressions {
			expressions := req.PathExpression.MergeExpressions(expr)
			for _, expression := range expressions {
				matchedPaths, diags := req.Config.PathMatches(ctx, expression)
				if diags.HasError() {
					continue
				}
				for _, mp := range matchedPaths {
					if !mp.Equal(req.Path) {
						fieldNames = append(fieldNames, mp.String())
					}
				}
			}
		}

		// 添加当前字段
		fieldNames = append(fieldNames, req.Path.String())

		res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			req.Path,
			fmt.Sprintf("必须至少指定以下字段之一: %v", fieldNames),
		))
	}
}

func (v validatorAtLeastOneOf) Description(_ context.Context) string {
	return "至少一个字段必须被指定"
}

func (v validatorAtLeastOneOf) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorAtLeastOneOf) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	validateReq := atLeastOneOfValidatorRequest{
		Config:         request.Config,
		ConfigValue:    request.ConfigValue,
		Path:           request.Path,
		PathExpression: request.PathExpression,
	}
	validateResp := &atLeastOneOfValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	response.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAtLeastOneOf) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := atLeastOneOfValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &atLeastOneOfValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAtLeastOneOf) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	validateReq := atLeastOneOfValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &atLeastOneOfValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAtLeastOneOf) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := atLeastOneOfValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &atLeastOneOfValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAtLeastOneOf) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := atLeastOneOfValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &atLeastOneOfValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAtLeastOneOf) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := atLeastOneOfValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &atLeastOneOfValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v validatorAtLeastOneOf) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := atLeastOneOfValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &atLeastOneOfValidatorResponse{}
	v.Validate(ctx, validateReq, validateResp)
	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
