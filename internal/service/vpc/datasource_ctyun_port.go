package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunNetworkInterfaces() datasource.DataSource {
	return &ctyunNetworkInterfaces{}
}

type ctyunNetworkInterfaces struct {
	meta *common.CtyunMetadata
}

func (c *ctyunNetworkInterfaces) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ports"
}

func (c *ctyunNetworkInterfaces) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10225195`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "VPC ID，用于过滤指定VPC下的网卡",
			},
			"device_id": schema.StringAttribute{
				Optional:    true,
				Description: "设备ID，用于过滤绑定到指定设备的网卡",
			},
			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Description: "子网ID，用于过滤指定子网下的网卡",
			},
			"page_no": schema.Int64Attribute{
				Optional:    true,
				Description: "页码，从1开始，默认为1",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "每页记录数，取值范围1-50，默认为10",
			},
			"network_interfaces": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "网卡ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "网卡名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "网卡描述",
						},
						"port_id": schema.StringAttribute{
							Computed:    true,
							Description: "网卡ID",
						},
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "子网ID",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "VPC ID",
						},
						"mac_address": schema.StringAttribute{
							Computed:    true,
							Description: "MAC地址",
						},
						"primary_private_ip": schema.StringAttribute{
							Computed:    true,
							Description: "主私有IP地址",
						},
						"secondary_private_ips": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "辅助私有IP地址列表",
						},
						"ipv6_addresses": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "IPv6地址列表",
						},
						"security_group_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "安全组ID列表",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "绑定的实例ID",
						},
						"instance_type": schema.StringAttribute{
							Computed:    true,
							Description: "实例类型",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "网卡状态",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunNetworkInterfaces) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CtyunNetworkInterfacesConfig
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

	// 构造查询请求
	request := &ctvpc.CtvpcVpcListPortRequest{
		RegionID: regionId,
		PageNo:   1,
		PageSize: 10,
	}

	// 设置可选参数
	if !config.VpcId.IsNull() {
		vpcId := config.VpcId.ValueString()
		request.VpcID = &vpcId
	}

	if !config.DeviceId.IsNull() {
		deviceId := config.DeviceId.ValueString()
		request.DeviceID = &deviceId
	}

	if !config.SubnetId.IsNull() {
		subnetId := config.SubnetId.ValueString()
		request.SubnetID = &subnetId
	}

	if !config.PageNo.IsNull() {
		request.PageNo = int32(config.PageNo.ValueInt64())
	}

	if !config.PageSize.IsNull() {
		pageSize := int32(config.PageSize.ValueInt64())
		if pageSize > 50 {
			pageSize = 50
		}
		request.PageSize = pageSize
	}

	// 调用API查询网卡列表
	response, err := c.meta.Apis.SdkCtVpcApis.CtvpcVpcListPortApi.Do(ctx, c.meta.SdkCredential, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"获取网卡列表失败",
			fmt.Sprintf("获取网卡列表时发生错误: %s", err.Error()),
		)
		return
	}

	// 处理响应数据
	var networkInterfaces []CtyunNetworkInterfaces
	for _, port := range response.ReturnObj {
		networkInterface := CtyunNetworkInterfaces{
			Id:                 types.StringPointerValue(port.NetworkInterfaceID),
			Name:               types.StringPointerValue(port.NetworkInterfaceName),
			Description:        types.StringPointerValue(port.Description),
			NetworkInterfaceId: types.StringPointerValue(port.NetworkInterfaceID),
			SubnetId:           types.StringPointerValue(port.SubnetID),
			VpcId:              types.StringPointerValue(port.VpcID),
			MacAddress:         types.StringPointerValue(port.MacAddress),
			PrimaryPrivateIp:   types.StringPointerValue(port.PrimaryPrivateIp),
			InstanceId:         types.StringPointerValue(port.InstanceID),
			InstanceType:       types.StringPointerValue(port.InstanceType),
			Status:             types.StringPointerValue(port.AdminStatus),
		}

		// 处理辅助私有IP
		if port.SecondaryPrivateIps != nil {
			secondaryIps := make([]types.String, len(port.SecondaryPrivateIps))
			for i, ip := range port.SecondaryPrivateIps {
				if ip != nil {
					secondaryIps[i] = types.StringValue(*ip)
				}
			}
			networkInterface.SecondaryPrivateIps = secondaryIps
		}

		// 处理IPv6地址
		if port.Ipv6Addresses != nil {
			ipv6Addresses := make([]types.String, len(port.Ipv6Addresses))
			for i, addr := range port.Ipv6Addresses {
				if addr != nil {
					ipv6Addresses[i] = types.StringValue(*addr)
				}
			}
			networkInterface.Ipv6Addresses = ipv6Addresses
		}

		// 处理安全组ID
		if port.SecurityGroupIds != nil {
			securityGroupIds := make([]types.String, len(port.SecurityGroupIds))
			for i, sgId := range port.SecurityGroupIds {
				if sgId != nil {
					securityGroupIds[i] = types.StringValue(*sgId)
				}
			}
			networkInterface.SecurityGroupIds = securityGroupIds
		}

		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	config.NetworkInterfaces = networkInterfaces
	config.RegionId = types.StringValue(regionId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *ctyunNetworkInterfaces) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunNetworkInterfaces struct {
	Id                  types.String   `tfsdk:"id"`
	Name                types.String   `tfsdk:"name"`
	Description         types.String   `tfsdk:"description"`
	NetworkInterfaceId  types.String   `tfsdk:"port_id"`
	SubnetId            types.String   `tfsdk:"subnet_id"`
	VpcId               types.String   `tfsdk:"vpc_id"`
	MacAddress          types.String   `tfsdk:"mac_address"`
	PrimaryPrivateIp    types.String   `tfsdk:"primary_private_ip"`
	SecondaryPrivateIps []types.String `tfsdk:"secondary_private_ips"`
	Ipv6Addresses       []types.String `tfsdk:"ipv6_addresses"`
	SecurityGroupIds    []types.String `tfsdk:"security_group_ids"`
	InstanceId          types.String   `tfsdk:"instance_id"`
	InstanceType        types.String   `tfsdk:"instance_type"`
	Status              types.String   `tfsdk:"status"`
}

type CtyunNetworkInterfacesConfig struct {
	RegionId          types.String             `tfsdk:"region_id"`
	VpcId             types.String             `tfsdk:"vpc_id"`
	DeviceId          types.String             `tfsdk:"device_id"`
	SubnetId          types.String             `tfsdk:"subnet_id"`
	PageNo            types.Int64              `tfsdk:"page_no"`
	PageSize          types.Int64              `tfsdk:"page_size"`
	NetworkInterfaces []CtyunNetworkInterfaces `tfsdk:"network_interfaces"`
}
