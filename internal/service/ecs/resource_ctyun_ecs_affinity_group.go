package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

var (
	_ resource.Resource                = &ctyunEcsAffinityGroup{}
	_ resource.ResourceWithConfigure   = &ctyunEcsAffinityGroup{}
	_ resource.ResourceWithImportState = &ctyunEcsAffinityGroup{}
)

type ctyunEcsAffinityGroup struct {
	meta *common.CtyunMetadata
}

func NewCtyunEcsAffinityGroup() resource.Resource {
	return &ctyunEcsAffinityGroup{}
}

func (c *ctyunEcsAffinityGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_affinity_group"
}

type CtyunEcsAffinityGroupConfig struct {
	ID                  types.String `tfsdk:"id"`
	AffinityGroupID     types.String `tfsdk:"affinity_group_id"`
	RegionID            types.String `tfsdk:"region_id"`
	AffinityGroupName   types.String `tfsdk:"affinity_group_name"`
	AffinityGroupPolicy types.String `tfsdk:"affinity_group_policy"`
}

func (c *ctyunEcsAffinityGroup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10597693`,
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
			"affinity_group_id": schema.StringAttribute{
				Computed:    true,
				Description: "云主机组ID",
			},
			"affinity_group_name": schema.StringAttribute{
				Required:    true,
				Description: "云主机组名称，满足以下规则：长度在1-64个字符，只能由中文、英文字母、数字、下划线_、中划线-、点.组成 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_.-]+$`), "不满足云主机组名称要求"),
				},
			},
			"affinity_group_policy": schema.StringAttribute{
				Required:    true,
				Description: "云主机组策略类型，取值范围：<br />anti-affinity（强制反亲和性），<br />affinity（强制亲和性），<br />soft-anti-affinity（反亲和性），<br />soft-affinity（亲和性)，<br />power-anti-affinity（电力反亲和性)",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EcsGroupAntiAffinityPolicy, business.EcsGroupAffinityPolicy, business.EcsGroupSoftAntiAffinity, business.EcsGroupSoftAffinity, business.EcsGroupPowerAntiAffinity),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (c *ctyunEcsAffinityGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcsAffinityGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建
	groupID, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.AffinityGroupID = types.StringValue(groupID)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEcsAffinityGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsAffinityGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEcsAffinityGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunEcsAffinityGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunEcsAffinityGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.updateName(ctx, plan, state)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEcsAffinityGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsAffinityGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEcsAffinityGroup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [groupID],[regionID]
func (c *ctyunEcsAffinityGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEcsAffinityGroupConfig
	var groupID, regionID string
	err = terraform_extend.Split(request.ID, &groupID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.AffinityGroupID = types.StringValue(groupID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建云主机组
func (c *ctyunEcsAffinityGroup) create(ctx context.Context, plan CtyunEcsAffinityGroupConfig) (groupID string, err error) {
	params := &ctecs2.CtecsCreateAffinityGroupV41Request{
		RegionID:          plan.RegionID.ValueString(),
		AffinityGroupName: plan.AffinityGroupName.ValueString(),
		PolicyType:        business.EcSGroupPolicyMap[plan.AffinityGroupPolicy.ValueString()],
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsCreateAffinityGroupV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	groupID = resp.ReturnObj.AffinityGroupID
	return
}

// getAndMerge 查询主机组
func (c *ctyunEcsAffinityGroup) getAndMerge(ctx context.Context, plan *CtyunEcsAffinityGroupConfig) (err error) {
	params := &ctecs2.CtecsListAffinityGroupV41Request{
		RegionID:        plan.RegionID.ValueString(),
		AffinityGroupID: plan.AffinityGroupID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsListAffinityGroupV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if len(resp.ReturnObj.Results) != 1 || resp.ReturnObj.Results[0].AffinityGroupPolicy == nil {
		err = common.InvalidReturnObjResultsError
		return
	}
	plan.AffinityGroupName = types.StringValue(resp.ReturnObj.Results[0].AffinityGroupName)
	plan.AffinityGroupPolicy = types.StringValue(resp.ReturnObj.Results[0].AffinityGroupPolicy.PolicyTypeName)
	plan.ID = plan.AffinityGroupID
	return
}

// updateName 修改主机组名称
func (c *ctyunEcsAffinityGroup) updateName(ctx context.Context, plan, state CtyunEcsAffinityGroupConfig) (err error) {
	if plan.AffinityGroupName.Equal(state.AffinityGroupName) {
		return
	}
	params := &ctecs2.CtecsUpdateAffinityGroupRequest{
		RegionID:          state.RegionID.ValueString(),
		AffinityGroupID:   state.AffinityGroupID.ValueString(),
		AffinityGroupName: plan.AffinityGroupName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsUpdateAffinityGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// delete 删除主机组
func (c *ctyunEcsAffinityGroup) delete(ctx context.Context, plan CtyunEcsAffinityGroupConfig) (err error) {
	params := &ctecs2.CtecsDeleteAffinityGroupRequest{
		RegionID:        plan.RegionID.ValueString(),
		AffinityGroupID: plan.AffinityGroupID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsDeleteAffinityGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}
