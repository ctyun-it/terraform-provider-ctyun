package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunSubnets{}
	_ datasource.DataSourceWithConfigure = &ctyunSubnets{}
)

type ctyunSubnets struct {
	meta *common.CtyunMetadata
}

func NewCtyunSubnets() datasource.DataSource {
	return &ctyunSubnets{}
}

func (c *ctyunSubnets) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_subnets"
}

type CtyunSubnetsModel struct {
	SubnetID          types.String `tfsdk:"subnet_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	VpcID             types.String `tfsdk:"vpc_id"`
	AvailabilityZones types.Set    `tfsdk:"availability_zones"`
	RouteTableID      types.String `tfsdk:"route_table_id"`
	NetworkAclID      types.String `tfsdk:"network_acl_id"`
	CIDR              types.String `tfsdk:"cidr"`
	GatewayIP         types.String `tfsdk:"gateway_ip"`
	Start             types.String `tfsdk:"start"`
	End               types.String `tfsdk:"end"`
	AvailableIPCount  types.Int32  `tfsdk:"available_ipcount"`
	Ipv6Enabled       types.Int32  `tfsdk:"ipv6_enabled"`
	EnableIpv6        types.Bool   `tfsdk:"enable_ipv6"`
	Ipv6CIDR          types.String `tfsdk:"ipv6_cidr"`
	Ipv6Start         types.String `tfsdk:"ipv6_start"`
	Ipv6End           types.String `tfsdk:"ipv6_end"`
	Ipv6GatewayIP     types.String `tfsdk:"ipv6_gateway_ip"`
	DnsList           types.Set    `tfsdk:"dns_list"`
	NtpList           types.Set    `tfsdk:"ntp_list"`
	Type              types.Int32  `tfsdk:"type"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

type CtyunSubnetsConfig struct {
	RegionID types.String `tfsdk:"region_id"`
	VpcID    types.String `tfsdk:"vpc_id"`
	PageNo   types.Int32  `tfsdk:"page_no"`
	PageSize types.Int32  `tfsdk:"page_size"`
	SubnetID types.String `tfsdk:"subnet_id"`

	CurrentCount types.Int32         `tfsdk:"current_count"`
	TotalCount   types.Int32         `tfsdk:"total_count"`
	TotalPage    types.Int32         `tfsdk:"total_page"`
	Subnets      []CtyunSubnetsModel `tfsdk:"subnets"`
}

func (c *ctyunSubnets) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026755/10197656**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "虚拟私有云ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "分页查询时每页的行数，最大值为50，默认值为10。",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Description: "子网ID",
			},
			"current_count": schema.Int32Attribute{
				Computed:    true,
				Description: "分页查询时每页的行数。",
			},
			"total_count": schema.Int32Attribute{
				Computed:    true,
				Description: "总数。",
			},
			"total_page": schema.Int32Attribute{
				Computed:    true,
				Description: "总页数。",
			},
			"subnets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "subnetID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "VpcID",
						},
						"availability_zones": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "子网所在的可用区名",
						},
						"route_table_id": schema.StringAttribute{
							Computed:    true,
							Description: "子网路由表ID",
						},
						"network_acl_id": schema.StringAttribute{
							Computed:    true,
							Description: "子网aclID",
						},
						"cidr": schema.StringAttribute{
							Computed:    true,
							Description: "子网网段，掩码范围为16-28位",
						},
						"gateway_ip": schema.StringAttribute{
							Computed:    true,
							Description: "子网网关",
						},
						"start": schema.StringAttribute{
							Computed:    true,
							Description: "子网网段起始IP",
						},
						"end": schema.StringAttribute{
							Computed:    true,
							Description: "子网网段结束IP",
						},
						"available_ipcount": schema.Int64Attribute{
							Computed:    true,
							Description: "子网内可用IPv4数目",
						},
						"ipv6_enabled": schema.Int64Attribute{
							Computed:    true,
							Description: "是否配置了ipv6网段",
						},
						"enable_ipv6": schema.BoolAttribute{
							Computed:    true,
							Description: "是否开启ipv6",
						},
						"ipv6_cidr": schema.StringAttribute{
							Computed:    true,
							Description: "子网Ipv6网段，掩码范围为16-28位",
						},
						"ipv6_start": schema.StringAttribute{
							Computed:    true,
							Description: "子网内可用的起始IPv6地址",
						},
						"ipv6_end": schema.StringAttribute{
							Computed:    true,
							Description: "子网内可用的结束IPv6地址",
						},
						"ipv6_gateway_ip": schema.StringAttribute{
							Computed:    true,
							Description: "v6网关地址",
						},
						"dns_list": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "DNS服务器地址:默认为空；必须为正确的IPv4格式；重新触发DHCP后生效，最大数组长度为4",
						},
						"ntp_list": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "NTP服务器地址:默认为空，必须为正确的域名或IPv4格式；重新触发DHCP后生效，最大数组长度为4",
						},
						"type": schema.Int64Attribute{
							Computed:    true,
							Description: "子网类型:当前仅支持：0（普通子网）,1（裸金属子网）",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
					},
				},
			},
		}}
}

func (c *ctyunSubnets) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunSubnetsConfig
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
	params := &ctvpc.CtvpcNewSubnetListRequest{
		RegionID: regionId,
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	subnetId := config.SubnetID.ValueString()
	vpcId := config.VpcID.ValueString()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if subnetId != "" {
		params.SubnetID = &subnetId
	}
	if vpcId != "" {
		params.VpcID = &vpcId
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewSubnetListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Subnets = []CtyunSubnetsModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, s := range resp.ReturnObj.Subnets {
		item := CtyunSubnetsModel{
			SubnetID:         utils.SecStringValue(s.SubnetID),
			Name:             utils.SecStringValue(s.Name),
			Description:      utils.SecStringValue(s.Description),
			VpcID:            utils.SecStringValue(s.VpcID),
			RouteTableID:     utils.SecStringValue(s.RouteTableID),
			NetworkAclID:     utils.SecStringValue(s.NetworkAclID),
			CIDR:             utils.SecStringValue(s.CIDR),
			GatewayIP:        utils.SecStringValue(s.GatewayIP),
			Start:            utils.SecStringValue(s.Start),
			End:              utils.SecStringValue(s.End),
			AvailableIPCount: types.Int32Value(s.AvailableIPCount),
			Ipv6Enabled:      types.Int32Value(s.Ipv6Enabled),
			EnableIpv6:       utils.SecBoolValue(s.EnableIpv6),
			Ipv6CIDR:         utils.SecStringValue(s.Ipv6CIDR),
			Ipv6Start:        utils.SecStringValue(s.Ipv6Start),
			Ipv6End:          utils.SecStringValue(s.Ipv6End),
			Ipv6GatewayIP:    utils.SecStringValue(s.Ipv6GatewayIP),
			Type:             types.Int32Value(s.RawType),
			UpdatedAt:        utils.SecStringValue(s.UpdatedAt),
		}
		item.AvailabilityZones, _ = types.SetValueFrom(ctx, types.StringType, s.AvailabilityZones)
		item.DnsList, _ = types.SetValueFrom(ctx, types.StringType, s.DnsList)
		item.NtpList, _ = types.SetValueFrom(ctx, types.StringType, s.NtpList)
		config.Subnets = append(config.Subnets, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunSubnets) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
