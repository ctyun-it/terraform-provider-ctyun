package ecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

func NewCtyunEcsFlavors() datasource.DataSource {
	return &ctyunEcsFlavors{}
}

type ctyunEcsFlavors struct {
	meta *common.CtyunMetadata
}

func (c *ctyunEcsFlavors) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ecs_flavors"
}

func (c *ctyunEcsFlavors) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10118193**`,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "规格类型，取值范围：CPU、CPU_C3、CPU_C6、CPU_C7、CPU_c7ne、CPU_C8、CPU_D3、CPU_FC1、CPU_FM1、CPU_FS1、CPU_HC1、CPU_HM1、CPU_HS1、CPU_IP3、CPU_IR3、CPU_IP3_2、CPU_IR3_2、CPU_KC1、CPU_KM1、CPU_KS1、CPU_M2、CPU_M3、CPU_M6、CPU_M7、CPU_M8、CPU_S2、CPU_S3、CPU_S6、CPU_S7、CPU_S8、CPU_s8r、GPU_N_V100_V_FMGQ、GPU_N_V100_V、GPU_N_V100S_V、GPU_N_V100S_V_FMGQ、GPU_N_T4_V、GPU_N_G7_V、GPU_N_V100、GPU_N_V100_SHIPINYUN、GPU_N_V100_SUANFA、GPU_N_P2V_RENMIN、GPU_N_V100S、GPU_N_T4、GPU_N_T4_AIJISUAN、GPU_N_T4_ASR、GPU_N_T4_JX、GPU_N_T4_SHIPINYUN、GPU_N_T4_SUANFA、GPU_N_T4_YUNYOUXI、GPU_N_PI7、GPU_N_P8A、GPU_A_PAK1、GPU_C_PCH1，支持类型会随着功能升级增加",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EcsFlavorTypes...),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "规格名称，例如：pi7.4xlarge.4",
			},
			"cpu": schema.Int64Attribute{
				Optional:    true,
				Description: "VCPU个数",
			},
			"ram": schema.Int64Attribute{
				Optional:    true,
				Description: "内存大小，单位为GB",
			},
			"arch": schema.StringAttribute{
				Optional:    true,
				Description: "指令集架构",
			},
			"series": schema.StringAttribute{
				Optional:    true,
				Description: "云主机规格系列，规格系列说明：S（通用型），C（计算增强型），M（内存优化型），HS（海光通用型），HC（海光计算增强型），HM（海光内存优化型），FS（飞腾通用型），FC（飞腾计算增强型），FM（飞腾内存优化型），KS（鲲鹏通用型），KC（鲲鹏计算增强型），P（GPU计算加速型），G（GPU图像加速基础型），IP3（超高IO型）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EcsFlavorSeries...),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区id，如果不填则默认使用provider ctyun中的az_name或环境变量中的CTYUN_AZ_NAME",
			},
			"flavors": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机规格ID",
						},
						"series_name": schema.StringAttribute{
							Computed:    true,
							Description: "规格系列名称，参照参数flavorSeries说明",
						},
						"cpu_info": schema.StringAttribute{
							Computed:    true,
							Description: "cpu架构",
						},
						"base_bandwidth": schema.Float64Attribute{
							Computed:    true,
							Description: "基准带宽",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机规格名称，类似pi7.4xlarge.4入参，精确检索",
						},
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "规格类型",
						},
						"series": schema.StringAttribute{
							Optional:    true,
							Description: "规格系列",
						},
						"nic_multi_queue": schema.Int64Attribute{
							Optional:    true,
							Description: "网卡多队列数目",
						},
						"pps": schema.Int64Attribute{
							Optional:    true,
							Description: "最大收发包限制",
						},
						"cpu": schema.Int64Attribute{
							Optional:    true,
							Description: "VCPU个数",
						},
						"ram": schema.Int64Attribute{
							Optional:    true,
							Description: "内存大小，单位为GB",
						},
						"bandwidth": schema.Float64Attribute{
							Optional:    true,
							Description: "带宽",
						},
						"gpu_vendor": schema.StringAttribute{
							Optional:    true,
							Description: "GPU厂商",
						},
						"video_mem_size": schema.Int64Attribute{
							Optional:    true,
							Description: "GPU显存大小",
						},
						"gpu_type": schema.StringAttribute{
							Optional:    true,
							Description: "GPU类型，取值范围：T4、V100、V100S、A10、A100、atlas 300i pro、mlu370-s4，支持类型会随着功能升级增加",
						},
						"gpu_count": schema.Int64Attribute{
							Optional:    true,
							Description: "GPU设备数量",
						},
						"available": schema.BoolAttribute{
							Optional:    true,
							Description: "是否可用（true：可用；false：不可用，已售罄）",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEcsFlavors) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CtyunEcsFlavorsConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionId.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		msg := "regionId不能为空"
		resp.Diagnostics.AddError(msg, msg)
		return
	}
	azName := c.meta.GetExtraIfEmpty(config.AzName.ValueString(), common.ExtraAzName)

	ecsFlavorListResponse, err := c.meta.Apis.CtEcsApis.EcsFlavorListApi.Do(ctx, c.meta.Credential, &ctecs.EcsFlavorListRequest{
		RegionId:     regionId,
		AzName:       azName,
		FlavorType:   config.Type.ValueString(),
		FlavorName:   config.Name.ValueString(),
		FlavorCpu:    int(config.Cpu.ValueInt64()),
		FlavorRam:    int(config.Ram.ValueInt64()),
		FlavorArch:   config.Arch.ValueString(),
		FlavorSeries: strings.ToLower(config.Series.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	var flavors []CtyunEcsFlavorsFlavorsConfig
	for _, f := range ecsFlavorListResponse.FlavorList {
		flavors = append(flavors, CtyunEcsFlavorsFlavorsConfig{
			Id:               types.StringValue(f.FlavorId),
			FlavorSeriesName: types.StringValue(f.FlavorSeriesName),
			CpuInfo:          types.StringValue(f.CpuInfo),
			BaseBandwidth:    types.Float64Value(f.BaseBandwidth),
			Name:             types.StringValue(f.FlavorName),
			Type:             types.StringValue(f.FlavorType),
			Series:           types.StringValue(strings.ToUpper(f.FlavorSeries)),
			NicMultiQueue:    types.Int64Value(int64(f.NicMultiQueue)),
			Pps:              types.Int64Value(int64(f.Pps)),
			Cpu:              types.Int64Value(int64(f.FlavorCpu)),
			Ram:              types.Int64Value(int64(f.FlavorRam)),
			Bandwidth:        types.Float64Value(f.Bandwidth),
			GpuVendor:        types.StringValue(f.GpuVendor),
			VideoMemSize:     types.Int64Value(int64(f.VideoMemSize)),
			GpuType:          types.StringValue(f.GpuType),
			GpuCount:         types.Int64Value(int64(f.GpuCount)),
			Available:        types.BoolValue(f.Available),
		})
	}
	config.Flavors = flavors
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *ctyunEcsFlavors) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunEcsFlavorsFlavorsConfig struct {
	Id               types.String  `tfsdk:"id"`
	FlavorSeriesName types.String  `tfsdk:"series_name"`
	CpuInfo          types.String  `tfsdk:"cpu_info"`
	BaseBandwidth    types.Float64 `tfsdk:"base_bandwidth"`
	Name             types.String  `tfsdk:"name"`
	Type             types.String  `tfsdk:"type"`
	Series           types.String  `tfsdk:"series"`
	NicMultiQueue    types.Int64   `tfsdk:"nic_multi_queue"`
	Pps              types.Int64   `tfsdk:"pps"`
	Cpu              types.Int64   `tfsdk:"cpu"`
	Ram              types.Int64   `tfsdk:"ram"`
	Bandwidth        types.Float64 `tfsdk:"bandwidth"`
	GpuVendor        types.String  `tfsdk:"gpu_vendor"`
	VideoMemSize     types.Int64   `tfsdk:"video_mem_size"`
	GpuType          types.String  `tfsdk:"gpu_type"`
	GpuCount         types.Int64   `tfsdk:"gpu_count"`
	Available        types.Bool    `tfsdk:"available"`
}

type CtyunEcsFlavorsConfig struct {
	Type     types.String                   `tfsdk:"type"`
	Name     types.String                   `tfsdk:"name"`
	Cpu      types.Int64                    `tfsdk:"cpu"`
	Ram      types.Int64                    `tfsdk:"ram"`
	Arch     types.String                   `tfsdk:"arch"`
	Series   types.String                   `tfsdk:"series"`
	RegionId types.String                   `tfsdk:"region_id"`
	AzName   types.String                   `tfsdk:"az_name"`
	Flavors  []CtyunEcsFlavorsFlavorsConfig `tfsdk:"flavors"`
}
