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
	var key string
	switch d.ctyunMetadataExtraKey {
	case common.ExtraRegionId:
		key = "region_id"
	case common.ExtraAzName:
		key = "az_name"
	case common.ExtraProjectId:
		key = "project_id"
	default:
		panic("invalid extra key")
	}
	if value == "" && d.mustAcquire {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			fmt.Sprintf("%s值为空，且ctyun_provider与环境变量均未设置", key),
			fmt.Sprintf("%s值为空，且ctyun_provider与环境变量均未设置", key))
		return
	}
	resp.PlanValue = types.StringValue(value)
}
