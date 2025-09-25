package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunEcsAffinityGroupAssociation{}
	_ resource.ResourceWithConfigure   = &ctyunEcsAffinityGroupAssociation{}
	_ resource.ResourceWithImportState = &ctyunEcsAffinityGroupAssociation{}
)

type ctyunEcsAffinityGroupAssociation struct {
	meta *common.CtyunMetadata
}

func NewCtyunEcsAffinityGroupAssociation() resource.Resource {
	return &ctyunEcsAffinityGroupAssociation{}
}

func (c *ctyunEcsAffinityGroupAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_affinity_group_association"
}

type CtyunEcsAffinityGroupAssociationConfig struct {
	ID              types.String `tfsdk:"id"`
	RegionID        types.String `tfsdk:"region_id"`
	InstanceID      types.String `tfsdk:"instance_id"`
	AffinityGroupID types.String `tfsdk:"affinity_group_id"`
}

func (c *ctyunEcsAffinityGroupAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10597685**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
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
					validator2.UUID(),
				},
			},
			"affinity_group_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机组ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
		},
	}
}

func (c *ctyunEcsAffinityGroupAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcsAffinityGroupAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeAssociation(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	err = c.associate(ctx, plan)
	if err != nil {
		return
	}
	err = c.checkAfterAssociation(ctx, plan)
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEcsAffinityGroupAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsAffinityGroupAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "未关联") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEcsAffinityGroupAssociation) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunEcsAffinityGroupAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsAffinityGroupAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeDissociate(ctx, state)
	if err != nil {
		return
	}
	err = c.dissociate(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDissociation(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEcsAffinityGroupAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [instanceID],[groupID],[regionID]
func (c *ctyunEcsAffinityGroupAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEcsAffinityGroupAssociationConfig
	var instanceID, groupID, regionID string
	err = terraform_extend.Split(request.ID, &instanceID, &groupID, &regionID)
	if err != nil {
		return
	}

	cfg.InstanceID = types.StringValue(instanceID)
	cfg.AffinityGroupID = types.StringValue(groupID)
	cfg.RegionID = types.StringValue(regionID)

	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeAssociate 绑定前检查
func (c *ctyunEcsAffinityGroupAssociation) checkBeforeAssociation(ctx context.Context, plan CtyunEcsAffinityGroupAssociationConfig) (err error) {
	instanceID, groupID, regionID := plan.InstanceID.ValueString(), plan.AffinityGroupID.ValueString(), plan.RegionID.ValueString()
	params := &ctecs2.CtecsAffinityGroupbindInstanceCheckV41Request{
		RegionID:        regionID,
		InstanceID:      instanceID,
		AffinityGroupID: groupID,
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsAffinityGroupbindInstanceCheckV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if resp.ReturnObj.NeedMigrate != 0 {
		err = fmt.Errorf("云主机 %s 需要迁移后才可加入云主机组 %s", instanceID, groupID)
		return
	}
	err = business.NewEcsService(c.meta).CheckEcsStatus(ctx, instanceID, regionID)
	return
}

// associate 将云主机加入主机组
func (c *ctyunEcsAffinityGroupAssociation) associate(ctx context.Context, plan CtyunEcsAffinityGroupAssociationConfig) (err error) {
	params := &ctecs2.CtecsAffinityGroupbindInstanceV41Request{
		RegionID:        plan.RegionID.ValueString(),
		InstanceID:      plan.InstanceID.ValueString(),
		AffinityGroupID: plan.AffinityGroupID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsAffinityGroupbindInstanceV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}

	return
}

// checkAfterAssociation 绑定后检查
func (c *ctyunEcsAffinityGroupAssociation) checkAfterAssociation(ctx context.Context, plan CtyunEcsAffinityGroupAssociationConfig) (err error) {
	var executeSuccessFlag bool
	var bindID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			bindID, err = c.getEcsAffinityGroup(ctx, plan)
			if err != nil {
				return false
			}
			if bindID == plan.AffinityGroupID.ValueString() {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		return fmt.Errorf("云主机 %s 和云主机组 %s 未关联", plan.InstanceID.ValueString(), plan.AffinityGroupID.ValueString())
	}
	return nil
}

// checkBeforeDissociate 解绑前检查
func (c *ctyunEcsAffinityGroupAssociation) checkBeforeDissociate(ctx context.Context, plan CtyunEcsAffinityGroupAssociationConfig) (err error) {
	instanceID, groupID, regionID := plan.InstanceID.ValueString(), plan.AffinityGroupID.ValueString(), plan.RegionID.ValueString()
	bindID, err := c.getEcsAffinityGroup(ctx, plan)
	if err != nil {
		return
	}
	if bindID != groupID {
		err = fmt.Errorf("云主机 %s 和云主机组 %s 未关联", instanceID, groupID)
		return
	}
	err = business.NewEcsService(c.meta).CheckEcsStatus(ctx, instanceID, regionID)
	return
}

// dissociate 解绑云主机组
func (c *ctyunEcsAffinityGroupAssociation) dissociate(ctx context.Context, plan CtyunEcsAffinityGroupAssociationConfig) (err error) {
	params := &ctecs2.CtecsAffinityGroupUnbindInstanceV41Request{
		RegionID:        plan.RegionID.ValueString(),
		InstanceID:      plan.InstanceID.ValueString(),
		AffinityGroupID: plan.AffinityGroupID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsAffinityGroupUnbindInstanceV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// checkAfterDissociation 解绑后检查
func (c *ctyunEcsAffinityGroupAssociation) checkAfterDissociation(ctx context.Context, plan CtyunEcsAffinityGroupAssociationConfig) (err error) {
	var executeSuccessFlag bool
	var bindID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			bindID, err = c.getEcsAffinityGroup(ctx, plan)
			if err != nil {
				return false
			}
			if bindID != plan.AffinityGroupID.ValueString() {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		return fmt.Errorf("云主机 %s 和云主机组 %s 解绑失败", plan.InstanceID.ValueString(), plan.AffinityGroupID.ValueString())
	}
	return nil
}

// getAndMerge 查询绑定关系
func (c *ctyunEcsAffinityGroupAssociation) getAndMerge(ctx context.Context, plan *CtyunEcsAffinityGroupAssociationConfig) (err error) {
	instanceID, groupID, regionID := plan.InstanceID.ValueString(), plan.AffinityGroupID.ValueString(), plan.RegionID.ValueString()
	bindID, err := c.getEcsAffinityGroup(ctx, *plan)
	if err != nil {
		return
	}
	if bindID != groupID {
		err = fmt.Errorf("云主机 %s 和云主机组 %s 未关联", instanceID, groupID)
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", instanceID, groupID, regionID))
	return
}

// getEcsAffinityGroup 查询云主机绑定的云主机组
func (c *ctyunEcsAffinityGroupAssociation) getEcsAffinityGroup(ctx context.Context, plan CtyunEcsAffinityGroupAssociationConfig) (groupID string, err error) {
	params := &ctecs2.CtecsGetAffinityGroupV41Request{
		RegionID:   plan.RegionID.ValueString(),
		InstanceID: plan.InstanceID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsGetAffinityGroupV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		if resp.ErrorCode == common.EcsAffinityGroupNotBound { // 没绑定主机组，返回空groupID
			err = nil
			return
		}
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	groupID = resp.ReturnObj.AffinityGroupID
	return
}
