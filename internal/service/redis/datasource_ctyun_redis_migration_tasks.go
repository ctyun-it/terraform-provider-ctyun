package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgdcs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &ctyunRedisMigrationTasks{}
)

type ctyunRedisMigrationTasks struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisMigrationTasks() datasource.DataSource {
	return &ctyunRedisMigrationTasks{}
}

func (c *ctyunRedisMigrationTasks) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_migration_tasks"
}

type CtyunRedisMigrationTasksModel struct {
	ID       types.String `tfsdk:"id"`
	RegionId types.String `tfsdk:"region_id"`
	PageNum  types.Int32  `tfsdk:"page_num"`
	PageSize types.Int32  `tfsdk:"page_size"`
	Status   types.String `tfsdk:"status"`

	// 详情字段
	DataSyncCountInfo      *DataSyncCountInfoModel `tfsdk:"data_sync_count_info"`
	SourceProgressInfoList types.List              `tfsdk:"source_progress_info_list"`

	// 列表字段
	List  types.List  `tfsdk:"list"`
	Total types.Int32 `tfsdk:"total"`
	Size  types.Int32 `tfsdk:"size"`
}

type DataSyncCountInfoModel struct {
	ReadCount  types.Int64   `tfsdk:"read_count"`
	ReadOps    types.Float64 `tfsdk:"read_ops"`
	WriteCount types.Int64   `tfsdk:"write_count"`
	WriteOps   types.Float64 `tfsdk:"write_ops"`
}

type SourceProgressInfoModel struct {
	Address       types.String  `tfsdk:"address"`
	SyncPercent   types.Float64 `tfsdk:"sync_percent"`
	AofOffsetDiff types.Int64   `tfsdk:"aof_offset_diff"`
	State         types.String  `tfsdk:"state"`
}

func (c *ctyunRedisMigrationTasks) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10518385`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "任务ID，指定时查询任务详情",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"page_num": schema.Int32Attribute{
				Optional:    true,
				Description: "页码（范围：> 0）",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "数量（范围：[1,100]）",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "查询指定的任务状态，可选值：0：所有状态（默认），1：运行中，2：成功，3：失败",
				Validators: []validator.String{
					stringvalidator.OneOf("0", "1", "2", "3"),
				},
			},
			// 详情相关字段
			"data_sync_count_info": schema.SingleNestedAttribute{
				Description: "命令同步信息汇总",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"read_count": schema.Int64Attribute{
						Computed:    true,
						Description: "从源端读取的总命令数",
					},
					"read_ops": schema.Float64Attribute{
						Computed:    true,
						Description: "从源端读取的命令总OPS",
					},
					"write_count": schema.Int64Attribute{
						Computed:    true,
						Description: "发送给目标的总命令数",
					},
					"write_ops": schema.Float64Attribute{
						Computed:    true,
						Description: "发送给目标的命令总OPS",
					},
				},
			},
			"source_progress_info_list": schema.ListNestedAttribute{
				Description: "同步进度信息",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Computed:    true,
							Description: "源端ip:port",
						},
						"sync_percent": schema.Float64Attribute{
							Computed:    true,
							Description: "全量数据的同步百分比（不包含增量）",
						},
						"aof_offset_diff": schema.Int64Attribute{
							Computed:    true,
							Description: "AOF偏移差距（增量同步时有）",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "同步阶段",
						},
					},
				},
			},
			// 列表相关字段
			"list": schema.ListNestedAttribute{
				Description: "任务列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user_id": schema.StringAttribute{
							Computed:    true,
							Description: "用户ID",
						},
						"tenant_id": schema.StringAttribute{
							Computed:    true,
							Description: "租户ID",
						},
						"task_id": schema.StringAttribute{
							Computed:    true,
							Description: "任务ID",
						},
						"source_spu_inst_id": schema.StringAttribute{
							Computed:    true,
							Description: "源库实例ID",
						},
						"target_spu_inst_id": schema.StringAttribute{
							Computed:    true,
							Description: "目标库实例ID",
						},
						"type": schema.Int32Attribute{
							Computed:    true,
							Description: "类型",
						},
						"status": schema.Int32Attribute{
							Computed:    true,
							Description: "任务状态。0：初始态，1：运行中，2：成功，3：失败",
						},
						"run_step": schema.Int32Attribute{
							Computed:    true,
							Description: "迁移进度。1：全量同步中，2：增量同步中",
						},
						"create_time": schema.Int64Attribute{
							Computed:    true,
							Description: "创建时间戳（秒）",
						},
						"complete_time": schema.Int64Attribute{
							Computed:    true,
							Description: "结束时间戳（秒，=-1时表示时间未知）",
						},
					},
				},
			},
			"total": schema.Int32Attribute{
				Computed:    true,
				Description: "总数",
			},
			"size": schema.Int32Attribute{
				Computed:    true,
				Description: "本次返回数量",
			},
		},
	}
}

func (c *ctyunRedisMigrationTasks) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunRedisMigrationTasks) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunRedisMigrationTasksModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionId.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.RegionId = types.StringValue(regionId)

	// 根据输入参数决定调用哪个API
	if !config.ID.IsNull() && !config.ID.IsUnknown() {
		// 查询任务详情
		err = c.getTaskDetail(ctx, &config)
		if err != nil {
			return
		}
	} else {
		// 查询任务列表
		err = c.getTaskList(ctx, &config)
		if err != nil {
			return
		}
	}

	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

// getTaskDetail 查询任务详情
func (c *ctyunRedisMigrationTasks) getTaskDetail(ctx context.Context, data *CtyunRedisMigrationTasksModel) (err error) {
	params := &ctgdcs2.Dcs2GetTaskProgressDetailInfoRequest{
		RegionId: data.RegionId.ValueString(),
		TaskId:   data.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2GetTaskProgressDetailInfoApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 处理命令同步信息汇总
	if resp.ReturnObj.DataSyncCountInfo != nil {
		data.DataSyncCountInfo = &DataSyncCountInfoModel{
			ReadCount:  types.Int64Value(resp.ReturnObj.DataSyncCountInfo.ReadCount),
			ReadOps:    types.Float64Value(resp.ReturnObj.DataSyncCountInfo.ReadOps),
			WriteCount: types.Int64Value(resp.ReturnObj.DataSyncCountInfo.WriteCount),
			WriteOps:   types.Float64Value(resp.ReturnObj.DataSyncCountInfo.WriteOps),
		}
	}

	// 处理同步进度信息
	if len(resp.ReturnObj.SourceProgressInfoList) > 0 {
		progressInfoModels := make([]SourceProgressInfoModel, 0, len(resp.ReturnObj.SourceProgressInfoList))
		for _, info := range resp.ReturnObj.SourceProgressInfoList {
			progressInfoModels = append(progressInfoModels, SourceProgressInfoModel{
				Address:       types.StringValue(info.Address),
				SyncPercent:   types.Float64Value(info.SyncPercent),
				AofOffsetDiff: types.Int64Value(info.AofOffsetDiff),
				State:         types.StringValue(info.State),
			})
		}

		progressInfoObjType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"address":         types.StringType,
				"sync_percent":    types.Float64Type,
				"aof_offset_diff": types.Int64Type,
				"state":           types.StringType,
			},
		}

		progressInfoList, diags := types.ListValueFrom(ctx, progressInfoObjType, progressInfoModels)
		if diags.HasError() {
			err = fmt.Errorf("failed to set source_progress_info_list: %v", diags)
			return
		}
		data.SourceProgressInfoList = progressInfoList
	} else {
		data.SourceProgressInfoList = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"address":         types.StringType,
				"sync_percent":    types.Float64Type,
				"aof_offset_diff": types.Int64Type,
				"state":           types.StringType,
			},
		})
	}

	return
}

// getTaskList 查询任务列表
func (c *ctyunRedisMigrationTasks) getTaskList(ctx context.Context, data *CtyunRedisMigrationTasksModel) (err error) {
	params := &ctgdcs2.Dcs2ListTaskInfoRequest{
		RegionId: data.RegionId.ValueString(),
	}

	// 设置分页参数
	if !data.PageNum.IsNull() && !data.PageNum.IsUnknown() {
		params.PageNum = fmt.Sprintf("%d", data.PageNum.ValueInt32())
	} else {
		params.PageNum = "1" // 默认第一页
	}

	if !data.PageSize.IsNull() && !data.PageSize.IsUnknown() {
		params.PageSize = fmt.Sprintf("%d", data.PageSize.ValueInt32())
	} else {
		params.PageSize = "10" // 默认每页10条
	}

	// 设置状态参数
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status := data.Status.ValueString()
		params.Status = &status
	} else {
		status := "0" // 默认查询所有状态
		params.Status = &status
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2ListTaskInfoApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 设置分页信息
	data.Total = types.Int32Value(resp.ReturnObj.Total)
	data.Size = types.Int32Value(resp.ReturnObj.Size)

	// 处理列表数据
	if len(resp.ReturnObj.List) > 0 {
		// 定义列表项的结构
		taskListModels := make([]types.Object, 0, len(resp.ReturnObj.List))

		// 定义列表项的属性类型
		listObjType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"user_id":            types.StringType,
				"tenant_id":          types.StringType,
				"task_id":            types.StringType,
				"source_spu_inst_id": types.StringType,
				"target_spu_inst_id": types.StringType,
				"type":               types.Int32Type,
				"status":             types.Int32Type,
				"run_step":           types.Int32Type,
				"create_time":        types.Int64Type,
				"complete_time":      types.Int64Type,
			},
		}

		for _, task := range resp.ReturnObj.List {
			// 构建每个任务项的属性值
			attrValues := map[string]attr.Value{
				"user_id":            stringValueOrDefault(task.UserId, ""),
				"tenant_id":          stringValueOrDefault(task.TenantId, ""),
				"task_id":            stringValueOrDefault(task.TaskId, ""),
				"source_spu_inst_id": stringValueOrDefault(task.SourceSpuInstId, ""),
				"target_spu_inst_id": stringValueOrDefault(task.TargetSpuInstId, ""),
				"type":               types.Int32Value(task.RawType),
				"status":             types.Int32Value(task.Status),
				"run_step":           types.Int32Value(task.RunStep),
				"create_time":        types.Int64Value(task.CreateTime),
				"complete_time":      types.Int64Value(task.CompleteTime),
			}

			// 创建对象
			taskObj, diags := types.ObjectValue(listObjType.AttrTypes, attrValues)
			if diags.HasError() {
				err = fmt.Errorf("failed to create task list object: %v", diags)
				return
			}

			taskListModels = append(taskListModels, taskObj)
		}

		// 转换taskListModels为[]attr.Value
		attrValues := make([]attr.Value, len(taskListModels))
		for i, obj := range taskListModels {
			attrValues[i] = obj
		}

		// 创建列表
		taskList, diags := types.ListValue(listObjType, attrValues)
		if diags.HasError() {
			err = fmt.Errorf("failed to set task list: %v", diags)
			return
		}
		data.List = taskList
	} else {
		// 如果没有数据，创建空列表
		listObjType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"user_id":            types.StringType,
				"tenant_id":          types.StringType,
				"task_id":            types.StringType,
				"source_spu_inst_id": types.StringType,
				"target_spu_inst_id": types.StringType,
				"type":               types.Int32Type,
				"status":             types.Int32Type,
				"run_step":           types.Int32Type,
				"create_time":        types.Int64Type,
				"complete_time":      types.Int64Type,
			},
		}
		data.List = types.ListNull(types.ObjectType{AttrTypes: listObjType.AttrTypes})
	}

	return
}

// stringValueOrDefault 返回字符串值或默认值
func stringValueOrDefault(s *string, defaultValue string) types.String {
	if s != nil {
		return types.StringValue(*s)
	}
	return types.StringValue(defaultValue)
}
