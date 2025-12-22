package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
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
	_ resource.Resource                = &ctyunEcsPortAssociation{}
	_ resource.ResourceWithConfigure   = &ctyunEcsPortAssociation{}
	_ resource.ResourceWithImportState = &ctyunEcsPortAssociation{}
)

type ctyunEcsPortAssociation struct {
	meta *common.CtyunMetadata
}

func NewCtyunEcsPortAssociation() resource.Resource {
	return &ctyunEcsPortAssociation{}
}

func (c *ctyunEcsPortAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_port_association"
}

type CtyunEcsPortAssociationConfig struct {
	ID         types.String `tfsdk:"id"`
	RegionID   types.String `tfsdk:"region_id"`
	InstanceID types.String `tfsdk:"instance_id"`
	PortID     types.String `tfsdk:"port_id"`
	AzName     types.String `tfsdk:"az_name"`
	ProjectID  types.String `tfsdk:"project_id"`
}

func (c *ctyunEcsPortAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10225195`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识符",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.Any(
						validator2.UUID(),
					),
				},
			},
			"port_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性网卡ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.Any(
						validator2.PortValidate(),
					),
				},
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称，如果不填则默认使用provider ctyun中的az_name或环境变量中的CTYUN_AZ_NAME",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraAzName, true),
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
		},
	}
}

func (c *ctyunEcsPortAssociation) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunEcsPortAssociation) Create(ctx context.Context, req resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcsPortAssociationConfig
	response.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, &plan)
	if err != nil {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	// 设置ID
	plan.ID = types.StringValue(generateEcsPortAssociationId(plan.RegionID.ValueString(), plan.InstanceID.ValueString(), plan.PortID.ValueString()))

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEcsPortAssociation) create(ctx context.Context, plan *CtyunEcsPortAssociationConfig) (err error) {
	// 绑定弹性网卡到云主机
	attachRequest := &ctecs.CtecsPortsAttachInstanceV41Request{
		ClientToken:        uuid.NewString(),
		RegionID:           plan.RegionID.ValueString(),
		ProjectID:          plan.ProjectID.ValueString(),
		AzName:             plan.AzName.ValueString(),
		NetworkInterfaceID: plan.PortID.ValueString(),
		InstanceID:         plan.InstanceID.ValueString(),
		InstanceType:       3, // 3-虚拟机
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsPortsAttachInstanceV41Api.Do(ctx, c.meta.SdkCredential, attachRequest)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

func (c *ctyunEcsPortAssociation) Read(ctx context.Context, req resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state *CtyunEcsPortAssociationConfig
	response.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.getAndMerge(ctx, state)

	if err != nil {
		if strings.Contains(err.Error(), "弹性网卡未绑定") {
			err = nil
			// 如果弹性网卡未绑定，则从状态中移除该资源
			response.State.RemoveResource(ctx)
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEcsPortAssociation) getAndMerge(ctx context.Context, state *CtyunEcsPortAssociationConfig) (err error) {
	regionId := c.meta.GetExtraIfEmpty(state.RegionID.ValueString(), common.ExtraRegionId)
	// 查询云主机详情以确认弹性网卡是否仍然绑定
	describeRequest := &ctecs.CtecsDetailsInstanceV41Request{
		RegionID:   regionId,
		InstanceID: state.InstanceID.ValueString(),
	}

	describeResponse, err := c.meta.Apis.SdkCtEcsApis.CtecsDetailsInstanceV41Api.Do(ctx, c.meta.SdkCredential, describeRequest)
	if err != nil {
		return
	}

	// 检查弹性网卡是否仍然绑定到该云主机
	if describeResponse.ReturnObj != nil {
		for _, networkCard := range describeResponse.ReturnObj.NetworkCardList {
			if networkCard != nil && utils.SecString(networkCard.NetworkCardID) == state.PortID.ValueString() {
				return nil
			}
		}
	}

	return fmt.Errorf("弹性网卡未绑定")

}

func (c *ctyunEcsPortAssociation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// 由于所有属性都需要替换资源，实际上不会执行更新操作
}

func (c *ctyunEcsPortAssociation) Delete(ctx context.Context, req resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsPortAssociationConfig
	response.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEcsPortAssociation) destroy(ctx context.Context, state CtyunEcsPortAssociationConfig) (err error) {
	regionId := c.meta.GetExtraIfEmpty(state.RegionID.ValueString(), common.ExtraRegionId)

	// 解绑弹性网卡
	detachRequest := &ctecs.CtecsPortsDetachInstanceV41Request{
		ClientToken:        uuid.NewString(),
		RegionID:           regionId,
		NetworkInterfaceID: state.PortID.ValueString(),
		InstanceID:         state.InstanceID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsPortsDetachInstanceV41Api.Do(ctx, c.meta.SdkCredential, detachRequest)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

func (c *ctyunEcsPortAssociation) ImportState(ctx context.Context, req resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [instanceId],[networkInterfaceId],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunEcsPortAssociationConfig
	var regionId, instanceId, networkInterfaceId string

	if strings.Count(req.ID, common.ImportSeparator) == 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(req.ID, &instanceId, &networkInterfaceId)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(req.ID, &instanceId, &networkInterfaceId, &regionId)
		if err != nil {
			return
		}
	}
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	if instanceId == "" {
		err = fmt.Errorf("instanceId不能为空")
		return
	}
	if networkInterfaceId == "" {
		err = fmt.Errorf("networkInterfaceId不能为空")
		return
	}

	cfg.ID = types.StringValue(req.ID)
	cfg.RegionID = types.StringValue(regionId)
	cfg.InstanceID = types.StringValue(instanceId)
	cfg.PortID = types.StringValue(networkInterfaceId)

	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}

	// 设置导入的属性
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)

}

func (c *ctyunEcsPortAssociation) checkBeforeCreate(ctx context.Context, c2 *CtyunEcsPortAssociationConfig) error {
	return nil
}

func generateEcsPortAssociationId(regionId, instanceId, networkInterfaceId string) string {
	return fmt.Sprintf("%s,%s,%s", instanceId, networkInterfaceId, regionId)
}
