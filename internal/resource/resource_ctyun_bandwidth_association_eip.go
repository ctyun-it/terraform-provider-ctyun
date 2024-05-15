package resource

import (
	"context"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctvpc"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "terraform-provider-ctyun/internal/extend/terraform/defaults"
)

type ctyunBandwidthAssociationEip struct {
	meta             *common.CtyunMetadata
	bandwidthService *business.BandwidthService
	eipService       *business.EipService
}

func NewCtyunBandwidthAssociationEip() resource.Resource {
	return &ctyunBandwidthAssociationEip{}
}

func (c *ctyunBandwidthAssociationEip) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_bandwidth_association_eip"
}

func (c *ctyunBandwidthAssociationEip) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026761/10030030**`,
		Attributes: map[string]schema.Attribute{
			"bandwidth_id": schema.StringAttribute{
				Required:    true,
				Description: "共享带宽id",
				Validators:  nil,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"eip_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性ip的id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目id，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
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
	err := c.bandwidthService.MustExist(ctx, plan.BandwidthId.ValueString(), plan.RegionId.ValueString(), plan.ProjectId.ValueString())
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

// 导入命令：terraform import [配置标识].[导入配置名称] [bandwidthId],[eipId],[regionId]
func (c *ctyunBandwidthAssociationEip) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunBandwidAssociationEipConfig
	var bandwidthId, eipId, regionId string
	err := terraform_extend.Split(request.ID, &bandwidthId, &eipId, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.BandwidthId = types.StringValue(bandwidthId)
	cfg.EipId = types.StringValue(eipId)
	cfg.RegionId = types.StringValue(regionId)

	instance, err := c.getAndMergeBandwidthAssociationEip(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
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
		return nil, err
	}
	if len(result.Eips) == 0 {
		return nil, nil
	}
	for _, eip := range result.Eips {
		if eip.EipId == cfg.EipId.ValueString() {
			cfg.EipId = types.StringValue(eip.EipId)
			break
		}
	}
	return &cfg, nil
}

type CtyunBandwidAssociationEipConfig struct {
	BandwidthId types.String `tfsdk:"bandwidth_id"`
	EipId       types.String `tfsdk:"eip_id"`
	ProjectId   types.String `tfsdk:"project_id"`
	RegionId    types.String `tfsdk:"region_id"`
}
