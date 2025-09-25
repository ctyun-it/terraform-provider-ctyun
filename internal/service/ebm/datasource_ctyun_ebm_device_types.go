package ebm

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebm"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEbmDeviceTypes{}
	_ datasource.DataSourceWithConfigure = &ctyunEbmDeviceTypes{}
)

type ctyunEbmDeviceTypes struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbmDeviceTypes() datasource.DataSource {
	return &ctyunEbmDeviceTypes{}
}

func (c *ctyunEbmDeviceTypes) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebm_device_types"
}

type CtyunEbmDeviceTypesModel struct {
	ID                      types.Int32  `tfsdk:"id"`
	DeviceType              types.String `tfsdk:"device_type"`
	CpuModel                types.String `tfsdk:"cpu_model"`
	NvmeVolumeType          types.String `tfsdk:"nvme_volume_type"`
	NameZh                  types.String `tfsdk:"name_zh"`
	NvmeVolumeInterface     types.String `tfsdk:"nvme_volume_interface"`
	UpdateTime              types.String `tfsdk:"update_time"`
	SystemVolumeSize        types.Int32  `tfsdk:"system_volume_size"`
	SystemVolumeType        types.String `tfsdk:"system_volume_type"`
	CpuManufacturer         types.String `tfsdk:"cpu_manufacturer"`
	NameEn                  types.String `tfsdk:"name_en"`
	NicAmount               types.Int32  `tfsdk:"nic_amount"`
	NvmeVolumeAmount        types.Int32  `tfsdk:"nvme_volume_amount"`
	SmartNicExist           types.Bool   `tfsdk:"smart_nic_exist"`
	CpuFrequency            types.String `tfsdk:"cpu_frequency"`
	CpuThreadAmount         types.Int32  `tfsdk:"cpu_thread_amount"`
	SystemVolumeInterface   types.String `tfsdk:"system_volume_interface"`
	GpuManufacturer         types.String `tfsdk:"gpu_manufacturer"`
	DataVolumeType          types.String `tfsdk:"data_volume_type"`
	GpuModel                types.String `tfsdk:"gpu_model"`
	SystemVolumeAmount      types.Int32  `tfsdk:"system_volume_amount"`
	DataVolumeDescription   types.String `tfsdk:"data_volume_description"`
	GpuSize                 types.Int32  `tfsdk:"gpu_size"`
	MemAmount               types.Int32  `tfsdk:"mem_amount"`
	MemSize                 types.Int32  `tfsdk:"mem_size"`
	GpuAmount               types.Int32  `tfsdk:"gpu_amount"`
	SystemVolumeDescription types.String `tfsdk:"system_volume_description"`
	MemFrequency            types.Int32  `tfsdk:"mem_frequency"`
	AzName                  types.String `tfsdk:"az_name"`
	NvmeVolumeSize          types.Int32  `tfsdk:"nvme_volume_size"`
	CpuSockets              types.Int32  `tfsdk:"cpu_sockets"`
	CpuAmount               types.Int32  `tfsdk:"cpu_amount"`
	CreateTime              types.String `tfsdk:"create_time"`
	SupportCloud            types.Bool   `tfsdk:"support_cloud"`
	DataVolumeAmount        types.Int32  `tfsdk:"data_volume_amount"`
	NumaNodeAmount          types.Int32  `tfsdk:"numa_node_amount"`
	Region                  types.String `tfsdk:"region"`
	DataVolumeSize          types.Int32  `tfsdk:"data_volume_size"`
	DataVolumeInterface     types.String `tfsdk:"data_volume_interface"`
	NicRate                 types.Int32  `tfsdk:"nic_rate"`
	CloudBoot               types.Bool   `tfsdk:"cloud_boot"`
	EnableShadowVpc         types.Bool   `tfsdk:"enable_shadow_vpc"`
	ComputeIBAmount         types.Int32  `tfsdk:"compute_i_b_amount"`
	ComputeIBRate           types.Int32  `tfsdk:"compute_i_b_rate"`
	StorageIBAmount         types.Int32  `tfsdk:"storage_i_b_amount"`
	StorageIBRate           types.Int32  `tfsdk:"storage_i_b_rate"`
	ComputeRoCEAmount       types.Int32  `tfsdk:"compute_ro_c_e_amount"`
	ComputeRoCERate         types.Int32  `tfsdk:"compute_ro_c_e_rate"`
	StorageRoCEAmount       types.Int32  `tfsdk:"storage_ro_c_e_amount"`
	StorageRoCERate         types.Int32  `tfsdk:"storage_ro_c_e_rate"`
	Project                 types.String `tfsdk:"project"`
}

type CtyunEbmDeviceTypesConfig struct {
	DeviceType  types.String               `tfsdk:"device_type"`
	RegionID    types.String               `tfsdk:"region_id"`
	AzName      types.String               `tfsdk:"az_name"`
	DeviceTypes []CtyunEbmDeviceTypesModel `tfsdk:"device_types"`
}

func (c *ctyunEbmDeviceTypes) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027724/10754001`,
		Attributes: map[string]schema.Attribute{
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
			"device_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "套餐类型",
			},
			"device_types": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "套餐ID",
						},
						"device_type": schema.StringAttribute{
							Computed:    true,
							Description: "套餐类型",
						},
						"cpu_model": schema.StringAttribute{
							Computed:    true,
							Description: "cpu型号",
						},
						"nvme_volume_type": schema.StringAttribute{
							Computed:    true,
							Description: "NVME介质类型； 包含SSD、HDD",
						},
						"name_zh": schema.StringAttribute{
							Computed:    true,
							Description: "物理机中文名",
						},
						"nvme_volume_interface": schema.StringAttribute{
							Computed:    true,
							Description: "NVME接口类型；包含SAS、SATA、NVMe",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "最后更新时间",
						},
						"system_volume_size": schema.Int64Attribute{
							Computed:    true,
							Description: "系统盘单盘大小(GB)",
						},
						"system_volume_type": schema.StringAttribute{
							Computed:    true,
							Description: "系统盘介质类型； 包含SSD、HDD",
						},
						"cpu_manufacturer": schema.StringAttribute{
							Computed:    true,
							Description: "cpu厂商；Intel，AMD，Hygon，HiSilicon，Loongson等",
						},
						"name_en": schema.StringAttribute{
							Computed:    true,
							Description: "英文名",
						},
						"nic_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "网卡数",
						},
						"nvme_volume_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "NVME硬盘数量",
						},
						"smart_nic_exist": schema.BoolAttribute{
							Computed:    true,
							Description: "是否有智能网卡，true为弹性裸金属; false为标准裸金属",
						},
						"cpu_frequency": schema.StringAttribute{
							Computed:    true,
							Description: "cpu频率(G)",
						},
						"cpu_thread_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "单个cpu核超线程数量",
						},
						"system_volume_interface": schema.StringAttribute{
							Computed:    true,
							Description: "系统盘接口类型；包含SAS、SATA、NVMe",
						},
						"gpu_manufacturer": schema.StringAttribute{
							Computed:    true,
							Description: "GPU厂商；Nvidia，Huawei，Cambricon等",
						},
						"data_volume_type": schema.StringAttribute{
							Computed:    true,
							Description: "数据盘介质类型； 包含SSD、HDD",
						},
						"gpu_model": schema.StringAttribute{
							Computed:    true,
							Description: "GPU型号",
						},
						"system_volume_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "系统盘数量",
						},
						"data_volume_description": schema.StringAttribute{
							Computed:    true,
							Description: "数据盘描述",
						},
						"gpu_size": schema.Int64Attribute{
							Computed:    true,
							Description: "GPU显存",
						},
						"mem_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "内存数",
						},
						"mem_size": schema.Int64Attribute{
							Computed:    true,
							Description: "内存大小(G)",
						},
						"gpu_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "GPU数目",
						},
						"system_volume_description": schema.StringAttribute{
							Computed:    true,
							Description: "系统盘描述",
						},
						"mem_frequency": schema.Int64Attribute{
							Computed:    true,
							Description: "内存频率(MHz)",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区",
						},
						"nvme_volume_size": schema.Int64Attribute{
							Computed:    true,
							Description: "NVME硬盘数量",
						},
						"cpu_sockets": schema.Int64Attribute{
							Computed:    true,
							Description: "物理cpu数量",
						},
						"cpu_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "单个cpu核数",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"support_cloud": schema.BoolAttribute{
							Computed:    true,
							Description: "是否支持云盘",
						},
						"data_volume_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "数据盘数量",
						},
						"numa_node_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "单个cpu numa node数量",
						},
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "资源池",
						},
						"data_volume_size": schema.Int64Attribute{
							Computed:    true,
							Description: "数据盘单盘大小(GB)",
						},
						"data_volume_interface": schema.StringAttribute{
							Computed:    true,
							Description: "数据盘接口；包含SAS、SATA、NVMe",
						},
						"nic_rate": schema.Int64Attribute{
							Computed:    true,
							Description: "网卡传播速率(GE)",
						},
						"cloud_boot": schema.BoolAttribute{
							Computed:    true,
							Description: "是否支持云盘启动",
						},
						"enable_shadow_vpc": schema.BoolAttribute{
							Computed:    true,
							Description: "是否支持存储高速网络；如支持存储高速网络则会占用对应可用网卡数量",
						},
						"compute_i_b_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "计算ib网卡大小",
						},
						"compute_i_b_rate": schema.Int64Attribute{
							Computed:    true,
							Description: "计算ib网卡速率(GE)",
						},
						"storage_i_b_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "存储ib网卡大小",
						},
						"storage_i_b_rate": schema.Int64Attribute{
							Computed:    true,
							Description: "存储ib网卡速率(GE)",
						},
						"compute_ro_c_e_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "计算RoCE网卡大小",
						},
						"compute_ro_c_e_rate": schema.Int64Attribute{
							Computed:    true,
							Description: "计算RoCE网卡速率(GE)",
						},
						"storage_ro_c_e_amount": schema.Int64Attribute{
							Computed:    true,
							Description: "存储RoCE网卡大小",
						},
						"storage_ro_c_e_rate": schema.Int64Attribute{
							Computed:    true,
							Description: "存储RoCE网卡速率(GE)",
						},
						"project": schema.StringAttribute{
							Computed:    true,
							Description: "项目信息",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEbmDeviceTypes) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunEbmDeviceTypesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	azName := c.meta.GetExtraIfEmpty(config.AzName.ValueString(), common.ExtraAzName)
	if azName == "" {
		err = fmt.Errorf("azName不能为空")
		return
	}
	deviceType := config.DeviceType.ValueString()
	params := &ctebm.EbmDeviceTypeListRequest{
		RegionID:   regionId,
		AzName:     azName,
		DeviceType: &(deviceType),
	}

	resp, err := c.meta.Apis.CtEbmApis.EbmDeviceTypeListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	deviceTypes := []CtyunEbmDeviceTypesModel{}
	for _, f := range resp.ReturnObj.Results {
		deviceTypes = append(deviceTypes, CtyunEbmDeviceTypesModel{
			ID:                      types.Int32Value(f.Id),
			DeviceType:              utils.SecStringValue(f.DeviceType),
			CpuModel:                utils.SecStringValue(f.CpuModel),
			NvmeVolumeType:          utils.SecStringValue(f.NvmeVolumeType),
			NameZh:                  utils.SecStringValue(f.NameZh),
			NvmeVolumeInterface:     utils.SecStringValue(f.NvmeVolumeInterface),
			UpdateTime:              utils.SecStringValue(f.UpdateTime),
			SystemVolumeSize:        types.Int32Value(f.SystemVolumeSize),
			SystemVolumeType:        utils.SecStringValue(f.SystemVolumeType),
			CpuManufacturer:         utils.SecStringValue(f.CpuManufacturer),
			NameEn:                  utils.SecStringValue(f.NameEn),
			NicAmount:               types.Int32Value(f.NicAmount),
			NvmeVolumeAmount:        types.Int32Value(f.NvmeVolumeAmount),
			SmartNicExist:           utils.SecBoolValue(f.SmartNicExist),
			CpuFrequency:            utils.SecStringValue(f.CpuFrequency),
			CpuThreadAmount:         types.Int32Value(f.CpuThreadAmount),
			SystemVolumeInterface:   utils.SecStringValue(f.SystemVolumeInterface),
			GpuManufacturer:         utils.SecStringValue(f.GpuManufacturer),
			DataVolumeType:          utils.SecStringValue(f.DataVolumeType),
			GpuModel:                utils.SecStringValue(f.GpuModel),
			SystemVolumeAmount:      types.Int32Value(f.SystemVolumeAmount),
			DataVolumeDescription:   utils.SecStringValue(f.DataVolumeDescription),
			GpuSize:                 types.Int32Value(f.GpuSize),
			MemAmount:               types.Int32Value(f.MemAmount),
			MemSize:                 types.Int32Value(f.MemSize),
			GpuAmount:               types.Int32Value(f.GpuAmount),
			SystemVolumeDescription: utils.SecStringValue(f.SystemVolumeDescription),
			MemFrequency:            types.Int32Value(f.MemFrequency),
			AzName:                  utils.SecStringValue(f.AzName),
			NvmeVolumeSize:          types.Int32Value(f.NvmeVolumeSize),
			CpuSockets:              types.Int32Value(f.CpuSockets),
			CpuAmount:               types.Int32Value(f.CpuAmount),
			CreateTime:              utils.SecStringValue(f.CreateTime),
			SupportCloud:            utils.SecBoolValue(f.SupportCloud),
			DataVolumeAmount:        types.Int32Value(f.DataVolumeAmount),
			NumaNodeAmount:          types.Int32Value(f.NumaNodeAmount),
			Region:                  utils.SecStringValue(f.Region),
			DataVolumeSize:          types.Int32Value(f.DataVolumeSize),
			DataVolumeInterface:     utils.SecStringValue(f.DataVolumeInterface),
			NicRate:                 types.Int32Value(f.NicRate),
			CloudBoot:               utils.SecBoolValue(f.CloudBoot),
			EnableShadowVpc:         utils.SecBoolValue(f.EnableShadowVpc),
			ComputeIBAmount:         types.Int32Value(f.ComputeIBAmount),
			ComputeIBRate:           types.Int32Value(f.ComputeIBRate),
			StorageIBAmount:         types.Int32Value(f.StorageIBAmount),
			StorageIBRate:           types.Int32Value(f.StorageIBRate),
			ComputeRoCEAmount:       types.Int32Value(f.ComputeRoCEAmount),
			ComputeRoCERate:         types.Int32Value(f.ComputeRoCERate),
			StorageRoCEAmount:       types.Int32Value(f.StorageRoCEAmount),
			StorageRoCERate:         types.Int32Value(f.StorageRoCERate),
			Project:                 utils.SecStringValue(f.Project),
		})
	}
	config.DeviceTypes = deviceTypes
	config.RegionID = types.StringValue(regionId)
	config.AzName = types.StringValue(azName)
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEbmDeviceTypes) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
