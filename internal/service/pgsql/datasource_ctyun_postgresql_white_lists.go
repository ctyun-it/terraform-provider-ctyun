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
	_ datasource.DataSource              = &CtyunPostgresqlWhiteLists{}
	_ datasource.DataSourceWithConfigure = &CtyunPostgresqlWhiteLists{}
)

type CtyunPostgresqlWhiteLists struct {
	meta *common.CtyunMetadata
}

func NewCtyunPostgresqlWhiteLists() *CtyunPostgresqlWhiteLists {
	return &CtyunPostgresqlWhiteLists{}
}
func (c *CtyunPostgresqlWhiteLists) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPostgresqlWhiteLists) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_white_lists"
}

func (c *CtyunPostgresqlWhiteLists) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10161484",
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
			"ip": schema.StringAttribute{
				Optional:    true,
				Description: "白名单ip,支持模糊查询",
			},
			"white_list": schema.SetAttribute{
				Computed:    true,
				Description: "当前白名单列表",
				ElementType: types.StringType,
			},
		},
	}
}

func (c *CtyunPostgresqlWhiteLists) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPostgresqlWhiteListsConfig
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

	params := &pgsql.PgsqlGetWhiteListRequest{
		ProdInstId: config.InstID.ValueString(),
	}
	if !config.Ip.IsNull() {
		params.IP = config.Ip.ValueStringPointer()
	}

	header := &pgsql.PgsqlGetWhiteListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetWhiteListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例获取白名单ip失败，接口返回nil，请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return
	}

	ips, diags := types.SetValueFrom(ctx, types.StringType, resp.ReturnObj)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return
	}
	config.WhiteList = ips
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunPostgresqlWhiteListsConfig struct {
	InstID    types.String `tfsdk:"instance_id"`
	RegionID  types.String `tfsdk:"region_id"`
	ProjectID types.String `tfsdk:"project_id"`
	Ip        types.String `tfsdk:"ip"`
	WhiteList types.Set    `tfsdk:"white_list"`
}
