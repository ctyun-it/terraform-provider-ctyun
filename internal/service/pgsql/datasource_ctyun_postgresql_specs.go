package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ datasource.DataSource              = &CtyunPgsqlSpecs{}
	_ datasource.DataSourceWithConfigure = &CtyunPgsqlSpecs{}
)

type CtyunPgsqlSpecs struct {
	meta *common.CtyunMetadata
}

func NewCtyunPgsqlSpecs() *CtyunPgsqlSpecs {
	return &CtyunPgsqlSpecs{}
}
func (c *CtyunPgsqlSpecs) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPgsqlSpecs) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_specs"
}

func (c *CtyunPgsqlSpecs) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10034019/10167295**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"instance_series": schema.StringAttribute{
				Required:    true,
				Description: "实例规格，取值范围:S(通用型)， C(计算增强型)，M(内存增强型)",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目id",
			},
			"specs": schema.ListNestedAttribute{
				Computed:    true,
				Description: "产品规格列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"prod_id": schema.Int64Attribute{
							Computed:    true,
							Description: "产品id",
						},
						"prod_code": schema.StringAttribute{
							Computed:    true,
							Description: "产品编码",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ProdCode...),
							},
						},
						"prod_spec_name": schema.StringAttribute{
							Computed:    true,
							Description: "产品名称",
						},
						"prod_spec_desc": schema.StringAttribute{
							Computed:    true,
							Description: "产品描述",
						},
						"instance_desc": schema.StringAttribute{
							Computed:    true,
							Description: "实例描述",
						},
						"prod_version": schema.StringAttribute{
							Computed:    true,
							Description: "产品版本",
						},
						"host_spec": schema.StringAttribute{
							Computed:    true,
							Description: "主机规格",
						},
						"lvs_spec": schema.StringAttribute{
							Computed:    true,
							Description: "lvs规格",
						},
						"inst_spec_info_list": schema.ListNestedAttribute{
							Computed:    true,
							Description: "AZ支持的产品规格信息，以及规格代S6/S7",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"spec_id": schema.StringAttribute{
										Computed:    true,
										Description: "spec id",
									},
									"prod_performance_spec": schema.StringAttribute{
										Computed:    true,
										Description: "规格名称",
									},
									"az_list": schema.StringAttribute{
										Computed:    true,
										Description: "该规格支持的AZ列表",
									},
									"spec_name": schema.StringAttribute{
										Computed:    true,
										Description: "主机世代完整名称",
									},
									"cpu_type": schema.StringAttribute{
										Computed:    true,
										Description: "cpu类型",
									},
									"generation": schema.StringAttribute{
										Computed:    true,
										Description: "主机世代缩写",
									},
									"min_rate": schema.StringAttribute{
										Computed:    true,
										Description: "带宽下限",
									},
									"max_rate": schema.StringAttribute{
										Computed:    true,
										Description: "带宽上限",
									},
								},
							},
						},
						"prod_host_config": schema.ListNestedAttribute{
							Computed:    true,
							Description: "host实例",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"host_type_name": schema.StringAttribute{
										Computed:    true,
										Description: "类型名称",
									},
									"host_type": schema.StringAttribute{
										Computed:    true,
										Description: "节点类型，当类型为mysql时，取值范围master（主节点）、readnode（只读节点）",
									},
									"prod_performance_specs": schema.StringAttribute{
										Computed:    true,
										Description: "支持的性能指标规格列表",
									},
									"host_default_num": schema.Int32Attribute{
										Computed:    true,
										Description: "节点默认数量",
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

func (c *CtyunPgsqlSpecs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPgsqlSpecsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "1",
		ProdCode:     "POSTGRESQL",
		RegionID:     regionId,
		InstanceType: business.PgsqlInstanceSeriesDict[config.InstanceSeries.ValueString()],
	}
	headers := &mysql.TeledbMysqlSpecsRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		headers.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbMysqlSpecsApi.Do(ctx, c.meta.Credential, params, headers)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析产品规格
	returnObjData := resp.ReturnObj.Data
	var specs []CtyunPgsqlSpecInfoModel
	for _, specItem := range returnObjData {
		var spec CtyunPgsqlSpecInfoModel
		spec.ProdId = types.Int64Value(specItem.ProdId)
		spec.ProdCode = types.StringValue(specItem.ProdCode)
		spec.ProdSpecName = types.StringValue(specItem.ProdSpecName)
		spec.ProdSpecDesc = types.StringValue(specItem.ProdSpecDesc)
		spec.InstanceDesc = types.StringValue(specItem.InstanceDesc)
		spec.ProdVersion = types.StringValue(specItem.ProdVersion)
		spec.HostSpec = types.StringValue(specItem.HostSpec)
		spec.LvsSpec = types.StringValue(specItem.LvsSpec)
		//InstSpecInfoList解析
		var instSpecInfoList []InstSpecInfo
		for _, instSpecInfoItem := range specItem.InstSpecInfoList {
			var specInfo InstSpecInfo
			specInfo.SpecId = types.StringValue(instSpecInfoItem.SpecId)
			specInfo.ProdPerformanceSpec = types.StringValue(instSpecInfoItem.ProdPerformanceSpec)
			specInfo.SpecName = types.StringValue(instSpecInfoItem.SpecName)
			specInfo.CpuType = types.StringValue(instSpecInfoItem.CpuType)
			specInfo.Generation = types.StringValue(instSpecInfoItem.Generation)
			specInfo.MinRate = types.StringValue(instSpecInfoItem.MinRate)
			specInfo.MaxRate = types.StringValue(instSpecInfoItem.MaxRate)
			// 解析azList
			var azList []string
			for _, az := range instSpecInfoItem.AzList {
				azList = append(azList, az)
			}
			specInfo.AzList = types.StringValue(strings.Join(azList, ","))
			instSpecInfoList = append(instSpecInfoList, specInfo)
		}
		spec.InstSpecInfoList = instSpecInfoList
		// ProdHostConfig解析
		var prodHostConfigs []ProdHostConfig
		for _, prodHostConfigItem := range specItem.ProdHostConfig.HostInsts {
			var prodHostConfig ProdHostConfig
			prodHostConfig.HostTypeName = types.StringValue(prodHostConfigItem.HostTypeName)
			prodHostConfig.HostType = types.StringValue(prodHostConfigItem.HostType)
			prodHostConfig.HostDefaultNum = types.Int32Value(prodHostConfigItem.HostDefaultNum)
			var prodPerformanceSpeces []string
			for _, performance := range prodHostConfigItem.ProdPerformanceSpeces {
				prodPerformanceSpeces = append(prodPerformanceSpeces, performance)
			}
			prodHostConfig.ProdPerformanceSpecs = types.StringValue(strings.Join(prodPerformanceSpeces, ","))
			prodHostConfigs = append(prodHostConfigs, prodHostConfig)
		}
		spec.ProdHostConfig = prodHostConfigs
		specs = append(specs, spec)
	}
	config.Specs = specs
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunPgsqlSpecsConfig struct {
	RegionID       types.String              `tfsdk:"region_id"`
	InstanceSeries types.String              `tfsdk:"instance_series"`
	ProjectID      types.String              `tfsdk:"project_id"`
	Specs          []CtyunPgsqlSpecInfoModel `tfsdk:"specs"`
}
type CtyunPgsqlSpecInfoModel struct {
	ProdId           types.Int64      `tfsdk:"prod_id"`             // 产品id
	ProdCode         types.String     `tfsdk:"prod_code"`           // 产品编码
	ProdSpecName     types.String     `tfsdk:"prod_spec_name"`      // 产品名称
	ProdSpecDesc     types.String     `tfsdk:"prod_spec_desc"`      // 产品描述
	InstanceDesc     types.String     `tfsdk:"instance_desc"`       // 实例描述
	ProdVersion      types.String     `tfsdk:"prod_version"`        // 产品版本
	HostSpec         types.String     `tfsdk:"host_spec"`           // 主机规格
	LvsSpec          types.String     `tfsdk:"lvs_spec"`            // lvs规格
	InstSpecInfoList []InstSpecInfo   `tfsdk:"inst_spec_info_list"` // AZ支持的产品规格信息，以及规格代S6/S7
	ProdHostConfig   []ProdHostConfig `tfsdk:"prod_host_config"`    //host实例
}

type InstSpecInfo struct {
	SpecId              types.String `tfsdk:"spec_id"`               // 废弃
	ProdPerformanceSpec types.String `tfsdk:"prod_performance_spec"` // 规格名称
	AzList              types.String `tfsdk:"az_list"`               // 该规格支持的AZ列表
	SpecName            types.String `tfsdk:"spec_name"`             // 主机世代完整名称
	CpuType             types.String `tfsdk:"cpu_type"`              // cpu类型
	Generation          types.String `tfsdk:"generation"`            // 主机世代缩写
	MinRate             types.String `tfsdk:"min_rate"`              // 带宽下限
	MaxRate             types.String `tfsdk:"max_rate"`              // 带宽上限
}
type ProdHostConfig struct {
	HostTypeName         types.String `tfsdk:"host_type_name"`         // 类型名称
	HostType             types.String `tfsdk:"host_type"`              // 节点类型，当类型为mysql时，取值范围master（主节点）、readnode（只读节点）
	ProdPerformanceSpecs types.String `tfsdk:"prod_performance_specs"` // 支持的性能指标规格列表
	HostDefaultNum       types.Int32  `tfsdk:"host_default_num"`       // 节点默认数量
}
