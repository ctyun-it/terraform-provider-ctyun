package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/amqp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunRabbitmqSpecs{}
	_ datasource.DataSourceWithConfigure = &ctyunRabbitmqSpecs{}
)

type ctyunRabbitmqSpecs struct {
	meta *common.CtyunMetadata
}

func NewCtyunRabbitmqSpecs() datasource.DataSource {
	return &ctyunRabbitmqSpecs{}
}

func (c *ctyunRabbitmqSpecs) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_specs"
}

type CtyunRabbitmqSpecsSkuResItem struct {
	ResType  types.String                           `tfsdk:"res_type"`
	ResName  types.String                           `tfsdk:"res_name"`
	ResItems []CtyunRabbitmqSpecsSkuResItemResItems `tfsdk:"res_items"`
}

type CtyunRabbitmqSpecsSkuResItemResItems struct {
	CpuArch  types.String                               `tfsdk:"cpu_arch"`
	HostType types.String                               `tfsdk:"host_type"`
	Spec     []CtyunRabbitmqSpecsSkuResItemResItemsSpec `tfsdk:"spec"`
}

type CtyunRabbitmqSpecsSkuResItemResItemsSpec struct {
	SpecName    types.String `tfsdk:"spec_name"`
	Description types.String `tfsdk:"description"`
	Cpu         types.Int32  `tfsdk:"cpu"`
	Memory      types.Int32  `tfsdk:"memory"`
}

type CtyunRabbitmqSpecsSkuDiskItem struct {
	ResType  types.String `tfsdk:"res_type"`
	ResName  types.String `tfsdk:"res_name"`
	ResItems []string     `tfsdk:"res_items"`
}

type CtyunRabbitmqSpecsSku struct {
	ProdId   types.String                  `tfsdk:"prod_id"`
	ProdName types.String                  `tfsdk:"prod_name"`
	ProdCode types.String                  `tfsdk:"prod_code"`
	ResItem  CtyunRabbitmqSpecsSkuResItem  `tfsdk:"res_item"`
	DiskItem CtyunRabbitmqSpecsSkuDiskItem `tfsdk:"disk_item"`
}

type CtyunRabbitmqSpecsModel struct {
	ProdId   types.String            `tfsdk:"prod_id"`
	ProdName types.String            `tfsdk:"prod_name"`
	ProdCode types.String            `tfsdk:"prod_code"`
	Sku      []CtyunRabbitmqSpecsSku `tfsdk:"sku"`
}

type CtyunRabbitmqSpecsConfig struct {
	RegionID types.String              `tfsdk:"region_id"`
	Specs    []CtyunRabbitmqSpecsModel `tfsdk:"specs"`
}

func (c *ctyunRabbitmqSpecs) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029625/10032819`,
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

func (c *ctyunRabbitmqSpecs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRabbitmqSpecsConfig
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
	params := &amqp.AmqpProdDetailRequest{regionId}
	// 调用API
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpProdDetailApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Specs = []CtyunRabbitmqSpecsModel{}
	// 解析返回值
	for _, s := range resp.ReturnObj.Data.Series {
		item := CtyunRabbitmqSpecsModel{
			ProdId:   types.StringValue(s.ProdId),
			ProdName: types.StringValue(s.ProdName),
			ProdCode: types.StringValue(s.ProdCode),
		}
		for _, sk := range s.Sku {
			skItem := CtyunRabbitmqSpecsSku{
				ProdId:   types.StringValue(sk.ProdId),
				ProdName: types.StringValue(sk.ProdName),
				ProdCode: types.StringValue(sk.ProdCode),
				ResItem: CtyunRabbitmqSpecsSkuResItem{
					ResName: types.StringValue(sk.ResItem.ResName),
					ResType: types.StringValue(sk.ResItem.ResType),
				},
				DiskItem: CtyunRabbitmqSpecsSkuDiskItem{
					ResName:  types.StringValue(sk.DiskItem.ResName),
					ResType:  types.StringValue(sk.DiskItem.ResType),
					ResItems: sk.DiskItem.ResItems,
				},
			}
			for _, r := range sk.ResItem.ResItems {
				rItem := CtyunRabbitmqSpecsSkuResItemResItems{
					CpuArch:  types.StringValue(r.CpuArch),
					HostType: types.StringValue(r.HostType),
				}
				for _, sp := range r.Spec {
					spItem := CtyunRabbitmqSpecsSkuResItemResItemsSpec{
						SpecName:    types.StringValue(sp.SpecName),
						Description: types.StringValue(sp.Description),
						Cpu:         types.Int32Value(sp.Cpu),
						Memory:      types.Int32Value(sp.Memory),
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

func (c *ctyunRabbitmqSpecs) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
