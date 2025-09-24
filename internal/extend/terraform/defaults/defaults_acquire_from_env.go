package defaults

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AcquireFromGlobalString(ctyunMetadataExtraKey string, mustAcquire bool) defaults.String {
	return globalStringDefault{
		ctyunMetadataExtraKey: ctyunMetadataExtraKey,
		mustAcquire:           mustAcquire,
	}
}

type globalStringDefault struct {
	ctyunMetadataExtraKey string
	mustAcquire           bool
}

func (d globalStringDefault) Description(_ context.Context) string {
	return fmt.Sprintf("当此值为空时，默认取自ctyun_provider中配置或系统环境变量")
}

func (d globalStringDefault) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("当此值为空时，默认取自ctyun_provider中配置或系统环境变量")
}

func (d globalStringDefault) DefaultString(_ context.Context, req defaults.StringRequest, resp *defaults.StringResponse) {
	metadata := common.AcquireCtyunMetadata()
	value := metadata.GetExtra(d.ctyunMetadataExtraKey)
	if value == "" && d.mustAcquire {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"当前值为空，ctyun_provider未配置且环境变量未设置",
			"当前值为空，ctyun_provider未配置且环境变量未设置")
		return
	}
	resp.PlanValue = types.StringValue(value)
}
