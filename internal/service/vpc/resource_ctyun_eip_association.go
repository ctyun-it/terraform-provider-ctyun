package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunEipAssociation() resource.Resource {
	return &ctyunEipAssociation{}
}

type ctyunEipAssociation struct {
	meta       *common.CtyunMetadata
	eipService *business.EipService
	ecsService *business.EcsService
}

func (c *ctyunEipAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_eip_association"
}

func (c *ctyunEipAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026753/10219975`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "id",
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
			"association_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "绑定的实例类型：ECS云主机：vm，目前仅支持云主机vm，后续会补充更多可选项",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EipAssociationTypes...),
				},
				Default: stringdefault.StaticString(business.EipAssociationTypeVm),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "绑定对象的实例id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
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

func (c *ctyunEipAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunEipAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 校验eip必须存在
	err := c.eipService.MustExist(ctx, plan.EipId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 校验目标绑定对象存在
	switch plan.AssociationType.ValueString() {
	case business.EipAssociationTypeVm:
		err := c.ecsService.MustExist(ctx, plan.InstanceId.ValueString(), plan.RegionId.ValueString())
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	}

	associationType, err := business.EipAssociationTypeMap.FromOriginalScene(plan.AssociationType.ValueString(), business.EipAssociationTypeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	_, err2 := c.meta.Apis.CtVpcApis.EipAssociateApi.Do(ctx, c.meta.Credential, &ctvpc.EipAssociateRequest{
		RegionId:        regionId,
		ProjectId:       projectId,
		ClientToken:     uuid.NewString(),
		AssociationType: associationType.(int),
		EipId:           plan.EipId.ValueString(),
		AssociationId:   plan.InstanceId.ValueString(),
	})
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	instance, ctyunRequestError := c.getAndMergeEipAssociation(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEipAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunEipAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeEipAssociation(ctx, state)
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

func (c *ctyunEipAssociation) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {

}

func (c *ctyunEipAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunEipAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtVpcApis.EipDisassociateApi.Do(ctx, c.meta.Credential, &ctvpc.EipDisassociateRequest{
		EipId:       state.EipId.ValueString(),
		RegionId:    state.RegionId.ValueString(),
		ProjectId:   state.ProjectId.ValueString(),
		ClientToken: uuid.NewString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [eipId],[regionId]
func (c *ctyunEipAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunEipAssociationConfig
	var eipId, regionId string
	err := terraform_extend.Split(request.ID, &eipId, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.EipId = types.StringValue(eipId)
	cfg.RegionId = types.StringValue(regionId)

	instance, err := c.getAndMergeEipAssociation(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEipAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.eipService = business.NewEipService(meta)
	c.ecsService = business.NewEcsService(meta)
}

// getAndMergeEipAssociation 查询eip绑定关系
func (c *ctyunEipAssociation) getAndMergeEipAssociation(ctx context.Context, cfg CtyunEipAssociationConfig) (*CtyunEipAssociationConfig, error) {
	resp, err := c.meta.Apis.CtVpcApis.EipShowApi.Do(ctx, c.meta.Credential, &ctvpc.EipShowRequest{
		RegionId: cfg.RegionId.ValueString(),
		EipId:    cfg.EipId.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiEipNotFound {
			return nil, nil
		}
		return nil, err
	}
	if resp.AssociationType == "" || resp.AssociationId == "" {
		// 移除实例
		return nil, nil
	}
	associationType, err2 := business.EipAssociationTypeMap.ToOriginalScene(resp.AssociationType, business.EipAssociationTypeMapScene2)
	if err2 != nil {
		return nil, err2
	}
	cfg.AssociationType = types.StringValue(associationType.(string))
	cfg.InstanceId = types.StringValue(resp.AssociationId)
	cfg.ID = types.StringValue(fmt.Sprintf("%s,%s", cfg.EipId.ValueString(), cfg.RegionId.ValueString()))
	return &cfg, nil
}

type CtyunEipAssociationConfig struct {
	ID              types.String `tfsdk:"id"`
	EipId           types.String `tfsdk:"eip_id"`
	AssociationType types.String `tfsdk:"association_type"`
	InstanceId      types.String `tfsdk:"instance_id"`
	ProjectId       types.String `tfsdk:"project_id"`
	RegionId        types.String `tfsdk:"region_id"`
}
