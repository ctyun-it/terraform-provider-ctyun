package ecs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
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
云主机备份
*/

func NewCtyunEcsBackup() resource.Resource {
	return &ctyunEcsBackup{}
}

type ctyunEcsBackup struct {
	meta       *common.CtyunMetadata
	ecsService *business.EcsService
}

func (c *ctyunEcsBackup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_backup"
}

type CtyunEcsBackupConfig struct {
	Id                        types.String `tfsdk:"id"`
	RegionID                  types.String `tfsdk:"region_id"`
	InstanceID                types.String `tfsdk:"instance_id"`
	InstanceBackupName        types.String `tfsdk:"name"`
	InstanceBackupDescription types.String `tfsdk:"instance_backup_description"`
	RepositoryID              types.String `tfsdk:"repository_id"`

	// 返回字段
	InstanceBackupStatus types.String `tfsdk:"instance_backup_status"`
	InstanceName         types.String `tfsdk:"instance_name"`
	RepositoryName       types.String `tfsdk:"repository_name"`
	DiskTotalSize        types.Int64  `tfsdk:"disk_total_size"`
	UsedSize             types.Int64  `tfsdk:"used_size"`
	CreatedTime          types.String `tfsdk:"created_time"`
	BackupType           types.String `tfsdk:"backup_type"`
	FullBackup           types.Bool   `tfsdk:"full_backup"`
	ProjectID            types.String `tfsdk:"project_id"`
}

func (c *ctyunEcsBackup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026751/10033761**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "云主机备份ID",
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
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "云主机备份名称。满足以下规则：长度为2-63字符，头尾不支持输入空格。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 63),
				},
			},
			"instance_backup_description": schema.StringAttribute{
				Optional:    true,
				Description: "云主机备份描述，字符长度不超过256字符。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"repository_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机备份存储库ID",
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

			// 返回字段
			"project_id": schema.StringAttribute{
				Computed:    true,
				Description: "企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看创建企业项目了解如何创建企业项目。注：默认值为\"0\"",
			},
			"instance_backup_status": schema.StringAttribute{
				Computed:    true,
				Description: "备份状态，取值范围：CREATING: 备份创建中, ACTIVE: 可用， RESTORING: 备份恢复中，DELETING: 删除中，EXPIRED：到期，ERROR：错误",
			},
			"instance_name": schema.StringAttribute{
				Computed:    true,
				Description: "云主机名称",
			},
			"repository_name": schema.StringAttribute{
				Computed:    true,
				Description: "云主机备份存储库名称",
			},
			"disk_total_size": schema.Int64Attribute{
				Computed:    true,
				Description: "云盘总容量大小，单位为GB",
			},
			"used_size": schema.Int64Attribute{
				Computed:    true,
				Description: "云硬盘备份已使用大小",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"backup_type": schema.StringAttribute{
				Computed:    true,
				Description: "备份类型",
			},
		},
	}
}

func (c *ctyunEcsBackup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcsBackupConfig
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

func (c *ctyunEcsBackup) checkCreate(ctx context.Context, plan CtyunEcsBackupConfig) error {
	// 1.云主机和备份存储库必须存在
	err := c.ecsService.MustExist(ctx, plan.InstanceID.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return err
	}
	// 2.云主机的状态必须处于开机或关机状态下
	status, err := c.ecsService.GetEcsStatus(ctx, plan.InstanceID.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	allowedStatuses := map[string]bool{
		business.EcsStatusRunning: true,
		business.EcsStatusStopped: true,
	}

	if !allowedStatuses[status] {
		return fmt.Errorf("云主机状态无效(当前:%s)，仅允许在%s或%s状态下创建云主机备份",
			status, business.EcsStatusRunning, business.EcsStatusStopped)
	}

	// 3.当备份存储库的容量小于0，已到期或已冻结这些情况下，备份存储库不可用 （创建接口返回 不需提前校验）
	//4.云主机盘限制：不可含有本地盘、共享盘、ISCSI磁盘模式盘
	//5. 云主机所挂全部盘的状态需要为"已挂载"
	return nil
}

// getAndMerge 查询
func (c *ctyunEcsBackup) getAndMerge(ctx context.Context, cfg *CtyunEcsBackupConfig) (err error) {
	params := &ctecs2.CtecsDetailsInstanceBackupV41Request{
		RegionID:         cfg.RegionID.ValueString(),
		InstanceBackupID: cfg.Id.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsDetailsInstanceBackupV41Api.Do(ctx, c.meta.SdkCredential, params)
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
	cfg.InstanceBackupName = types.StringValue(result.InstanceBackupName)
	cfg.InstanceBackupStatus = types.StringValue(result.InstanceBackupStatus)
	// 处理 instance_backup_description 字段，确保空值处理一致
	if result.InstanceBackupDescription == "" {
		cfg.InstanceBackupDescription = types.StringNull()
	} else {
		cfg.InstanceBackupDescription = types.StringValue(result.InstanceBackupDescription)
	}

	cfg.InstanceID = types.StringValue(result.InstanceID)
	cfg.InstanceName = types.StringValue(result.InstanceName)
	cfg.RepositoryID = types.StringValue(result.RepositoryID)
	cfg.RepositoryName = types.StringValue(result.RepositoryName)
	cfg.DiskTotalSize = types.Int64Value(int64(result.DiskTotalSize))
	cfg.UsedSize = types.Int64Value(result.UsedSize)
	cfg.CreatedTime = types.StringValue(result.CreatedTime)
	cfg.ProjectID = types.StringValue(result.ProjectID)
	cfg.BackupType = types.StringValue(result.BackupType)
	return
}

func (c *ctyunEcsBackup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsBackupConfig
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

func (c *ctyunEcsBackup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan, state CtyunEcsBackupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 准备更新参数
	params := &ctecs2.CtecsUpdateInstanceBackupV41Request{
		RegionID:                  state.RegionID.ValueString(),
		InstanceBackupID:          state.Id.ValueString(),
		InstanceBackupName:        plan.InstanceBackupName.ValueString(),
		InstanceBackupDescription: plan.InstanceBackupDescription.ValueString(),
	}

	// 调用API更新云主机备份
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsUpdateInstanceBackupV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 查询更新后的信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEcsBackup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsBackupConfig
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

func (c *ctyunEcsBackup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ecsService = business.NewEcsService(meta)
}

// create 创建
func (c *ctyunEcsBackup) create(ctx context.Context, plan *CtyunEcsBackupConfig) (id string, err error) {

	params := &ctecs2.CtecsCreateInstanceBackupV41Request{
		RegionID:                  plan.RegionID.ValueString(),
		InstanceID:                plan.InstanceID.ValueString(),
		InstanceBackupName:        plan.InstanceBackupName.ValueString(),
		InstanceBackupDescription: plan.InstanceBackupDescription.ValueString(),
		RepositoryID:              plan.RepositoryID.ValueString(),
		FullBackup:                plan.FullBackup.ValueBool(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsCreateInstanceBackupV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.Results == nil {
		err = common.InvalidReturnObjResultsError
		return
	}

	id = resp.ReturnObj.Results.InstanceBackupID
	return
}

// delete 删除
func (c *ctyunEcsBackup) delete(ctx context.Context, plan CtyunEcsBackupConfig) (err error) {
	if plan.Id.ValueString() == "" {
		return fmt.Errorf("instance backup ID is required for deletion")
	}
	params := &ctecs2.CtecsDeleteInstanceBackupRequest{
		RegionID:         plan.RegionID.ValueString(),
		InstanceBackupID: plan.Id.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsDeleteInstanceBackupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *ctyunEcsBackup) StartedLoop(ctx context.Context, state *CtyunEcsBackupConfig, loopCount ...int) (err error) {
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
			params := &ctecs2.CtecsDetailsInstanceBackupV41Request{
				RegionID:         state.RegionID.ValueString(),
				InstanceBackupID: state.Id.ValueString(),
			}
			// 调用API
			resp, err := c.meta.Apis.SdkCtEcsApis.CtecsDetailsInstanceBackupV41Api.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			runningStatus := strings.ToLower(resp.ReturnObj.InstanceBackupStatus)
			if runningStatus == business.EcsBackupRepoStatusActive {
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

func (c *ctyunEcsBackup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEcsBackupConfig
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
