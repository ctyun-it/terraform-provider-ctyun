package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebsbackup "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebsbackup"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEbsBackups{}
	_ datasource.DataSourceWithConfigure = &ctyunEbsBackups{}
)

type ctyunEbsBackups struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbsBackups() datasource.DataSource {
	return &ctyunEbsBackups{}
}

func (c *ctyunEbsBackups) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_backups"
}

type ctyunEbsBackupsModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	BackupStatus        types.String `tfsdk:"backup_status"`
	Description         types.String `tfsdk:"description"`
	DiskID              types.String `tfsdk:"disk_id"`
	DiskName            types.String `tfsdk:"disk_name"`
	RepositoryID        types.String `tfsdk:"repository_id"`
	RepositoryName      types.String `tfsdk:"repository_name"`
	RepositoryFreeze    types.Bool   `tfsdk:"repository_freeze"`
	DiskTotalSize       types.Int64  `tfsdk:"disk_size"`
	BackupSize          types.Int64  `tfsdk:"backup_size"`
	RestoreFinishedTime types.String `tfsdk:"restore_finished_time"`
	CreatedTime         types.String `tfsdk:"created_time"`
	FinishedTime        types.String `tfsdk:"finished_time"`
	UpdatedTime         types.String `tfsdk:"updated_time"`
	RestoredTime        types.String `tfsdk:"restored_time"`
	Encrypted           types.Bool   `tfsdk:"encrypted"`
	DiskType            types.String `tfsdk:"disk_type"`
	Paas                types.Bool   `tfsdk:"paas"`
	InstanceID          types.String `tfsdk:"instance_id"`
	InstanceName        types.String `tfsdk:"instance_name"`
	ProjectID           types.String `tfsdk:"project_id"`
	BackupType          types.String `tfsdk:"backup_type"`
}

type ctyunEbsBackupsConfig struct {
	RegionID     types.String           `tfsdk:"region_id"`
	PageNo       types.Int32            `tfsdk:"page_no"`
	PageSize     types.Int32            `tfsdk:"page_size"`
	DiskID       types.String           `tfsdk:"disk_id"`
	DiskName     types.String           `tfsdk:"disk_name"`
	RepositoryID types.String           `tfsdk:"repository_id"`
	Name         types.String           `tfsdk:"name"`
	QueryContent types.String           `tfsdk:"query_content"`
	BackupStatus types.String           `tfsdk:"backup_status"`
	ProjectID    types.String           `tfsdk:"project_id"`
	Backups      []ctyunEbsBackupsModel `tfsdk:"backups"`
}

func (c *ctyunEbsBackups) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026752/10037428**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数目，取值范围：[1, 50]，注：默认值为10",
			},
			"disk_id": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘ID",
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"disk_name": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘名称，模糊过滤",
			},
			"repository_id": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘备份存储库ID",
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘备份名称，模糊过滤",
			},
			"query_content": schema.StringAttribute{
				Optional:    true,
				Description: "该参数，可用于模糊过滤 云硬盘ID/云硬盘名称/备份ID/备份名称，即上述4个字段如果包含该参数的值，则会被过滤出来",
			},
			"backup_status": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘备份状态，取值范围：available：可用，error：错误，restoring：恢复中，creating：创建中，deleting：删除中，merging_backup：合并中，frozen：已冻结",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href=\"https://www.ctyun.cn/document/10026730/10238876\">创建企业项目</a>了解如何创建企业项目<br />注：默认值为\"0\"",
			},
			"backups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份名称",
						},
						"backup_status": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份状态，取值范围：available：可用，error：错误，restoring：恢复中，creating：创建中，deleting：删除中，merging_backup：合并中，frozen：已冻结",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份描述",
						},
						"disk_id": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘ID",
						},
						"disk_name": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘名称",
						},
						"repository_id": schema.StringAttribute{
							Computed:    true,
							Description: "备份存储库ID",
						},
						"repository_name": schema.StringAttribute{
							Computed:    true,
							Description: "备份存储库名称",
						},
						"repository_freeze": schema.BoolAttribute{
							Computed:    true,
							Description: "备份是否冻结",
						},
						"disk_size": schema.Int64Attribute{
							Computed:    true,
							Description: "云硬盘大小，单位GB",
						},
						"backup_size": schema.Int64Attribute{
							Computed:    true,
							Description: "云硬盘备份大小，单位Byte",
						},
						"restore_finished_time": schema.StringAttribute{
							Computed:    true,
							Description: "使用该云硬盘备份恢复完成时间",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "备份创建时间",
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
						"encrypted": schema.BoolAttribute{
							Computed:    true,
							Description: "云硬盘是否加密",
						},
						"disk_type": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘类型，取值范围：SATA：普通IO，SAS：高IO，SSD：超高IO，FAST-SSD：极速型SSD，XSSD-0、XSSD-1、XSSD-2：X系列云硬盘",
						},
						"paas": schema.BoolAttribute{
							Computed:    true,
							Description: "是否支持PAAS",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘挂载的云主机ID",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘挂载的云主机名称",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID",
						},
						"backup_type": schema.StringAttribute{
							Computed:    true,
							Description: "备份类型，取值范围：full-backup：全量备份，incremental-backup：增量备份",
						},
					},
				}},
		},
	}
}

func (c *ctyunEbsBackups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config ctyunEbsBackupsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}

	config.RegionID = types.StringValue(regionId)
	// 组装请求体
	params := &ctebsbackup.EbsbackupListBackupRequest{
		RegionID:     config.RegionID.ValueString(),
		PageNo:       config.PageNo.ValueInt32(),
		PageSize:     config.PageSize.ValueInt32(),
		DiskID:       config.DiskID.ValueString(),
		DiskName:     config.DiskName.ValueString(),
		RepositoryID: config.RepositoryID.ValueString(),
		BackupName:   config.Name.ValueString(),
		QueryContent: config.QueryContent.ValueString(),
		BackupStatus: config.BackupStatus.ValueString(),
		ProjectID:    config.ProjectID.ValueString(),
	}

	// 调用API
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupListBackupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	// 解析返回值
	config.Backups = []ctyunEbsBackupsModel{}
	for _, backup := range resp.ReturnObj.BackupList {
		item := ctyunEbsBackupsModel{
			Id:                  types.StringValue(backup.BackupID),
			Name:                types.StringValue(backup.BackupName),
			BackupStatus:        types.StringValue(backup.BackupStatus),
			Description:         types.StringValue(backup.Description),
			DiskID:              types.StringValue(backup.DiskID),
			DiskName:            types.StringValue(backup.DiskName),
			RepositoryID:        types.StringValue(backup.RepositoryID),
			RepositoryName:      types.StringValue(backup.RepositoryName),
			RepositoryFreeze:    types.BoolValue(*backup.Freeze),
			DiskTotalSize:       types.Int64Value(int64(backup.DiskSize)),
			BackupSize:          types.Int64Value(int64(backup.BackupSize)),
			RestoreFinishedTime: types.StringValue(fmt.Sprintf("%d", backup.RestoreFinishedTime)),
			CreatedTime:         types.StringValue(fmt.Sprintf("%d", backup.CreatedTime)),
			FinishedTime:        types.StringValue(fmt.Sprintf("%d", backup.FinishedTime)),
			UpdatedTime:         types.StringValue(fmt.Sprintf("%d", backup.UpdatedTime)),
			RestoredTime:        types.StringValue(fmt.Sprintf("%d", backup.RestoredTime)),
			Encrypted:           types.BoolValue(*backup.Encrypted),
			DiskType:            types.StringValue(backup.DiskType),
			Paas:                types.BoolValue(*backup.Paas),
			InstanceID:          types.StringValue(backup.InstanceID),
			InstanceName:        types.StringValue(backup.InstanceName),
			ProjectID:           types.StringValue(backup.ProjectID),
			BackupType:          types.StringValue(backup.BackupType),
		}

		config.Backups = append(config.Backups, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEbsBackups) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
