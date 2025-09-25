package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027696/10169293`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"ebs_id": schema.StringAttribute{
				Required:    true,
				Description: "磁盘id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机id，多可用区资源池下，云硬盘和云主机必须在同个az才能支持挂载",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
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

func (c *ctyunEbsAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbsAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.ebsService.MustExist(ctx, plan.EbsId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		return
	}
	err = c.ecsService.CheckEcsStatus(ctx, plan.InstanceId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		return
	}

	regionId := plan.RegionId.ValueString()
	resp, err := c.meta.Apis.CtEbsApis.EbsAssociateApi.Do(ctx, c.meta.Credential, &ctebs.EbsAssociateRequest{
		RegionId:   regionId,
		DiskId:     plan.EbsId.ValueString(),
		InstanceId: plan.InstanceId.ValueString(),
	})
	if err != nil {
		return
	}
	helper := business.NewGeneralJobHelper(c.meta.Apis.CtEcsApis.JobShowApi)
	_, err = helper.JobLoop(ctx, c.meta.Credential, regionId, resp.DiskJobId)
	if err != nil {
		return
	}

	plan.RegionId = types.StringValue(regionId)
	instance, err := c.getAndMergeEbsAssociationEcs(ctx, plan)
	if err != nil {
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
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.ecsService.CheckEcsStatus(ctx, state.InstanceId.ValueString(), state.RegionId.ValueString())
	if err != nil {
		return
	}
	resp, err := c.meta.Apis.CtEbsApis.EbsDisassociateApi.Do(ctx, c.meta.Credential, &ctebs.EbsDisassociateRequest{
		RegionId: state.RegionId.ValueString(),
		DiskId:   state.EbsId.ValueString(),
	})
	if err != nil {
		return
	}
	helper := business.NewGeneralJobHelper(c.meta.Apis.CtEcsApis.JobShowApi)
	_, err = helper.JobLoop(ctx, c.meta.Credential, state.RegionId.ValueString(), resp.DiskJobId)
	if err != nil {
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [diskId],[ecsId],[regionId]
func (c *ctyunEbsAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEbsAssociationConfig
	var diskId, ecsId, regionId string
	err = terraform_extend.Split(request.ID, &diskId, &ecsId, &regionId)
	if err != nil {
		return
	}

	cfg.EbsId = types.StringValue(diskId)
	cfg.InstanceId = types.StringValue(ecsId)
	cfg.RegionId = types.StringValue(regionId)

	var instance *CtyunEbsAssociationConfig
	instance, err = c.getAndMergeEbsAssociationEcs(ctx, cfg)
	if err != nil {
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
	cfg.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", cfg.EbsId.ValueString(), cfg.InstanceId.ValueString(), cfg.RegionId.ValueString()))
	return &cfg, err
}

type CtyunEbsAssociationConfig struct {
	ID         types.String `tfsdk:"id"`
	EbsId      types.String `tfsdk:"ebs_id"`
	InstanceId types.String `tfsdk:"instance_id"`
	RegionId   types.String `tfsdk:"region_id"`
}
