package mongodb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &CtyunMongodbAccounts{}
	_ datasource.DataSourceWithConfigure = &CtyunMongodbAccounts{}
)

func NewCtyunMongodbAccounts() datasource.DataSource {
	return &CtyunMongodbAccounts{}
}

type CtyunMongodbAccounts struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMongodbAccounts) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_accounts"
}

func (c *CtyunMongodbAccounts) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10089535`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源唯一标识",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MongoDB实例ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "当前页",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "单页记录条数",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"accounts": schema.ListNestedAttribute{
				Computed:    true,
				Description: "账号列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "账号名称",
						},
						"database": schema.StringAttribute{
							Computed:    true,
							Description: "数据库名称",
						},
						"roles": schema.ListNestedAttribute{
							Computed:    true,
							Description: "角色列表",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"db": schema.StringAttribute{
										Computed:    true,
										Description: "数据库名称",
									},
									"role": schema.StringAttribute{
										Computed:    true,
										Description: "角色",
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

func (c *CtyunMongodbAccounts) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMongodbAccounts) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MongodbAccountsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeReq := &mongodb.MongodbDescribeAccountsRequest{
		ProdInstId: data.InstanceID.ValueString(),
		PageNow:    data.PageNow.ValueInt32(),
		PageSize:   data.PageSize.ValueInt32(),
	}

	headers := &mongodb.MongodbDescribeAccountsRequestHeaders{
		RegionID: data.RegionID.ValueString(),
	}
	if !data.ProjectID.IsNull() {
		headers.ProjectID = data.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "查询MongoDB账号列表", map[string]interface{}{
		"instance_id": data.InstanceID.ValueString(),
	})

	describeResp, err := c.meta.Apis.SdkMongodbApis.MongodbDescribeAccountsApi.Do(ctx, c.meta.Credential, describeReq, headers)
	if err != nil {
		resp.Diagnostics.AddError("查询MongoDB账号列表失败", err.Error())
		return
	}

	if describeResp.StatusCode != 800 {
		if describeResp.Message != nil {
			resp.Diagnostics.AddError("查询MongoDB账号列表失败", fmt.Sprintf("API返回错误: %s", *describeResp.Message))
		} else {
			resp.Diagnostics.AddError("查询MongoDB账号列表失败", fmt.Sprintf("API返回错误，状态码: %d", describeResp.StatusCode))
		}
		return
	}

	if describeResp.ReturnObj == nil {
		resp.Diagnostics.AddError("查询MongoDB账号列表失败", "API返回空结果")
		return
	}

	// 转换API响应数据到Terraform模型
	var accounts []MongodbAccountModel
	for _, account := range describeResp.ReturnObj.List {
		accountModel := MongodbAccountModel{
			Name:     types.StringValue(account.User),
			Database: types.StringValue(account.DB),
		}

		// 转换角色信息
		var roles []MongodbAccountRoleModel
		for _, role := range account.Roles {
			roles = append(roles, MongodbAccountRoleModel{
				DB:   types.StringValue(role.DB),
				Role: types.StringValue(role.Role),
			})
		}
		accountModel.Roles = roles

		accounts = append(accounts, accountModel)
	}

	data.Accounts = accounts
	data.ID = types.StringValue(data.InstanceID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type MongodbAccountRoleModel struct {
	DB   types.String `tfsdk:"db"`
	Role types.String `tfsdk:"role"`
}

type MongodbAccountModel struct {
	Name     types.String              `tfsdk:"name"`
	Database types.String              `tfsdk:"database"`
	Roles    []MongodbAccountRoleModel `tfsdk:"roles"`
}

type MongodbAccountsDataSourceModel struct {
	ID         types.String          `tfsdk:"id"`
	InstanceID types.String          `tfsdk:"instance_id"`
	RegionID   types.String          `tfsdk:"region_id"`
	ProjectID  types.String          `tfsdk:"project_id"`
	PageNow    types.Int32           `tfsdk:"page_no"`
	PageSize   types.Int32           `tfsdk:"page_size"`
	Accounts   []MongodbAccountModel `tfsdk:"accounts"`
}
