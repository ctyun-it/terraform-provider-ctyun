package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgdcs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &ctyunRedisMigrationTask{}
	_ resource.ResourceWithConfigure   = &ctyunRedisMigrationTask{}
	_ resource.ResourceWithImportState = &ctyunRedisMigrationTask{}
)

type ctyunRedisMigrationTask struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisMigrationTask() resource.Resource {
	return &ctyunRedisMigrationTask{}
}

func (c *ctyunRedisMigrationTask) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_migration_task"
}

type CtyunRedisMigrationTaskConfig struct {
	ID           types.String   `tfsdk:"id"`
	RegionId     types.String   `tfsdk:"region_id"`
	SyncMode     types.Int32    `tfsdk:"sync_mode"`
	ConflictMode types.Int32    `tfsdk:"conflict_mode"`
	Status       types.Int32    `tfsdk:"status"`
	CreateTime   types.String   `tfsdk:"create_time"`
	CompleteTime types.String   `tfsdk:"complete_time"`
	SourceDbInfo InstanceDbInfo `tfsdk:"source_db_info"`
	TargetDbInfo InstanceDbInfo `tfsdk:"target_db_info"`
	OperateType  types.Int32    `tfsdk:"operate_type"`
}

type InstanceDbInfo struct {
	SpuInstId       types.String `tfsdk:"spu_inst_id"`
	IpAddr          types.String `tfsdk:"ip_addr"`
	OriginalCluster types.Bool   `tfsdk:"original_cluster"`
	AccountName     types.String `tfsdk:"account_name"`
	Password        types.String `tfsdk:"password"`
}

func (c *ctyunRedisMigrationTask) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10518385`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "任务ID",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
			"sync_mode": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "同步模式 1：全量同步+增量同步 2：全量同步",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.OneOf(1, 2),
				},
			},
			"conflict_mode": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "数据冲突时的处理办法 1：中断迁移 2：跳过目标key，继续执行 3：覆盖目标key，继续执行",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.OneOf(1, 2, 3),
				},
			},
			"status": schema.Int32Attribute{
				Computed:    true,
				Description: "任务状态 0：初始态 1：运行中 2：成功 3：失败",
				Validators: []validator.Int32{
					int32validator.OneOf(0, 1, 2, 3),
				},
			},
			"operate_type": schema.Int32Attribute{
				Optional:    true,
				Description: "操作类型，可选值：2：结束运行中的任务，3：删除成功或者失败的任务记录 支持更新",
				Validators: []validator.Int32{
					int32validator.OneOf(2, 3),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间(秒)",
			},
			"complete_time": schema.StringAttribute{
				Computed:    true,
				Description: "结束时间（秒，-1表示时间未知）",
			},
			"source_db_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"spu_inst_id": schema.StringAttribute{
						Required:    true,
						Description: "实例ID",
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"ip_addr": schema.StringAttribute{
						Required:    true,
						Description: "连接地址",
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"original_cluster": schema.BoolAttribute{
						Optional:    true,
						Description: "是否是原生cluster集群 输入实例ID可不填，否则必填。",
					},
					"account_name": schema.StringAttribute{
						Required:    true,
						Description: "数据库账号",
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"password": schema.StringAttribute{
						Required:    true,
						Description: "数据库密码",
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
						Sensitive: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				Required:    true,
				Description: "源数据库信息",
			},
			"target_db_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"spu_inst_id": schema.StringAttribute{
						Required:    true,
						Description: "实例ID",
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"ip_addr": schema.StringAttribute{
						Required:    true,
						Description: "连接地址",
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"original_cluster": schema.BoolAttribute{
						Optional:    true,
						Description: "是否是原生cluster集群",
					},
					"account_name": schema.StringAttribute{
						Required:    true,
						Description: "数据库账号",
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"password": schema.StringAttribute{
						Required:    true,
						Description: "数据库密码",
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
						Sensitive: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				Required:    true,
				Description: "目标数据库信息",
			},
		},
	}
}

func (c *ctyunRedisMigrationTask) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRedisMigrationTaskConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建迁移任务
	id, err := c.createMigrationTask(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(id)
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRedisMigrationTask) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisMigrationTaskConfig
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

func (c *ctyunRedisMigrationTask) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunRedisMigrationTaskConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// state中的
	var state CtyunRedisMigrationTaskConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 检查 operate_type 字段，执行相应的操作
	if !plan.OperateType.IsNull() && !plan.OperateType.IsUnknown() {
		operateType := plan.OperateType.ValueInt32()
		switch operateType {
		case 2: // 结束运行中的任务
			if state.Status.ValueInt32() == 1 { // 只有运行中的任务才能结束
				err = c.operateTransferTask(ctx, state, 2)
				if err != nil {
					return
				}
			}
		case 3: // 删除成功或者失败的任务记录
			if state.Status.ValueInt32() == 2 || state.Status.ValueInt32() == 3 { // 成功或失败的任务
				err = c.operateTransferTask(ctx, state, 3)
				if err != nil {
					return
				}
				// 删除成功后从状态中移除资源
				response.State.RemoveResource(ctx)
				return
			}
		}
	}

	state.OperateType = plan.OperateType
	// 更新后查询最新状态
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	// 迁移任务不支持更新，直接返回
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisMigrationTask) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisMigrationTaskConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除迁移任务记录
	err = c.deleteMigrationTask(ctx, state)
	if err != nil {
		return
	}
	response.State.RemoveResource(ctx)
}

func (c *ctyunRedisMigrationTask) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunRedisMigrationTask) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [regionID]/[taskID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()

	var cfg CtyunRedisMigrationTaskConfig

	var regionId, taskId string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		taskId = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &taskId, &regionId)
		if err != nil {
			return
		}
	}

	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if taskId == "" {
		err = fmt.Errorf("taskID不能为空")
		return
	}
	cfg.RegionId = types.StringValue(regionId)
	cfg.ID = types.StringValue(taskId)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// createMigrationTask 创建迁移任务
func (c *ctyunRedisMigrationTask) createMigrationTask(ctx context.Context, plan CtyunRedisMigrationTaskConfig) (id string, err error) {
	// 构建创建迁移任务请求
	sourceDbInfo := &ctgdcs2.Dcs2CreateTransferTaskSourceDbInfoRequest{}
	targetDbInfo := &ctgdcs2.Dcs2CreateTransferTaskTargetDbInfoRequest{}

	sourceDbInfo.SpuInstId = plan.SourceDbInfo.SpuInstId.ValueString()
	sourceDbInfo.IpAddr = plan.SourceDbInfo.IpAddr.ValueString()
	sourceDbInfo.AccountName = plan.SourceDbInfo.AccountName.ValueString()
	sourceDbInfo.Password = plan.SourceDbInfo.Password.ValueString()
	// 处理 original_cluster 字段：如果提供了实例ID，则该字段可选，否则必填
	if !plan.SourceDbInfo.OriginalCluster.IsNull() {
		val := plan.SourceDbInfo.OriginalCluster.ValueBool()
		sourceDbInfo.OriginalCluster = &val
	} else if plan.SourceDbInfo.SpuInstId.IsNull() {
		// 如果没有提供实例ID，original_cluster 是必填的
		err = fmt.Errorf("source_db_info.original_cluster is required when spu_inst_id is not provided")
		return
	}

	targetDbInfo.SpuInstId = plan.TargetDbInfo.SpuInstId.ValueString()
	targetDbInfo.IpAddr = plan.TargetDbInfo.IpAddr.ValueString()
	targetDbInfo.AccountName = plan.TargetDbInfo.AccountName.ValueString()
	targetDbInfo.Password = plan.TargetDbInfo.Password.ValueString()

	// 处理 original_cluster 字段：如果提供了实例ID，则该字段可选，否则必填
	if !plan.TargetDbInfo.OriginalCluster.IsNull() {
		val := plan.TargetDbInfo.OriginalCluster.ValueBool()
		targetDbInfo.OriginalCluster = &val
	} else if plan.TargetDbInfo.SpuInstId.IsNull() {
		// 如果没有提供实例ID，original_cluster 是必填的
		err = fmt.Errorf("target_db_info.original_cluster is required when spu_inst_id is not provided")
		return
	}

	params := &ctgdcs2.Dcs2CreateTransferTaskRequest{
		RegionId:     plan.RegionId.ValueString(),
		SyncMode:     plan.SyncMode.ValueInt32(),
		ConflictMode: plan.ConflictMode.ValueInt32(),
		SourceDbInfo: sourceDbInfo,
		TargetDbInfo: targetDbInfo,
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2CreateTransferTaskApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	plan.ID = types.StringValue(resp.ReturnObj.TaskId)
	return plan.ID.ValueString(), nil
}

// operateTransferTask 结束运行中的迁移任务
func (c *ctyunRedisMigrationTask) operateTransferTask(ctx context.Context, state CtyunRedisMigrationTaskConfig, operateType int32) (err error) {
	params := &ctgdcs2.Dcs2OperateTransferTaskRequest{
		RegionId:    state.RegionId.ValueString(),
		TaskId:      state.ID.ValueString(),
		OperateType: operateType,
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2OperateTransferTaskApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// deleteMigrationTask 删除迁移任务
func (c *ctyunRedisMigrationTask) deleteMigrationTask(ctx context.Context, state CtyunRedisMigrationTaskConfig) (err error) {
	params := &ctgdcs2.Dcs2DeleteTransferTaskRequest{
		RegionId:   state.RegionId.ValueString(),
		TaskIdList: []string{state.ID.ValueString()},
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DeleteTransferTaskApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// getAndMerge 从远端查询迁移任务信息
func (c *ctyunRedisMigrationTask) getAndMerge(ctx context.Context, state *CtyunRedisMigrationTaskConfig) (err error) {
	// 调用API查询迁移任务详情
	params := &ctgdcs2.Dcs2GetTransferTaskInfoRequest{
		RegionId: state.RegionId.ValueString(),
		TaskId:   state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2GetTransferTaskInfoApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 更新state中的信息
	state.ID = types.StringValue(resp.ReturnObj.TaskId)
	state.Status = types.Int32Value(resp.ReturnObj.Status)

	state.CreateTime = types.StringValue(utils.FromUnixToUTC(resp.ReturnObj.CreateTime))
	state.CompleteTime = types.StringValue(utils.FromUnixToUTC(resp.ReturnObj.CompleteTime))
	state.SyncMode = types.Int32Value(resp.ReturnObj.SyncMode)
	state.ConflictMode = types.Int32Value(resp.ReturnObj.ConflictMode)
	// 处理源数据库信息
	if resp.ReturnObj.Detail != nil && resp.ReturnObj.Detail.SourceDbInfo != nil {
		sourceDbInfo := resp.ReturnObj.Detail.SourceDbInfo
		state.SourceDbInfo.SpuInstId = types.StringValue(sourceDbInfo.SpuInstId)
		state.SourceDbInfo.IpAddr = types.StringValue(sourceDbInfo.IpAddr)
		state.SourceDbInfo.AccountName = types.StringValue(sourceDbInfo.AccountName)
		// 处理 OriginalCluster 字段
		if sourceDbInfo.OriginalCluster != nil {
			state.SourceDbInfo.OriginalCluster = types.BoolValue(*sourceDbInfo.OriginalCluster)
		} else {
			state.SourceDbInfo.OriginalCluster = types.BoolNull()
		}
	}

	// 处理目标数据库信息
	if resp.ReturnObj.Detail != nil && resp.ReturnObj.Detail.TargetDbInfo != nil {
		targetDbInfo := resp.ReturnObj.Detail.TargetDbInfo
		state.TargetDbInfo.SpuInstId = types.StringValue(targetDbInfo.SpuInstId)
		state.TargetDbInfo.IpAddr = types.StringValue(targetDbInfo.IpAddr)
		state.TargetDbInfo.AccountName = types.StringValue(targetDbInfo.AccountName)
		// 处理 OriginalCluster 字段
		if targetDbInfo.OriginalCluster != nil {
			state.TargetDbInfo.OriginalCluster = types.BoolValue(*targetDbInfo.OriginalCluster)
		} else {
			state.TargetDbInfo.OriginalCluster = types.BoolNull()
		}
	}

	return
}
