package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunExpressConnectionRegionPeer{}
	_ datasource.DataSourceWithConfigure = &CtyunExpressConnectionRegionPeer{}
)

type CtyunExpressConnectionRegionPeer struct {
	meta *common.CtyunMetadata
}

func NewCtyunExpressConnectionRegionPeers() datasource.DataSource {
	return &CtyunExpressConnectionRegionPeer{}
}
func (c *CtyunExpressConnectionRegionPeer) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunExpressConnectionRegionPeer) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ec_region_peers"
}

func (c *CtyunExpressConnectionRegionPeer) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038250",
		Attributes: map[string]schema.Attribute{
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"packet_id": schema.StringAttribute{
				Optional:    true,
				Description: "带宽包ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cgw_id": schema.StringAttribute{
				Optional:    true,
				Description: "云网关ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"region_peers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "跨域连接ID",
						},
						"ec_id": schema.StringAttribute{
							Computed:    true,
							Description: "云间高速ID",
						},
						"packet_id": schema.StringAttribute{
							Computed:    true,
							Description: "带宽包ID",
						},
						"packet_name": schema.StringAttribute{
							Computed:    true,
							Description: "带宽包名称",
						},
						"peer_name": schema.StringAttribute{
							Computed:    true,
							Description: "跨域连接名称",
						},
						"src_cgw_id": schema.StringAttribute{
							Computed:    true,
							Description: "本端网关ID",
						},
						"dst_cgw_id": schema.StringAttribute{
							Computed:    true,
							Description: "对端网关ID",
						},
						"src_cgw_name": schema.StringAttribute{
							Computed:    true,
							Description: "本端网关名称",
						},
						"dst_cgw_name": schema.StringAttribute{
							Computed:    true,
							Description: "对端网关名称",
						},
						"src_region_id": schema.StringAttribute{
							Computed:    true,
							Description: "本端资源池ID",
						},
						"dst_region_id": schema.StringAttribute{
							Computed:    true,
							Description: "对端资源池ID",
						},
						"src_region_name": schema.StringAttribute{
							Computed:    true,
							Description: "本端资源池名称",
						},
						"dst_region_name": schema.StringAttribute{
							Computed:    true,
							Description: "对端资源池名称",
						},
						"peer_type": schema.Int32Attribute{
							Computed:    true,
							Description: "互通类型（1：境内，2：跨境（中国大陆-亚太），3：境外（亚太），4：定制）",
						},
						"rate": schema.Int32Attribute{
							Computed:    true,
							Description: "带宽值（MB）",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态（creating：加载中，running：已连接， removing：卸载中，expired：过期）",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间，为UTC格式",
						},
					},
				},
				Description: "跨域连接列表",
			},
		},
	}
}

func (c *CtyunExpressConnectionRegionPeer) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunEcRegionPeersConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &ec.EcRegionPeerListRequest{
		EcID: config.EcID.ValueString(),
	}
	if !config.PacketID.IsNull() {
		params.PacketID = config.PacketID.ValueStringPointer()
	}
	if !config.CgwID.IsNull() {
		params.CgwID = config.CgwID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkEcApis.EcRegionPeerListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询云间高速(id=%s)的跨域连接列表失败，接口返回nil，请联系研发确认问题原因！", config.EcID.ValueString())
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	var regionPeers []CtyunEcRegionPeerModel
	for _, peerItem := range resp.ReturnObj.Results {
		var regionPeer CtyunEcRegionPeerModel
		regionPeer.EcID = types.StringValue(*peerItem.EcID)
		regionPeer.ID = types.StringValue(*peerItem.PeerID)
		regionPeer.PacketID = types.StringValue(*peerItem.PacketID)
		regionPeer.PacketName = types.StringValue(*peerItem.PacketName)
		regionPeer.SrcCgwID = types.StringValue(*peerItem.SrcCgwID)
		regionPeer.DstCgwID = types.StringValue(*peerItem.DstCgwID)
		regionPeer.SrcCgwName = types.StringValue(*peerItem.SrcCgwName)
		regionPeer.DstCgwName = types.StringValue(*peerItem.DstCgwName)
		regionPeer.SrcRegionID = types.StringValue(*peerItem.SrcDcID)
		regionPeer.DstRegionID = types.StringValue(*peerItem.DstDcID)
		regionPeer.SrcRegionName = types.StringValue(*peerItem.SrcDcName)
		regionPeer.DstRegionName = types.StringValue(*peerItem.DstDcName)
		regionPeer.PeerType = types.Int32Value(*peerItem.PeerType)
		regionPeer.Rate = types.Int32Value(*peerItem.Rate)
		regionPeer.Status = types.StringValue(*peerItem.Status)
		regionPeer.UpdateTime = types.StringValue(utils.FromBJTimeToUTCZ(*peerItem.UpdateDate))

		regionPeers = append(regionPeers, regionPeer)
	}
	config.RegionPeers = regionPeers
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunEcRegionPeerModel struct {
	ID            types.String `tfsdk:"id"`
	EcID          types.String `tfsdk:"ec_id"`
	PacketID      types.String `tfsdk:"packet_id"`
	PacketName    types.String `tfsdk:"packet_name"`
	PeerName      types.String `tfsdk:"peer_name"`
	SrcCgwID      types.String `tfsdk:"src_cgw_id"`
	DstCgwID      types.String `tfsdk:"dst_cgw_id"`
	SrcCgwName    types.String `tfsdk:"src_cgw_name"`
	DstCgwName    types.String `tfsdk:"dst_cgw_name"`
	SrcRegionID   types.String `tfsdk:"src_region_id"`
	DstRegionID   types.String `tfsdk:"dst_region_id"`
	SrcRegionName types.String `tfsdk:"src_region_name"`
	DstRegionName types.String `tfsdk:"dst_region_name"`
	PeerType      types.Int32  `tfsdk:"peer_type"`
	Rate          types.Int32  `tfsdk:"rate"`
	Status        types.String `tfsdk:"status"`
	UpdateTime    types.String `tfsdk:"update_time"`
}

type CtyunEcRegionPeersConfig struct {
	EcID        types.String             `tfsdk:"ec_id"`
	PacketID    types.String             `tfsdk:"packet_id"`
	CgwID       types.String             `tfsdk:"cgw_id"`
	RegionPeers []CtyunEcRegionPeerModel `tfsdk:"region_peers"`
}
