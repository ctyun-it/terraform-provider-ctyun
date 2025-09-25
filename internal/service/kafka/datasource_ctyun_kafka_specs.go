package kafka

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunKafkaSpecs{}
	_ datasource.DataSourceWithConfigure = &ctyunKafkaSpecs{}
)

type ctyunKafkaSpecs struct {
	meta *common.CtyunMetadata
}

func NewCtyunKafkaSpecs() datasource.DataSource {
	return &ctyunKafkaSpecs{}
}

func (c *ctyunKafkaSpecs) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_specs"
}

type CtyunKafkaSpecsSkuResItem struct {
	ResType  types.String                        `tfsdk:"res_type"`
	ResName  types.String                        `tfsdk:"res_name"`
	ResItems []CtyunKafkaSpecsSkuResItemResItems `tfsdk:"res_items"`
}

type CtyunKafkaSpecsSkuResItemResItems struct {
	CpuArch  types.String                            `tfsdk:"cpu_arch"`
	HostType types.String                            `tfsdk:"host_type"`
	Spec     []CtyunKafkaSpecsSkuResItemResItemsSpec `tfsdk:"spec"`
}

type CtyunKafkaSpecsSkuResItemResItemsSpec struct {
	SpecName     types.String `tfsdk:"spec_name"`
	Description  types.String `tfsdk:"description"`
	Tps          types.Int32  `tfsdk:"tps"`
	MaxPartition types.Int32  `tfsdk:"max_partition"`
	Flow         types.Int32  `tfsdk:"flow"`
	Cpu          types.Int32  `tfsdk:"cpu"`
	Memory       types.Int32  `tfsdk:"memory"`
}

type CtyunKafkaSpecsSkuDiskItem struct {
	ResType  types.String `tfsdk:"res_type"`
	ResName  types.String `tfsdk:"res_name"`
	ResItems []string     `tfsdk:"res_items"`
}

type CtyunKafkaSpecsSku struct {
	ProdId   types.String               `tfsdk:"prod_id"`
	ProdName types.String               `tfsdk:"prod_name"`
	ProdCode types.String               `tfsdk:"prod_code"`
	ResItem  CtyunKafkaSpecsSkuResItem  `tfsdk:"res_item"`
	DiskItem CtyunKafkaSpecsSkuDiskItem `tfsdk:"disk_item"`
}

type CtyunKafkaSpecsModel struct {
	ProdId   types.String         `tfsdk:"prod_id"`
	ProdName types.String         `tfsdk:"prod_name"`
	ProdCode types.String         `tfsdk:"prod_code"`
	Sku      []CtyunKafkaSpecsSku `tfsdk:"sku"`
}

type CtyunKafkaSpecsConfig struct {
	RegionID types.String           `tfsdk:"region_id"`
	Specs    []CtyunKafkaSpecsModel `tfsdk:"specs"`
}

func (c *ctyunKafkaSpecs) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10029624/10030704**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"specs": schema.ListNestedAttribute{
				Description: "产品系列信息",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"prod_id": schema.StringAttribute{
							Description: "产品系列id",
							Computed:    true,
						},
						"prod_name": schema.StringAttribute{
							Description: "产品系列名称",
							Computed:    true,
						},
						"prod_code": schema.StringAttribute{
							Description: "产品系列编码",
							Computed:    true,
						},
						"sku": schema.ListNestedAttribute{
							Description: "产品系列详情",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"prod_id": schema.StringAttribute{
										Description: "产品系列详情id",
										Computed:    true,
									},
									"prod_name": schema.StringAttribute{
										Description: "产品系列详情名称",
										Computed:    true,
									},
									"prod_code": schema.StringAttribute{
										Description: "产品系列详情编码",
										Computed:    true,
									},
									"res_item": schema.SingleNestedAttribute{
										Description: "主机信息",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"res_type": schema.StringAttribute{
												Description: "资源类型：ecs",
												Computed:    true,
											},
											"res_name": schema.StringAttribute{
												Description: "资源名称：云服务器",
												Computed:    true,
											},
											"res_items": schema.ListNestedAttribute{
												Description: "主机规格信息",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"cpu_arch": schema.StringAttribute{
															Description: "cpu架构",
															Computed:    true,
														},
														"host_type": schema.StringAttribute{
															Description: "主机类型",
															Computed:    true,
														},
														"spec": schema.ListNestedAttribute{
															Description: "主机规格列表",
															Computed:    true,
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{
																	"spec_name": schema.StringAttribute{
																		Description: "产品规格名称",
																		Computed:    true,
																	},
																	"description": schema.StringAttribute{
																		Description: "产品规格描述",
																		Computed:    true,
																	},
																	"tps": schema.Int32Attribute{
																		Description: "单个代理TPS",
																		Computed:    true,
																	},
																	"max_partition": schema.Int32Attribute{
																		Description: "单个代理最大分区数",
																		Computed:    true,
																	},
																	"flow": schema.Int32Attribute{
																		Description: "单个代理流量规格",
																		Computed:    true,
																	},
																	"cpu": schema.Int64Attribute{
																		Description: "cpu核心数",
																		Computed:    true,
																	},
																	"memory": schema.Int64Attribute{
																		Description: "内存大小",
																		Computed:    true,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"disk_item": schema.SingleNestedAttribute{
										Description: "磁盘信息",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"res_type": schema.StringAttribute{
												Description: "资源类型：ebs",
												Computed:    true,
											},
											"res_name": schema.StringAttribute{
												Description: "资源名称：磁盘",
												Computed:    true,
											},
											"res_items": schema.ListAttribute{
												Description: "磁盘类型",
												Computed:    true,
												ElementType: types.StringType,
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

func (c *ctyunKafkaSpecs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunKafkaSpecsConfig
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
	params := &ctgkafka.CtgkafkaProdDetailRequest{
		RegionId: config.RegionID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaProdDetailApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	config.Specs = []CtyunKafkaSpecsModel{}
	for _, s := range resp.ReturnObj.Data.Series {
		item := CtyunKafkaSpecsModel{
			ProdId:   types.StringValue(s.ProdId),
			ProdName: types.StringValue(s.ProdName),
			ProdCode: types.StringValue(s.ProdCode),
		}
		for _, sk := range s.Sku {
			skItem := CtyunKafkaSpecsSku{
				ProdId:   types.StringValue(sk.ProdId),
				ProdName: types.StringValue(sk.ProdName),
				ProdCode: types.StringValue(sk.ProdCode),
				ResItem: CtyunKafkaSpecsSkuResItem{
					ResName: types.StringValue(sk.ResItem.ResName),
					ResType: types.StringValue(sk.ResItem.ResType),
				},
				DiskItem: CtyunKafkaSpecsSkuDiskItem{
					ResName:  types.StringValue(sk.DiskItem.ResName),
					ResType:  types.StringValue(sk.DiskItem.ResType),
					ResItems: sk.DiskItem.ResItems,
				},
			}
			for _, r := range sk.ResItem.ResItems {
				rItem := CtyunKafkaSpecsSkuResItemResItems{
					CpuArch:  types.StringValue(r.CpuArch),
					HostType: types.StringValue(r.HostType),
				}
				for _, sp := range r.Spec {
					spItem := CtyunKafkaSpecsSkuResItemResItemsSpec{
						SpecName:     types.StringValue(sp.SpecName),
						Description:  types.StringValue(sp.Description),
						Tps:          types.Int32Value(sp.Tps),
						MaxPartition: types.Int32Value(sp.MaxPartition),
						Flow:         types.Int32Value(sp.Flow),
						Cpu:          types.Int32Value(sp.Cpu),
						Memory:       types.Int32Value(sp.Memory),
					}
					rItem.Spec = append(rItem.Spec, spItem)
				}
				skItem.ResItem.ResItems = append(skItem.ResItem.ResItems, rItem)
			}
			item.Sku = append(item.Sku, skItem)
		}
		config.Specs = append(config.Specs, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunKafkaSpecs) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
