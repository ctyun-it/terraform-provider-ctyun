package mongodb

import (
	"context"
	"fmt"
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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &CtyunMongodbBackupResource{}
	_ resource.ResourceWithConfigure   = &CtyunMongodbBackupResource{}
	_ resource.ResourceWithImportState = &CtyunMongodbBackupResource{}
)

func NewCtyunMongodbBackupResource() resource.Resource {
	return &CtyunMongodbBackupResource{}
}

// CtyunMongodbBackupResource defines the resource implementation.
type CtyunMongodbBackupResource struct {
	meta *common.CtyunMetadata
}

// CtyunMongodbBackupConfig describes the resource data model.
type CtyunMongodbBackupConfig struct {
	ID          types.String `tfsdk:"id"`
	RegionID    types.String `tfsdk:"region_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	InstanceID  types.String `tfsdk:"instance_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *CtyunMongodbBackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_backup"
}

func (r *CtyunMongodbBackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10089535`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源ID，格式为region_id:inst_id:backup_name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "备份名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "备份描述",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
					validator2.Desc(),
				},
			},
		},
	}
}

func (r *CtyunMongodbBackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	r.meta = meta
}

func (r *CtyunMongodbBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMongodbBackupConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = r.checkBeforeCreate(ctx, &plan)
	if err != nil {
		return
	}
	err = r.create(ctx, &plan)
	if err != nil {
		return
	}
	err = r.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	// 保存数据到Terraform状态
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CtyunMongodbBackupResource) Read(ctx context.Context, req resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMongodbBackupConfig
	response.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = r.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *CtyunMongodbBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// 备份资源不支持更新
}

func (r *CtyunMongodbBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMongodbBackupConfig

	// Read Terraform prior state plan into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = r.delete(ctx, plan)
	if err != nil {
		return
	}
}

func (r *CtyunMongodbBackupResource) delete(ctx context.Context, plan CtyunMongodbBackupConfig) (err error) {
	// 删除备份
	deleteReq := &mongodb.MongodbDeleteBackupRequest{
		ProdInstId: plan.InstanceID.ValueString(),
		BackupId:   plan.ID.ValueString(),
	}

	header := &mongodb.MongodbDeleteBackupRequestHeaders{
		RegionID: plan.RegionID.ValueString(),
	}

	if !plan.ProjectID.IsNull() {
		header.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	resp, err := r.meta.Apis.SdkMongodbApis.MongodbDeleteBackupApi.Do(ctx, r.meta.Credential, deleteReq, header)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	return
}

func (c *CtyunMongodbBackupResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[instID],[projectID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunMongodbBackupConfig
	var name, regionId, projectId, instId string

	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		projectId = c.meta.GetExtraIfEmpty(projectId, common.ExtraProjectId)
		err = terraform_extend.Split(request.ID, &name, &instId)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &name, &instId, &projectId, &regionId)
		if err != nil {
			return
		}
	}
	if name == "" {
		err = fmt.Errorf("name 不能为空")
		return
	}
	if instId == "" {
		err = fmt.Errorf("instID 不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID 不能为空")
		return
	}
	cfg.RegionID = types.StringValue(regionId)
	cfg.ProjectID = types.StringValue(projectId)
	cfg.Name = types.StringValue(name)
	cfg.InstanceID = types.StringValue(instId)
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)

}

func (r *CtyunMongodbBackupResource) checkBeforeCreate(ctx context.Context, c *CtyunMongodbBackupConfig) (err error) {
	return

}
func (r *CtyunMongodbBackupResource) create(ctx context.Context, plan *CtyunMongodbBackupConfig) (err error) {
	// Create backup
	createReq := &mongodb.MongodbCreateBackupRequest{
		ProdInstId: plan.InstanceID.ValueString(),
	}

	createReq.BackupName = plan.Name.ValueStringPointer()

	if !plan.Description.IsNull() {
		createReq.Description = plan.Description.ValueStringPointer()
	}

	header := &mongodb.MongodbCreateBackupRequestHeaders{
		RegionID: plan.RegionID.ValueString(),
	}

	if !plan.ProjectID.IsNull() {
		header.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "开始创建MongoDB备份", map[string]interface{}{
		"instance_id": plan.InstanceID.ValueString(),
	})

	resp, err := r.meta.Apis.SdkMongodbApis.MongodbCreateBackupApi.Do(ctx, r.meta.Credential, createReq, header)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	//plan.ID = types.StringValue(plan.InstanceID.ValueString() + "," + plan.BackupName.ValueString())

	return
}
func (r *CtyunMongodbBackupResource) getAndMerge(ctx context.Context, plan *CtyunMongodbBackupConfig) (err error) {
	// 获取备份信息
	describeReq := &mongodb.MongodbDescribeBackupsRequest{
		ProdInstId: plan.InstanceID.ValueString(),
		BackupName: plan.Name.ValueStringPointer(),
		PageNow:    1,
		PageSize:   1000,
	}

	header := &mongodb.MongodbDescribeBackupsRequestHeaders{
		RegionID: plan.RegionID.ValueString(),
	}

	if !plan.ProjectID.IsNull() {
		header.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	resp, err := r.meta.Apis.SdkMongodbApis.MongodbDescribeBackupsApi.Do(ctx, r.meta.Credential, describeReq, header)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	// 查找备份信息
	var backupInfo *mongodb.MongodbBackupInfo
	for _, item := range resp.ReturnObj.List {
		if item.BackupName == plan.Name.ValueString() {
			backupInfo = &item
			break
		}
	}
	if backupInfo == nil {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%d", backupInfo.BackupId))
	if backupInfo.BackupName != "" {
		plan.Name = types.StringValue(backupInfo.BackupName)
	}
	// 添加 Description 字段的空指针检查
	if backupInfo.Description != nil {
		plan.Description = types.StringValue(*backupInfo.Description)
	}
	return
}
