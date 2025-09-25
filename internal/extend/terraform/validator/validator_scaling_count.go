package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ScalingCountError  = "弹性伸缩组实例个数有误"
	minCountError      = "min_count输入有误，应当小于等于max_count"
	expectedCountError = "expected_count输入有误，取值区间为[min_count, max_count]"
)

type validatorScalingCount struct {
}

func (v validatorScalingCount) Description(ctx context.Context) string {
	return ScalingCountError
}

func (v validatorScalingCount) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorScalingCount) ValidateInt32(ctx context.Context, request validator.Int32Request, response *validator.Int32Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	config := request.Config
	maxCount, err := geInt32Attribute(config, path.Root("max_count"))
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	minCount, err := geInt32Attribute(config, path.Root("min_count"))
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	expectedCount, err := geInt32Attribute(config, path.Root("expected_count"))
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// minCount <= maxCount。 expectedCount >= minCount && expectedCount <= maxCount
	if minCount.ValueInt32() > maxCount.ValueInt32() {
		response.Diagnostics.AddError(minCountError, minCountError)
		return
	}
	if expectedCount.IsNull() || expectedCount.IsUnknown() {
		return
	}
	if expectedCount.ValueInt32() < minCount.ValueInt32() || expectedCount.ValueInt32() > maxCount.ValueInt32() {
		response.Diagnostics.AddError(expectedCountError, expectedCountError)
		return
	}
}

func geInt32Attribute(config tfsdk.Config, attrPath path.Path) (types.Int32, error) {
	var value types.Int32
	diags := config.GetAttribute(context.Background(), attrPath, &value)
	if diags.HasError() {
		return types.Int32Value(-1), fmt.Errorf("failed to get attribute: %v", diags.Errors())
	}
	return value, nil
}

func ScalingCountValidate() validator.Int32 {
	return validatorScalingCount{}
}
