package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPgsqlBackups{}
	_ datasource.DataSourceWithConfigure = &ctyunPgsqlBackups{}
)

type ctyunPgsqlBackups struct {
	meta *common.CtyunMetadata
}

func NewCtyunPgsqlBackups() datasource.DataSource {
	return &ctyunPgsqlBackups{}
}
func (c *ctyunPgsqlBackups) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunPgsqlBackups) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_backups"
}

func (c *ctyunPgsqlBackups) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10160072",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目ID",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "PostgreSQL实例ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "备份名称过滤条件",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "分页页码，默认为1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数，默认为10",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "备份类型过滤条件（auto：自动备份，manual：手动备份，recovery：恢复备份）",
				Validators: []validator.String{
					stringvalidator.OneOf("auto", "manual", "recovery"),
				},
			},
			"backups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "备份ID",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "实例ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "备份名称",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "备份类型（auto：自动备份，manual：手动备份，recovery：恢复备份）",
						},
						"result": schema.StringAttribute{
							Computed:    true,
							Description: "备份结果（success：成功，failed：失败，ing：运行中）",
						},
						"start_time": schema.StringAttribute{
							Computed:    true,
							Description: "备份开始时间",
						},
						"end_time": schema.StringAttribute{
							Computed:    true,
							Description: "备份结束时间",
						},
						"data_len": schema.StringAttribute{
							Computed:    true,
							Description: "数据长度（格式化）",
						},
						"compress_size": schema.StringAttribute{
							Computed:    true,
							Description: "备份压缩大小（格式化）",
						},
					},
				},
				Description: "PostgreSQL备份列表",
			},
		},
	}
}

func (c *ctyunPgsqlBackups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPostgresqlBackupsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)
	respList, err := c.getPgsqlBackupList(ctx, config)
	if err != nil {
		return
	}
	var postgresqlBackups []PostgresqlBackupsInfoList
	for _, backupItem := range respList.ReturnObj.List {
		var backup PostgresqlBackupsInfoList
		backup.ID = types.Int64Value(backupItem.Id)
		backup.ProdInstId = types.StringValue(backupItem.ProdInstId)
		backup.BackupName = types.StringValue(backupItem.BackupName)
		backup.BackupType = types.StringValue(business.PgsqlBackupTypeMapConv[backupItem.Type])
		backup.BackupResult = types.StringValue(business.PgsqlBackupResultMapConv[backupItem.Result])
		backup.StartTime = types.StringValue(utils.FromBJTimeToUTCZ(backupItem.StartTime))
		backup.EndTime = types.StringValue(utils.FromBJTimeToUTCZ(backupItem.EndTime))
		backup.DataLen = types.StringValue(backupItem.DataLen)
		backup.BackupCompressSize = types.StringValue(backupItem.BackupCompressSize)
		postgresqlBackups = append(postgresqlBackups, backup)
	}
	config.PostgresqlBackups = postgresqlBackups
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunPgsqlBackups) getPgsqlBackupList(ctx context.Context, config CtyunPostgresqlBackupsConfig) (*pgsql.PgsqlGetBackupListResponse, error) {
	params := &pgsql.PgsqlGetBackupListRequest{
		ProdInstId: config.InstID.ValueString(),
		PageNum:    1,
		PageSize:   10,
		BackupName: config.Name.ValueStringPointer(),
	}
	header := &pgsql.PgsqlGetBackupListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetBackupListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询postgresql实例(id=%s)的备份集(name=%s)失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

type PostgresqlBackupsInfoList struct {
	ID                 types.Int64  `tfsdk:"id"`
	ProdInstId         types.String `tfsdk:"instance_id"`
	BackupName         types.String `tfsdk:"name"`
	BackupType         types.String `tfsdk:"type"`
	BackupResult       types.String `tfsdk:"result"`
	StartTime          types.String `tfsdk:"start_time"`
	EndTime            types.String `tfsdk:"end_time"`
	DataLen            types.String `tfsdk:"data_len"`
	BackupCompressSize types.String `tfsdk:"compress_size"`
}

type CtyunPostgresqlBackupsConfig struct {
	RegionID          types.String                `tfsdk:"region_id"`
	ProjectID         types.String                `tfsdk:"project_id"`
	InstID            types.String                `tfsdk:"instance_id"`
	Name              types.String                `tfsdk:"name"`
	PageNo            types.Int32                 `tfsdk:"page_no"`
	PageSize          types.Int32                 `tfsdk:"page_size"`
	BackupType        types.String                `tfsdk:"type"`
	PostgresqlBackups []PostgresqlBackupsInfoList `tfsdk:"backups"`
}
