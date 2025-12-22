package nat

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctnat"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPrivateNatTransitIpsDatasource{}
	_ datasource.DataSourceWithConfigure = &ctyunPrivateNatTransitIpsDatasource{}
)

type ctyunPrivateNatTransitIpsDatasource struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrivateNatTransitIps() datasource.DataSource {
	return &ctyunPrivateNatTransitIpsDatasource{}
}

func (c *ctyunPrivateNatTransitIpsDatasource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_nat_transit_ips"
}

func (c *ctyunPrivateNatTransitIpsDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
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
			"transit_ips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "中转IP列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Computed:    true,
							Description: "中转IP地址",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "中转IP状态: running代表运行中, freeze代表已冻结, expired代表已到期",
						},
						"is_default": schema.BoolAttribute{
							Computed:    true,
							Description: "是否为默认中转地址",
						},
						"snat_count": schema.Int64Attribute{
							Computed:    true,
							Description: "在使用此中转IP的snat数量",
						},
						"dnat_count": schema.Int64Attribute{
							Computed:    true,
							Description: "在使用此中转IP的dnat数量",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunPrivateNatTransitIpsDatasource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunPrivateNatTransitIpsConfig
	// 读取请求信息
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	// region_id不能为空
	if regionId == "" {
		msg := "regionID不能为空"
		response.Diagnostics.AddError(msg, msg)
		return
	}
	natGatewayId := config.NatGatewayID.ValueString()
	params := &ctnat.CtnatQueryPrivatenatIPRequest{
		RegionID:     regionId,
		NatGatewayID: natGatewayId,
	}

	resp, err := c.meta.Apis.SdkCtNatApis.CtnatQueryPrivatenatIPApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var transitIps []CtyunPrivateNatTransitIpModel
	for _, ip := range resp.ReturnObj {
		ipItem := CtyunPrivateNatTransitIpModel{
			Address:   types.StringValue(ip.Address),
			Status:    types.StringValue(ip.Status),
			SnatCount: types.Int64Value(int64(ip.SnatCnt)),
			DnatCount: types.Int64Value(int64(ip.DnarCnt)),
		}
		if ip.IsDefault != nil {
			ipItem.IsDefault = types.BoolValue(*ip.IsDefault)
		}

		transitIps = append(transitIps, ipItem)
	}

	config.RegionID = types.StringValue(regionId)
	config.NatGatewayID = types.StringValue(natGatewayId)
	config.TransitIps = transitIps
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunPrivateNatTransitIpsDatasource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunPrivateNatTransitIpsConfig struct {
	RegionID     types.String                    `tfsdk:"region_id"`
	NatGatewayID types.String                    `tfsdk:"nat_gateway_id"`
	TransitIps   []CtyunPrivateNatTransitIpModel `tfsdk:"transit_ips"`
}

type CtyunPrivateNatTransitIpModel struct {
	Address   types.String `tfsdk:"address"`    // 中转IP地址
	Status    types.String `tfsdk:"status"`     // 中转IP状态
	IsDefault types.Bool   `tfsdk:"is_default"` // 是否为默认中转地址
	SnatCount types.Int64  `tfsdk:"snat_count"` // 在使用此中转IP的snat数量
	DnatCount types.Int64  `tfsdk:"dnat_count"` // 在使用此中转IP的dnat数量
}
