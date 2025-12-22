package peer_connection

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CtyunVpcPeerConnections struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpcPeerConnections() datasource.DataSource {
	return &CtyunVpcPeerConnections{}
}

func (c *CtyunVpcPeerConnections) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunVpcPeerConnections) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_peer_connections"
}

func (c *CtyunVpcPeerConnections) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026760/10037873",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Description: "资源池id",
				Optional:    true,
			},
			"page_size": schema.Int32Attribute{
				Description: "当前页数据条数，默认为10，最大值为50",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"page_no": schema.Int32Attribute{
				Description: "页码，默认为1",
				Optional:    true,
			},
			"peer_connections": schema.ListNestedAttribute{
				Description: "对等连接列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "对等连接名称",
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Description: "对等连接 ID",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "对等连接状态：agree(同意建立对等连接) / pending(待同意建立对等链接)",
							Computed:    true,
						},
						"request_vpc_id": schema.StringAttribute{
							Description: "本端vpc id ",
							Computed:    true,
						},
						"request_vpc_cidr": schema.StringAttribute{
							Description: "本端vpc cidr",
							Computed:    true,
						},
						"request_vpc_name": schema.StringAttribute{
							Description: "本端vpc的名称",
							Computed:    true,
						},
						"accept_vpc_id": schema.StringAttribute{
							Description: "对端的vpc id",
							Computed:    true,
						},
						"accept_vpc_cidr": schema.StringAttribute{
							Description: "对端vpc的cidr",
							Computed:    true,
						},
						"accept_vpc_name": schema.StringAttribute{
							Description: "对端vpc的名称",
							Computed:    true,
						},
						"user_type": schema.StringAttribute{
							Description: "对等连接类型：current(同一个租户) / other(不同租户)",
							Computed:    true,
						},
						"accept_email": schema.StringAttribute{
							Description: "对端vpc账户的邮箱",
							Computed:    true,
						},
						"create_time": schema.StringAttribute{
							Description: "创建时间，为UTC格式",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (c *CtyunVpcPeerConnections) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunVpcPeerConnectionsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	config.RegionID = types.StringValue(regionId)
	params := &ctvpc.CtvpcListVpcPeerConnectionRequest{
		RegionID:   regionId,
		PageSize:   config.PageSize.ValueInt32(),
		PageNumber: config.PageNo.ValueInt32(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListVpcPeerConnectionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询对待连接列表失败，接口返回值为nil，请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 遍历读取vpc对等连接
	var peerConnections []CtyunVpcPeerConnectionModel
	for _, peerConnItem := range resp.ReturnObj {
		var peerConnection CtyunVpcPeerConnectionModel
		peerConnection.Name = types.StringValue(*peerConnItem.Name)
		peerConnection.ID = types.StringValue(*peerConnItem.InstanceID)
		peerConnection.Status = types.StringValue(*peerConnItem.Status)
		peerConnection.RequestVpcID = types.StringValue(*peerConnItem.RequestVpcID)
		peerConnection.RequestVpcCidr = types.StringValue(*peerConnItem.RequestVpcCidr)
		peerConnection.RequestVpcName = types.StringValue(*peerConnItem.RequestVpcName)
		peerConnection.AcceptVpcID = types.StringValue(*peerConnItem.AcceptVpcID)
		peerConnection.AcceptVpcName = types.StringValue(*peerConnItem.AcceptVpcName)
		peerConnection.AcceptVpcCidr = types.StringValue(*peerConnItem.AcceptVpcCidr)
		peerConnection.UserType = types.StringValue(*peerConnItem.UserType)
		peerConnection.AcceptEmail = types.StringValue(*peerConnItem.AcceptEmail)
		peerConnection.CreateTime = types.StringValue(*peerConnItem.CreationTime)
		peerConnections = append(peerConnections, peerConnection)
	}
	config.PeerConnections = peerConnections
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunVpcPeerConnectionModel struct {
	Name           types.String `tfsdk:"name"`
	ID             types.String `tfsdk:"id"`
	Status         types.String `tfsdk:"status"`
	RequestVpcID   types.String `tfsdk:"request_vpc_id"`
	RequestVpcCidr types.String `tfsdk:"request_vpc_cidr"`
	RequestVpcName types.String `tfsdk:"request_vpc_name"`
	AcceptVpcID    types.String `tfsdk:"accept_vpc_id"`
	AcceptVpcCidr  types.String `tfsdk:"accept_vpc_cidr"`
	AcceptVpcName  types.String `tfsdk:"accept_vpc_name"`
	UserType       types.String `tfsdk:"user_type"`
	AcceptEmail    types.String `tfsdk:"accept_email"`
	CreateTime     types.String `tfsdk:"create_time"`
}

type CtyunVpcPeerConnectionsConfig struct {
	RegionID        types.String                  `tfsdk:"region_id"`
	PageSize        types.Int32                   `tfsdk:"page_size"`
	PageNo          types.Int32                   `tfsdk:"page_no"`
	PeerConnections []CtyunVpcPeerConnectionModel `tfsdk:"peer_connections"`
}
