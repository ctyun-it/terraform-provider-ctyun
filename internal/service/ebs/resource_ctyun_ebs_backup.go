package ebs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebsbackup "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebsbackup"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

/*
云硬盘备份
*/

func NewCtyunEbsBackup() resource.Resource {
	return &ctyunEbsBackup{}
}

type ctyunEbsBackup struct {
	meta       *common.CtyunMetadata
	ebsService *business.EbsService
}

func (c *ctyunEbsBackup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_backup"
}

type CtyunEbsBackupConfig struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	RegionID     types.String `tfsdk:"region_id"`
	Description  types.String `tfsdk:"description"`
	DiskID       types.String `tfsdk:"disk_id"`
	RepositoryID types.String `tfsdk:"repository_id"`
	FullBackup   types.Bool   `tfsdk:"full_backup"`

	// 返回字段
	InstanceName        types.String `tfsdk:"instance_name"`
	RepositoryName      types.String `tfsdk:"repository_name"`
	CreatedTime         types.String `tfsdk:"created_time"`
	ProjectID           types.String `tfsdk:"project_id"`
	BackupStatus        types.String `tfsdk:"backup_status"`
	DiskSize            types.Int64  `tfsdk:"disk_size"`
	BackupSize          types.Int64  `tfsdk:"backup_size"`
	UpdatedTime         types.String `tfsdk:"updated_time"`
	FinishedTime        types.String `tfsdk:"finished_time"`
	RestoredTime        types.String `tfsdk:"restored_time"`
	RestoreFinishedTime types.String `tfsdk:"restore_finished_time"`
	Freeze              types.Bool   `tfsdk:"freeze"`
	Encrypted           types.Bool   `tfsdk:"encrypted"`
	DiskType            types.String `tfsdk:"disk_type"`
	Paas                types.Bool   `tfsdk:"paas"`
	InstanceID          types.String `tfsdk:"instance_id"`
}

func (c *ctyunEbsBackup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026752/10037428`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "云硬盘备份ID",
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
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"disk_id": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘备份名称。满足以下规则：长度为2-63字符，头尾不支持输入空格",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 63),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "云硬盘备份描述",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"repository_id": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘备份存储库ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"full_backup": schema.BoolAttribute{
				Optional:    true,
				Description: "是否启用全量备份，取值范围：true：是，false：否。若启用该参数，则此次备份的类型为全量备份。注：只有4.0资源池支持该参数。",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},

			// 返回参数
			"backup_status": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘备份状态，该接口会返回creating状态",
			},
			"disk_size": schema.Int64Attribute{
				Computed:    true,
				Description: "云硬盘大小，单位GB",
			},
			"backup_size": schema.Int64Attribute{
				Computed:    true,
				Description: "云硬盘备份大小，单位Byte",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "备份更新时间",
			},
			"finished_time": schema.StringAttribute{
				Computed:    true,
				Description: "备份完成时间",
			},
			"restored_time": schema.StringAttribute{
				Computed:    true,
				Description: "使用该云硬盘备份恢复数据时间",
			},
			"restore_finished_time": schema.StringAttribute{
				Computed:    true,
				Description: "使用该云硬盘备份恢复完成时间",
			},
			"freeze": schema.BoolAttribute{
				Computed:    true,
				Description: "备份是否冻结",
			},
			"encrypted": schema.BoolAttribute{
				Computed:    true,
				Description: "云硬盘是否加密",
			},
			"disk_type": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘类型，取值范围为：SATA：普通IO。SAS：高IO。SSD：超高IO。FAST-SSD：极速型SSD。XSSD-0、XSSD-1、XSSD-2：X系列云硬盘",
			},
			"paas": schema.BoolAttribute{
				Computed:    true,
				Description: "是否支持PAAS",
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘挂载的云主机ID",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Description: "企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看创建企业项目了解如何创建企业项目。注：默认值为\"0\"",
			},
			"instance_name": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘挂载的云主机名称",
			},
			"repository_name": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘备份存储库名称",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
		},
	}
}

func (c *ctyunEbsBackup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbsBackupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 校验创建动作的前置条件
	err = c.checkCreate(ctx, plan)
	if err != nil {
		return
	}

	// 实际创建
	id, err := c.create(ctx, &plan)
	if err != nil {
		return
	}
	plan.Id = types.StringValue(id)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)

	//轮询状态为可用状态
	err = c.StartedLoop(ctx, &plan)
	if err != nil {
		return
	}
	// 查询信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}
func (c *ctyunEbsBackup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

// getAndMerge 查询
func (c *ctyunEbsBackup) getAndMerge(ctx context.Context, cfg *CtyunEbsBackupConfig) (err error) {

	params := &ctebsbackup.EbsbackupShowBackupRequest{
		RegionID: cfg.RegionID.ValueString(),
		BackupID: cfg.Id.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupShowBackupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	//资源返回内容更新
	result := resp.ReturnObj
	cfg.Name = types.StringValue(result.BackupName)
	cfg.BackupStatus = types.StringValue(result.BackupStatus)
	// 处理 description 字段，确保空值处理一致
	if result.Description == "" {
		cfg.Description = types.StringNull()
	} else {
		cfg.Description = types.StringValue(result.Description)
	}

	cfg.DiskID = types.StringValue(result.DiskID)
	cfg.RepositoryID = types.StringValue(result.RepositoryID)
	cfg.RepositoryName = types.StringValue(result.RepositoryName)
	cfg.DiskSize = types.Int64Value(int64(result.DiskSize))
	cfg.BackupSize = types.Int64Value(int64(result.BackupSize))
	cfg.CreatedTime = types.StringValue(fmt.Sprintf("%d", result.CreatedTime))
	cfg.ProjectID = types.StringValue(result.ProjectID)

	// 新增字段处理
	cfg.UpdatedTime = types.StringValue(fmt.Sprintf("%d", result.UpdatedTime))
	cfg.FinishedTime = types.StringValue(fmt.Sprintf("%d", result.FinishedTime))
	cfg.RestoredTime = types.StringValue(fmt.Sprintf("%d", result.RestoredTime))
	cfg.RestoreFinishedTime = types.StringValue(fmt.Sprintf("%d", result.RestoreFinishedTime))
	cfg.Freeze = types.BoolValue(*result.Freeze)
	cfg.Encrypted = types.BoolValue(*result.Encrypted)
	cfg.DiskType = types.StringValue(result.DiskType)
	cfg.Paas = types.BoolValue(*result.Paas)
	cfg.InstanceID = types.StringValue(result.InstanceID)
	cfg.InstanceName = types.StringValue(result.InstanceName)
	return
}

func (c *ctyunEbsBackup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsBackupConfig
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

func (c *ctyunEbsBackup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsBackupConfig
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

func (c *ctyunEbsBackup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ebsService = business.NewEbsService(meta)
}

// create 创建
func (c *ctyunEbsBackup) create(ctx context.Context, cfg *CtyunEbsBackupConfig) (id string, err error) {

	// 处理全量备份参数
	var fullBackup *bool
	if !cfg.FullBackup.IsNull() && !cfg.FullBackup.IsUnknown() {
		fullBackupVal := cfg.FullBackup.ValueBool()
		fullBackup = &fullBackupVal
	}
	params := &ctebsbackup.EbsbackupCreateBackupRequest{
		RegionID:     cfg.RegionID.ValueString(),
		DiskID:       cfg.DiskID.ValueString(),
		BackupName:   cfg.Name.ValueString(),
		Description:  cfg.Description.ValueString(),
		RepositoryID: cfg.RepositoryID.ValueString(),
		FullBackup:   fullBackup,
	}
	// 调用API
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupCreateBackupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	id = resp.ReturnObj.BackupID
	return
}

func (c *ctyunEbsBackup) checkCreate(ctx context.Context, plan CtyunEbsBackupConfig) error {
	// 1.云硬盘必须存在
	err := c.ebsService.MustExist(ctx, plan.DiskID.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return err
	}
	// 2.云硬盘必须处于“未挂载”或“已挂载”状态 ——资料描述错误 应该是available状态
	resp, err := c.ebsService.GetEbsInfo(ctx, plan.DiskID.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return err
	}
	if resp.DiskStatus != "available" {
		return fmt.Errorf("云硬盘必须处于“未挂载”或“已挂载”状态")
	}

	// 3.云硬盘备份存储库的状态为可用状态允许备份 (创建接口会返回 不需提前校验)

	return nil
}

// delete 删除
func (c *ctyunEbsBackup) delete(ctx context.Context, plan CtyunEbsBackupConfig) (err error) {
	if plan.Id.ValueString() == "" {
		return fmt.Errorf("instance backup ID is required for deletion")
	}
	params := &ctebsbackup.EbsbackupDeleteEbsBackupRequest{
		RegionID: plan.RegionID.ValueString(),
		BackupID: plan.Id.ValueString(),
	}
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupDeleteEbsBackupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}
func (c *ctyunEbsBackup) StartedLoop(ctx context.Context, state *CtyunEbsBackupConfig, loopCount ...int) (err error) {
	count := 30
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			params := &ctebsbackup.EbsbackupShowBackupRequest{
				RegionID: state.RegionID.ValueString(),
				BackupID: state.Id.ValueString(),
			}
			// 调用API
			resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupShowBackupApi.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			runningStatus := strings.ToLower(resp.ReturnObj.BackupStatus)
			if runningStatus == business.EcsSnapshotStatusAvailable {
				return false
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未到达启动状态！")
	}
	return
}

func (c *ctyunEbsBackup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEbsBackupConfig
	var id, regionID string
	err = terraform_extend.Split(request.ID, &id, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.Id = types.StringValue(id)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}
