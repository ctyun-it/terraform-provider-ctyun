package resource

import (
	"context"
	"errors"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctiam"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
)

func NewCtyunPolicyAssociationUserGroup() resource.Resource {
	return &ctyunPolicyAssociationUserGroup{}
}

type ctyunPolicyAssociationUserGroup struct {
	meta             *common.CtyunMetadata
	userGroupService *business.UserGroupService
}

func (c *ctyunPolicyAssociationUserGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_iam_policy_association_user_group"
}

func (c *ctyunPolicyAssociationUserGroup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10345725/10409392**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "绑定关系id",
			},
			"user_group_id": schema.StringAttribute{
				Required:    true,
				Description: "用户组id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_id": schema.StringAttribute{
				Required:    true,
				Description: "策略id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，当授权的策略为资源池级别时必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(""),
			},
		},
	}
}

func (c *ctyunPolicyAssociationUserGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunPolicyAssociationUserGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	rangeType, err := c.checkAndGetRangeType(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	err = c.checkUserGroup(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	regionIds := []string{}
	if plan.RegionId.ValueString() != "" {
		regionIds = append(regionIds, plan.RegionId.ValueString())
	}
	resp, err := c.meta.Apis.CtIamApis.PolicyAttachUserGroupApi.Do(ctx, c.meta.Credential, &ctiam.PolicyAttachUserGroupRequest{
		UserGroupId: plan.UserGroupId.ValueString(),
		PolicyIds:   []string{plan.PolicyId.ValueString()},
		RangeType:   rangeType,
		RegionIds:   regionIds,
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.StringValue(resp.PrivilegeMessage[0].PrivilegeId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeIamPolicyAssociationUserGroup(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunPolicyAssociationUserGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunPolicyAssociationUserGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeIamPolicyAssociationUserGroup(ctx, state)
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

func (c *ctyunPolicyAssociationUserGroup) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {

}

func (c *ctyunPolicyAssociationUserGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunPolicyAssociationUserGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtIamApis.PolicyInvalidUserGroupApi.Do(ctx, c.meta.Credential, &ctiam.PolicyInvalidUserGroupRequest{
		PrivilegeId: state.Id.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [privilegeId]
func (c *ctyunPolicyAssociationUserGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunPolicyAssociationUserGroupConfig
	var privilegeId string
	err := terraform_extend.Split(request.ID, &privilegeId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(privilegeId)

	instance, err := c.getAndMergeIamPolicyAssociationUserGroup(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunPolicyAssociationUserGroup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.userGroupService = business.NewUserGroupService(meta)
}

// checkUserGroup 校验用户组是否存在
func (c *ctyunPolicyAssociationUserGroup) checkUserGroup(ctx context.Context, cfg CtyunPolicyAssociationUserGroupConfig) error {
	return c.userGroupService.MustExist(ctx, cfg.UserGroupId.ValueString())
}

// checkRangeType 校验并返回范围
func (c *ctyunPolicyAssociationUserGroup) checkAndGetRangeType(ctx context.Context, cfg CtyunPolicyAssociationUserGroupConfig) (string, error) {
	resp, err := c.meta.Apis.CtIamApis.PolicyGetApi.Do(ctx, c.meta.Credential, &ctiam.PolicyGetRequest{
		PolicyId: cfg.PolicyId.ValueString(),
	})
	if err != nil {
		return "", err
	}
	policy, err2 := business.PolicyRangeMap.ToOriginalScene(resp.PolicyRange, business.PolicyRangeMapScene1)
	if err2 != nil {
		return "", err
	}
	switch policy.(string) {
	case business.PolicyRangeRegion:
		// 如果是资源池级别的范围，用户又没有填写资源池id，报错
		if cfg.RegionId.IsNull() || cfg.RegionId.IsUnknown() {
			return "", errors.New("策略：" + cfg.PolicyId.ValueString() + "为资源池级别的范围，授权时必须填写资源池id")
		}
	case business.PolicyRangeGlobal:
		// 如果是全局级别的范围，用户又填写资源池id，报错
		if !cfg.RegionId.IsNull() && !cfg.RegionId.IsUnknown() && cfg.RegionId.ValueString() != "" {
			return "", errors.New("策略：" + cfg.PolicyId.ValueString() + "为全局级别的范围，授权时不能填写资源池id")
		}
	}
	key, err2 := business.PolicyRangeMap.Map(resp.PolicyRange, business.PolicyRangeMapScene1, business.PolicyRangeMapScene2)
	if err2 != nil {
		return "", err
	}
	return key.(string), nil
}

// getAndMergeIamPolicyAssociationUserGroup 查询授权信息
func (c *ctyunPolicyAssociationUserGroup) getAndMergeIamPolicyAssociationUserGroup(ctx context.Context, cfg CtyunPolicyAssociationUserGroupConfig) (*CtyunPolicyAssociationUserGroupConfig, error) {
	resp, err := c.meta.Apis.CtIamApis.PrivilegeGetApi.Do(ctx, c.meta.Credential, &ctiam.PrivilegeGetRequest{
		PrivilegeId: cfg.Id.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.CtiamNoPrivilege {
			return nil, nil
		}
		return nil, err
	}

	regionId := ""
	if resp.RegionId != "No ProjectID" {
		regionId = resp.RegionId
	}

	cfg.Id = types.StringValue(resp.PrivilegeId)
	cfg.UserGroupId = types.StringValue(resp.Id)
	cfg.PolicyId = types.StringValue(resp.PolicyId)
	cfg.RegionId = types.StringValue(regionId)
	return &cfg, nil
}

type CtyunPolicyAssociationUserGroupConfig struct {
	Id          types.String `tfsdk:"id"`
	UserGroupId types.String `tfsdk:"user_group_id"`
	PolicyId    types.String `tfsdk:"policy_id"`
	RegionId    types.String `tfsdk:"region_id"`
}
