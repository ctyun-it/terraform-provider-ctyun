package planmodifier

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	EcsStatusError = "云主机配置更新时，status需要为stopped"
)

type checkValueWhenChange struct {
	targetPath    path.Path
	expectedValue string
}

func (c checkValueWhenChange) Description(ctx context.Context) string {
	return EcsStatusError
}

func (c checkValueWhenChange) MarkdownDescription(ctx context.Context) string {
	return EcsStatusError
}

func (c checkValueWhenChange) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	// 创建时无需校验
	if request.State.Raw.IsNull() && !request.Plan.Raw.IsNull() {
		return
	}
	// 删除时无需校验
	if request.Plan.Raw.IsNull() && !request.State.Raw.IsNull() {
		return
	}
	var planExpectedValue types.String
	// 获取对比值
	response.Diagnostics.Append(request.Plan.GetAttribute(ctx, c.targetPath, &planExpectedValue)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 表示发生变配
	if !request.PlanValue.IsNull() && !request.PlanValue.IsUnknown() {
		if !request.PlanValue.Equal(request.StateValue) && planExpectedValue.ValueString() != c.expectedValue {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"参数校验失败",
				fmt.Sprintf("%s当发生更新时，%s预期值为%s", request.Path.String(), c.targetPath.String(), c.expectedValue),
			)
		}
	}
}

func (c checkValueWhenChange) PlanModifyInt64(ctx context.Context, request planmodifier.Int64Request, response *planmodifier.Int64Response) {
	// 创建时无需校验
	if request.State.Raw.IsNull() && !request.Plan.Raw.IsNull() {
		return
	}
	// 删除时无需校验
	if request.Plan.Raw.IsNull() && !request.State.Raw.IsNull() {
		return
	}
	var planExpectedValue types.String
	response.Diagnostics.Append(request.Plan.GetAttribute(ctx, c.targetPath, &planExpectedValue)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 表示发生变配
	if !request.PlanValue.IsNull() && !request.PlanValue.IsUnknown() {
		if !request.PlanValue.Equal(request.StateValue) && planExpectedValue.ValueString() != c.expectedValue {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"参数校验失败",
				fmt.Sprintf("%s当发生更新时，%s预期值为%s", request.Path.String(), c.targetPath.String(), c.expectedValue),
			)
		}
	}
}

func CheckValueWhenChangeString(targetPath path.Path, expectedValue string) planmodifier.String {
	return &checkValueWhenChange{targetPath: targetPath, expectedValue: expectedValue}
}

func CheckValueWhenChangeInt64(targetPath path.Path, expectedValue string) planmodifier.Int64 {
	return &checkValueWhenChange{targetPath: targetPath, expectedValue: expectedValue}
}
