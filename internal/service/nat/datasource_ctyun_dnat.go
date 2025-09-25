package nat

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunDNatDatasource{}
	_ datasource.DataSourceWithConfigure = &ctyunDNatDatasource{}
)

type ctyunDNatDatasource struct {
	meta *common.CtyunMetadata
}

func NewCtyunDNats() datasource.DataSource {
	return &ctyunDNatDatasource{}
}

func (c *ctyunDNatDatasource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_nat_dnats"
}

func (c *ctyunDNatDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10166345`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填这默认使用provider ctyun总region_id 或者环境变量",
			},
			"nat_gateway_id": schema.StringAttribute{
				Required:    true,
				Description: "要查询的NAT网关的ID",
			},
			"dnats": schema.ListNestedAttribute{
				Computed:    true,
				Description: "dnats列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述信息",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "dnatID 值",
						},
						"dnat_id": schema.StringAttribute{
							Computed:    true,
							Description: "dnatID 值",
						},
						"ip_expire_time": schema.StringAttribute{
							Computed:    true,
							Description: "ip到期时间",
						},
						"extend_id": schema.StringAttribute{
							Computed:    true,
							Description: "弹性 IP id",
						},
						"external_ip": schema.StringAttribute{
							Computed:    true,
							Description: "弹性 IP 地址",
						},
						"external_port": schema.Int64Attribute{
							Computed:    true,
							Description: "外部访问端口",
							Validators: []validator.Int64{
								int64validator.Between(1, 1024),
							},
						},
						"internal_port": schema.Int64Attribute{
							Computed:    true,
							Description: "内部访问端口",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"internal_ip": schema.StringAttribute{
							Computed:    true,
							Description: "内网 IP 地址",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "TCP:转发TCP协议的报文 UDP：转发UDP协议的报文",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "运行状态: ACTIVE / FREEZING / CREATING",
							Validators: []validator.String{
								stringvalidator.OneOf(business.DNatStates...),
							},
						},
						"virtual_machine_display_name": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟机展示名称",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟机id",
						},
						"virtual_machine_name": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟机名称",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunDNatDatasource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunDNatConfig
	// 读取请求信息
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	// region_id必不能为空
	if regionId == "" {
		msg := "regionID不能为空"
		response.Diagnostics.AddError(msg, msg)
		return
	}
	natGatewayId := config.NatGateWayID.ValueString()
	params := &ctvpc.CtvpcListDnatEntriesRequest{
		RegionID:     regionId,
		NatGatewayID: natGatewayId,
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListDnatEntriesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var dnats []CtyunDNatModel
	for _, dnat := range resp.ReturnObj {
		dnatItem := CtyunDNatModel{
			ID:           utils.SecStringValue(dnat.Id),
			DNatID:       utils.SecStringValue(dnat.DNatID),
			Description:  utils.SecStringValue(dnat.Description),
			CreateTime:   utils.SecStringValue(dnat.CreationTime),
			IpExpireTime: utils.SecStringValue(dnat.IpExpireTime),
			ExternalIp:   utils.SecStringValue(dnat.ExternalIp),
			InternalIP:   utils.SecStringValue(dnat.InternalIp),
			//Protocol:                  utils.SecStringValue(dnat.Protocol),
			VirtualMachineId:          utils.SecStringValue(dnat.VirtualMachineID),
			VirtualMachineName:        utils.SecStringValue(dnat.VirtualMachineName),
			VirtualMachineDisplayName: utils.SecStringValue(dnat.VirtualMachineDisplayName),
			//State:                     utils.SecStringValue(dnat.State),
			//ExternalPort:              types.Int64Value(int64(dnat.ExternalPort)),
			//InternalPort:              types.Int64Value(int64(dnat.InternalPort)),
		}
		protocol := utils.SecStringValue(dnat.Protocol)
		if c.contains(business.DNatProtocols, protocol.ValueString()) {
			dnatItem.Protocol = protocol
		}
		state := utils.SecStringValue(dnat.State)
		if c.contains(business.DNatStatus, state.ValueString()) {
			dnatItem.State = state
		}
		externalPort := types.Int64Value(int64(dnat.ExternalPort))
		if c.isPort(externalPort) {
			dnatItem.ExternalPort = externalPort
		}
		internalPort := types.Int64Value(int64(dnat.InternalPort))
		if c.isPort(internalPort) {
			dnatItem.InternalPort = internalPort
		}

		dnats = append(dnats, dnatItem)
	}

	config.RegionID = types.StringValue(regionId)
	config.NatGateWayID = types.StringValue(natGatewayId)
	config.Dnats = dnats
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}
func (c *ctyunDNatDatasource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// contains判断string字符串切片中是否包含某个字符串
func (c *ctyunDNatDatasource) contains(slice []string, item string) bool {
	if item == "" || len(item) == 0 {
		return false
	}
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isPort 判断参数是否符合端口0~65535范围
func (c *ctyunDNatDatasource) isPort(port types.Int64) bool {
	if port.IsNull() {
		return false
	}

	if port.ValueInt64() >= 0 && port.ValueInt64() <= 65535 {
		return true
	}
	return false
}

type CtyunDNatConfig struct {
	RegionID     types.String     `tfsdk:"region_id"`
	NatGateWayID types.String     `tfsdk:"nat_gateway_id"`
	Dnats        []CtyunDNatModel `tfsdk:"dnats"`
}

type CtyunDNatModel struct {
	CreateTime                types.String `tfsdk:"create_time"`                  //创建时间
	Description               types.String `tfsdk:"description"`                  //描述信息
	ID                        types.String `tfsdk:"id"`                           //dnatID 值
	DNatID                    types.String `tfsdk:"dnat_id"`                      //dnatID 值
	IpExpireTime              types.String `tfsdk:"ip_expire_time"`               //ip到期时间
	ExtendID                  types.String `tfsdk:"extend_id"`                    //弹性 IP id
	ExternalIp                types.String `tfsdk:"external_ip"`                  //弹性 IP 地址
	ExternalPort              types.Int64  `tfsdk:"external_port"`                //外部访问端口
	InternalPort              types.Int64  `tfsdk:"internal_port"`                //内部访问端口
	InternalIP                types.String `tfsdk:"internal_ip"`                  //内网 IP 地址
	Protocol                  types.String `tfsdk:"protocol"`                     //TCP:转发TCP协议的报文 UDP：转发UDP协议的报文
	State                     types.String `tfsdk:"state"`                        //运行状态: ACTIVE / FREEZING / CREATING
	VirtualMachineDisplayName types.String `tfsdk:"virtual_machine_display_name"` //虚拟机展示名称
	VirtualMachineId          types.String `tfsdk:"instance_id"`                  //虚拟机id
	VirtualMachineName        types.String `tfsdk:"virtual_machine_name"`         //虚拟机名称
}
