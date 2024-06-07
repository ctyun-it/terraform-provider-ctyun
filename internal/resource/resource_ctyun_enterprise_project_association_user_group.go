package resource

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"sync"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
)

func NewCtyunEnterpriseProjectAssociationUserGroup() resource.Resource {
	return &ctyunEnterpriseProjectAssociationUserGroup{}
}

type ctyunEnterpriseProjectAssociationUserGroup struct {
	meta                     *common.CtyunMetadata
	enterpriseProjectService *business.EnterpriseProjectService
	userGroupService         *business.UserGroupService
	policyService            *business.PolicyService
}

func (c *ctyunEnterpriseProjectAssociationUserGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_enterprise_project_association_user_group"
}

func (c *ctyunEnterpriseProjectAssociationUserGroup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10345725/10356399**`,
		Attributes: map[string]schema.Attribute{
			"enterprise_project_id": schema.StringAttribute{
				Required:    true,
				Description: "企业项目id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_group_id": schema.StringAttribute{
				Required:    true,
				Description: "用户组id",
			},
			"policy_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: "策略id列表",
				ElementType: types.StringType,
				Validators:  []validator.Set{},
				Default:     setdefault.StaticValue(types.SetValueMust(basetypes.StringType{}, []attr.Value{})),
			},
		},
	}
}

func (c *ctyunEnterpriseProjectAssociationUserGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunEnterpriseProjectAssociationUserGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.check(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	err = c.attachGroupPolicy(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEnterpriseProjectAssociationUserGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunEnterpriseProjectAssociationUserGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.CtIamApis.EnterpriseProjectGetPolicyApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectGetPolicyRequest{
		GroupId:   state.UserGroupId.ValueString(),
		ProjectId: state.EnterpriseProjectId.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	if len(resp.List) == 0 {
		response.State.RemoveResource(ctx)
		return
	}

	var rawPolicyId []attr.Value
	for _, l := range resp.List {
		rawPolicyId = append(rawPolicyId, types.StringValue(l.Id))
	}
	state.PolicyIds = types.SetValueMust(types.StringType, rawPolicyId)
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *ctyunEnterpriseProjectAssociationUserGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan CtyunEnterpriseProjectAssociationUserGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)

	err := c.check(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	err = c.attachGroupPolicy(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEnterpriseProjectAssociationUserGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunEnterpriseProjectAssociationUserGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtIamApis.EnterpriseProjectRemoveGroupApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectRemoveGroupRequest{
		ProjectId: state.EnterpriseProjectId.ValueString(),
		GroupIds:  []string{state.UserGroupId.ValueString()},
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

func (c *ctyunEnterpriseProjectAssociationUserGroup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.enterpriseProjectService = business.NewEnterpriseProjectService(meta)
	c.userGroupService = business.NewUserGroupService(meta)
	c.policyService = business.NewPolicyService(meta)
}

// attachGroupPolicy 绑定用户组和策略的关系
func (c *ctyunEnterpriseProjectAssociationUserGroup) attachGroupPolicy(ctx context.Context, cfg CtyunEnterpriseProjectAssociationUserGroupConfig) error {
	var policyIds []string
	cfg.PolicyIds.ElementsAs(ctx, &policyIds, true)
	var err error
	if len(policyIds) == 0 {
		// 先删除，再进行更新
		_, err := c.meta.Apis.CtIamApis.EnterpriseProjectRemoveGroupApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectRemoveGroupRequest{
			ProjectId: cfg.EnterpriseProjectId.ValueString(),
			GroupIds:  []string{cfg.UserGroupId.ValueString()},
		})
		if err != nil {
			return err
		}

		_, err = c.meta.Apis.CtIamApis.EnterpriseProjectAssignmentToGroupApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectAssignmentToGroupRequest{
			GroupIds:  []string{cfg.UserGroupId.ValueString()},
			ProjectId: cfg.EnterpriseProjectId.ValueString(),
		})
	} else {
		_, err = c.meta.Apis.CtIamApis.EnterpriseProjectSetGroupPolicyApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectSetGroupPolicyRequest{
			GroupId:   cfg.UserGroupId.ValueString(),
			ProjectId: cfg.EnterpriseProjectId.ValueString(),
			PloyIds:   policyIds,
		})
	}
	return err
}

// check 校验
func (c *ctyunEnterpriseProjectAssociationUserGroup) check(ctx context.Context, cfg CtyunEnterpriseProjectAssociationUserGroupConfig) error {
	// 校验企业项目
	err := c.enterpriseProjectService.MustExist(ctx, cfg.EnterpriseProjectId.ValueString())
	if err != nil {
		return err
	}

	// 校验用户组
	err = c.userGroupService.MustExist(ctx, cfg.UserGroupId.ValueString())
	if err != nil {
		return err
	}

	// 校验策略
	var wg sync.WaitGroup
	var lock sync.Mutex

	policyIds := []types.String{}
	cfg.PolicyIds.ElementsAs(ctx, &policyIds, true)
	var es []error

	for _, id := range policyIds {
		wg.Add(1)
		go func(value string) {
			defer wg.Done()
			err = c.policyService.MustExist(ctx, value)
			if err != nil {
				lock.Lock()
				defer lock.Unlock()
				es = append(es, err)
			}

		}(id.ValueString())
	}
	wg.Wait()
	if len(es) > 0 {
		return errors.Join(es...)
	}

	return nil
}

type CtyunEnterpriseProjectAssociationUserGroupConfig struct {
	EnterpriseProjectId types.String `tfsdk:"enterprise_project_id"`
	UserGroupId         types.String `tfsdk:"user_group_id"`
	PolicyIds           types.Set    `tfsdk:"policy_ids"`
}
