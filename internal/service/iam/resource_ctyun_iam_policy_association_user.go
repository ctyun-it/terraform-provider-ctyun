package iam

import (
	"context"
	"errors"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunPolicyAssociationUser() resource.Resource {
	return &ctyunPolicyAssociationUser{}
}

type ctyunPolicyAssociationUser struct {
	meta        *common.CtyunMetadata
	userService *business.UserService
}

func (c *ctyunPolicyAssociationUser) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_iam_policy_association_user"
}

func (c *ctyunPolicyAssociationUser) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10345725/10409392`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "绑定关系id",
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "用户id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"policy_id": schema.StringAttribute{
				Required:    true,
				Description: "策略id",
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
				Description: "资源池ID，当授权的策略为资源池级别时必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
		},
	}
}

func (c *ctyunPolicyAssociationUser) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunPolicyAssociationUserConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	rangeType, err := c.checkAndGetRangeType(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	err = c.checkUser(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	regionIds := []string{}
	if plan.RegionId.ValueString() != "" {
		regionIds = append(regionIds, plan.RegionId.ValueString())
	}
	resp, err := c.meta.Apis.CtIamApis.PolicyAttachUserApi.Do(ctx, c.meta.Credential, &ctiam.PolicyAttachUserRequest{
		UserId:    plan.UserId.ValueString(),
		PolicyIds: []string{plan.PolicyId.ValueString()},
		RangeType: rangeType,
		RegionIds: regionIds,
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

	instance, ctyunRequestError := c.getAndMergeIamPolicyAssociationUser(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunPolicyAssociationUser) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunPolicyAssociationUserConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	instance, err := c.getAndMergeIamPolicyAssociationUser(ctx, state)
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

func (c *ctyunPolicyAssociationUser) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {

}

func (c *ctyunPolicyAssociationUser) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunPolicyAssociationUserConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	_, err := c.meta.Apis.CtIamApis.PolicyInvalidUserGroupApi.Do(ctx, c.meta.Credential, &ctiam.PolicyInvalidUserGroupRequest{
		PrivilegeId: state.Id.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [privilegeId]
func (c *ctyunPolicyAssociationUser) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunPolicyAssociationUserConfig
	var privilegeId string
	err := terraform_extend.Split(request.ID, &privilegeId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(privilegeId)

	instance, err := c.getAndMergeIamPolicyAssociationUser(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunPolicyAssociationUser) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.userService = business.NewUserService(meta)
}

// checkRangeType 校验并返回范围
func (c *ctyunPolicyAssociationUser) checkAndGetRangeType(ctx context.Context, cfg CtyunPolicyAssociationUserConfig) (string, error) {
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
		// 如果是资源池级别的范围，用户又没有填写资源池ID，报错
		if cfg.RegionId.IsNull() || cfg.RegionId.IsUnknown() {
			return "", errors.New("策略：" + cfg.PolicyId.ValueString() + "为资源池级别的范围，授权时必须填写资源池ID")
		}
	case business.PolicyRangeGlobal:
		// 如果是全局级别的范围，用户又填写资源池ID，报错
		if !cfg.RegionId.IsNull() && !cfg.RegionId.IsUnknown() && cfg.RegionId.ValueString() != "" {
			return "", errors.New("策略：" + cfg.PolicyId.ValueString() + "为全局级别的范围，授权时不能填写资源池ID")
		}
	}
	key, err2 := business.PolicyRangeMap.Map(resp.PolicyRange, business.PolicyRangeMapScene1, business.PolicyRangeMapScene2)
	if err2 != nil {
		return "", err
	}
	return key.(string), nil
}

// checkUser 校验用户是否存在
func (c *ctyunPolicyAssociationUser) checkUser(ctx context.Context, cfg CtyunPolicyAssociationUserConfig) error {
	return c.userService.MustExist(ctx, cfg.UserId.ValueString())
}

// getAndMergeIamPolicyAssociationUser 查询授权信息
func (c *ctyunPolicyAssociationUser) getAndMergeIamPolicyAssociationUser(ctx context.Context, cfg CtyunPolicyAssociationUserConfig) (*CtyunPolicyAssociationUserConfig, error) {
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
	cfg.UserId = types.StringValue(resp.Id)
	cfg.PolicyId = types.StringValue(resp.PolicyId)
	cfg.RegionId = types.StringValue(regionId)
	return &cfg, nil
}

type CtyunPolicyAssociationUserConfig struct {
	Id       types.String `tfsdk:"id"`
	UserId   types.String `tfsdk:"user_id"`
	PolicyId types.String `tfsdk:"policy_id"`
	RegionId types.String `tfsdk:"region_id"`
}
