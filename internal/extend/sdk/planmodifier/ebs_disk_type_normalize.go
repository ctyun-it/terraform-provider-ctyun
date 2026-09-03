package planmodifier

import (
	"context"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ebsDiskTypeNormalize 将用户输入的磁盘类型统一映射为 API 侧格式。
// 例如 "sata" → "SATA"，"ssd-genric" → "SSD-genric"，"x-entry" → "X-Entry"。
// 如果用户输入已经是 API 侧格式，则保持不变。
type ebsDiskTypeNormalize struct {
	fromOriginalScene func(source any, targetScene utils.Scene) (any, error)
	scene             utils.Scene
}

// EbsDiskTypeNormalize 返回一个 String PlanModifier，在 Plan 阶段将用户输入的磁盘类型
// 通过 fromOriginalScene 映射为 API 侧格式。映射失败时（说明输入已经是 API 侧格式）保持原值。
func EbsDiskTypeNormalize(fromOriginalScene func(source any, targetScene utils.Scene) (any, error), scene utils.Scene) planmodifier.String {
	return &ebsDiskTypeNormalize{
		fromOriginalScene: fromOriginalScene,
		scene:             scene,
	}
}

func (d ebsDiskTypeNormalize) Description(_ context.Context) string {
	return "将磁盘类型统一映射为 API 侧格式"
}

func (d ebsDiskTypeNormalize) MarkdownDescription(_ context.Context) string {
	return d.Description(nil)
}

func (d ebsDiskTypeNormalize) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	mapped, err := d.fromOriginalScene(req.PlanValue.ValueString(), d.scene)
	if err != nil {
		// 映射失败，说明输入已经是 API 侧格式，保持原值
		return
	}
	resp.PlanValue = types.StringValue(mapped.(string))
}
