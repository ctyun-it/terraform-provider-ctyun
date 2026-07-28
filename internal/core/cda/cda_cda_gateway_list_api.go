package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaGatewayListApi
/* 查询账户下已创建的云专线(Cloud Dedicated Access)网关。 */
type CdaCdaGatewayListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaGatewayListApi(client *core.CtyunClient) *CdaCdaGatewayListApi {
	return &CdaCdaGatewayListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/gateway/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaGatewayListApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaGatewayListRequest) (*CdaCdaGatewayListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaGatewayListRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CdaCdaGatewayListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaGatewayListRequest struct {
	PageNo      int32   `json:"pageNo"`                /*  页数  */
	PageSize    int32   `json:"pageSize"`              /*  每页数据量  */
	Account     string  `json:"account"`               /*  天翼云客户邮箱  */
	ProjectID   *string `json:"projectID,omitempty"`   /*  项目ID  */
	RegionID    *string `json:"regionID,omitempty"`    /*  regionID  */
	GatewayName *string `json:"gatewayName,omitempty"` /*  专线网关名字  */
}

type CdaCdaGatewayListResponse struct {
	StatusCode            *int32                                `json:"statusCode"`            /*  返回状态码(800为成功，900为失败)  */
	Message               *string                               `json:"message"`               /*  失败时的错误描述，一般为英文描述  */
	Description           *string                               `json:"description"`           /*  失败时的错误描述，一般为中文描述  */
	ErrorCode             *string                               `json:"errorCode"`             /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail           *CdaCdaGatewayListErrorDetailResponse `json:"errorDetail"`           /*  错误明细  */
	Count                 *int32                                `json:"count"`                 /*  数量  */
	Fuid                  *string                               `json:"fuid"`                  /*  专线网关ID  */
	VrfName               *string                               `json:"vrfName"`               /*  专线网关名字  */
	AccessPoint           *string                               `json:"accessPoint"`           /*  接入点  */
	Account               *string                               `json:"account"`               /*  天翼云账号  */
	IsSwConfig            *bool                                 `json:"isSwConfig"`            /*  是否下发配置到交换机  */
	LineList              []*string                             `json:"lineList"`              /*  绑定的物理专线列表  */
	ResourcePool          *string                               `json:"resourcePool"`          /*  资源池ID  */
	ResourcePoolName      *string                               `json:"resourcePoolName"`      /*  资源池名字  */
	IsAutomation          *bool                                 `json:"isAutomation"`          /*  是否自动化  */
	ProjectIDcs           *string                               `json:"projectIDcs"`           /*  租户ID  */
	LgcreateTime          *string                               `json:"lgcreateTime"`          /*  创建时间  */
	FuserLastUpdated      *string                               `json:"fuserLastUpdated"`      /*  上次更新时间  */
	DeleteTime            *string                               `json:"deleteTime"`            /*  删除时间  */
	IpVersion             *string                               `json:"ipVersion"`             /*  本参数表示包周期类型。<br>取值范围：<br>IPV4<br>IPV6<br>DUALSTACK  */
	GatewayName           *string                               `json:"gatewayName"`           /*  专线网关名字  */
	DstCidr               []*string                             `json:"dstCidr"`               /*  目的IPV4地址列表  */
	DstCidrV6             []*string                             `json:"dstCidrV6"`             /*  目的IPV6地址列表  */
	SRId                  *string                               `json:"SRId"`                  /*  静态路由ID  */
	Srid                  *string                               `json:"srid"`                  /*  静态路由ID  */
	RemoteGatewayIp       *string                               `json:"remoteGatewayIp"`       /*  下一跳，即物理专线的远端互联ip  */
	Priority              *int32                                `json:"priority"`              /*  优先级  */
	Track                 *int32                                `json:"track"`                 /*  0为关闭，1为开启  */
	MultiPathNumberIpv6   *string                               `json:"multiPathNumberIpv6"`   /*  BGP-IPv6多路功能序号(负载线路数)  */
	NetworkCidr           []*string                             `json:"networkCidr"`           /*  客户侧子网列表(IPv4)  */
	MultiPathIpv6         *bool                                 `json:"multiPathIpv6"`         /*  是否开启BGP-IPv6多路功能  */
	NetworkCidrV6         []*string                             `json:"networkCidrV6"`         /*  客户侧子网列表(IPv6)  */
	MultiPathType         *string                               `json:"multiPathType"`         /*  Bgp多路功能类型(IBGP/EBGP)  */
	BGPId                 *string                               `json:"BGPId"`                 /*  BGP ID  */
	MultiPath             *bool                                 `json:"multiPath"`             /*  是否开启BGP-IPv4多路功能  */
	MultiPathNumber       *string                               `json:"multiPathNumber"`       /*  BGP-IPv4多路功能序号(负载线路数)  */
	BGPNeighbor           *string                               `json:"BGPNeighbor"`           /*  BGP邻居名称  */
	BGPIP                 *string                               `json:"BGPIP"`                 /*  BGP邻居IP(物理专线的远端互联IP)  */
	LineName              *string                               `json:"lineName"`              /*  物理专线名称  */
	PeerAS                *string                               `json:"peerAS"`                /*  Peer AS号  */
	BGPKey                *string                               `json:"BGPKey"`                /*  BGP密钥  */
	Bfd                   *bool                                 `json:"bfd"`                   /*  是否打开bfd 功能  */
	CdaId                 *string                               `json:"cdaId"`                 /*  专线ID  */
	VpcId                 *string                               `json:"vpcId"`                 /*  VPC  ID  */
	DcType                *string                               `json:"dcType"`                /*  本参数表示资源池类型。<br>取值范围：<br>MAZ<br>CNP  */
	AzNameIpv6            *string                               `json:"azNameIpv6"`            /*  IPv6 VPC的可用区名字  */
	CdaIdV6               *string                               `json:"cdaIdV6"`               /*  IPv6 VPC的专线ID  */
	VpcSubnet             *string                               `json:"vpcSubnet"`             /*  VPC子网  */
	EmazId                *string                               `json:"emazId"`                /*  多可用区ID  */
	VirtualBandwidth      *int32                                `json:"virtualBandwidth"`      /*  虚拟带宽  */
	VpcName               *string                               `json:"vpcName"`               /*  VPC名字  */
	EmazIdV6              *string                               `json:"emazIdV6"`              /*  多可用区ID(IPv6)  */
	GuestGatewayList      []*string                             `json:"guestGatewayList"`      /*  客户侧网关列表  */
	AzName                *string                               `json:"azName"`                /*  可用区名字  */
	VpcSubnetIpv6         *string                               `json:"vpcSubnetIpv6"`         /*  VPC IPv6子网  */
	VpcNetworkSegment     *string                               `json:"vpcNetworkSegment"`     /*  VPC网段  */
	CtUserId              *string                               `json:"ctUserId"`              /*  天翼云用户ID  */
	VpcNetworkSegmentIpv6 *string                               `json:"vpcNetworkSegmentIpv6"` /*  VPC IPv6网段  */
	GuestGatewayListIpv6  []*string                             `json:"guestGatewayListIpv6"`  /*  客户侧IPv6网关列表  */
	Error                 *string                               `json:"error"`                 /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaGatewayListErrorDetailResponse struct{}
