package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaBgpRouteUpdateApi
/* 更新BGP动态路由 */
type CdaCdaBgpRouteUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaBgpRouteUpdateApi(client *core.CtyunClient) *CdaCdaBgpRouteUpdateApi {
	return &CdaCdaBgpRouteUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/bgp-route/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaBgpRouteUpdateApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaBgpRouteUpdateRequest) (*CdaCdaBgpRouteUpdateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaBgpRouteUpdateRequest
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
	var resp CdaCdaBgpRouteUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaBgpRouteUpdateRequest struct {
	BGPID            string                                    `json:"BGPID"`                      /*  BGP路由ID  */
	IpVersion        string                                    `json:"ipVersion"`                  /*  本参数表示包周期类型。<br>取值范围：<br>IPV4<br>IPV6<br>DUALSTACK  */
	NetworkCidr      []*string                                 `json:"networkCidr,omitempty"`      /*  客户侧子网列表(IPv4)  */
	NetworkCidrV6    []*string                                 `json:"networkCidrV6,omitempty"`    /*  客户侧子网列表(IPv6)  */
	BGPList          []*CdaCdaBgpRouteUpdateBGPListRequest     `json:"BGPList,omitempty"`          /*  IPv4类型的BGP列表  */
	BGPIpv6List      []*CdaCdaBgpRouteUpdateBGPIpv6ListRequest `json:"BGPIpv6List,omitempty"`      /*  IPv6类型的BGP列表  */
	MultiPath        *bool                                     `json:"multiPath,omitempty"`        /*  是否开启多路功能  */
	MultiPathNum     *string                                   `json:"multiPathNum,omitempty"`     /*  Bgp多路功能序号(负载线路数)  */
	MultiPathType    *string                                   `json:"multiPathType,omitempty"`    /*  Bgp多路功能类型(IBGP/EBGP)  */
	MultiPathIpv6    *bool                                     `json:"multiPathIpv6,omitempty"`    /*  是否开启BGP-IPv6多路功能  */
	MultiPathNumIpv6 *string                                   `json:"multiPathNumIpv6,omitempty"` /*  BGP-IPv6多路功能序号(负载线路数)  */
}

type CdaCdaBgpRouteUpdateBGPListRequest struct {
	BGPNeighbor   string    `json:"BGPNeighbor"`             /*  BGP邻居名称  */
	BGPIP         string    `json:"BGPIP"`                   /*  BGP邻居IP(物理专线的远端互联IP)  */
	LineID        string    `json:"lineID"`                  /*  物理专线ID  */
	PeerAS        string    `json:"peerAS"`                  /*  Peer AS号  */
	BGPKey        *string   `json:"BGPKey,omitempty"`        /*  BGP密钥  */
	Bfd           bool      `json:"bfd"`                     /*  是否打开bfd 功能  */
	NetworkCidr   []*string `json:"networkCidr,omitempty"`   /*  BGP下的客户侧子网列表(IPv4)  */
	NetworkCidrV6 []*string `json:"networkCidrV6,omitempty"` /*  BGP下的客户侧子网列表(IPv6)  */
}

type CdaCdaBgpRouteUpdateBGPIpv6ListRequest struct {
	BGPNeighbor   string    `json:"BGPNeighbor"`             /*  BGP邻居名称  */
	BGPIP         string    `json:"BGPIP"`                   /*  BGP邻居IP(物理专线的远端互联IP)  */
	LineID        string    `json:"lineID"`                  /*  物理专线ID  */
	PeerAS        string    `json:"peerAS"`                  /*  Peer AS号  */
	BGPKey        *string   `json:"BGPKey,omitempty"`        /*  BGP密钥  */
	Bfd           bool      `json:"bfd"`                     /*  是否打开bfd 功能  */
	NetworkCidr   []*string `json:"networkCidr,omitempty"`   /*  BGP下的客户侧子网列表(IPv4)  */
	NetworkCidrV6 []*string `json:"networkCidrV6,omitempty"` /*  BGP下的客户侧子网列表(IPv6)  */
}

type CdaCdaBgpRouteUpdateResponse struct {
	StatusCode  *int32                                   `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaBgpRouteUpdateReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaBgpRouteUpdateErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaBgpRouteUpdateReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为空  */
	ErrorCode *string `json:"errorCode"` /*  错误代码，成功为空  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为空  */
}

type CdaCdaBgpRouteUpdateErrorDetailResponse struct{}
