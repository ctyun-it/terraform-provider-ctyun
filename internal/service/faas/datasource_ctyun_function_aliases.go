package faas

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunFunctionAliases{}
	_ datasource.DataSourceWithConfigure = &CtyunFunctionAliases{}
)

func NewCtyunFunctionAliases() datasource.DataSource {
	return &CtyunFunctionAliases{}
}

type CtyunFunctionAliases struct {
	meta *common.CtyunMetadata
}

func (c *CtyunFunctionAliases) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function_aliases"
}

func (c *CtyunFunctionAliases) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询函数别名列表", "函数工作流（FunctionGraph）", "https://www.ctyun.cn/document/10355289"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源唯一标识",
			},
			"function_name": schema.StringAttribute{
				Required:    true,
				Description: "函数名称",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID，如果不填则默认使用 provider ctyun 中的 region_id 或环境变量中的 CTYUN_REGION_ID",
			},
			"aliases": schema.ListNestedAttribute{
				Computed:    true,
				Description: "函数别名列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"alias_name": schema.StringAttribute{
							Computed:    true,
							Description: "别名名称",
						},
						"version_id": schema.StringAttribute{
							Computed:    true,
							Description: "主版本 ID",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "关于别名的描述",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"gray_type": schema.Int32Attribute{
							Computed:    true,
							Description: "灰度类型，当前支持：1、按百分比随机灰度",
						},
						"gray_version_id": schema.StringAttribute{
							Computed:    true,
							Description: "灰度版本 ID",
						},
						"gray_weight": schema.Int32Attribute{
							Computed:    true,
							Description: "切流的比例。范围是 [0-100]",
						},
						"gray": schema.ListNestedAttribute{
							Computed:    true,
							Description: "灰度版本的配置",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"version_id": schema.StringAttribute{
										Computed:    true,
										Description: "灰度版本 ID",
									},
									"raw_type": schema.Int32Attribute{
										Computed:    true,
										Description: "灰度类型，当前支持：1、按百分比随机灰度",
									},
									"config": schema.ListNestedAttribute{
										Computed:    true,
										Description: "对应类型的配置",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"weight": schema.Int32Attribute{
													Computed:    true,
													Description: "切流的比例。范围是 [0-100]",
												},
											},
										},
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

func (c *CtyunFunctionAliases) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunctionAliases) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CtyunFunctionAliasesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 获取地区 ID
	regionId := c.meta.GetExtraIfEmpty(data.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		resp.Diagnostics.AddError("region_id 不能为空", "")
		return
	}
	data.RegionID = types.StringValue(regionId)

	// TODO: 当有 ListAliases API 时，在这里调用
	// 目前由于没有列表 API，我们只能返回空列表
	// 如果后续需要查询特定别名，可以通过参数传递 alias_name 来查询单个别名

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type FunctionAliasModel struct {
	AliasName     types.String             `tfsdk:"alias_name"`
	VersionID     types.String             `tfsdk:"version_id"`
	Description   types.String             `tfsdk:"description"`
	CreateTime    types.String             `tfsdk:"create_time"`
	UpdateTime    types.String             `tfsdk:"update_time"`
	GrayType      types.Int32              `tfsdk:"gray_type"`
	GrayVersionID types.String             `tfsdk:"gray_version_id"`
	GrayWeight    types.Int32              `tfsdk:"gray_weight"`
	Gray          []FunctionAliasGrayModel `tfsdk:"gray"`
}

type FunctionAliasGrayModel struct {
	VersionID types.String                   `tfsdk:"version_id"`
	RawType   types.Int32                    `tfsdk:"raw_type"`
	Config    []FunctionAliasGrayConfigModel `tfsdk:"config"`
}

type FunctionAliasGrayConfigModel struct {
	Weight types.Int32 `tfsdk:"weight"`
}

type CtyunFunctionAliasesDataSourceModel struct {
	ID           types.String         `tfsdk:"id"`
	FunctionName types.String         `tfsdk:"function_name"`
	RegionID     types.String         `tfsdk:"region_id"`
	Aliases      []FunctionAliasModel `tfsdk:"aliases"`
}
