package iam

import (
	"context"
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

func NewCtyunEnterpriseProject() resource.Resource {
	return &ctyunEnterpriseProject{}
}

type ctyunEnterpriseProject struct {
	meta *common.CtyunMetadata
}

func (c *ctyunEnterpriseProject) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_enterprise_project"
}

func (c *ctyunEnterpriseProject) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10345725/10358242`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "企业项目名称，长度为1-32",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 32),
				},
			},
			// "status": schema.StringAttribute{
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "企业项目状态，enable：启用，disable：禁用",
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf(business.EnterpriseProjectStatuses...),
			// 	},
			// 	Default: stringdefault.StaticString(business.EnterpriseProjectStatusEnable),
			// },
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目描述，长度最大为64",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(64),
				},
				Default: stringdefault.StaticString(""),
			},
		},
	}
}

func (c *ctyunEnterpriseProject) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunEnterpriseProjectConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.CtIamApis.EnterpriseProjectCreateApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectCreateRequest{
		ProjectName: plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.StringValue(resp.ProjectId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// if plan.Status.ValueString() != business.EnterpriseProjectStatusEnable {
	// 	err := c.changeStatus(ctx, resp.ProjectId, business.EnterpriseProjectStatusEnable)
	// 	if err != nil {
	// 		response.Diagnostics.AddError(err.Error(), err.Error())
	// 		return
	// 	}
	// }

	instance, ctyunRequestError := c.getAndMergeEnterpriseProject(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEnterpriseProject) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunEnterpriseProjectConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeEnterpriseProject(ctx, state)
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

func (c *ctyunEnterpriseProject) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunEnterpriseProjectConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan CtyunEnterpriseProjectConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtIamApis.EnterpriseProjectUpdateApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectUpdateRequest{
		Id:          state.Id.ValueString(),
		ProjectName: plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// if !plan.Status.Equal(state.Status) {
	// 	err := c.changeStatus(ctx, state.Id.ValueString(), plan.Status.ValueString())
	// 	if err != nil {
	// 		response.Diagnostics.AddError(err.Error(), err.Error())
	// 		return
	// 	}
	// }

	instance, ctyunRequestError := c.getAndMergeEnterpriseProject(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEnterpriseProject) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunEnterpriseProjectConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.changeStatus(ctx, state.Id.ValueString(), business.EnterpriseProjectStatusDisable)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [enterpriseProjectId]
func (c *ctyunEnterpriseProject) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunEnterpriseProjectConfig
	var enterpriseProjectId string
	err := terraform_extend.Split(request.ID, &enterpriseProjectId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(enterpriseProjectId)

	instance, err := c.getAndMergeEnterpriseProject(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEnterpriseProject) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getAndMergeEnterpriseProject 查询企业项目
func (c *ctyunEnterpriseProject) getAndMergeEnterpriseProject(ctx context.Context, cfg CtyunEnterpriseProjectConfig) (*CtyunEnterpriseProjectConfig, error) {
	resp, err := c.meta.Apis.CtIamApis.EnterpriseProjectGetApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectGetRequest{
		Id: cfg.Id.ValueString(),
	})
	if err != nil {
		return nil, err
	}

	if resp.Id == "" {
		return nil, nil
	}

	// status, err2 := business.EnterpriseProjectStatusMap.ToOriginalScene(resp.Status, business.EnterpriseProjectSceneRequest)
	// if err2 != nil {
	// 	return nil, err2
	// }
	cfg.Name = types.StringValue(resp.ProjectName)
	cfg.Description = types.StringValue(resp.Description)
	// cfg.Status = types.StringValue(status.(string))
	return &cfg, nil
}

// changeStatus 改变状态
func (c *ctyunEnterpriseProject) changeStatus(ctx context.Context, projectId string, statusTo string) error {
	//status, err := business.EnterpriseProjectStatusMap.FromOriginalScene(statusTo, business.EnterpriseProjectSceneRequest)
	//if err != nil {
	//	return err
	//}
	_, err := c.meta.Apis.CtIamApis.EnterpriseProjectStatusUpdateApi.Do(ctx, c.meta.Credential, &ctiam.EnterpriseProjectStatusUpdateRequest{
		ProjectId: projectId,
		Status:    3,
	})
	if err != nil {
		return err
	}
	return nil
}

type CtyunEnterpriseProjectConfig struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	// Status      types.String `tfsdk:"status"`
}
