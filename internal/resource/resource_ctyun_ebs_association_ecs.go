package resource

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "terraform-provider-ctyun/internal/extend/terraform/defaults"
)

func NewCtyunEbsAssociation() resource.Resource {
	return &ctyunEbsAssociation{}
}

type ctyunEbsAssociation struct {
	meta       *common.CtyunMetadata
	ecsService *business.EcsService
	ebsService *business.EbsService
}

func (c *ctyunEbsAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_association_ecs"
}

func (c *ctyunEbsAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027696/10169293**`,
		Attributes: map[string]schema.Attribute{
			"ebs_id": schema.StringAttribute{
				Required:    true,
				Description: "磁盘id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机id，多可用区资源池下，云硬盘和云主机必须在同个az才能支持挂载",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

func (c *ctyunEbsAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunEbsAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.ebsService.MustExist(ctx, plan.EbsId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	err = c.ecsService.MustExist(ctx, plan.InstanceId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	regionId := plan.RegionId.ValueString()
	resp, err := c.meta.Apis.CtEbsApis.EbsAssociateApi.Do(ctx, c.meta.Credential, &ctebs.EbsAssociateRequest{
		RegionId:   regionId,
		DiskId:     plan.EbsId.ValueString(),
		InstanceId: plan.InstanceId.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	helper := business.NewGeneralJobHelper(c.meta.Apis.CtEcsApis.JobShowApi)
	_, err2 := helper.JobLoop(ctx, c.meta.Credential, regionId, resp.DiskJobId)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	plan.RegionId = types.StringValue(regionId)
	instance, ctyunRequestError := c.getAndMergeEbsAssociationEcs(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEbsAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunEbsAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeEbsAssociationEcs(ctx, state)
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

func (c *ctyunEbsAssociation) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {

}

func (c *ctyunEbsAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunEbsAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.CtEbsApis.EbsDisassociateApi.Do(ctx, c.meta.Credential, &ctebs.EbsDisassociateRequest{
		RegionId: state.RegionId.ValueString(),
		DiskId:   state.EbsId.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	helper := business.NewGeneralJobHelper(c.meta.Apis.CtEcsApis.JobShowApi)
	_, err2 := helper.JobLoop(ctx, c.meta.Credential, state.RegionId.ValueString(), resp.DiskJobId)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [diskId],[ecsId],[regionId]
func (c *ctyunEbsAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunEbsAssociationConfig
	var diskId, ecsId, regionId string
	err := terraform_extend.Split(request.ID, &diskId, &ecsId, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.EbsId = types.StringValue(diskId)
	cfg.InstanceId = types.StringValue(ecsId)
	cfg.RegionId = types.StringValue(regionId)

	instance, err := c.getAndMergeEbsAssociationEcs(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEbsAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ebsService = business.NewEbsService(meta)
	c.ecsService = business.NewEcsService(meta)
}

// getAndMergeEbsAssociationEcs 查询绑定关系
func (c *ctyunEbsAssociation) getAndMergeEbsAssociationEcs(ctx context.Context, cfg CtyunEbsAssociationConfig) (*CtyunEbsAssociationConfig, error) {
	resp, err := c.meta.Apis.CtEbsApis.EbsShowApi.Do(ctx, c.meta.Credential, &ctebs.EbsShowRequest{
		RegionId: cfg.RegionId.ValueString(),
		DiskId:   cfg.EbsId.ValueString(),
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Attachments) == 0 {
		return nil, nil
	}
	for _, each := range resp.Attachments {
		if each.InstanceId == cfg.InstanceId.ValueString() {
			cfg.InstanceId = types.StringValue(each.InstanceId)
			break
		}
	}
	return &cfg, err
}

type CtyunEbsAssociationConfig struct {
	EbsId      types.String `tfsdk:"ebs_id"`
	InstanceId types.String `tfsdk:"instance_id"`
	RegionId   types.String `tfsdk:"region_id"`
}
