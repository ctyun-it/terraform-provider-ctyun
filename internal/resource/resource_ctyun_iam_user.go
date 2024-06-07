package resource

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
)

func NewCtyunIamUser() resource.Resource {
	return &ctyunIamUser{}
}

type ctyunIamUser struct {
	meta *common.CtyunMetadata
}

func (c *ctyunIamUser) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_iam_user"
}

func (c *ctyunIamUser) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "用户id",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "账户id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "用户名，长度为4到32位",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(4, 32),
				},
			},
			"password": schema.StringAttribute{
				Required:    true,
				Description: "密码，密码必须包含数字大小写字母，密码长度必须在8-26位之间，密码必须包含特殊字符：$./,;~!@#%_$^*?+{}[-]",
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(8, 26),
				},
			},
			"phone": schema.StringAttribute{
				Required:    true,
				Description: "手机号",
				Validators: []validator.String{
					validator2.Phone(),
				},
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "登录邮箱",
				Validators: []validator.String{
					validator2.Email(),
				},
			},
			"user_group_ids": schema.SetAttribute{
				Required:    true,
				Description: "用户组id，用户加入的目标安全组id，创建用户时至少加入一个用户组",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "备注，长度最大为64",
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(64),
				},
			},
		},
	}
}

func (c *ctyunIamUser) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunIamUserConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var userGroupIds []types.String
	plan.UserGroupIds.ElementsAs(ctx, &userGroupIds, true)
	var ugs []ctiam.UserGroup
	for _, id := range userGroupIds {
		ugs = append(ugs, ctiam.UserGroup{Id: id.ValueString()})
	}
	resp, err := c.meta.Apis.CtIamApis.UserCreateApi.Do(ctx, c.meta.Credential, &ctiam.UserCreateRequest{
		LoginEmail:         plan.Email.ValueString(),
		MobilePhone:        plan.Phone.ValueString(),
		UserName:           plan.Name.ValueString(),
		Remark:             plan.Description.ValueString(),
		GeneratePassword:   false,
		LoginResetPassword: false,
		SourcePassword:     plan.Password.ValueString(),
		Groups:             ugs,
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.StringValue(resp.UserId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeIamUser(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunIamUser) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunIamUserConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeIamUser(ctx, state)
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

func (c *ctyunIamUser) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunIamUserConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan CtyunIamUserConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtIamApis.UserUpdateApi.Do(ctx, c.meta.Credential, &ctiam.UserUpdateRequest{
		UserId:      state.Id.ValueString(),
		Remark:      plan.Description.ValueString(),
		LoginEmail:  plan.Email.ValueString(),
		MobilePhone: plan.Phone.ValueString(),
		UserName:    plan.Name.ValueString(),
		Prohibit:    0,
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 更新密码
	if !plan.Password.Equal(state.Password) {
		_, err := c.meta.Apis.CtIamApis.UserResetPasswordApi.Do(ctx, c.meta.Credential, &ctiam.UserResetPasswordRequest{
			UserId:      state.Id.ValueString(),
			OldPassword: state.Password.ValueString(),
			NewPassword: plan.Password.ValueString(),
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	}

	// 更新用户组
	err2 := c.updateUserGroup(ctx, state, plan)
	if err != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	instance, ctyunRequestError := c.getAndMergeIamUser(ctx, state)
	instance.Password = plan.Password
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunIamUser) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunIamUserConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtIamApis.UserInvalidApi.Do(ctx, c.meta.Credential, &ctiam.UserInvalidRequest{
		UserId: state.Id.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [iamUserId]
func (c *ctyunIamUser) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunIamUserConfig
	var iamUserId string
	err := terraform_extend.Split(request.ID, &iamUserId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(iamUserId)

	instance, err := c.getAndMergeIamUser(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunIamUser) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// updateSecurityGroup 更新用户组
func (c *ctyunIamUser) updateUserGroup(ctx context.Context, state CtyunIamUserConfig, plan CtyunIamUserConfig) error {
	var mapping = make(map[string]struct{})
	var userGroups []types.String
	state.UserGroupIds.ElementsAs(ctx, &userGroups, true)
	for _, group := range userGroups {
		mapping[group.ValueString()] = struct{}{}
	}

	// 需要新加入的用户组id
	var addUserGroups []string
	plan.UserGroupIds.ElementsAs(ctx, &userGroups, true)
	for _, group := range userGroups {
		groupStr := group.ValueString()
		_, ok := mapping[groupStr]
		if ok {
			delete(mapping, groupStr)
		} else {
			addUserGroups = append(addUserGroups, groupStr)
		}
	}
	if len(addUserGroups) != 0 {
		_, err := c.meta.Apis.CtIamApis.UserAttachUserGroupApi.Do(ctx, c.meta.Credential, &ctiam.UserAttachUserGroupRequest{
			UserId:   state.Id.ValueString(),
			GroupIds: addUserGroups,
		})
		if err != nil {
			return err
		}
	}

	// 剩余的是离开的用户组
	var removeUserGroups []string
	for key := range mapping {
		removeUserGroups = append(removeUserGroups, key)
	}
	if len(removeUserGroups) != 0 {
		_, err := c.meta.Apis.CtIamApis.UserRemoveUserGroupApi.Do(ctx, c.meta.Credential, &ctiam.UserRemoveUserGroupRequest{
			UserId:   state.Id.ValueString(),
			GroupIds: removeUserGroups,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// getAndMergeIamUser 查询iam用户
func (c *ctyunIamUser) getAndMergeIamUser(ctx context.Context, cfg CtyunIamUserConfig) (*CtyunIamUserConfig, error) {
	resp, err := c.meta.Apis.CtIamApis.UserGetApi.Do(ctx, c.meta.Credential, &ctiam.UserGetRequest{
		UserId: cfg.Id.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.CtiamNoPermission {
			return nil, nil
		}
		return nil, err
	}
	groups := []types.String{}
	for _, group := range resp.Groups {
		groups = append(groups, types.StringValue(group.Id))
	}
	ugs, _ := types.SetValueFrom(ctx, types.StringType, groups)

	cfg.AccountId = types.StringValue(resp.AccountId)
	cfg.Name = types.StringValue(resp.UserName)
	cfg.Phone = types.StringValue(resp.MobilePhone)
	cfg.Email = types.StringValue(resp.LoginEmail)
	cfg.UserGroupIds = ugs
	cfg.Description = types.StringValue(resp.Remark)
	return &cfg, nil
}

type CtyunIamUserConfig struct {
	Id           types.String `tfsdk:"id"`
	AccountId    types.String `tfsdk:"account_id"`
	Name         types.String `tfsdk:"name"`
	Password     types.String `tfsdk:"password"`
	Phone        types.String `tfsdk:"phone"`
	Email        types.String `tfsdk:"email"`
	UserGroupIds types.Set    `tfsdk:"user_group_ids"`
	Description  types.String `tfsdk:"description"`
}
