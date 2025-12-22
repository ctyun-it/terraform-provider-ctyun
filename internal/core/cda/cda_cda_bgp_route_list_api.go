package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaBgpRouteListApi
/* 查询用户专线网关下的BGP动态路由 */
type CdaCdaBgpRouteListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaBgpRouteListApi(client *core.CtyunClient) *CdaCdaBgpRouteListApi {
	return &CdaCdaBgpRouteListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/bgp-route/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaBgpRouteListApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaBgpRouteListRequest) (*CdaCdaBgpRouteListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaBgpRouteListRequest
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
	var resp CdaCdaBgpRouteListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaBgpRouteListRequest struct {
	GatewayName string `json:"gatewayName"` /*  专线网关名字  */
	Account     string `json:"account"`     /*  天翼云客户邮箱  */
}

type CdaCdaBgpRouteListResponse struct {
	StatusCode          *int32                                 `json:"statusCode"`          /*  返回状态码(800为成功，900为失败)  */
	Message             *string                                `json:"message"`             /*  失败时的错误描述，一般为英文描述  */
	Description         *string                                `json:"description"`         /*  失败时的错误描述，一般为中文描述  */
	ErrorCode           *string                                `json:"errorCode"`           /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail         *CdaCdaBgpRouteListErrorDetailResponse `json:"errorDetail"`         /*  错误明细  */
	MultiPathNumberIpv6 *string                                `json:"multiPathNumberIpv6"` /*  BGP-IPv6多路功能序号(负载线路数)  */
	NetworkCidr         []*string                              `json:"networkCidr"`         /*  客户侧子网列表(IPv4)  */
	GatewayName         *string                                `json:"gatewayName"`         /*  专线网关名字  */
	MultiPathIpv6       *bool                                  `json:"multiPathIpv6"`       /*  是否开启BGP-IPv6多路功能  */
	NetworkCidrV6       []*string                              `json:"networkCidrV6"`       /*  客户侧子网列表(IPv6)  */
	MultiPathType       *string                                `json:"multiPathType"`       /*  Bgp多路功能类型(IBGP/EBGP)  */
	IpVersion           *string                                `json:"ipVersion"`           /*  本参数表示包周期类型。<br>取值范围：<br>IPV4<br>IPV6<br>DUALSTACK  */
	BGPID               *string                                `json:"BGPID"`               /*  BGP ID  */
	MultiPath           *bool                                  `json:"multiPath"`           /*  是否开启BGP-IPv4多路功能  */
	MultiPathNumber     *string                                `json:"multiPathNumber"`     /*  BGP-IPv4多路功能序号(负载线路数)  */
	BGPNeighbor         *string                                `json:"BGPNeighbor"`         /*  BGP邻居名称  */
	BGPIP               *string                                `json:"BGPIP"`               /*  BGP邻居IP(物理专线的远端互联IP)  */
	LineName            *string                                `json:"lineName"`            /*  物理专线名称  */
	PeerAS              *string                                `json:"peerAS"`              /*  Peer AS号  */
	BGPKey              *string                                `json:"BGPKey"`              /*  BGP密钥  */
	Bfd                 *bool                                  `json:"bfd"`                 /*  是否打开bfd 功能  */
	Error               *string                                `json:"error"`               /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaBgpRouteListErrorDetailResponse struct{}
