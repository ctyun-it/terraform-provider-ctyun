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
	_ datasource.DataSource              = &ctyunRedisSpecs{}
	_ datasource.DataSourceWithConfigure = &ctyunRedisSpecs{}
)

type ctyunRedisSpecs struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisSpecs() datasource.DataSource {
	return &ctyunRedisSpecs{}
}

func (c *ctyunRedisSpecs) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_specs"
}

type CtyunRedisSpecsSeriesInfo struct {
	Version           types.String                           `tfsdk:"version"`              /*  版本类型<li>BASIC：基础版<li>PLUS：增强版<li>Classic：经典版<li>Capacity：容量型  */
	SeriesCode        types.String                           `tfsdk:"series_code"`          /*  产品系列编码  */
	SeriesName        types.String                           `tfsdk:"series_name"`          /*  产品系列名称  */
	SeriesId          types.Int64                            `tfsdk:"series_id"`            /*  产品系列ID  */
	EngineTypeItems   []string                               `tfsdk:"engine_type_items"`    /*  引擎版本  */
	MemSizeItems      []string                               `tfsdk:"mem_size_items"`       /*  内存容量可选值(GB)<br>说明：version为Classic和Capacity有值  */
	ShardMemSizeItems []string                               `tfsdk:"shard_mem_size_items"` /*  单分片内存可选值(GB)。<br>说明：version为BASIC和PLUS有值  */
	ResItems          []CtyunRedisSpecsSeriesInfoListResItem `tfsdk:"res_items"`
}

type CtyunRedisSpecsSeriesInfoListResItem struct {
	ResType types.String `tfsdk:"res_type"` /*  资源类型<li>ecs：云服务器<li>ebs：磁盘  */
	ResName types.String `tfsdk:"res_name"` /*  资源名称  */
	Items   []string     `tfsdk:"items"`    /*  资源类型可选值<br>说明：以实际返回为准<br><br>云服务器<li>S7：通用型<li>C7：计算型<li>M7：内存型<li>HS1：海光通用型<li>HC1：海光计算增强型<li>KS1：鲲鹏通用型<li>KC1：鲲鹏计算增强型  <br><br>磁盘<li>SATA：普通IO<li>SAS：高IO<li>SSD：超高IO<li>FAST-SSD：极速型SSD  */
}

type CtyunRedisSpecsMirror struct {
	AttrVal  types.String `tfsdk:"attr_val"`  /*  操作系统  */
	AttrName types.String `tfsdk:"attr_name"` /*  操作系统名称  */
	Status   types.Int32  `tfsdk:"status"`    /*  状态,1：正常，其他表示异常  */
}

type CtyunRedisSpecsConfig struct {
	RegionID types.String                `tfsdk:"region_id"`
	Mirrors  []CtyunRedisSpecsMirror     `tfsdk:"mirrors"`
	Series   []CtyunRedisSpecsSeriesInfo `tfsdk:"series_infos"`
}

func (c *ctyunRedisSpecs) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/11030280`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"series_infos": schema.ListNestedAttribute{
				Computed:    true,
				Description: "系列信息",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"res_items": schema.ListNestedAttribute{
							Computed:    true,
							Description: "资源类型信息",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"res_type": schema.StringAttribute{
										Computed:    true,
										Description: "资源类型<li>ecs：云服务器<li>ebs：磁盘",
									},
									"res_name": schema.StringAttribute{
										Computed:    true,
										Description: "资源名称",
									},
									"items": schema.ListAttribute{
										ElementType: types.StringType,
										Computed:    true,
										Description: "资源类型可选值<br>说明：以实际返回为准<br><br>云服务器<li>S7：通用型<li>C7：计算型<li>M7：内存型<li>HS1：海光通用型<li>HC1：海光计算增强型<li>KS1：鲲鹏通用型<li>KC1：鲲鹏计算增强型  <br><br>磁盘<li>SATA：普通IO<li>SAS：高IO<li>SSD：超高IO<li>FAST-SSD：极速型SSD",
									},
								},
							},
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "版本类型<li>BASIC：基础版<li>PLUS：增强版<li>Classic：经典版<li>Capacity：容量型",
						},
						"series_code": schema.StringAttribute{
							Computed:    true,
							Description: "产品系列编码",
						},
						"series_name": schema.StringAttribute{
							Computed:    true,
							Description: "状series_name",
						},
						"series_id": schema.Int64Attribute{
							Computed:    true,
							Description: "版本类型<li>BASIC：基础版<li>PLUS：增强版<li>Classic：经典版<li>Capacity：容量型",
						},
						"engine_type_items": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "引擎版本",
						},
						"mem_size_items": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "内存容量可选值(GB)<br>说明：version为Classic和Capacity有值",
						},
						"shard_mem_size_items": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "单分片内存可选值(GB)。<br>说明：version为BASIC和PLUS有值",
						},
					},
				},
			},
			"mirrors": schema.ListNestedAttribute{
				Computed:    true,
				Description: "镜像信息",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attr_val": schema.StringAttribute{
							Computed:    true,
							Description: "操作系统",
						},
						"attr_name": schema.StringAttribute{
							Computed:    true,
							Description: "操作系统名称",
						},
						"status": schema.Int32Attribute{
							Computed:    true,
							Description: "状态，1为正常",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRedisSpecs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRedisSpecsConfig
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
	// 组装请求体
	params := &dcs2.Dcs2DescribeAvailableResourceRequest{
		RegionId: regionId,
	}
	// 调用API
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeAvailableResourceApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.Mirrors = []CtyunRedisSpecsMirror{}
	for _, m := range resp.ReturnObj.MirrorArray {
		item := CtyunRedisSpecsMirror{
			AttrVal:  types.StringValue(m.AttrVal),
			AttrName: types.StringValue(m.AttrName),
			Status:   types.Int32Value(m.Status),
		}

		config.Mirrors = append(config.Mirrors, item)
	}
	config.Series = []CtyunRedisSpecsSeriesInfo{}
	for _, s := range resp.ReturnObj.SeriesInfoList {
		item := CtyunRedisSpecsSeriesInfo{
			Version:           types.StringValue(s.Version),
			SeriesCode:        types.StringValue(s.SeriesCode),
			SeriesName:        types.StringValue(s.SeriesName),
			SeriesId:          types.Int64Value(s.SeriesId),
			EngineTypeItems:   s.EngineTypeItems,
			MemSizeItems:      s.MemSizeItems,
			ShardMemSizeItems: s.ShardMemSizeItems,
		}
		for _, r := range s.ResItems {
			res := CtyunRedisSpecsSeriesInfoListResItem{
				ResType: types.StringValue(r.ResType),
				ResName: types.StringValue(r.ResName),
				Items:   r.Items,
			}
			item.ResItems = append(item.ResItems, res)
		}
		config.Series = append(config.Series, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRedisSpecs) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
