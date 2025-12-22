package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunMysqlBackups{}
	_ datasource.DataSourceWithConfigure = &ctyunMysqlBackups{}
)

type ctyunMysqlBackups struct {
	meta *common.CtyunMetadata
}

func NewCtyunMysqlBackups() datasource.DataSource {
	return &ctyunMysqlBackups{}
}
func (c *ctyunMysqlBackups) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunMysqlBackups) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_backups"
}

func (c *ctyunMysqlBackups) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10098797",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "MySQL实例ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目ID",
			},
			"instance_name": schema.StringAttribute{
				Optional:    true,
				Description: "MySQL实例名称",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "备份名称",
			},
			"id": schema.Int64Attribute{
				Optional:    true,
				Description: "备份集ID",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数，默认10",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，默认1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"backups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "实例ID",
						},
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "备份ID",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "实例名称",
						},
						"records": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"backup_record_id": schema.Int64Attribute{
										Computed:    true,
										Description: "备份记录ID",
									},
									"backup_task_id": schema.Int64Attribute{
										Computed:    true,
										Description: "备份任务ID",
									},
									"task_id": schema.StringAttribute{
										Computed:    true,
										Description: "任务ID",
									},
									"backup_name": schema.StringAttribute{
										Computed:    true,
										Description: "备份名称",
									},
									"outer_prod_inst_id": schema.StringAttribute{
										Computed:    true,
										Description: "外部产品实例ID",
									},
									"prod_inst_name": schema.StringAttribute{
										Computed:    true,
										Description: "产品实例名称",
									},
									"description": schema.StringAttribute{
										Computed:    true,
										Description: "备份描述",
									},
									"storage_type": schema.StringAttribute{
										Computed:    true,
										Description: "存储类型",
									},
									"op_type": schema.StringAttribute{
										Computed:    true,
										Description: "操作类型（auto/manual）",
									},
									"task_type": schema.StringAttribute{
										Computed:    true,
										Description: "任务类型（full/incr）",
									},
									"task_status": schema.Int32Attribute{
										Computed:    true,
										Description: "任务状态",
									},
									"backed_up_data_size": schema.Int64Attribute{
										Computed:    true,
										Description: "备份数据大小（字节）",
									},
									"backed_up_data_size_human": schema.StringAttribute{
										Computed:    true,
										Description: "备份数据大小（格式化）",
									},
									"backup_start_time": schema.StringAttribute{
										Computed:    true,
										Description: "备份开始时间",
									},
									"backup_end_time": schema.StringAttribute{
										Computed:    true,
										Description: "备份结束时间",
									},
									"disabled": schema.BoolAttribute{
										Computed:    true,
										Description: "是否禁用",
									},
								},
							},
							Description: "备份记录列表",
						},
					},
				},
				Description: "备份列表",
			},
		},
	}
}

func (c *ctyunMysqlBackups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMysqlBackupsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	params := &mysql.TeledbGetBackupListRequest{
		OuterProdInstId: config.InstID.ValueString(),
		PageNow:         1,
		PageSize:        10,
	}
	header := &mysql.TeledbGetBackupListRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: regionId,
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	if !config.InstName.IsNull() {
		params.ProdInstName = config.InstName.ValueStringPointer()
	}
	if !config.Name.IsNull() {
		params.BackupName = config.Name.ValueStringPointer()
	}
	if !config.BackupID.IsNull() {
		params.BlockId = config.BackupID.ValueInt64Pointer()
	}
	if !config.PageNo.IsNull() {
		params.PageNow = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetBackupListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询mysql实例id=%s备份集列表失败，接口返回nil，请与研发联系确认问题原因。", config.InstID.ValueString())
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("get mysql backup list failed, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var backupList []BackupModel
	for _, backupItem := range resp.ReturnObj.List {
		var backupInfo BackupModel
		backupInfo.InstId = types.StringValue(backupItem.OuterProdInstId)
		backupInfo.BackupID = types.Int64Value(backupItem.BlockId)
		backupInfo.InstName = types.StringValue(backupItem.ProdInstName)
		var records []BackupRecordModel
		for _, recordItem := range backupItem.Records {
			var recordInfo BackupRecordModel
			recordInfo.BackupRecordId = types.Int64Value(recordItem.BackupRecordId)
			recordInfo.TaskId = types.StringValue(recordItem.TaskId)
			recordInfo.BackupTaskId = types.Int64Value(recordItem.BackupTaskId)
			recordInfo.BackupName = types.StringValue(recordItem.BackupName)
			recordInfo.OuterProdInstId = types.StringValue(recordItem.OuterProdInstId)
			recordInfo.ProdInstName = types.StringValue(recordItem.ProdInstName)
			recordInfo.Description = types.StringValue(recordItem.Description)
			recordInfo.OpType = types.StringValue(recordItem.OpType)
			recordInfo.TaskType = types.StringValue(recordItem.TaskType)
			recordInfo.TaskStatus = types.Int32Value(recordItem.TaskStatus)
			recordInfo.BackedUpDataSize = types.Int64Value(recordItem.BackedUpDataSize)
			recordInfo.BackedUpDataSizeHuman = types.StringValue(recordItem.BackedUpDataSizeHuman)
			recordInfo.BackupStartTime = types.StringValue(utils.FromBJTimeToUTCZ(recordItem.BackupStartTime))
			recordInfo.BackupEndTime = types.StringValue(utils.FromBJTimeToUTCZ(recordItem.BackupEndTime))
			recordInfo.Disabled = types.BoolValue(recordItem.Disabled)
			records = append(records, recordInfo)
		}
		backupInfo.Records = records
		backupList = append(backupList, backupInfo)
	}

	config.BackupList = backupList
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

}

type BackupRecordModel struct {
	BackupRecordId        types.Int64  `tfsdk:"backup_record_id"`          // 备份记录id
	BackupTaskId          types.Int64  `tfsdk:"backup_task_id"`            // 备份任务id
	TaskId                types.String `tfsdk:"task_id"`                   // 任务id
	BackupName            types.String `tfsdk:"backup_name"`               // 备份名称
	OuterProdInstId       types.String `tfsdk:"outer_prod_inst_id"`        // 外部实例id
	ProdInstName          types.String `tfsdk:"prod_inst_name"`            // 实例名称
	Description           types.String `tfsdk:"description"`               // 备份描述
	StorageType           types.String `tfsdk:"storage_type"`              // 存储类型（s3/disk/region_s3）
	OpType                types.String `tfsdk:"op_type"`                   // 操作类型（auto/manual）
	TaskType              types.String `tfsdk:"task_type"`                 // 备份类型（full/incr）
	TaskStatus            types.Int32  `tfsdk:"task_status"`               // 任务状态（100/101/102/1/-1）
	BackedUpDataSize      types.Int64  `tfsdk:"backed_up_data_size"`       // 备份大小（字节）
	BackedUpDataSizeHuman types.String `tfsdk:"backed_up_data_size_human"` // 展示大小（自动适配单位）
	BackupStartTime       types.String `tfsdk:"backup_start_time"`         // 备份开始时间
	BackupEndTime         types.String `tfsdk:"backup_end_time"`           // 备份结束时间
	Disabled              types.Bool   `tfsdk:"disabled"`                  // 禁用备份
}
type BackupModel struct {
	InstId   types.String        `tfsdk:"instance_id"`
	BackupID types.Int64         `tfsdk:"id"`
	InstName types.String        `tfsdk:"instance_name"` // 实例名称
	Records  []BackupRecordModel `tfsdk:"records"`       // 备份记录集合 (元素类型为 BackupRecordModel)
}

type CtyunMysqlBackupsConfig struct {
	RegionID   types.String  `tfsdk:"region_id"`
	InstID     types.String  `tfsdk:"instance_id"`
	InstName   types.String  `tfsdk:"instance_name"`
	ProjectID  types.String  `tfsdk:"project_id"`
	Name       types.String  `tfsdk:"name"`
	BackupID   types.Int64   `tfsdk:"id"`
	PageSize   types.Int32   `tfsdk:"page_size"`
	PageNo     types.Int32   `tfsdk:"page_no"`
	BackupList []BackupModel `tfsdk:"backups"`
}
