package iam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunIamUserGroup() resource.Resource {
	return &ctyunIamUserGroup{}
}

type ctyunIamUserGroup struct {
	meta *common.CtyunMetadata
}

func (c *ctyunIamUserGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_iam_user_group"
}

func (c *ctyunIamUserGroup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10345725/10355805**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "用户组id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "用户组名称，长度1-32位",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 32),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "用户组描述，长度最大为64",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(64),
				},
			},
		},
	}
}

func (c *ctyunIamUserGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunIamUserGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.CtIamApis.UserGroupCreateApi.Do(ctx, c.meta.Credential, &ctiam.UserGroupCreateRequest{
		GroupName:  plan.Name.ValueString(),
		GroupIntro: plan.Description.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.StringValue(resp.Id)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err = c.meta.Apis.CtIamApis.EnterpriseProjectAssignmentToGroupApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectAssignmentToGroupRequest{
		ProjectId: "0",
		GroupIds:  []string{resp.Id},
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	instance, ctyunRequestError := c.getAndMergeIamUserGroup(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunIamUserGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunIamUserGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeIamUserGroup(ctx, state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, &instance)...)
}

func (c *ctyunIamUserGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunIamUserGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan CtyunIamUserGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtIamApis.UserGroupUpdateApi.Do(ctx, c.meta.Credential, &ctiam.UserGroupUpdateRequest{
		Id:         state.Id.ValueString(),
		GroupName:  plan.Name.ValueString(),
		GroupIntro: plan.Description.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	instance, ctyunRequestError := c.getAndMergeIamUserGroup(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunIamUserGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunIamUserGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtIamApis.UserGroupInvalidApi.Do(ctx, c.meta.Credential, &ctiam.UserGroupInvalidRequest{
		GroupId: state.Id.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [iamUserGroupId]
func (c *ctyunIamUserGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunIamUserGroupConfig
	var iamUserGroupId string
	err := terraform_extend.Split(request.ID, &iamUserGroupId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(iamUserGroupId)

	instance, err := c.getAndMergeIamUserGroup(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunIamUserGroup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getAndMergeIamUserGroup 查询iam用户组
func (c *ctyunIamUserGroup) getAndMergeIamUserGroup(ctx context.Context, cfg CtyunIamUserGroupConfig) (*CtyunIamUserGroupConfig, error) {
	resp, err := c.meta.Apis.CtIamApis.UserGroupGetApi.Do(ctx, c.meta.Credential, &ctiam.UserGroupGetRequest{
		GroupId: cfg.Id.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.CtiamNoPermission {
			return nil, nil
		}
		return nil, err
	}
	cfg.Name = types.StringValue(resp.GroupName)
	cfg.Description = types.StringValue(resp.GroupIntro)
	return &cfg, nil
}

type CtyunIamUserGroupConfig struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}
