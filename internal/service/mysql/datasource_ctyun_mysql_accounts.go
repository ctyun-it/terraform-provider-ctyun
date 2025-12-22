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
	_ datasource.DataSource              = &ctyunMysqlAccounts{}
	_ datasource.DataSourceWithConfigure = &ctyunMysqlAccounts{}
)

type ctyunMysqlAccounts struct {
	meta *common.CtyunMetadata
}

func NewCtyunMysqlAccounts() datasource.DataSource {
	return &ctyunMysqlAccounts{}
}
func (c *ctyunMysqlAccounts) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunMysqlAccounts) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_accounts"
}

func (c *ctyunMysqlAccounts) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10133363",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MySQL实例ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
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
				Description: "数据库账号名称",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"accounts": schema.ListNestedAttribute{
				Computed:    true,
				Description: "MySQL账户权限列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "数据库账号名称",
						},
						"schema_privilege_list": schema.ListNestedAttribute{
							Computed:    true,
							Description: "数据库权限列表",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"grant_schema": schema.StringAttribute{
										Computed:    true,
										Description: "授权数据库名称",
									},
									"privilege": schema.StringAttribute{
										Computed:    true,
										Description: "权限类型",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunMysqlAccounts) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMysqlAccountsConfig
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
	accounts, err := c.getMysqlAccountInfo(ctx, &config)
	if err != nil {
		return
	}
	var mysqlAccounts []CtyunMysqlAccountPrivilegeModel
	for _, accountItem := range accounts {
		var mysqlAccount CtyunMysqlAccountPrivilegeModel
		mysqlAccount.AccountName = types.StringValue(accountItem.AccountName)
		var schemaPrivilegeList []SchemaPrivilegeModel
		for _, privilege := range accountItem.SchemaPrivilegeVOList {
			var schemaPrivilege SchemaPrivilegeModel
			schemaPrivilege.Privilege = types.StringValue(c.getPrivilege(privilege))
			schemaPrivilege.GrantSchema = types.StringValue(privilege.GrantSchema)
			schemaPrivilegeList = append(schemaPrivilegeList, schemaPrivilege)
		}
		mysqlAccount.SchemaPrivilegeList = schemaPrivilegeList
		mysqlAccounts = append(mysqlAccounts, mysqlAccount)
	}
	config.MysqlAccounts = mysqlAccounts
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunMysqlAccounts) getMysqlAccountInfo(ctx context.Context, config *CtyunMysqlAccountsConfig) ([]mysql.TeledbGetAccountInfoResponseReturnObj, error) {
	params := &mysql.TeledbGetAccountInfoRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbGetAccountInfoRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetAccountInfoApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询mysql实例(id=%s)的用户权限列表失败", config.InstID.ValueString())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("get mysql account failed, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	var accountPrivilegeList []mysql.TeledbGetAccountInfoResponseReturnObj
	if !config.Name.IsNull() {
		for _, accountPrivilege := range resp.ReturnObj {
			accountName := accountPrivilege.AccountName
			if accountName == config.Name.ValueString() {
				accountPrivilegeList = append(accountPrivilegeList, accountPrivilege)
				return accountPrivilegeList, nil
			}
		}
	} else {
		return resp.ReturnObj, nil
	}

	return nil, fmt.Errorf("mysql实例(id=%s)不存在name=%s的权限配置", config.InstID.ValueString(), config.Name.ValueString())
}

func (c *ctyunMysqlAccounts) getPrivilege(item mysql.SchemaPrivilegeVO) string {
	if *item.DDLPrivilege {
		return business.MysqlSchemaPrivilegeDDL
	} else if *item.ReadOnly {
		return business.MysqlSchemaPrivilegeReadOnly
	} else if *item.DMLPrivilege {
		return business.MysqlSchemaPrivilegeDML
	}
	return ""
}

type SchemaPrivilegeModel struct {
	GrantSchema types.String `tfsdk:"grant_schema"` // 数据库schema
	Privilege   types.String `tfsdk:"privilege"`
}

type CtyunMysqlAccountPrivilegeModel struct {
	AccountName         types.String           `tfsdk:"name"`
	SchemaPrivilegeList []SchemaPrivilegeModel `tfsdk:"schema_privilege_list"`
}
type CtyunMysqlAccountsConfig struct {
	InstID        types.String                      `tfsdk:"instance_id"`
	ProjectID     types.String                      `tfsdk:"project_id"`
	RegionID      types.String                      `tfsdk:"region_id"`
	Name          types.String                      `tfsdk:"name"`
	MysqlAccounts []CtyunMysqlAccountPrivilegeModel `tfsdk:"accounts"`
}
