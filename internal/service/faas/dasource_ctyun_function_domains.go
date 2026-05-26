package faas

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/cf"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunFunctionDomains{}
	_ datasource.DataSourceWithConfigure = &CtyunFunctionDomains{}
)

type CtyunFunctionDomains struct {
	meta *common.CtyunMetadata
}

func NewCtyunFunctionDomains() datasource.DataSource {
	return &CtyunFunctionDomains{}
}

func (c *CtyunFunctionDomains) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunctionDomains) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_function_domains"
}

func (c *CtyunFunctionDomains) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询自定义域名列表", "函数工作流（FunctionGraph）", "https://www.ctyun.cn/document/10355289"),
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID，默认使用 provider 配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"page_index": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，默认为 -1",
				Validators: []validator.Int32{
					int32validator.AtLeast(-1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页大小，默认为 -1",
				Validators: []validator.Int32{
					int32validator.AtLeast(-1),
				},
			},
			"search_key": schema.StringAttribute{
				Optional:    true,
				Description: "模糊查询的关键字，默认为空",
			},
			"domains": schema.ListNestedAttribute{
				Computed:    true,
				Description: "自定义域名列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain_name": schema.StringAttribute{
							Computed:    true,
							Description: "自定义域名",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "协议类型",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"domain_status": schema.StringAttribute{
							Computed:    true,
							Description: "域名备案状态",
						},
						"cname_valid": schema.BoolAttribute{
							Computed:    true,
							Description: "CNAME 是否有效",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"cert_name": schema.StringAttribute{
							Computed:    true,
							Description: "证书名称",
						},
						"routes": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"function_name": schema.StringAttribute{
										Computed:    true,
										Description: "函数名称",
									},
									"path": schema.StringAttribute{
										Computed:    true,
										Description: "请求路径",
									},
									"qualifier": schema.StringAttribute{
										Computed:    true,
										Description: "函数版本",
									},
									"enable_jwt": schema.Int32Attribute{
										Computed:    true,
										Description: "是否启用 JWT",
									},
									"methods": schema.ListAttribute{
										ElementType: types.StringType,
										Computed:    true,
										Description: "请求方法列表",
									},
									"function_id": schema.Int32Attribute{
										Computed:    true,
										Description: "函数 ID",
									},
									"function_unique_name": schema.StringAttribute{
										Computed:    true,
										Description: "函数唯一名称",
									},
								},
							},
							Description: "路由映射列表",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunFunctionDomains) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunFunctionDomainsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("region_id 不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)

	params := &cf.CfListCustomDomainsRequest{
		RegionId: config.RegionID.ValueString(),
	}

	if !config.PageIndex.IsNull() {
		pageIndex := config.PageIndex.ValueInt32()
		params.PageIndex = &pageIndex
	}

	if !config.PageSize.IsNull() {
		pageSize := config.PageSize.ValueInt32()
		params.PageSize = &pageSize
	}

	if !config.SearchKey.IsNull() {
		searchKey := config.SearchKey.ValueString()
		params.SearchKey = &searchKey
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfListCustomDomainsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询自定义域名列表失败！接口返回 nil，请联系研发确认问题原因！")
		return
	} else if *resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	var domains []CtyunFunctionDomainInfoModel
	for _, domainItem := range resp.ReturnObj.Data {
		var domain CtyunFunctionDomainInfoModel
		domain.DomainName = types.StringPointerValue(domainItem.DomainName)
		domain.Protocol = types.StringPointerValue(domainItem.Protocol)
		domain.Description = types.StringPointerValue(domainItem.Description)
		domain.DomainStatus = types.StringPointerValue(domainItem.DomainStatus)
		domain.CnameValid = types.BoolPointerValue(domainItem.CnameValid)
		domain.CreatedAt = types.StringPointerValue(domainItem.CreatedAt)
		domain.UpdatedAt = types.StringPointerValue(domainItem.UpdatedAt)

		if domainItem.CertConfig != nil {
			domain.CertName = types.StringPointerValue(domainItem.CertConfig.CertName)
		}

		// 处理 routes
		if domainItem.RouteConfig != nil && len(domainItem.RouteConfig.Routes) > 0 {
			routes := make([]CtyunFunctionDomainRouteInfoModel, 0, len(domainItem.RouteConfig.Routes))
			for _, routeItem := range domainItem.RouteConfig.Routes {
				methods := make([]types.String, 0)
				if routeItem.Methods != nil {
					for _, method := range routeItem.Methods {
						methods = append(methods, types.StringPointerValue(method))
					}
				}
				route := CtyunFunctionDomainRouteInfoModel{
					FunctionName:       types.StringPointerValue(routeItem.FunctionName),
					Path:               types.StringPointerValue(routeItem.Path),
					Qualifier:          types.StringPointerValue(routeItem.Qualifier),
					EnableJwt:          types.Int32PointerValue(routeItem.EnableJwt),
					Methods:            methods,
					FunctionID:         types.Int32PointerValue(routeItem.FunctionId),
					FunctionUniqueName: types.StringPointerValue(routeItem.FunctionUniqueName),
				}
				routes = append(routes, route)
			}
			domain.Routes = routes
		}

		domains = append(domains, domain)
	}

	config.Domains = domains
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunFunctionDomainRouteInfoModel struct {
	FunctionName       types.String   `tfsdk:"function_name"`
	Path               types.String   `tfsdk:"path"`
	Qualifier          types.String   `tfsdk:"qualifier"`
	EnableJwt          types.Int32    `tfsdk:"enable_jwt"`
	Methods            []types.String `tfsdk:"methods"`
	FunctionID         types.Int32    `tfsdk:"function_id"`
	FunctionUniqueName types.String   `tfsdk:"function_unique_name"`
}

type CtyunFunctionDomainInfoModel struct {
	DomainName   types.String                        `tfsdk:"domain_name"`
	Protocol     types.String                        `tfsdk:"protocol"`
	Description  types.String                        `tfsdk:"description"`
	DomainStatus types.String                        `tfsdk:"domain_status"`
	CnameValid   types.Bool                          `tfsdk:"cname_valid"`
	CreatedAt    types.String                        `tfsdk:"created_at"`
	UpdatedAt    types.String                        `tfsdk:"updated_at"`
	CertName     types.String                        `tfsdk:"cert_name"`
	Routes       []CtyunFunctionDomainRouteInfoModel `tfsdk:"routes"`
}

type CtyunFunctionDomainsConfig struct {
	RegionID  types.String                   `tfsdk:"region_id"`
	PageIndex types.Int32                    `tfsdk:"page_index"`
	PageSize  types.Int32                    `tfsdk:"page_size"`
	SearchKey types.String                   `tfsdk:"search_key"`
	Domains   []CtyunFunctionDomainInfoModel `tfsdk:"domains"`
}
