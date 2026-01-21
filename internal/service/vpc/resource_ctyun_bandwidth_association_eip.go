package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	explanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/planmodifier"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &ctyunBandwidthAssociationEip{}
	_ resource.ResourceWithConfigure   = &ctyunBandwidthAssociationEip{}
	_ resource.ResourceWithImportState = &ctyunBandwidthAssociationEip{}
)

type ctyunBandwidthAssociationEip struct {
	meta             *common.CtyunMetadata
	name             string
	bandwidthService *business.BandwidthService
	eipService       *business.EipService
}

func NewCtyunBandwidthAssociationEip() resource.Resource {
	return &ctyunBandwidthAssociationEip{}
}

func (c *ctyunBandwidthAssociationEip) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_bandwidth_association_eip"
	c.name = response.TypeName
}

func (c *ctyunBandwidthAssociationEip) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("BANDWIDTH", "https://www.ctyun.cn/document/10026761/10030030"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "id",
			},
			"bandwidth_id": schema.StringAttribute{
				Required:    true,
				Description: "共享带宽id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"eip_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性ip的id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.EipValidate(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:           true,
				Computed:           true,
				DeprecationMessage: "废弃字段，请不要指定",
				Description:        "企业项目ID",
				PlanModifiers: []planmodifier.String{
					explanmodifier.Project(),
				},
				Validators: []validator.String{
					validator2.Project(),
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

func (c *ctyunBandwidthAssociationEip) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunBandwidAssociationEipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 校验带宽必须存在
	err := c.bandwidthService.MustExist(ctx, plan.BandwidthId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	// 校验弹性ip必须存在
	err = c.eipService.MustExist(ctx, plan.EipId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	regionId := plan.RegionId.ValueString()
	_, err = c.meta.Apis.CtVpcApis.BandwidthAssociateEipApi.Do(ctx, c.meta.Credential, &ctvpc.BandwidthAssociateEipRequest{
		RegionId:    regionId,
		ClientToken: uuid.NewString(),
		BandwidthId: plan.BandwidthId.ValueString(),
		EipIds:      []string{plan.EipId.ValueString()},
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.RegionId = types.StringValue(regionId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	instance, err := c.getAndMergeBandwidthAssociationEip(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunBandwidthAssociationEip) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunBandwidAssociationEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeBandwidthAssociationEip(ctx, state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunBandwidthAssociationEip) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {

}

func (c *ctyunBandwidthAssociationEip) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunBandwidAssociationEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtVpcApis.BandwidthDisassociateEipApi.Do(ctx, c.meta.Credential, &ctvpc.BandwidthDisassociateEipRequest{
		RegionId:    state.RegionId.ValueString(),
		ClientToken: uuid.NewString(),
		EipIds:      []string{state.EipId.ValueString()},
		BandwidthId: state.BandwidthId.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

func (c *ctyunBandwidthAssociationEip) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := c.name + "导入失败：" + err.Error()
			detail := "导入命令：terraform import " + c.name + ".[导入配置名称] [bandwidth_id],[eip_id],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunBandwidAssociationEipConfig
	var bandwidthID, eipID, regionID string
	cnt := strings.Count(request.ID, common.ImportSeparator)
	switch cnt {
	case 0:
		err = fmt.Errorf("bandwidth_id和eip_id必须输入")
		return
	case 1:
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &bandwidthID, &eipID)
		if err != nil {
			return
		}
	default:
		err = terraform_extend.Split(request.ID, &bandwidthID, &eipID, &regionID)
		if err != nil {
			return
		}
	}
	if bandwidthID == "" {
		err = fmt.Errorf("bandwidth_id不能为空")
		return
	}
	if eipID == "" {
		err = fmt.Errorf("eip_id不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}

	cfg.BandwidthId = types.StringValue(bandwidthID)
	cfg.EipId = types.StringValue(eipID)
	cfg.RegionId = types.StringValue(regionID)
	instance, err := c.getAndMergeBandwidthAssociationEip(ctx, cfg)
	if err != nil {
		return
	}
	if instance == nil {
		err = common.ResourceNotExistError
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunBandwidthAssociationEip) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.bandwidthService = business.NewBandwidthService(meta)
	c.eipService = business.NewEipService(meta)
}

// getAndMergeBandwidthAssociationEip 查询绑定关系
func (c *ctyunBandwidthAssociationEip) getAndMergeBandwidthAssociationEip(ctx context.Context, cfg CtyunBandwidAssociationEipConfig) (*CtyunBandwidAssociationEipConfig, error) {
	result, err := c.meta.Apis.CtVpcApis.BandwidthDescribeApi.Do(ctx, c.meta.Credential, &ctvpc.BandwidthDescribeRequest{
		RegionId:    cfg.RegionId.ValueString(),
		BandwidthId: cfg.BandwidthId.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiSharedbandwidthNotFound {
			return nil, nil
		}
		return nil, err
	}
	if len(result.Eips) == 0 {
		return nil, nil
	}
	var bind bool
	for _, eip := range result.Eips {
		if eip.EipId == cfg.EipId.ValueString() {
			cfg.EipId = types.StringValue(eip.EipId)
			bind = true
			break
		}
	}
	if !bind {
		return nil, nil
	}
	if cfg.ProjectId.IsUnknown() {
		cfg.ProjectId = types.StringNull()
	}
	cfg.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", cfg.BandwidthId.ValueString(), cfg.EipId.ValueString(), cfg.RegionId.ValueString()))
	return &cfg, nil
}

type CtyunBandwidAssociationEipConfig struct {
	ID          types.String `tfsdk:"id"`
	BandwidthId types.String `tfsdk:"bandwidth_id"`
	EipId       types.String `tfsdk:"eip_id"`
	ProjectId   types.String `tfsdk:"project_id"`
	RegionId    types.String `tfsdk:"region_id"`
}
