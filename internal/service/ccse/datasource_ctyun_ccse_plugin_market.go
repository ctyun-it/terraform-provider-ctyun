package ccse

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/crs"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunCcsePluginMarket{}
	_ datasource.DataSourceWithConfigure = &ctyunCcsePluginMarket{}
)

type ctyunCcsePluginMarket struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcsePluginMarket() datasource.DataSource {
	return &ctyunCcsePluginMarket{}
}

func (c *ctyunCcsePluginMarket) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_plugin_market"
}

type CtyunCcsePluginMarketModel struct {
	Namespace        types.String `tfsdk:"namespace"`
	ChartName        types.String `tfsdk:"chart_name"`
	ImageUrl         types.String `tfsdk:"image_url"`
	ImageUrlInternal types.String `tfsdk:"image_url_internal"`
}

type CtyunCcsePluginVersion struct {
	ChartVersion types.String `tfsdk:"chart_version"`
	Size         types.String `tfsdk:"size"`
	Description  types.String `tfsdk:"description"`
}

type CtyunCcsePluginMarketConfig struct {
	RegionID     types.String                 `tfsdk:"region_id"`
	Total        types.Int32                  `tfsdk:"total"`
	Size         types.Int32                  `tfsdk:"size"`
	Current      types.Int32                  `tfsdk:"current"`
	PageNo       types.Int32                  `tfsdk:"page_no"`
	PageSize     types.Int32                  `tfsdk:"page_size"`
	ChartName    types.String                 `tfsdk:"chart_name"`
	ChartVersion types.String                 `tfsdk:"chart_version"`
	ValuesType   types.String                 `tfsdk:"values_type"`
	Values       types.String                 `tfsdk:"values"`
	Records      []CtyunCcsePluginMarketModel `tfsdk:"records"`
	Versions     []CtyunCcsePluginVersion     `tfsdk:"versions"`
}

func (c *ctyunCcsePluginMarket) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10083472/10656137`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页数据量大小，支持范围1-10",
				Validators: []validator.Int32{
					int32validator.Between(1, 10),
				},
			},
			"total": schema.Int32Attribute{
				Computed:    true,
				Description: "插件总数",
			},
			"size": schema.Int32Attribute{
				Computed:    true,
				Description: "每页条数",
			},
			"current": schema.Int32Attribute{
				Computed:    true,
				Description: "当前页码",
			},
			"chart_name": schema.StringAttribute{
				Optional:    true,
				Description: "插件名称",
			},
			"chart_version": schema.StringAttribute{
				Optional:    true,
				Description: "插件版本号，必须传递chart_name时才有效",
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("chart_name")),
				},
			},
			"values_type": schema.StringAttribute{
				Optional:    true,
				Description: "values类型，支持YAML或JSON",
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("chart_name")),
					stringvalidator.AlsoRequires(path.MatchRoot("chart_version")),
					stringvalidator.OneOf("YAML", "JSON"),
				},
			},
			"values": schema.StringAttribute{
				Computed:    true,
				Description: "values，需要同时填写chart_name、chart_version、values_type时才可以获取",
			},
			"records": schema.ListNestedAttribute{
				Computed:    true,
				Description: "插件列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"namespace": schema.StringAttribute{
							Computed:    true,
							Description: "命名空间",
						},
						"chart_name": schema.StringAttribute{
							Computed:    true,
							Description: "插件名称",
						},
						"image_url": schema.StringAttribute{
							Computed:    true,
							Description: "公网地址",
						},
						"image_url_internal": schema.StringAttribute{
							Computed:    true,
							Description: "内网地址",
						},
					},
				},
			},
			"versions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "插件版本列表，当指定chart_name时才可以获取",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"chart_version": schema.StringAttribute{
							Computed:    true,
							Description: "插件版本号",
						},
						"size": schema.StringAttribute{
							Computed:    true,
							Description: "插件大小",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "插件版本描述",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunCcsePluginMarket) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunCcsePluginMarketConfig
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
	err = c.getRecordsAndMerge(ctx, &config)
	if err != nil {
		return
	}
	err = c.getVersionsAndMerge(ctx, &config)
	if err != nil {
		return
	}
	err = c.getValuesAndMerge(ctx, &config)
	if err != nil {
		return
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunCcsePluginMarket) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunCcsePluginMarket) getRecordsAndMerge(ctx context.Context, config *CtyunCcsePluginMarketConfig) (err error) {
	// 组装请求体
	params := &crs.CrsListPluginRequest{
		RegionIdHeader: config.RegionID.ValueString(),
		RegionIdParam:  config.RegionID.ValueString(),
		PageNum:        config.PageNo.ValueInt32(),
		PageSize:       config.PageSize.ValueInt32(),
		RepositoryName: config.ChartName.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCrsApis.CrsListPluginApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 解析返回值
	config.Records = []CtyunCcsePluginMarketModel{}
	config.Total = types.Int32Value(resp.ReturnObj.Total)
	config.Current = types.Int32Value(resp.ReturnObj.Current)
	config.Size = types.Int32Value(resp.ReturnObj.Size)
	for _, r := range resp.ReturnObj.Records {
		item := CtyunCcsePluginMarketModel{
			Namespace:        types.StringValue(r.NamespaceName),
			ChartName:        types.StringValue(r.RepositoryName),
			ImageUrl:         types.StringValue(r.ImageUrl),
			ImageUrlInternal: types.StringValue(r.ImageUrlInternal),
		}
		config.Records = append(config.Records, item)
	}
	return
}

func (c *ctyunCcsePluginMarket) getVersionsAndMerge(ctx context.Context, config *CtyunCcsePluginMarketConfig) (err error) {
	if config.ChartName.ValueString() == "" {
		return
	}
	params := &crs.CrsListTagRequest{
		RegionIdHeader: config.RegionID.ValueString(),
		RegionIdParam:  config.RegionID.ValueString(),
		RepositoryName: config.ChartName.ValueString(),
		NamespaceName:  "plugins",
		TagName:        config.ChartVersion.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCrsApis.CrsListTagApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 解析返回值
	config.Versions = []CtyunCcsePluginVersion{}
	for _, r := range resp.ReturnObj.Records {
		item := CtyunCcsePluginVersion{
			ChartVersion: types.StringValue(r.Name),
			Size:         types.StringValue(r.Size),
			Description:  types.StringValue(r.Description),
		}
		config.Versions = append(config.Versions, item)
	}
	return

}

func (c *ctyunCcsePluginMarket) getValuesAndMerge(ctx context.Context, config *CtyunCcsePluginMarketConfig) (err error) {
	if config.ValuesType.IsNull() {
		return
	}
	params := &crs.CrsGetValuesRequest{
		RegionIdHeader: config.RegionID.ValueString(),
		RegionIdParam:  config.RegionID.ValueString(),
		RepositoryName: config.ChartName.ValueString(),
		NamespaceName:  "plugins",
		TagName:        config.ChartVersion.ValueString(),
		RawType:        config.ValuesType.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCrsApis.CrsGetValuesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == "" {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 解析返回值
	config.Values = types.StringValue(resp.ReturnObj)
	return
}
