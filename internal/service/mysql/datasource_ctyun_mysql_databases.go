package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunMysqlDatabases{}
	_ datasource.DataSourceWithConfigure = &ctyunMysqlDatabases{}
)

type ctyunMysqlDatabases struct {
	meta *common.CtyunMetadata
}

func NewCtyunMysqlDatabases() datasource.DataSource {
	return &ctyunMysqlDatabases{}
}
func (c *ctyunMysqlDatabases) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunMysqlDatabases) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_databases"
}

func (c *ctyunMysqlDatabases) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10140487",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MySQL实例ID",
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
			"databases": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "数据库名称",
						},
						"user_vo_list": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_name": schema.StringAttribute{
										Computed:    true,
										Description: "数据库账号名称",
									},
									"read_only": schema.BoolAttribute{
										Computed:    true,
										Description: "是否只读权限",
									},
									"select_priv": schema.BoolAttribute{
										Computed:    true,
										Description: "SELECT权限",
									},
									"insert_priv": schema.BoolAttribute{
										Computed:    true,
										Description: "INSERT权限",
									},
								},
							},
							Description: "数据库用户权限列表",
						},
					},
				},
				Description: "MySQL数据库列表",
			},
		},
	}
}

func (c *ctyunMysqlDatabases) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMysqlDatabasesConfig
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
	resp, err := c.getMysqlDatabasesInfo(ctx, config)
	var databases []CtyunMysqlDatabaseModel
	for _, schemaItem := range resp {
		var database CtyunMysqlDatabaseModel
		database.DBName = types.StringValue(schemaItem.GrantSchema)
		var userVOList []CtyunMysqlUserDatabasePrivilegeModel
		for _, privilegeItem := range schemaItem.UserVOList {
			var privilege CtyunMysqlUserDatabasePrivilegeModel
			privilege.InsertPriv = types.BoolValue(business.PrivilegeMap[privilegeItem.InsertPriv])
			privilege.SelectPriv = types.BoolValue(business.PrivilegeMap[privilegeItem.SelectPriv])
			privilege.AccountName = types.StringValue(privilegeItem.AccountName)
			privilege.ReadOnly = types.BoolValue(privilegeItem.ReadOnly)
			userVOList = append(userVOList, privilege)
		}
		database.UserVOList = userVOList
		databases = append(databases, database)
	}
	config.MysqlDatabases = databases
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunMysqlDatabases) getMysqlDatabasesInfo(ctx context.Context, config CtyunMysqlDatabasesConfig) ([]mysql.TeledbGetDatabaseSchemaResponseReturnObj, error) {
	params := &mysql.TeledbGetDatabaseSchemaRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbGetDatabaseSchemaRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetDatabaseSchemaApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询mysql实例(id=%s)的database schema列表失败，接口返回nil。请联系研发确认问题原因！", config.InstID.ValueString())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil || len(resp.ReturnObj) < 1 {
		err = common.InvalidReturnObjError
		return nil, err
	}

	return resp.ReturnObj, nil
}

type CtyunMysqlUserDatabasePrivilegeModel struct {
	AccountName types.String `tfsdk:"account_name"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	SelectPriv  types.Bool   `tfsdk:"select_priv"`
	InsertPriv  types.Bool   `tfsdk:"insert_priv"`
}

type CtyunMysqlDatabaseModel struct {
	DBName     types.String                           `tfsdk:"name"`
	UserVOList []CtyunMysqlUserDatabasePrivilegeModel `tfsdk:"user_vo_list"`
}

type CtyunMysqlDatabasesConfig struct {
	InstID         types.String              `tfsdk:"instance_id"`
	ProjectID      types.String              `tfsdk:"project_id"`
	RegionID       types.String              `tfsdk:"region_id"`
	MysqlDatabases []CtyunMysqlDatabaseModel `tfsdk:"databases"`
}
