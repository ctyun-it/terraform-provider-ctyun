package nat

import (
	"context"
	"fmt"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctnat"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPrivateDnatDatasource{}
	_ datasource.DataSourceWithConfigure = &ctyunPrivateDnatDatasource{}
)

type ctyunPrivateDnatDatasource struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrivateDnats() datasource.DataSource {
	return &ctyunPrivateDnatDatasource{}
}

func (c *ctyunPrivateDnatDatasource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_nat_dnats"
}

func (c *ctyunPrivateDnatDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
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
				Description: "要查询的私网NAT网关的ID",
			},
			"dnats": schema.ListNestedAttribute{
				Computed:    true,
				Description: "私网dnats列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
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
						"external_ip": schema.StringAttribute{
							Computed:    true,
							Description: "中转 IP 地址",
						},
						"external_port": schema.Int64Attribute{
							Computed:    true,
							Description: "外部访问端口",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
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
						"port_id": schema.StringAttribute{
							Computed:    true,
							Description: "网卡ID",
						},
						"port_name": schema.StringAttribute{
							Computed:    true,
							Description: "网卡名称",
						},
						"device_id": schema.StringAttribute{
							Computed:    true,
							Description: "网卡对应的设备ID",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "协议: tcp/udp",
							Validators: []validator.String{
								stringvalidator.OneOf("tcp", "udp"),
							},
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "运行状态: running代表运行中, freeze代表已冻结, expired代表已到期",
							Validators: []validator.String{
								stringvalidator.OneOf("running", "freeze", "expired"),
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunPrivateDnatDatasource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunDataPrivateDnatConfig
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
	params := &ctnat.CtnatQueryPrivatenatDnatRequest{
		RegionID:     regionId,
		NatGatewayID: natGatewayId,
		PageNo:       1,
		PageSize:     50,
	}

	resp, err := c.meta.Apis.SdkCtNatApis.CtnatQueryPrivatenatDnatApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var dnats []CtyunPrivateDnatModel
	for _, dnat := range resp.ReturnObj {
		dnatItem := CtyunPrivateDnatModel{
			ID:           types.StringValue(dnat.DnatID),
			DNatID:       types.StringValue(dnat.DnatID),
			Description:  types.StringValue(dnat.Description),
			CreatedAt:    types.StringValue(dnat.CreatedAt),
			ExternalIP:   types.StringValue(dnat.ExternalIP),
			InternalIP:   types.StringValue(dnat.InternalIP),
			PortID:       types.StringValue(dnat.PortID),
			PortName:     types.StringValue(dnat.PortName),
			DeviceID:     types.StringValue(dnat.DeviceID),
			Protocol:     types.StringValue(dnat.Protocol),
			State:        types.StringValue(dnat.State),
			ExternalPort: types.Int64Value(int64(dnat.ExternalPort)),
			InternalPort: types.Int64Value(int64(dnat.InternalPort)),
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

func (c *ctyunPrivateDnatDatasource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunDataPrivateDnatConfig struct {
	RegionID     types.String            `tfsdk:"region_id"`
	NatGateWayID types.String            `tfsdk:"nat_gateway_id"`
	Dnats        []CtyunPrivateDnatModel `tfsdk:"dnats"`
}

type CtyunPrivateDnatModel struct {
	ID           types.String `tfsdk:"id"`            //dnatID 值
	DNatID       types.String `tfsdk:"dnat_id"`       //dnatID 值
	ExternalIP   types.String `tfsdk:"external_ip"`   /*  中转IP  */
	ExternalPort types.Int64  `tfsdk:"external_port"` /*  外部端口  */
	InternalIP   types.String `tfsdk:"internal_ip"`   /*  内部IP  */
	InternalPort types.Int64  `tfsdk:"internal_port"` /*  内部端口  */
	PortID       types.String `tfsdk:"port_id"`       /*  对应的网卡ID  */
	PortName     types.String `tfsdk:"port_name"`     /*  网卡名称  */
	DeviceID     types.String `tfsdk:"device_id"`     /*  网卡对应的设备ID  */
	Protocol     types.String `tfsdk:"protocol"`      /*  协议: tcp/udp  */
	State        types.String `tfsdk:"state"`         /*  DNAT状态: running代表运行中, freeze代表已冻结, expired代表已到期  */
	CreatedAt    types.String `tfsdk:"create_time"`   /*  创建时间  */
	Description  types.String `tfsdk:"description"`   /*  描述  */
}
