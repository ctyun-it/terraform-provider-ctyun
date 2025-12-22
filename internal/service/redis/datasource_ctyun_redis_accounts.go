package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunRedisAccounts{}
	_ datasource.DataSourceWithConfigure = &ctyunRedisAccounts{}
)

type ctyunRedisAccounts struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisAccounts() datasource.DataSource {
	return &ctyunRedisAccounts{}
}

func (c *ctyunRedisAccounts) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_accounts"
}

type CtyunRedisAccountModel struct {
	Name        types.String `tfsdk:"name"`
	Privilege   types.String `tfsdk:"privilege"`
	Description types.String `tfsdk:"description"`
}

type CtyunRedisAccountsConfig struct {
	RegionID   types.String             `tfsdk:"region_id"`
	InstanceId types.String             `tfsdk:"instance_id"`
	Accounts   []CtyunRedisAccountModel `tfsdk:"accounts"`
}

func (c *ctyunRedisAccounts) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10403139`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
			},
			"accounts": schema.ListNestedAttribute{
				Computed:    true,
				Description: "账户列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "账户名称",
						},
						"privilege": schema.StringAttribute{
							Computed:    true,
							Description: "账户权限",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "账户描述信息",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRedisAccounts) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRedisAccountsConfig
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
	instanceId := config.InstanceId.ValueString()
	if instanceId == "" {
		err = fmt.Errorf("instanceId不能为空")
		return
	}

	// 组装请求体
	params := &dcs2.Dcs2DescribeAccountsRequest{
		RegionId:   regionId,
		ProdInstId: instanceId,
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeAccountsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	config.Accounts = []CtyunRedisAccountModel{}
	for _, account := range resp.ReturnObj.Rows {
		item := CtyunRedisAccountModel{
			Name:        types.StringValue(account.Name),
			Privilege:   types.StringValue(account.RawType),
			Description: types.StringValue(account.AccountDescription),
		}
		config.Accounts = append(config.Accounts, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRedisAccounts) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
