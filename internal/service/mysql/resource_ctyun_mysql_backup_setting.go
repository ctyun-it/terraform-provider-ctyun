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
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunMysqlBackupSetting{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlBackupSetting{}
	_ resource.ResourceWithImportState = &CtyunMysqlBackupSetting{}
)

type CtyunMysqlBackupSetting struct {
	meta         *common.CtyunMetadata
	mysqlService *business.MysqlService
}

func (c *CtyunMysqlBackupSetting) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_backup_setting"
}
func NewCtyunMysqlBackupSetting() resource.Resource {
	return &CtyunMysqlBackupSetting{}
}

func (c *CtyunMysqlBackupSetting) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.mysqlService = business.NewMysqlService(c.meta)
}

func (c *CtyunMysqlBackupSetting) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [instanceID],[projectID],[regionID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunMysqlBackupSettingConfig
	var regionId, projectId, instId string

	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		instId = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &instId, &projectId, &regionId)
		if err != nil {
			return
		}
	}
	if instId == "" {
		err = fmt.Errorf("instanceID 不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID 不能为空")
	}

	cfg.RegionID = types.StringValue(regionId)
	cfg.ProjectID = types.StringValue(projectId)
	cfg.InstID = types.StringValue(instId)
	cfg.ID = types.StringValue(fmt.Sprintf("%s-setting", instId))
	err = c.getAndMergeMysqlBackupSetting(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *CtyunMysqlBackupSetting) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
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
			"storage_day": schema.Int32Attribute{
				Required:    true,
				Description: "备份数据保留天数，最少保留1天。支持更新。",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"frequency_backup": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "高频备份 true=开启，false=关闭。默认关闭，支持更新。",
			},
			"frequency_backup_unit_time": schema.Int64Attribute{
				Optional:    true,
				Description: "高频备份频率 单位: 秒，最小为1小时，即3600。支持更新",
				Validators: []validator.Int64{
					int64validator.AtLeast(3600),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("frequency_backup"),
						types.BoolValue(false),
					),
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("frequency_backup"),
						types.BoolValue(true),
					),
				},
			},
			"allow_earliest_time": schema.StringAttribute{
				Required:    true,
				Description: "允许最早开始备份时间 默认：00:00",
				Validators: []validator.String{
					validator2.BackupTimeValidator(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"trigger_days_of_week": schema.SetAttribute{
				Required:    true,
				ElementType: types.Int32Type,
				Description: "全备触发星期，单个元素取值范围为1~7，1代表星期一，7代表星期日，以此类推。支持更新",
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 7),
					setvalidator.ValueInt32sAre(int32validator.Between(1, 7)),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
			},
		},
	}
}

func (c *CtyunMysqlBackupSetting) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlBackupSettingConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 修改备份设置
	err = c.updateMysqlBackupSettingConfig(ctx, &plan)
	if err != nil {
		return
	}

	// 获取更新后的备份设置
	err = c.getAndMergeMysqlBackupSetting(ctx, &plan)
	if err != nil {
		return
	}
	// id赋值， instance id + "-" + settings
	plan.ID = types.StringValue(fmt.Sprintf("%s-setting", plan.InstID.ValueString()))
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlBackupSetting) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMysqlBackupSettingConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMysqlBackupSetting(ctx, &state)
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

func (c *CtyunMysqlBackupSetting) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunMysqlBackupSettingConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMysqlBackupSettingConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.updateMysqlBackupSettingConfig(ctx, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeMysqlBackupSetting(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlBackupSetting) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunMysqlBackupSetting) updateMysqlBackupSettingConfig(ctx context.Context, config *CtyunMysqlBackupSettingConfig) error {

	err := c.mysqlService.WaitInstanceStatus(ctx, config.InstID.ValueString(), config.ProjectID.ValueString(), config.RegionID.ValueString(),
		business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted)
	if err != nil {
		return err
	}
	var triggerWeek []int32
	diags := config.TriggerDaysOfWeek.ElementsAs(ctx, &triggerWeek, false)
	if diags.HasError() {
		err := errors.New(diags[0].Detail())
		return err
	}
	week, err := c.processWeek(ctx, triggerWeek, business.MysqlBackupSettingConfigWeek)
	if err != nil {
		return err
	}
	params := &mysql.TeledbUpdateBackupSettingRequest{
		ExpiredTime:       c.transSec(config.StorageDay.ValueInt32()),
		FrequencyBackup:   config.FrequencyBackup.ValueBool(),
		AllowEarliestTime: config.AllowEarliestTime.ValueString() + ":00",
		OuterProdInstId:   config.InstID.ValueString(),
		TriggerDaysOfWeek: week,
	}
	if !config.FrequencyBackupUnitTime.IsNull() {
		params.FrequencyBackupUnitTime = config.FrequencyBackupUnitTime.ValueInt64Pointer()
	}
	header := &mysql.TeledbUpdateBackupSettingRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateBackupSettingApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新mysql实例(id=%s)备份设置有误，接口返回nil。请联系研发确认问题原因", config.InstID.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

func (c *CtyunMysqlBackupSetting) transSec(storageDay int32) int64 {
	return int64(storageDay * 24 * 3600)
}

func (c *CtyunMysqlBackupSetting) transDay(sec int64) int32 {
	return int32(sec / 24 / 3600)
}

func (c *CtyunMysqlBackupSetting) processWeek(ctx context.Context, weeks []int32, weekCorr map[int32]int32) ([]int32, error) {
	var outputWeeks []int32

	for _, day := range weeks {
		newDay := weekCorr[day]
		outputWeeks = append(outputWeeks, newDay)
	}
	return outputWeeks, nil
}

func (c *CtyunMysqlBackupSetting) getAndMergeMysqlBackupSetting(ctx context.Context, config *CtyunMysqlBackupSettingConfig) error {
	resp, err := c.getMysqlBackupSettingInfo(ctx, config)
	if err != nil {
		return err
	}
	// 解析备份设置
	returnObj := resp.ReturnObj
	config.StorageDay = types.Int32Value(c.transDay(returnObj.ExpiredTime))
	config.FrequencyBackup = types.BoolValue(returnObj.FrequencyBackup)
	if config.FrequencyBackup.ValueBool() {
		config.FrequencyBackupUnitTime = types.Int64Value(returnObj.FrequencyBackupUnitTime)
	}

	config.AllowEarliestTime = types.StringValue(extractHourMinute(returnObj.AllowEarliestTime))
	week, err := c.processWeek(ctx, returnObj.TriggerDaysOfWeek, business.MysqlBackupSettingConfigWeekRev)
	if err != nil {
		return err
	}
	weekTF, diags := types.SetValueFrom(ctx, types.Int32Type, week)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return err
	}
	config.TriggerDaysOfWeek = weekTF

	return nil
}

func (c *CtyunMysqlBackupSetting) getMysqlBackupSettingInfo(ctx context.Context, config *CtyunMysqlBackupSettingConfig) (*mysql.TeledbGetBackupSettingDetailResponse, error) {
	params := &mysql.TeledbGetBackupSettingDetailRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbGetBackupSettingDetailRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetBackupSettingDetailApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取mysql(id=%s)备份设置信息失败，接口返回nil。请联系研发确认问题原因！", config.InstID.ValueString())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func extractHourMinute(timeStr string) string {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 3 {
		return timeStr // 或返回错误
	}
	return parts[0] + ":" + parts[1]
}

type CtyunMysqlBackupSettingConfig struct {
	InstID                  types.String `tfsdk:"instance_id"`
	ProjectID               types.String `tfsdk:"project_id"`
	RegionID                types.String `tfsdk:"region_id"`
	StorageDay              types.Int32  `tfsdk:"storage_day"`
	FrequencyBackup         types.Bool   `tfsdk:"frequency_backup"`
	FrequencyBackupUnitTime types.Int64  `tfsdk:"frequency_backup_unit_time"` // 高频备份频率 单位: 秒，最小为1小时，即3600
	AllowEarliestTime       types.String `tfsdk:"allow_earliest_time"`
	TriggerDaysOfWeek       types.Set    `tfsdk:"trigger_days_of_week"`
	ID                      types.String `tfsdk:"id"`
}
