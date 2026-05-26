package rocketmq

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/rocketmq"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunRocketmqSpecs{}
	_ datasource.DataSourceWithConfigure = &ctyunRocketmqSpecs{}
)

type ctyunRocketmqSpecs struct {
	meta *common.CtyunMetadata
}

func NewCtyunRocketmqSpecs() datasource.DataSource {
	return &ctyunRocketmqSpecs{}
}

func (c *ctyunRocketmqSpecs) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rocketmq_specs"
}

type CtyunRocketmqSpecsModel struct {
	FlavorID      types.String   `tfsdk:"flavor_id"`
	SpecName      types.String   `tfsdk:"spec_name"`
	FlavorType    types.String   `tfsdk:"flavor_type"`
	FlavorName    types.String   `tfsdk:"flavor_name"`
	CpuNum        types.Int32    `tfsdk:"cpu_num"`
	MemSize       types.Int32    `tfsdk:"mem_size"`
	MultiQueue    types.Int32    `tfsdk:"multi_queue"`
	Pps           types.Int32    `tfsdk:"pps"`
	BandwidthBase types.Float64  `tfsdk:"bandwidth_base"`
	BandwidthMax  types.Float64  `tfsdk:"bandwidth_max"`
	CpuArch       types.String   `tfsdk:"cpu_arch"`
	Series        types.String   `tfsdk:"series"`
	AzList        []types.String `tfsdk:"az_list"`
	SkuProdId     types.String   `tfsdk:"sku_prod_id"`
}

type CtyunRocketmqSpecsConfig struct {
	RegionID types.String              `tfsdk:"region_id"`
	Specs    []CtyunRocketmqSpecsModel `tfsdk:"specs"`
}

func (c *ctyunRocketmqSpecs) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询 RocketMQ 可用的规格", "分布式消息服务 RocketMQ", ""),
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池 ID",
			},
			"specs": schema.ListNestedAttribute{
				Description: "规格详情列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"flavor_id": schema.StringAttribute{
							Description: "规格 id",
							Computed:    true,
						},
						"spec_name": schema.StringAttribute{
							Description: "规格名称",
							Computed:    true,
						},
						"flavor_type": schema.StringAttribute{
							Description: "规格类型",
							Computed:    true,
						},
						"flavor_name": schema.StringAttribute{
							Description: "规格类型名称",
							Computed:    true,
						},
						"cpu_num": schema.Int32Attribute{
							Description: "cpu 核数",
							Computed:    true,
						},
						"mem_size": schema.Int32Attribute{
							Description: "内存大小，单位 G",
							Computed:    true,
						},
						"multi_queue": schema.Int32Attribute{
							Description: "多队列数",
							Computed:    true,
						},
						"pps": schema.Int32Attribute{
							Description: "网络最大收发包能力 (万 PPS)",
							Computed:    true,
						},
						"bandwidth_base": schema.Float64Attribute{
							Description: "基准带宽 (Gbps)",
							Computed:    true,
						},
						"bandwidth_max": schema.Float64Attribute{
							Description: "最大带宽 (Gbps)",
							Computed:    true,
						},
						"cpu_arch": schema.StringAttribute{
							Description: "cpu 架构（x86、arm）",
							Computed:    true,
						},
						"series": schema.StringAttribute{
							Description: "系列",
							Computed:    true,
						},
						"az_list": schema.ListAttribute{
							Description: "支持的 az 名称列表",
							Computed:    true,
							ElementType: types.StringType,
						},
						"sku_prod_id": schema.StringAttribute{
							Description: "产品 id",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRocketmqSpecs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRocketmqSpecsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}

	config.RegionID = types.StringValue(regionId)
	// 组装请求体
	params := &rocketmq.RocketmqProdDetailRequest{RegionId: regionId}
	// 调用 API
	resp, err := c.meta.Apis.RocketmqApis.RocketmqProdDetailApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	config.Specs = []CtyunRocketmqSpecsModel{}
	// 解析返回值
	for _, flavor := range resp.ReturnObj.Data {
		item := CtyunRocketmqSpecsModel{
			FlavorID:      types.StringValue(flavor.FlavorID),
			SpecName:      types.StringValue(flavor.SpecName),
			FlavorType:    types.StringValue(flavor.FlavorType),
			FlavorName:    types.StringValue(flavor.FlavorName),
			CpuNum:        types.Int32Value(flavor.CpuNum),
			MemSize:       types.Int32Value(flavor.MemSize),
			MultiQueue:    types.Int32Value(flavor.MultiQueue),
			Pps:           types.Int32Value(flavor.Pps),
			BandwidthBase: types.Float64Value(float64(flavor.BandwidthBase)),
			BandwidthMax:  types.Float64Value(float64(flavor.BandwidthMax)),
			CpuArch:       types.StringValue(flavor.CpuArch),
			Series:        types.StringValue(flavor.Series),
			AzList:        make([]types.String, len(flavor.AzList)),
			SkuProdId:     types.StringValue(flavor.SkuProdId),
		}
		for i, az := range flavor.AzList {
			item.AzList[i] = types.StringValue(az)
		}
		config.Specs = append(config.Specs, item)
	}
	// 保存到 state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRocketmqSpecs) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
