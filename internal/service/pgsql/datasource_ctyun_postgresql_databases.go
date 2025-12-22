package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPostgresqlDatabases{}
	_ datasource.DataSourceWithConfigure = &ctyunPostgresqlDatabases{}
)

type ctyunPostgresqlDatabases struct {
	meta *common.CtyunMetadata
}

func NewCtyunPostgresqlDatabases() datasource.DataSource {
	return &ctyunPostgresqlDatabases{}
}
func (c *ctyunPostgresqlDatabases) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunPostgresqlDatabases) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_databases"
}

func (c *ctyunPostgresqlDatabases) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10159978",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "PostgreSQL实例ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目ID",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "数据库名称过滤条件",
			},
			"databases": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "实例ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "数据库名称",
						},
						"charset_name": schema.StringAttribute{
							Computed:    true,
							Description: "字符集名称",
						},
						"charset_collate": schema.StringAttribute{
							Computed:    true,
							Description: "排序规则",
						},
						"charset_type": schema.StringAttribute{
							Computed:    true,
							Description: "字符集类型",
						},
						"owner": schema.StringAttribute{
							Computed:    true,
							Description: "数据库所有者",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "数据库描述",
						},
					},
				},
				Description: "PostgreSQL数据库列表",
			},
		},
	}
}

func (c *ctyunPostgresqlDatabases) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPostgresqlDatabases
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
	databaseRespList, err := c.getDatabaseList(ctx, config)
	if err != nil {
		return
	}
	var postgresqlDatabases []CtyunPostgresqlDatabaseModel
	for _, databaseItem := range databaseRespList {
		var databaseInfo CtyunPostgresqlDatabaseModel
		databaseInfo.InstID = types.StringValue(databaseItem.ProdInstId)
		databaseInfo.DBName = types.StringValue(databaseItem.DBName)
		databaseInfo.CharSetName = types.StringValue(databaseItem.DBEncoding)
		databaseInfo.CharSetCollate = types.StringValue(databaseItem.DBCollate)
		databaseInfo.CharSetType = types.StringValue(databaseItem.DbType)
		databaseInfo.Owner = types.StringValue(databaseItem.DBOwner)
		databaseInfo.Description = types.StringValue(databaseItem.DBDescription)
		postgresqlDatabases = append(postgresqlDatabases, databaseInfo)
	}
	config.PostgresqlDatabases = postgresqlDatabases
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunPostgresqlDatabases) getDatabaseList(ctx context.Context, config CtyunPostgresqlDatabases) ([]pgsql.PgsqlGetDatabaseSchemaResponseReturnObj, error) {
	params := &pgsql.PgsqlGetDatabaseSchemaRequest{
		ProdInstId: config.InstID.ValueString(),
	}
	if !config.Name.IsNull() {
		params.DBName = config.Name.ValueStringPointer()
	}

	header := &pgsql.PgsqlGetDatabaseSchemaRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetDatabaseSchemaApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)查询数据库详情，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj, nil

}

type CtyunPostgresqlDatabaseModel struct {
	InstID         types.String `tfsdk:"instance_id"`
	DBName         types.String `tfsdk:"name"`
	CharSetName    types.String `tfsdk:"charset_name"`
	CharSetCollate types.String `tfsdk:"charset_collate"`
	CharSetType    types.String `tfsdk:"charset_type"`
	Owner          types.String `tfsdk:"owner"`
	Description    types.String `tfsdk:"description"`
}

type CtyunPostgresqlDatabases struct {
	InstID              types.String                   `tfsdk:"instance_id"`
	ProjectID           types.String                   `tfsdk:"project_id"`
	RegionID            types.String                   `tfsdk:"region_id"`
	Name                types.String                   `tfsdk:"name"`
	PostgresqlDatabases []CtyunPostgresqlDatabaseModel `tfsdk:"databases"`
}
