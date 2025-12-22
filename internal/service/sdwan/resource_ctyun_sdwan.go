package sdwan

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sdwan"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"time"
)

var (
	_ resource.Resource                = &CtyunSdwan{}
	_ resource.ResourceWithConfigure   = &CtyunSdwan{}
	_ resource.ResourceWithImportState = &CtyunSdwan{}
)

func NewCtyunSdwan() resource.Resource {
	return &CtyunSdwan{}
}

type CtyunSdwan struct {
	meta *common.CtyunMetadata
}

type CtyunSdwanConfig struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	Desc      types.String `tfsdk:"description"`
}

func (c *CtyunSdwan) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdwan"
}

func (c *CtyunSdwan) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10035786/10035852`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, true),
				Validators: []validator.String{
					validator2.Project(),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "SD-WAN名称 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "SD-WAN描述 支持更新",
				Validators: []validator.String{
					validator2.Desc(),
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(`^\S*$`), "不能含有空格"),
				},
			},
		},
	}
}

func (c *CtyunSdwan) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSdwan) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunSdwanConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunSdwan) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSdwanConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 查询远端确认资源是否存在
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunSdwan) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunSdwanConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunSdwan) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSdwanConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunSdwan) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID]"
			resp.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunSdwanConfig
	config.ID = types.StringValue(req.ID)
	// 查询远端
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (c *CtyunSdwan) create(ctx context.Context, plan *CtyunSdwanConfig) (err error) {
	createReq := &sdwan.SdwanCreateSdwanRequest{
		SdwanName: plan.Name.ValueString(),
	}
	if plan.ProjectID.IsNull() || plan.ProjectID.IsUnknown() || plan.ProjectID.ValueString() == "" {
		createReq.ProjectID = "0"
	} else {
		createReq.ProjectID = plan.ProjectID.ValueString()
	}

	if !plan.Desc.IsNull() {

		createReq.Description = plan.Desc.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanCreateSdwanApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	// 等待创建完成
	time.Sleep(3 * time.Second)
	return
}

func (c *CtyunSdwan) getAndMerge(ctx context.Context, plan *CtyunSdwanConfig) (err error) {
	listReq := &sdwan.SdwanListSdwanRequest{
		PageNo:   1,
		PageSize: 1000,
	}
	if plan.ID.ValueStringPointer() != nil {
		listReq.SdwanID = plan.ID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkSdwanApis.SdwanListSdwanApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	// 查找对应的SD-WAN
	var found bool
	if plan.ID.ValueString() != "" {
		for _, sdwanItem := range resp.ReturnObj.Result {
			if sdwanItem.SdwanID != nil && *sdwanItem.SdwanID == plan.ID.ValueString() {
				plan.Name = types.StringValue(*sdwanItem.SdwanName)
				if sdwanItem.Description != nil {
					plan.Desc = types.StringValue(*sdwanItem.Description)
				}
				found = true
				break
			}
		}

		if !found {
			return common.ResourceNotExistError
		}
		return
	} else if plan.Name.ValueString() != "" {
		for _, sdwanItem := range resp.ReturnObj.Result {
			if sdwanItem.SdwanName != nil && *sdwanItem.SdwanName == plan.Name.ValueString() {
				plan.ID = types.StringValue(*sdwanItem.SdwanID)
				if sdwanItem.Description != nil {
					plan.Desc = types.StringValue(*sdwanItem.Description)
				}
				found = true
				break
			}
		}
		if !found {
			return common.ResourceNotExistError
		}
		return
	}

	return
}

func (c *CtyunSdwan) update(ctx context.Context, plan *CtyunSdwanConfig) (err error) {
	updateReq := &sdwan.SdwanUpdateSdwanRequest{
		SdwanID:   plan.ID.ValueString(),
		SdwanName: plan.Name.ValueStringPointer(),
	}

	if !plan.Desc.IsNull() {

		updateReq.Description = plan.Desc.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanUpdateSdwanApi.Do(ctx, c.meta.SdkCredential, updateReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	return
}

func (c *CtyunSdwan) delete(ctx context.Context, state *CtyunSdwanConfig) (err error) {
	deleteReq := &sdwan.SdwanDeleteSdwanRequest{
		SdwanID: state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanDeleteSdwanApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	// 等待删除完成
	time.Sleep(3 * time.Second)

	return
}
