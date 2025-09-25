package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEcsInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunEcsInstances{}
)

type ctyunEcsInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunEcsInstances() datasource.DataSource {
	return &ctyunEcsInstances{}
}

func (c *ctyunEcsInstances) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_instances"
}

type CtyunEcsInstancesImage struct {
	ImageID   types.String `tfsdk:"image_id"`
	ImageName types.String `tfsdk:"image_name"`
}

type CtyunEcsInstancesFlavor struct {
	FlavorID     types.String `tfsdk:"flavor_id"`
	FlavorName   types.String `tfsdk:"flavor_name"`
	FlavorCPU    types.Int32  `tfsdk:"flavor_cpu"`
	FlavorRAM    types.Int32  `tfsdk:"flavor_ram"`
	GpuType      types.String `tfsdk:"gpu_type"`
	GpuCount     types.Int32  `tfsdk:"gpu_count"`
	GpuVendor    types.String `tfsdk:"gpu_vendor"`
	VideoMemSize types.Int32  `tfsdk:"video_mem_size"`
}
type CtyunEcsInstancesAffinityGroup struct {
	Policy            types.String `tfsdk:"policy"`
	AffinityGroupName types.String `tfsdk:"affinity_group_name"`
	AffinityGroupID   types.String `tfsdk:"affinity_group_id"`
}

type CtyunEcsInstancesAddress struct {
	VpcName     types.String                   `tfsdk:"vpc_name"`
	AddressList []CtyunEcsInstancesAddressList `tfsdk:"address_list"`
}

type CtyunEcsInstancesAddressList struct {
	Addr       types.String `tfsdk:"addr"`
	Version    types.Int32  `tfsdk:"version"`
	Type       types.String `tfsdk:"type"`
	IsMaster   types.Bool   `tfsdk:"is_master"`
	MacAddress types.String `tfsdk:"mac_address"`
}

type CtyunEcsInstancesSecGroupList struct {
	SecurityGroupName types.String `tfsdk:"security_group_name"`
	SecurityGroupID   types.String `tfsdk:"security_group_id"`
}

type CtyunEcsInstancesVipInfoList struct {
	VipID          types.String `tfsdk:"vip_id"`
	VipAddress     types.String `tfsdk:"vip_address"`
	VipBindNicIP   types.String `tfsdk:"vip_bind_nic_ip"`
	VipBindNicIPv6 types.String `tfsdk:"vip_bind_nic_ipv6"`
	NicID          types.String `tfsdk:"nic_id"`
}

type CtyunEcsInstancesNetworkInfo struct {
	SubnetID  types.String `tfsdk:"subnet_id"`
	IpAddress types.String `tfsdk:"ip_address"`
}

type CtyunEcsInstancesModel struct {
	AzName              types.String                    `tfsdk:"az_name"`
	AzDisplayName       types.String                    `tfsdk:"az_display_name"`
	ExpiredTime         types.String                    `tfsdk:"expired_time"`
	CreatedTime         types.String                    `tfsdk:"created_time"`
	ProjectID           types.String                    `tfsdk:"project_id"`
	AttachedVolumes     []string                        `tfsdk:"attached_volumes"`
	InstanceID          types.String                    `tfsdk:"instance_id"`
	ID                  types.String                    `tfsdk:"id"`
	DisplayName         types.String                    `tfsdk:"display_name"`
	InstanceName        types.String                    `tfsdk:"instance_name"`
	OsType              types.Int32                     `tfsdk:"os_type"`
	InstanceDescription types.String                    `tfsdk:"instance_description"`
	InstanceStatus      types.String                    `tfsdk:"instance_status"`
	OnDemand            types.Bool                      `tfsdk:"on_demand"`
	KeypairName         types.String                    `tfsdk:"keypair_name"`
	Addresses           []CtyunEcsInstancesAddress      `tfsdk:"addresses"`
	SecGroupList        []CtyunEcsInstancesSecGroupList `tfsdk:"sec_group_list"`
	VipInfoList         []CtyunEcsInstancesVipInfoList  `tfsdk:"vip_info_list"`
	NetworkInfo         []CtyunEcsInstancesNetworkInfo  `tfsdk:"network_info"`
	AffinityGroup       CtyunEcsInstancesAffinityGroup  `tfsdk:"affinity_group"`
	Image               CtyunEcsInstancesImage          `tfsdk:"image"`
	Flavor              CtyunEcsInstancesFlavor         `tfsdk:"flavor"`
	DelegateName        types.String                    `tfsdk:"delegate_name"`
	DeletionProtection  types.Bool                      `tfsdk:"deletion_protection"`
}

type CtyunEcsInstancesConfig struct {
	RegionID       types.String             `tfsdk:"region_id"`
	AzName         types.String             `tfsdk:"az_name"`
	InstanceName   types.String             `tfsdk:"instance_name"`
	InstanceIDList types.String             `tfsdk:"instance_id_list"`
	ProjectID      types.String             `tfsdk:"project_id"`
	PageNo         types.Int32              `tfsdk:"page_no"`
	PageSize       types.Int32              `tfsdk:"page_size"`
	Instances      []CtyunEcsInstancesModel `tfsdk:"instances"`
}

func (c *ctyunEcsInstances) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730**`,
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
			"instance_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "云主机名称，精准匹配",
			},
			"instance_id_list": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "uuid列表，支持逗号分割",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "每页记录数目，取值范围：[1,50]，注：默认值为10",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID",
			},
			"instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"az_display_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区展示名称",
						},
						"attached_volumes": schema.ListAttribute{
							Description: "关联的云硬盘ID",
							Computed:    true,
							ElementType: types.StringType,
						},
						"addresses": schema.ListNestedAttribute{
							Computed:    true,
							Description: "网络地址信息",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"vpc_name": schema.StringAttribute{
										Computed:    true,
										Description: "VPC",
									},
									"address_list": schema.ListNestedAttribute{
										Computed:    true,
										Description: "网络地址信息",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"addr": schema.StringAttribute{
													Computed:    true,
													Description: "IP地址",
												},
												"version": schema.Int32Attribute{
													Computed:    true,
													Description: "IP版本",
												},
												"type": schema.StringAttribute{
													Computed:    true,
													Description: "网络类型，取值范围：<br/>fixed（内网），<br/>floating（弹性公网）",
												},
												"is_master": schema.BoolAttribute{
													Computed:    true,
													Description: "是否为主网卡",
												},
												"mac_address": schema.StringAttribute{
													Computed:    true,
													Description: "mac地址",
												},
											},
										},
									},
								},
							},
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机ID，值与instance_id相同",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机ID",
						},
						"display_name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机显示名称",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机名称",
						},
						"os_type": schema.Int32Attribute{
							Computed:    true,
							Description: "操作系统类型，取值范围：<br/>1（linux），<br/>2（windows），<br/>3（redhat），<br/>4（ubuntu），<br/>5（centos），<br/>6（oracle）",
						},
						"instance_status": schema.StringAttribute{
							Computed:    true,
							Description: "云主机状态",
						},
						"expired_time": schema.StringAttribute{
							Computed:    true,
							Description: "到期时间",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"sec_group_list": schema.ListNestedAttribute{
							Computed:    true,
							Description: "安全组信息列表",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"security_group_name": schema.StringAttribute{
										Computed:    true,
										Description: "安全组名称",
									},
									"security_group_id": schema.StringAttribute{
										Computed:    true,
										Description: "安全组id",
									},
								},
							},
						},
						"vip_info_list": schema.ListNestedAttribute{
							Computed:    true,
							Description: "虚拟IP信息列表",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"vip_id": schema.StringAttribute{
										Computed:    true,
										Description: "虚拟IP的ID",
									},
									"vip_address": schema.StringAttribute{
										Computed:    true,
										Description: "虚拟IP地址",
									},
									"vip_bind_nic_ip": schema.StringAttribute{
										Computed:    true,
										Description: "虚拟IP绑定的网卡对应IPv4地址",
									},
									"vip_bind_nic_ipv6": schema.StringAttribute{
										Computed:    true,
										Description: "虚拟IP绑定的网卡对应IPv6地址",
									},
									"nic_id": schema.StringAttribute{
										Computed:    true,
										Description: "网卡ID",
									},
								},
							},
						},
						"affinity_group": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "云主机组信息",
							Attributes: map[string]schema.Attribute{
								"affinity_group_id": schema.StringAttribute{
									Computed:    true,
									Description: "云主机组名称ID",
								},
								"affinity_group_name": schema.StringAttribute{
									Computed:    true,
									Description: "云主机组名称",
								},
								"policy": schema.StringAttribute{
									Computed:    true,
									Description: "云主机组策略类型，取值范围：<br />anti-affinity（强制反亲和性），<br />affinity（强制亲和性），<br />soft-anti-affinity（反亲和性），<br />soft-affinity（亲和性)，<br />power-anti-affinity（电力反亲和性)",
								},
							},
						},
						"image": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "镜像信息",
							Attributes: map[string]schema.Attribute{
								"image_id": schema.StringAttribute{
									Computed:    true,
									Description: "镜像id",
								},
								"image_name": schema.StringAttribute{
									Computed:    true,
									Description: "镜像名称",
								},
							},
						},
						"flavor": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "云主机规格信息",
							Attributes: map[string]schema.Attribute{
								"flavor_id": schema.StringAttribute{
									Computed:    true,
									Description: "规格ID",
								},
								"flavor_name": schema.StringAttribute{
									Computed:    true,
									Description: "规格名称",
								},
								"flavor_cpu": schema.Int32Attribute{
									Computed:    true,
									Description: "规格CPU",
								},
								"flavor_ram": schema.Int32Attribute{
									Computed:    true,
									Description: "规格RAM",
								},
								"gpu_type": schema.StringAttribute{
									Computed:    true,
									Description: "GPU类型",
								},
								"gpu_count": schema.Int32Attribute{
									Computed:    true,
									Description: "GPU数量",
								},
								"gpu_vendor": schema.StringAttribute{
									Computed:    true,
									Description: "GPU名称",
								},
								"video_mem_size": schema.Int32Attribute{
									Computed:    true,
									Description: "显存大小",
								},
							},
						},
						"on_demand": schema.BoolAttribute{
							Computed:    true,
							Description: "付费方式，取值范围：<br/>true（按量付费）;<br/>false（包周期）",
						},
						"keypair_name": schema.StringAttribute{
							Computed:    true,
							Description: "密钥对名称",
						},
						"network_info": schema.ListNestedAttribute{
							Computed:    true,
							Description: "网络信息",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"subnet_id": schema.StringAttribute{
										Computed:    true,
										Description: "子网ID",
									},
									"ip_address": schema.StringAttribute{
										Computed:    true,
										Description: "IP地址",
									},
								},
							},
						},
						"delegate_name": schema.StringAttribute{
							Computed:    true,
							Description: "委托名称，注：委托绑定目前仅支持多可用区类型资源池，非可用区资源池为空字符串",
						},
						"deletion_protection": schema.BoolAttribute{
							Computed:    true,
							Description: "是否开启实例删除保护",
						},
						"instance_description": schema.StringAttribute{
							Computed:    true,
							Description: "云主机描述信息",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEcsInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunEcsInstancesConfig
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
	config.Instances = []CtyunEcsInstancesModel{}
	config.AzName = types.StringValue(azName)
	config.RegionID = types.StringValue(regionId)
	// 组装请求体
	params := &ctecs.CtecsDescribeInstancesRequest{
		RegionID:       regionId,
		AzName:         azName,
		ProjectID:      config.ProjectID.ValueString(),
		PageNo:         config.PageNo.ValueInt32(),
		PageSize:       config.PageSize.ValueInt32(),
		InstanceName:   config.InstanceName.ValueString(),
		InstanceIDList: config.InstanceIDList.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsDescribeInstancesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	for _, ecs := range resp.ReturnObj.Results {
		if ecs == nil {
			continue
		}
		item := CtyunEcsInstancesModel{
			AzName:              types.StringValue(ecs.AzName),
			AzDisplayName:       types.StringValue(ecs.AzDisplayName),
			ExpiredTime:         types.StringValue(ecs.ExpiredTime),
			CreatedTime:         types.StringValue(ecs.CreatedTime),
			ProjectID:           types.StringValue(ecs.ProjectID),
			AttachedVolumes:     ecs.AttachedVolume,
			InstanceID:          types.StringValue(ecs.InstanceID),
			ID:                  types.StringValue(ecs.InstanceID),
			DisplayName:         types.StringValue(ecs.DisplayName),
			InstanceName:        types.StringValue(ecs.InstanceName),
			OsType:              types.Int32Value(ecs.OsType),
			InstanceDescription: types.StringValue(ecs.InstanceDescription),
			InstanceStatus:      types.StringValue(ecs.InstanceStatus),
			OnDemand:            utils.SecBoolValue(ecs.OnDemand),
			KeypairName:         types.StringValue(ecs.KeypairName),
			DelegateName:        types.StringValue(ecs.DelegateName),
			DeletionProtection:  utils.SecBoolValue(ecs.DeletionProtection),
		}
		if ecs.Image != nil {
			item.Image = CtyunEcsInstancesImage{
				ImageID:   types.StringValue(ecs.Image.ImageID),
				ImageName: types.StringValue(ecs.Image.ImageName),
			}
		}
		if ecs.Flavor != nil {
			item.Flavor = CtyunEcsInstancesFlavor{
				FlavorID:   types.StringValue(ecs.Flavor.FlavorID),
				FlavorName: types.StringValue(ecs.Flavor.FlavorName),
				FlavorCPU:  types.Int32Value(ecs.Flavor.FlavorCPU),
				FlavorRAM:  types.Int32Value(ecs.Flavor.FlavorRAM),
				GpuCount:   types.Int32Value(ecs.Flavor.GpuCount),
			}
		}
		if ecs.AffinityGroup != nil {
			item.AffinityGroup = CtyunEcsInstancesAffinityGroup{
				Policy:            types.StringValue(ecs.AffinityGroup.Policy),
				AffinityGroupID:   types.StringValue(ecs.AffinityGroup.AffinityGroupID),
				AffinityGroupName: types.StringValue(ecs.AffinityGroup.AffinityGroupName),
			}
		}
		for _, addr := range ecs.Addresses {
			t := CtyunEcsInstancesAddress{
				VpcName: types.StringValue(addr.VpcName),
			}
			for _, l := range addr.AddressList {
				t.AddressList = append(t.AddressList, CtyunEcsInstancesAddressList{
					Addr:       types.StringValue(l.Addr),
					Version:    types.Int32Value(l.Version),
					Type:       types.StringValue(l.RawType),
					IsMaster:   utils.SecBoolValue(l.IsMaster),
					MacAddress: types.StringValue(l.MacAddress),
				})
			}
			item.Addresses = append(item.Addresses, t)
		}
		for _, sg := range ecs.SecGroupList {
			t := CtyunEcsInstancesSecGroupList{
				SecurityGroupID:   types.StringValue(sg.SecurityGroupID),
				SecurityGroupName: types.StringValue(sg.SecurityGroupName),
			}
			item.SecGroupList = append(item.SecGroupList, t)
		}
		for _, vip := range ecs.VipInfoList {
			t := CtyunEcsInstancesVipInfoList{
				VipID:          types.StringValue(vip.VipID),
				VipAddress:     types.StringValue(vip.VipAddress),
				VipBindNicIP:   types.StringValue(vip.VipBindNicIP),
				VipBindNicIPv6: types.StringValue(vip.VipBindNicIPv6),
				NicID:          types.StringValue(vip.NicID),
			}
			item.VipInfoList = append(item.VipInfoList, t)
		}
		for _, nic := range ecs.NetworkInfo {
			t := CtyunEcsInstancesNetworkInfo{
				SubnetID:  types.StringValue(nic.SubnetID),
				IpAddress: types.StringValue(nic.IpAddress),
			}
			item.NetworkInfo = append(item.NetworkInfo, t)
		}
		config.Instances = append(config.Instances, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEcsInstances) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
