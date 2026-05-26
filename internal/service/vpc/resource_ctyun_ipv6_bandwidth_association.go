package vpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"

	"strings"

	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ctyunIPv6BandwidthAssociation{}
	_ resource.ResourceWithConfigure   = &ctyunIPv6BandwidthAssociation{}
	_ resource.ResourceWithImportState = &ctyunIPv6BandwidthAssociation{}
)

type ctyunIPv6BandwidthAssociation struct {
	meta                 *common.CtyunMetadata
	name                 string
	ipv6Service          *business.IPv6Service
	ipv6BandwidthService *business.IPv6BandwidthService
}

func NewCtyunIPv6BandwidthAssociation() resource.Resource {
	return &ctyunIPv6BandwidthAssociation{}
}

func (c *ctyunIPv6BandwidthAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ipv6_bandwidth_association"
	c.name = response.TypeName
}

type CtyunIPv6BandwidAssociationConfig struct {
	ID              types.String `tfsdk:"id"`
	IPv6BandwidthId types.String `tfsdk:"ipv6_bandwidth_id"`
	IPv6            types.String `tfsdk:"ipv6"`
	RegionId        types.String `tfsdk:"region_id"`
}

func (c *ctyunIPv6BandwidthAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理共享带宽和弹性IP的绑定关系", "共享流量包（SDP，Shared Data Package）", "https://www.ctyun.cn/document/10026761/10030030"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "id",
			},
			"ipv6_bandwidth_id": schema.StringAttribute{
				Required:    true,
				Description: "IPv6带宽id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"ipv6": schema.StringAttribute{
				Required:    true,
				Description: "IPv6 地址",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
		},
	}
}

func (c *ctyunIPv6BandwidthAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunIPv6BandwidAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 校验ipv6必须存在
	err := c.ipv6Service.MustExist(ctx, plan.IPv6.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	// 校验ipv6带宽必须存在
	err = c.ipv6BandwidthService.MustExist(ctx, plan.IPv6BandwidthId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	regionId := plan.RegionId.ValueString()
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcAddIPv6ToIPv6BandwidthApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcAddIPv6ToIPv6BandwidthRequest{
		RegionID:    regionId,
		BandwidthID: plan.IPv6BandwidthId.ValueString(),
		Ip:          plan.IPv6.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	err = c.getAndMergeIPv6BandwidthAssociationIPv6(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunIPv6BandwidthAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunIPv6BandwidAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.getAndMergeIPv6BandwidthAssociationIPv6(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *ctyunIPv6BandwidthAssociation) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {

}

func (c *ctyunIPv6BandwidthAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunIPv6BandwidAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcRemoveIPv6FromIPv6BandwidthApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcRemoveIPv6FromIPv6BandwidthRequest{
		RegionID:    state.RegionId.ValueString(),
		BandwidthID: state.IPv6BandwidthId.ValueString(),
		Ip:          state.IPv6.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

func (c *ctyunIPv6BandwidthAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [ipv6_bandwidth_id],[ipv6],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunIPv6BandwidAssociationConfig
	var bandwidthID, ipv6, regionID string
	cnt := strings.Count(request.ID, common.ImportSeparator)
	switch cnt {
	case 0:
		err = fmt.Errorf("bandwidth_id和eip_id必须输入")
		return
	case 1:
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &bandwidthID, &ipv6)
		if err != nil {
			return
		}
	default:
		err = terraform_extend.Split(request.ID, &bandwidthID, &ipv6, &regionID)
		if err != nil {
			return
		}
	}
	if bandwidthID == "" {
		err = fmt.Errorf("bandwidth_id不能为空")
		return
	}
	if ipv6 == "" {
		err = fmt.Errorf("eip_id不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}

	cfg.IPv6BandwidthId = types.StringValue(bandwidthID)
	cfg.IPv6 = types.StringValue(ipv6)
	cfg.RegionId = types.StringValue(regionID)
	err = c.getAndMergeIPv6BandwidthAssociationIPv6(ctx, &cfg)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunIPv6BandwidthAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ipv6Service = business.NewIPv6Service(meta)
	c.ipv6BandwidthService = business.NewIPv6BandwidthService(meta)
}

// getAndMergeIPv6BandwidthAssociationIPv6 查询绑定关系
func (c *ctyunIPv6BandwidthAssociation) getAndMergeIPv6BandwidthAssociationIPv6(ctx context.Context, cfg *CtyunIPv6BandwidAssociationConfig) error {
	result, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowIPv6BandwidthApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowIPv6BandwidthRequest{
		RegionID:    cfg.RegionId.ValueString(),
		BandwidthID: cfg.IPv6BandwidthId.ValueString(),
	})
	if err != nil {
		return err
	} else if utils.SecString(result.ErrorCode) == common.OpenapiIPv6BandwidthNotFound {
		return fmt.Errorf("带宽 %s 不存在", cfg.IPv6BandwidthId.ValueString())
	} else if result.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	if utils.SecString(result.ReturnObj.IpAddress) != cfg.IPv6.ValueString() {
		return common.ResourceNotExistError
	}

	cfg.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", cfg.IPv6BandwidthId.ValueString(), cfg.IPv6.ValueString(), cfg.RegionId.ValueString()))
	return nil
}
