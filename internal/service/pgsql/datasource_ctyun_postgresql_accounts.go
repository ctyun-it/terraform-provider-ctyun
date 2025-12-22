package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPgsqlAccounts{}
	_ datasource.DataSourceWithConfigure = &ctyunPgsqlAccounts{}
)

type ctyunPgsqlAccounts struct {
	meta *common.CtyunMetadata
}

func NewCtyunPgsqlAccounts() datasource.DataSource {
	return &ctyunPgsqlAccounts{}
}
func (c *ctyunPgsqlAccounts) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunPgsqlAccounts) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_accounts"
}

func (c *ctyunPgsqlAccounts) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10161317",
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
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "PostgreSQL实例的账户名称过滤条件",
			},
			"accounts": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "账户名称",
						},
						"rol_super": schema.BoolAttribute{
							Computed:    true,
							Description: "是否具有超级用户权限",
						},
						"rol_inherit": schema.BoolAttribute{
							Computed:    true,
							Description: "用户是否自动继承其所属角色的权限",
						},
						"rol_create_role": schema.BoolAttribute{
							Computed:    true,
							Description: "用户是否支持创建其他子用户",
						},
						"rol_create_db": schema.BoolAttribute{
							Computed:    true,
							Description: "用户是否可以创建数据库",
						},
						"rol_can_login": schema.BoolAttribute{
							Computed:    true,
							Description: "用户是否可以登录数据库",
						},
						"rol_conn_limit": schema.Int32Attribute{
							Computed:    true,
							Description: "用户连接数限制（-1表示无限制）",
						},
						"rol_by_pass_rls": schema.BoolAttribute{
							Computed:    true,
							Description: "用户是否绕过每个行级安全策略",
						},
					},
				},
				Description: "PostgreSQL账户列表",
			},
		},
	}
}

func (c *ctyunPgsqlAccounts) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPgsqlAccounts
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
	accountRespList, err := c.getPgsqlAccountList(ctx, config)
	if err != nil {
		return
	}

	var postgresqlAccounts []PgsqlAccountInfoModel
	for _, accountItem := range accountRespList {
		var accountInfo PgsqlAccountInfoModel
		accountInfo.AccountName = types.StringValue(accountItem.Username)
		accountInfo.RolConnLimit = types.Int32Value(accountItem.RolConnLimit)
		accountInfo.RolInherit = types.BoolValue(accountItem.RolInherit)
		accountInfo.RolCreateRole = types.BoolValue(accountItem.RolCreateRole)
		accountInfo.RolCreateDB = types.BoolValue(accountItem.RolCreateDB)
		accountInfo.RolCanLogin = types.BoolValue(accountItem.RolCanLogin)
		accountInfo.RolByPassRls = types.BoolValue(accountItem.RolByPassRls)
		postgresqlAccounts = append(postgresqlAccounts, accountInfo)
	}
	config.PostgresqlAccounts = postgresqlAccounts
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunPgsqlAccounts) getPgsqlAccountList(ctx context.Context, config CtyunPgsqlAccounts) ([]pgsql.PgsqlGetAccountListResponseReturnObj, error) {
	params := &pgsql.PgsqlGetAccountListRequest{
		ProdInstId: config.InstID.ValueString(),
		PageNum:    1,
		PageSize:   10,
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}
	if !config.PageNo.IsNull() {
		params.PageNum = config.PageNo.ValueInt32()
	}

	if !config.Name.IsNull() {
		params.Username = config.Name.ValueStringPointer()
	}
	header := &pgsql.PgsqlGetAccountListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetAccountListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询postgresql实例(id=%s)的账户信息(name=%s)失败，接口返回nil。请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}

	return resp.ReturnObj.List, nil
}

type PgsqlAccountInfoModel struct {
	AccountName   types.String `tfsdk:"name"`
	RolSuper      types.Bool   `tfsdk:"rol_super"`
	RolInherit    types.Bool   `tfsdk:"rol_inherit"`
	RolCreateRole types.Bool   `tfsdk:"rol_create_role"`
	RolCreateDB   types.Bool   `tfsdk:"rol_create_db"`
	RolCanLogin   types.Bool   `tfsdk:"rol_can_login"`
	RolConnLimit  types.Int32  `tfsdk:"rol_conn_limit"`
	RolByPassRls  types.Bool   `tfsdk:"rol_by_pass_rls"`
}

type CtyunPgsqlAccounts struct {
	RegionID           types.String            `tfsdk:"region_id"`
	ProjectID          types.String            `tfsdk:"project_id"`
	InstID             types.String            `tfsdk:"instance_id"`
	PageNo             types.Int32             `tfsdk:"page_no"`
	PageSize           types.Int32             `tfsdk:"page_size"`
	Name               types.String            `tfsdk:"name"`
	PostgresqlAccounts []PgsqlAccountInfoModel `tfsdk:"accounts"`
}
