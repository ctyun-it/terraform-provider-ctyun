package ec

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	explanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/planmodifier"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &CtyunExpressConnect{}
	_ resource.ResourceWithConfigure   = &CtyunExpressConnect{}
	_ resource.ResourceWithImportState = &CtyunExpressConnect{}
)

func NewCtyunExpressConnect() resource.Resource {
	return &CtyunExpressConnect{}
}

type CtyunExpressConnect struct {
	meta *common.CtyunMetadata
	name string
}

type CtyunExpressConnectConfig struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Status      types.Int64  `tfsdk:"status"`
	CreateTime  types.String `tfsdk:"create_time"`
	ProjectID   types.String `tfsdk:"project_id"`
}

func (c *CtyunExpressConnect) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_express_connect"
	c.name = resp.TypeName
}

func (c *CtyunExpressConnect) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理云间高速实例", "云间高速（标准版）（CT-EC, Express Connect Standard）", "https://www.ctyun.cn/document/10026763/10038220"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "云间高速实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "名称 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "描述信息 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"status": schema.Int64Attribute{
				Computed:    true,
				Description: "运行状态，取值范围: 0:创建中 2:运行中 18:删除中 21:设置中 22:更新带宽中 24:更新中",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					explanmodifier.Project(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
		},
	}
}

func (c *CtyunExpressConnect) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunExpressConnect) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunExpressConnectConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunExpressConnect) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunExpressConnectConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			resp.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunExpressConnect) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunExpressConnectConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunExpressConnect) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunExpressConnectConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *CtyunExpressConnect) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [id]", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunExpressConnectConfig
	config.ID = types.StringValue(request.ID)

	// 查询远端
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunExpressConnect) checkBeforeCreate(ctx context.Context, c2 *CtyunExpressConnectConfig) (err error) {
	return nil
}
func (c *CtyunExpressConnect) create(ctx context.Context, plan *CtyunExpressConnectConfig) (err error) {
	// 创建云间高速实例
	createReq := &ec.EcEcCreateRequest{
		EcName: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		createReq.EcDescription = plan.Description.ValueStringPointer()
	}
	if !plan.ProjectID.IsNull() {
		createReq.ProjectID = plan.ProjectID.ValueStringPointer()
	}
	tflog.Info(ctx, "创建云间高速实例", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcCreateApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	plan.ID = types.StringValue(*resp.ReturnObj.EcID)
	plan.Status = types.Int64Value(int64(*resp.ReturnObj.Status))
	return
}
func (c *CtyunExpressConnect) getAndMerge(ctx context.Context, plan *CtyunExpressConnectConfig) (err error) {
	// 查询云间高速实例
	listReq := &ec.EcEcListRequest{
		EcID:     plan.ID.ValueStringPointer(),
		PageNo:   func() *int32 { i := int32(1); return &i }(),
		PageSize: func() *int32 { i := int32(1); return &i }(),
	}

	resp, err := c.meta.Apis.SdkEcApis.EcEcListApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	} else if len(resp.ReturnObj.Results) == 0 {
		return common.ResourceNotExistError
	}
	result := resp.ReturnObj.Results[0]
	plan.Name = types.StringValue(*result.EcName)
	plan.Status = types.Int64Value(int64(*result.Status))
	plan.CreateTime = types.StringValue(utils.FromBJTimeToUTCZ(utils.SecString(result.CreateDate)))
	plan.ProjectID = types.StringValue(*result.Project)
	if result.EcDescription != nil && *result.EcDescription != "" {
		plan.Description = types.StringValue(*result.EcDescription)
	}

	return
}
func (c *CtyunExpressConnect) update(ctx context.Context, plan CtyunExpressConnectConfig) (err error) {
	// 更新云间高速实例
	updateReq := &ec.EcEcUpdateRequest{
		EcID:   plan.ID.ValueString(),
		EcName: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		updateReq.EcDescription = plan.Description.ValueStringPointer()
	}

	tflog.Info(ctx, "更新云间高速实例", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcUpdateApi.Do(ctx, c.meta.SdkCredential, updateReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	return
}
func (c *CtyunExpressConnect) delete(ctx context.Context, state CtyunExpressConnectConfig) (err error) {
	// 删除云间高速实例
	deleteReq := &ec.EcEcDeleteRequest{
		EcID: state.ID.ValueString(),
	}

	tflog.Info(ctx, "删除云间高速实例", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcDeleteApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	return
}
