package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
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
	"regexp"
	"strconv"
	"time"
)

var (
	_ resource.Resource              = &CtyunMysqlBackupRecovery{}
	_ resource.ResourceWithConfigure = &CtyunMysqlBackupRecovery{}
)

type CtyunMysqlBackupRecovery struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlBackupRecovery) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_backup_recovery"
}
func NewCtyunMysqlBackupRecovery() resource.Resource {
	return &CtyunMysqlBackupRecovery{}
}

func (c *CtyunMysqlBackupRecovery) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlBackupRecovery) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10098797",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "mysql实例id",
				Validators: []validator.String{
					stringvalidator.LengthBetween(32, 32),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"src_instance_id": schema.StringAttribute{
				Required:    true,
				Description: "恢复的源mysql实例id",
				Validators: []validator.String{
					stringvalidator.LengthBetween(32, 32),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dst_instance_id": schema.StringAttribute{
				Required:    true,
				Description: "恢复的目标mysql实例id",
				Validators: []validator.String{
					stringvalidator.LengthBetween(32, 32),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"to_timepoint": schema.StringAttribute{
				Optional:    true,
				Description: "恢复到的时间点，格式为：YYYY-MM-DDTHH:MM:SSZ【taskId和to_timepoint不能同时为空，优先取to_timepoint】",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"必须符合格式：YYYY-MM-DDTHH:MM:SSZ",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"task_id": schema.StringAttribute{
				Optional:    true,
				Description: "用来恢复的备份任务集【task_id和to_timepoint不能同时为空，优先取to_timepoint】",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "备份任务id",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunMysqlBackupRecovery) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlBackupRecoveryConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	//
	err = c.CreateMysqlBackupRecovery(ctx, &plan)
	if err != nil {
		return
	}
	err = c.checkBackup(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后，获取mysql详情
	//err = c.getAndMergeMysqlBackupRecovery(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlBackupRecovery) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	return
}

func (c *CtyunMysqlBackupRecovery) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunMysqlBackupRecovery) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunMysqlBackupRecovery) CreateMysqlBackupRecovery(ctx context.Context, config *CtyunMysqlBackupRecoveryConfig) error {
	params := &mysql.TeledbCreateRecoveryJobRequest{
		SrcOuterProdInstId: config.SrcInstId.ValueString(),
		DstOuterProdInstId: config.DstInstId.ValueString(),
	}
	if config.TaskId.IsNull() && config.ToTimepoint.IsNull() {
		err := fmt.Errorf("task_id和to_timepoint不能同时为空")
		return err
	}
	if config.TaskId.IsNull() && !config.ToTimepoint.IsNull() {
		a := utils.FromRFC3339ToLocal(config.ToTimepoint.ValueString())
		params.ToTimepoint = &a
	}
	if !config.TaskId.IsNull() && config.ToTimepoint.IsNull() {
		params.TaskId = config.TaskId.ValueStringPointer()
	}
	if !config.TaskId.IsNull() && !config.ToTimepoint.IsNull() {
		a := utils.FromRFC3339ToLocal(config.ToTimepoint.ValueString())
		params.ToTimepoint = &a
	}
	header := &mysql.TeledbCreateRecoveryJobRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	// 备份恢复操作前，先轮询确认实例是running的
	err := c.StartedLoop(ctx, config, 30)
	if err != nil {
		return err
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCreateRecoveryJobApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建从备份恢复实例任务失败，源mysql实例id=%s，目的mysql实例id=%s", config.SrcInstId.ValueString(), config.DstInstId.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	if resp.ReturnObj.Data == "" {
		err = common.InvalidReturnObjError
		return err
	}
	id, err := strconv.ParseInt(resp.ReturnObj.Data, 10, 64)
	if err != nil {
		return err
	}
	config.ID = types.Int64Value(id)
	return nil
}

//func (c *CtyunMysqlBackupRecovery) getAndMergeMysqlBackupRecovery(ctx context.Context, config *CtyunMysqlBackupRecoveryConfig) error {
//	return nil
//}

func (c *CtyunMysqlBackupRecovery) BackupRecoveryLoop(ctx context.Context, config *CtyunMysqlBackupRecoveryConfig, loopCount ...int) error {
	var err error
	count := 30
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			recoveryList, err2 := c.getBackupRecoveryList(ctx, config)
			if err2 != nil {
				err = err2
				return false
			}
			status := recoveryList[0].TaskStatus
			switch status {
			case business.MysqlBackupRecoveryJobStatusING:
				return true
			case business.MysqlBackupRecoveryJobStatusTODO:
				return true
			case business.MysqlBackupRecoveryJobStatusPreCheck:
				return true
			case business.MysqlBackupRecoveryJobStatusSuccess:
				return false
			case business.MysqlBackupRecoveryJobStatusFail:
				return false
			}

			return false
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，数据库备份恢复仍未完成！")
	}
	return err
}

func (c *CtyunMysqlBackupRecovery) getBackupRecoveryList(ctx context.Context, config *CtyunMysqlBackupRecoveryConfig) ([]mysql.BrRecoveryRecordVo, error) {
	params := &mysql.TeledbGetBackupRecoveryListRequest{
		OuterProdInstId: config.InstID.ValueString(),
		ID:              config.ID.ValueInt64Pointer(),
		PageSize:        10,
		PageNow:         1,
	}
	header := &mysql.TeledbGetBackupRecoveryListRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetBackupRecoveryListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询mysql实例(id=%s)的备份恢复任务(id=%d)列表失败，接口返回nil，请联系研发确认问题原因。", config.InstID.ValueString(), config.ID.ValueInt64())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	recoveryList := resp.ReturnObj.List
	if len(recoveryList) == 0 {
		err = common.InvalidReturnObjError
		return nil, err
	}
	if len(recoveryList) > 1 {
		err = fmt.Errorf("查询mysql实例(id=%s)的备份恢复任务(id=%s)列表长度不为1。", config.InstID.ValueString(), config.ID)
		return nil, err
	}
	return recoveryList, nil
}

func (c *CtyunMysqlBackupRecovery) StartedLoop(ctx context.Context, state *CtyunMysqlBackupRecoveryConfig, loopCount ...int) (err error) {
	count := 30
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	var cnt int
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: state.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   state.InstID.ValueString(),
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp == nil {
				err = fmt.Errorf("查询mysql实例(id=%s)的详情失败，接口返回nil。请与研发联系确认问题原因。", state.InstID.ValueString())
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			runningStatus := resp.ReturnObj.ProdRunningStatus
			orderStatus := resp.ReturnObj.ProdOrderStatus
			// 若变配前，发现数据库已冻结，将其恢复
			if orderStatus == business.MysqlOrderStatusPause {
				err = fmt.Errorf("mysql实例(id=%s)当前处理冻结状态", state.InstID.ValueString())
				if err != nil {
					return false
				}
			}
			if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlRunningStatusStarted {
				// 有三次是start，才认为状态正常
				cnt++
				if cnt > 3 {
					return false
				}
			}
			if orderStatus == business.MysqlOrderStatusPause {
				err = errors.New("订单处于暂停状态，不可进行备份恢复操作")
				return false
			}
			if runningStatus == business.MysqlRunningStatusStopping || runningStatus == business.MysqlRunningStatusStopped {
				err = errors.New("主机处于关机状态，不可进行备份恢复操作")
				return false
			}

			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未到达运行状态！")
	}
	return
}

func (c *CtyunMysqlBackupRecovery) checkBackup(ctx context.Context, config *CtyunMysqlBackupRecoveryConfig) error {
	// 轮询查看备份恢复任务
	err := c.BackupRecoveryLoop(ctx, config, 60)
	if err != nil {
		return err
	}
	return err
}

type CtyunMysqlBackupRecoveryConfig struct {
	InstID      types.String `tfsdk:"instance_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	RegionID    types.String `tfsdk:"region_id"`
	SrcInstId   types.String `tfsdk:"src_instance_id"`
	DstInstId   types.String `tfsdk:"dst_instance_id"`
	ToTimepoint types.String `tfsdk:"to_timepoint"`
	TaskId      types.String `tfsdk:"task_id"`
	ID          types.Int64  `tfsdk:"id"`
}
