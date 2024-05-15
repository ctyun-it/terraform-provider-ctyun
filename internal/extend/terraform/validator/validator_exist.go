package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	existError = "不存在指定对象，请确认此对象值已经存在"
)

type CheckFunc[T any] func(ctx context.Context, obj T) bool

func StringExist(f CheckFunc[types.String]) validator.String {
	return &validatorStringExist{
		checkFunc: f,
	}
}

func SetExist(f CheckFunc[types.Set]) validator.Set {
	return &validatorSetExist{
		checkFunc: f,
	}
}

type validatorStringExist struct {
	checkFunc CheckFunc[types.String]
}

func (v validatorStringExist) Description(_ context.Context) string {
	return existError
}

func (v validatorStringExist) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorStringExist) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if !v.checkFunc(ctx, request.ConfigValue) {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			request.Path,
			existError,
			existError,
		))
		return
	}
}

type validatorSetExist struct {
	checkFunc CheckFunc[types.Set]
}

func (v validatorSetExist) Description(_ context.Context) string {
	return existError
}

func (v validatorSetExist) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorSetExist) ValidateSet(ctx context.Context, request validator.SetRequest, response *validator.SetResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if !v.checkFunc(ctx, request.ConfigValue) {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			request.Path,
			existError,
			existError,
		))
		return
	}
}
