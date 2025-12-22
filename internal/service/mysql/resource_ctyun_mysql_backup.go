package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunMysqlBackup{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlBackup{}
	_ resource.ResourceWithImportState = &CtyunMysqlBackup{}
)

type CtyunMysqlBackup struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlBackup) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_backup"
}
func NewCtyunMysqlBackup() resource.Resource {
	return &CtyunMysqlBackup{}
}

func (c *CtyunMysqlBackup) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlBackup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[instanceID],[projectID],[regionID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunMysqlBackupConfig
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
	cfg.InstID = types.StringValue(instId)
	cfg.ID = types.StringValue(fmt.Sprintf("%s,%s", name, instId))
	err = c.getAndMergeMysqlBackup(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *CtyunMysqlBackup) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10098797",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "mysql实例id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(32, 32),
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
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "备份名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "备份集备注",
				Validators: []validator.String{
					validator2.Desc(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"task_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "备份类型,默认全量物理备份 全量物理备份:full 全量逻辑备份:logic_full。逻辑备份支持资源池：（华北2、西安7），具体可查看文档：https://www.ctyun.cn/document/10033813/10902204",
				Default:     stringdefault.StaticString("full"),
				Validators: []validator.String{
					stringvalidator.OneOf("full", "logic_full"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunMysqlBackup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlBackupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开始创建备份
	plan.Name = types.StringValue(plan.InstID.ValueString() + strings.ReplaceAll(uuid.NewString(), "-", ""))
	err = c.createMysqlBackup(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.Name.ValueString(), plan.InstID.ValueString()))
	// 创建后，获取mysql详情
	err = c.getAndMergeMysqlBackup(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlBackup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMysqlBackupConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMysqlBackup(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlBackup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunMysqlBackup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunMysqlBackupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.deleteBackupSetAndFile(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunMysqlBackup) createMysqlBackup(ctx context.Context, config CtyunMysqlBackupConfig) error {
	params := &mysql.TeledbCreateBackupRequest{
		OuterProdInstId: config.InstID.ValueString(),
		BackupName:      config.Name.ValueString(),
		TaskType:        config.TaskType.ValueString(),
	}
	if !config.Description.IsNull() {
		params.Description = config.Description.ValueStringPointer()
	}
	header := &mysql.TeledbCreateBackupRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	//  mysql备份之前，确定状态running
	err := c.startedLoop(ctx, config, 60)
	if err != nil {
		return err
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCreateBackupApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建mysql实例id=%s备份失败，接口返回nil。请与研发联系确认问题原因。", config.InstID.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("create mysql backup error, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	return c.checkAfterCreate(ctx, config)
}

func (c *CtyunMysqlBackup) getAndMergeMysqlBackup(ctx context.Context, config *CtyunMysqlBackupConfig) error {
	resp, err := c.getBackupRecordList(ctx, *config)
	if err != nil {
		return err
	}
	backupList := resp.List
	if len(backupList) == 0 {
		err = fmt.Errorf("通过backupName=%s,mysql实例id=%s未查询到备份集", config.Name.ValueString(), config.InstID.ValueString())
		return err
	}
	if len(backupList) > 1 {
		err = fmt.Errorf("通过backupName=%s,mysql实例id=%s查询到多条备份集", config.Name.ValueString(), config.InstID.ValueString())
		return err
	}
	return nil
}

func (c *CtyunMysqlBackup) deleteBackupSetAndFile(ctx context.Context, config CtyunMysqlBackupConfig) error {
	// 1. 根据accountName 确认备集下的 records是否存在正在备份的任务
	// 2. 根据accountName获取blockId
	// 3. 通过blockId删除

	// 轮询确认该备份集下是否有正在备份的实例
	err := c.backupIngLoop(ctx, config, 60)
	if err != nil {
		return err
	}

	resp, err := c.getBackupRecordList(ctx, config)
	if err != nil {
		return err
	}
	backupList := resp.List
	if len(backupList) == 0 {
		err = fmt.Errorf("通过backupName(%s),mysql实例(id=%s)未查询到备份集", config.Name.ValueString(), config.InstID.ValueString())
		return err
	}
	if len(backupList) > 1 {
		err = fmt.Errorf("通过backupName(%s),mysql实例(id=%s)查询到多条备份集", config.Name.ValueString(), config.InstID.ValueString())
		return err
	}

	//if !flag {
	//	err = fmt.Errorf("mysql实例(id=%s)备份状态不符合预期，只有备份成功备份集才支持删除！", config.InstID.ValueString())
	//	return err
	//}
	params := &mysql.TeledbDeleteBackupRequest{
		OuterProdInstId: config.InstID.ValueString(),
		BlockID:         backupList[0].BlockId,
	}
	header := &mysql.TeledbDeleteBackupRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	deleteResp, err := c.meta.Apis.SdkCtMysqlApis.TeledbDeleteBackupApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if deleteResp == nil {
		err = fmt.Errorf("mysql实例(id=%s)删除备份集(backup_name=%s)失败，接口返回nil。请联系研发确认问题原因。", config.InstID.ValueString(), config.Name.ValueString())
		return err
	} else if deleteResp.StatusCode != 0 {
		err = fmt.Errorf("delete backup set error, API return error. Message: %s Error: %s", deleteResp.Message, *deleteResp.Error)
		return err
	}
	return nil
}

func (c *CtyunMysqlBackup) getBackupRecordList(ctx context.Context, config CtyunMysqlBackupConfig) (*mysql.TeledbGetBackupListResponseReturnObj, error) {
	params := &mysql.TeledbGetBackupListRequest{
		OuterProdInstId: config.InstID.ValueString(),
		BackupName:      config.Name.ValueStringPointer(),
		PageNow:         1,
		PageSize:        10,
	}

	header := &mysql.TeledbGetBackupListRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetBackupListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取mysql实例(id=%s)备份集(id=%s)信息及备份任务失败，接口返回nil，具体原因请联系研发确认！", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("get backup list error, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj, nil
}

func (c *CtyunMysqlBackup) getBackupRecordDetail(ctx context.Context, config *CtyunBackupCancelConfig) (*mysql.TeledbGetBackupRecordDetailResponse, error) {
	params := &mysql.TeledbGetBackupRecordDetailRequest{
		Id: fmt.Sprintf("%d", config.BackupRecordId.ValueInt64()),
	}
	header := &mysql.TeledbGetBackupRecordDetailRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetBackupRecordDetailApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询备份任务(record_id=%d)信息失败，接口返回nil。具体原因请联系研发确定。", config.BackupRecordId.ValueInt64())
		return nil, err
	} else if resp.StatusCode == 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	}
	if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}

	return resp, nil
}

// 确保该备份集下没有正在备份的任务
func (c *CtyunMysqlBackup) backupIngLoop(ctx context.Context, config CtyunMysqlBackupConfig, loopCount ...int) error {
	var err error
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.getBackupRecordList(ctx, config)
			if err2 != nil {
				err = err2
				return false
			}
			backupList := resp.List
			if len(backupList) == 0 {
				err = fmt.Errorf("通过backupName(%s),mysql实例(id=%s)未查询到备份集", config.Name.ValueString(), config.InstID.ValueString())
				return false
			}
			if len(backupList) > 1 {
				err = fmt.Errorf("通过backupName(%s),mysql实例(id=%s)查询到多条备份集", config.Name.ValueString(), config.InstID.ValueString())
				return false
			}

			backupRecordList := backupList[0]
			for _, record := range backupRecordList.Records {
				taskStatus := record.TaskStatus
				switch taskStatus {
				case business.MysqlBackupTaskStatusSuccess:
					continue
				case business.MysqlBackupTaskStatusWaitStart:
					return true
				case business.MysqlBackupTaskStatusSubmit:
					return true
				case business.MysqlBackupTaskStatusFailed:
					continue
				case business.MysqlBackupTaskStatusCancel:
					continue
				default:
					return true
				}
			}
			return false
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，备份任务仍未备份成功！")
	}
	return err
}

func (c *CtyunMysqlBackup) startedLoop(ctx context.Context, config CtyunMysqlBackupConfig, loopCount ...int) (err error) {
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
				OuterProdInstId: config.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   config.InstID.ValueString(),
				RegionID: config.RegionID.ValueString(),
			}
			if config.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = config.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
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
				err = fmt.Errorf("mysql实例已冻结，请先启用后再进行备份操作！")
				if err != nil {
					return false
				}
			}
			if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlRunningStatusStarted {
				// 有三次是start，才认为状态正常
				cnt++
				if cnt > 2 {
					return false
				}
			}
			if orderStatus == business.MysqlOrderStatusPause {
				err = errors.New("订单处于暂停状态，不可进行变更操作")
				return false
			}
			if runningStatus == business.MysqlRunningStatusStopping || runningStatus == business.MysqlRunningStatusStopped {
				err = errors.New("主机处于关机状态，不可进行变更操作")
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

type CtyunMysqlBackupConfig struct {
	InstID      types.String `tfsdk:"instance_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	RegionID    types.String `tfsdk:"region_id"`
	Name        types.String `tfsdk:"name"` // 备份名称在4位到64位之间，不区分大小写，可以包含中文、字母、数字、中划线或下划线，不能包含其他特殊字符
	Description types.String `tfsdk:"description"`
	TaskType    types.String `tfsdk:"task_type"`
	ID          types.String `tfsdk:"id"`
}

func (c *CtyunMysqlBackup) checkAfterCreate(ctx context.Context, config CtyunMysqlBackupConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var resp *mysql.TeledbGetBackupListResponseReturnObj
			resp, err = c.getBackupRecordList(ctx, config)
			if err != nil {
				return false
			}
			backupList := resp.List
			if len(backupList) == 0 {
				err = fmt.Errorf("通过backupName(%s),mysql实例(id=%s)未查询到备份集", config.Name.ValueString(), config.InstID.ValueString())
				return false
			}
			if len(backupList) > 1 {
				err = fmt.Errorf("通过backupName(%s),mysql实例(id=%s)查询到多条备份集", config.Name.ValueString(), config.InstID.ValueString())
				return false
			}

			backupRecordList := backupList[0]
			if len(backupRecordList.Records) > 0 && backupRecordList.Records[0].TaskId != "" && backupRecordList.Records[0].BackupRecordId != 0 {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		return fmt.Errorf("backupName(%s),mysql实例(id=%s)未完成", config.Name.ValueString(), config.InstID.ValueString())
	}
	return nil
}
