package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEcsBackups{}
	_ datasource.DataSourceWithConfigure = &ctyunEcsBackups{}
)

type ctyunEcsBackups struct {
	meta *common.CtyunMetadata
}

func NewCtyunEcsBackups() datasource.DataSource {
	return &ctyunEcsBackups{}
}

func (c *ctyunEcsBackups) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_backups"
}

type ctyunEcsBackupsModel struct {
	InstanceBackupID          types.String `tfsdk:"instance_backup_id"`
	InstanceBackupName        types.String `tfsdk:"name"`
	InstanceBackupStatus      types.String `tfsdk:"instance_backup_status"`
	InstanceBackupDescription types.String `tfsdk:"instance_backup_description"`
	InstanceID                types.String `tfsdk:"instance_id"`
	InstanceName              types.String `tfsdk:"instance_name"`
	RepositoryID              types.String `tfsdk:"repository_id"`
	RepositoryName            types.String `tfsdk:"repository_name"`
	RepositoryExpired         types.Bool   `tfsdk:"repository_expired"`
	RepositoryFreeze          types.Bool   `tfsdk:"repository_freeze"`
	DiskTotalSize             types.Int64  `tfsdk:"disk_total_size"`
	UsedSize                  types.Int64  `tfsdk:"used_size"`
	DiskCount                 types.Int64  `tfsdk:"disk_count"`
	RestoreFinishedTime       types.String `tfsdk:"restore_finished_time"`
	CreatedTime               types.String `tfsdk:"created_time"`
	FinishedTime              types.String `tfsdk:"finished_time"`
	ProjectID                 types.String `tfsdk:"project_id"`
	BackupType                types.String `tfsdk:"backup_type"`
}

type ctyunEcsBackupsConfig struct {
	RegionID             types.String           `tfsdk:"region_id"`
	PageNo               types.Int32            `tfsdk:"page_no"`
	PageSize             types.Int32            `tfsdk:"page_size"`
	InstanceID           types.String           `tfsdk:"instance_id"`
	RepositoryID         types.String           `tfsdk:"repository_id"`
	InstanceBackupID     types.String           `tfsdk:"instance_backup_id"`
	QueryContent         types.String           `tfsdk:"query_content"`
	InstanceBackupStatus types.String           `tfsdk:"instance_backup_status"`
	ProjectID            types.String           `tfsdk:"project_id"`
	Backups              []ctyunEcsBackupsModel `tfsdk:"backups"`
}

func (c *ctyunEcsBackups) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026751/10033761**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数目，取值范围：[1, 50]，注：默认值为10",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "云主机ID",
			},
			"repository_id": schema.StringAttribute{
				Optional:    true,
				Description: "云主机备份存储库ID",
			},
			"instance_backup_id": schema.StringAttribute{
				Optional:    true,
				Description: "云主机备份ID",
			},
			"query_content": schema.StringAttribute{
				Optional:    true,
				Description: "模糊匹配查询内容（匹配字段：instanceBackupName、instanceBackupID、instanceBackupStatus、instanceName）",
			},
			"instance_backup_status": schema.StringAttribute{
				Optional:    true,
				Description: "云主机备份状态",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看创建企业项目了解如何创建企业项目",
			},
			"backups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_backup_id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机备份ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机备份名称",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看创建企业项目了解如何创建企业项目",
						},
						"instance_backup_status": schema.StringAttribute{
							Computed:    true,
							Description: "云主机备份状态，取值范围：CREATING: 备份创建中, ACTIVE: 可用， RESTORING: 备份恢复中，DELETING: 删除中，EXPIRED：到期，ERROR：错误",
						},
						"instance_backup_description": schema.StringAttribute{
							Computed:    true,
							Description: "云主机备份描述",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机ID",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机名称",
						},
						"repository_id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机备份存储库ID",
						},
						"repository_name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机备份存储库名称",
						},
						"repository_expired": schema.BoolAttribute{
							Computed:    true,
							Description: "云主机备份存储库是否过期",
						},
						"repository_freeze": schema.BoolAttribute{
							Computed:    true,
							Description: "存储库是否冻结",
						},
						"disk_total_size": schema.Int64Attribute{
							Computed:    true,
							Description: "云盘总容量大小，单位为GB",
						},
						"used_size": schema.Int64Attribute{
							Computed:    true,
							Description: "磁盘备份已使用大小",
						},
						"disk_count": schema.Int64Attribute{
							Computed:    true,
							Description: "云硬盘数目",
						},
						"restore_finished_time": schema.StringAttribute{
							Computed:    true,
							Description: "备份恢复完成时间",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"finished_time": schema.StringAttribute{
							Computed:    true,
							Description: "完成时间",
						},
						"backup_type": schema.StringAttribute{
							Computed:    true,
							Description: "备份类型，取值范围：FULL：全量备份，INCREMENT：增量备份",
						},
					},
				}},
		},
	}
}

func (c *ctyunEcsBackups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config ctyunEcsBackupsConfig
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
	params := &ctecs2.CtecsListInstanceBackupV41Request{
		RegionID:             config.RegionID.ValueString(),
		PageNo:               config.PageNo.ValueInt32(),
		PageSize:             config.PageSize.ValueInt32(),
		InstanceID:           config.InstanceID.ValueString(),
		RepositoryID:         config.RepositoryID.ValueString(),
		InstanceBackupID:     config.InstanceBackupID.ValueString(),
		QueryContent:         config.QueryContent.ValueString(),
		InstanceBackupStatus: config.InstanceBackupStatus.ValueString(),
		ProjectID:            config.ProjectID.ValueString(),
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsListInstanceBackupV41Api.Do(ctx, c.meta.SdkCredential, params)
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
	config.Backups = []ctyunEcsBackupsModel{}
	for _, backup := range resp.ReturnObj.Results {
		item := ctyunEcsBackupsModel{
			InstanceBackupID:          types.StringValue(backup.InstanceBackupID),
			InstanceBackupName:        types.StringValue(backup.InstanceBackupName),
			InstanceBackupStatus:      types.StringValue(backup.InstanceBackupStatus),
			InstanceBackupDescription: types.StringValue(backup.InstanceBackupDescription),
			InstanceID:                types.StringValue(backup.InstanceID),
			InstanceName:              types.StringValue(backup.InstanceName),
			RepositoryID:              types.StringValue(backup.RepositoryID),
			RepositoryName:            types.StringValue(backup.RepositoryName),
			RepositoryExpired:         types.BoolValue(*backup.RepositoryExpired),
			RepositoryFreeze:          types.BoolValue(*backup.RepositoryFreeze),
			DiskTotalSize:             types.Int64Value(int64(backup.DiskTotalSize)),
			UsedSize:                  types.Int64Value(backup.UsedSize),
			DiskCount:                 types.Int64Value(int64(backup.DiskCount)),
			RestoreFinishedTime:       types.StringValue(backup.RestoreFinishedTime),
			CreatedTime:               types.StringValue(backup.CreatedTime),
			FinishedTime:              types.StringValue(backup.FinishedTime),
			ProjectID:                 types.StringValue(backup.ProjectID),
			BackupType:                types.StringValue(backup.BackupType),
		}

		config.Backups = append(config.Backups, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEcsBackups) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
