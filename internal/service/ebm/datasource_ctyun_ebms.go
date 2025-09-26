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

type ctyunEbms struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbms() datasource.DataSource {
	return &ctyunEbms{}
}

func (c *ctyunEbms) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ebms"
}

type CtyunEbmsModel struct {
	InstanceID           types.String               `tfsdk:"instance_id"`
	DeviceType           types.String               `tfsdk:"device_type"`
	InstanceName         types.String               `tfsdk:"instance_name"`
	Hostname             types.String               `tfsdk:"hostname"`
	SystemVolumeRaidUUID types.String               `tfsdk:"system_volume_raid_uuid"`
	DataVolumeRaidUUID   types.String               `tfsdk:"data_volume_raid_uuid"`
	ImageUUID            types.String               `tfsdk:"image_uuid"`
	OsTypeName           types.String               `tfsdk:"os_type_name"`
	VpcID                types.String               `tfsdk:"vpc_id"`
	PublicIP             types.String               `tfsdk:"public_ip"`
	PublicIPv6           types.String               `tfsdk:"public_ipv6"`
	Status               types.String               `tfsdk:"status"`
	NetworkCardList      []CtyunEbmsNetworkCardList `tfsdk:"network_card_list"`
	AttachedVolumes      []string                   `tfsdk:"attached_volumes"`
	Freezing             types.Bool                 `tfsdk:"freezing"`
	Expired              types.Bool                 `tfsdk:"expired"`
	CreateTime           types.String               `tfsdk:"create_time"`
	UpdatedTime          types.String               `tfsdk:"updated_time"`
	ExpiredTime          types.String               `tfsdk:"expired_time"`
}
type CtyunEbmsNetworkCardList struct {
	InterfaceID types.String `tfsdk:"interface_id"`
	PortUUID    types.String `tfsdk:"port_uuid"`
	FixedIP     types.String `tfsdk:"fixed_ip"`
	Master      types.Bool   `tfsdk:"master"`
	Ipv6        types.String `tfsdk:"ipv6"`
	SubnetID    types.String `tfsdk:"subnet_id"`
}

type CtyunEbmsConfig struct {
	RegionID       types.String     `tfsdk:"region_id"`
	AzName         types.String     `tfsdk:"az_name"`
	InstanceIDList types.String     `tfsdk:"instance_id_list"`
	Instances      []CtyunEbmsModel `tfsdk:"instances"`
}

func (c *ctyunEbms) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027724/10040106`,
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
			"instance_id_list": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "uuid列表，支持逗号分割",
			},
			"instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "物理机UUID",
						},
						"device_type": schema.StringAttribute{
							Computed:    true,
							Description: "设备类型",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "物理机实例展示名",
						},
						"hostname": schema.StringAttribute{
							Computed:    true,
							Description: "物理机主机名称(hostname)",
						},
						"system_volume_raid_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "本地系统盘raidid",
						},
						"data_volume_raid_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "数据盘raidid",
						},
						"image_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "镜像ID",
						},
						"os_type_name": schema.StringAttribute{
							Computed:    true,
							Description: "操作系统类型",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "主网卡网络ID",
						},
						"public_ip": schema.StringAttribute{
							Computed:    true,
							Description: "公网IPIPv4地址",
						},
						"public_ipv6": schema.StringAttribute{
							Computed:    true,
							Description: "公网IPv6地址",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "物理机状态",
						},

						"network_card_list": schema.ListNestedAttribute{
							Computed:    true,
							Description: "网卡信息",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"interface_id": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "物理机网卡id",
									},
									"port_uuid": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "弹性网卡id",
									},
									"fixed_ip": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "内网IPv4地址",
									},
									"master": schema.BoolAttribute{
										Required:    true,
										Description: "是否主节点(True代表主节点)",
									},
									"ipv6": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "内网IPv6地址",
									},
									"subnet_id": schema.StringAttribute{
										Required:    true,
										Description: "子网id",
									},
								},
							},
						},

						"attached_volumes": schema.ListAttribute{
							Description: "关联的云硬盘ID",
							Computed:    true,
							ElementType: types.StringType,
						},
						"freezing": schema.BoolAttribute{
							Computed:    true,
							Description: "是否冻结",
						},
						"expired": schema.BoolAttribute{
							Computed:    true,
							Description: "是否到期",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"updated_time": schema.StringAttribute{
							Computed:    true,
							Description: "最后更新时间",
						},
						"expired_time": schema.StringAttribute{
							Computed:    true,
							Description: "到期时间",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEbms) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunEbmsConfig
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

	// 组装请求体
	uuids := config.InstanceIDList.ValueString()
	params := &ctebm.EbmListInstanceV4plusRequest{
		RegionID:         regionId,
		AzName:           azName,
		InstanceUUIDList: &(uuids),
	}
	// 调用API
	resp, err := c.meta.Apis.CtEbmApis.EbmListInstanceV4plusApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	model := []CtyunEbmsModel{}
	for _, ebm := range resp.ReturnObj.Results {
		item := CtyunEbmsModel{
			InstanceID:           utils.SecStringValue(ebm.InstanceUUID),
			DeviceType:           utils.SecStringValue(ebm.DeviceType),
			Hostname:             utils.SecStringValue(ebm.InstanceName),
			InstanceName:         utils.SecStringValue(ebm.DisplayName),
			ImageUUID:            utils.SecStringValue(ebm.ImageID),
			SystemVolumeRaidUUID: utils.SecStringValue(ebm.SystemVolumeRaidID),
			DataVolumeRaidUUID:   utils.SecStringValue(ebm.DataVolumeRaidID),
			VpcID:                utils.SecStringValue(ebm.VpcID),
			PublicIP:             utils.SecStringValue(ebm.PublicIP),
			PublicIPv6:           utils.SecStringValue(ebm.PublicIPv6),
			Status:               utils.SecLowerStringValue(ebm.EbmState),
			OsTypeName:           utils.SecStringValue(ebm.OsTypeName),
			Freezing:             utils.SecBoolValue(ebm.Freezing),
			Expired:              utils.SecBoolValue(ebm.Expired),
			UpdatedTime:          utils.SecStringValue(ebm.UpdatedTime),
			CreateTime:           utils.SecStringValue(ebm.CreateTime),
			ExpiredTime:          utils.SecStringValue(ebm.ExpiredTime),
		}

		item.AttachedVolumes = utils.StrPointerArrayToStrArray(ebm.AttachedVolumes)

		// 处理网卡
		var networkCards []CtyunEbmsNetworkCardList
		for _, card := range ebm.Interfaces {
			networkCards = append(networkCards, CtyunEbmsNetworkCardList{
				PortUUID:    utils.SecStringValue(card.PortUUID),
				Master:      utils.SecBoolValue(card.Master),
				FixedIP:     utils.SecStringValue(card.Ipv4),
				Ipv6:        utils.SecStringValue(card.Ipv6),
				SubnetID:    utils.SecStringValue(card.SubnetUUID),
				InterfaceID: utils.SecStringValue(card.InterfaceUUID),
			})
		}
		item.NetworkCardList = networkCards
		model = append(model, item)
	}

	// 保存到state
	config.RegionID = types.StringValue(regionId)
	config.AzName = types.StringValue(azName)
	config.Instances = model
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEbms) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
