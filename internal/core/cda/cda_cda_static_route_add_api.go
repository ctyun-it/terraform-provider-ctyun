package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaStaticRouteAddApi
/* 创建静态路由 */
type CdaCdaStaticRouteAddApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaStaticRouteAddApi(client *core.CtyunClient) *CdaCdaStaticRouteAddApi {
	return &CdaCdaStaticRouteAddApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/static-route/add",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaStaticRouteAddApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaStaticRouteAddRequest) (*CdaCdaStaticRouteAddResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaStaticRouteAddRequest
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
	var resp CdaCdaStaticRouteAddResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaStaticRouteAddRequest struct {
	IpVersion   string                                  `json:"ipVersion"`           /*  本参数表示包周期类型。<br>取值范围：<br>IPV4<br>IPV6<br>DUALSTACK  */
	GatewayName string                                  `json:"gatewayName"`         /*  专线网关名字  */
	DstCidr     []*string                               `json:"dstCidr,omitempty"`   /*  目的IPV4地址列表  */
	DstCidrV6   []*string                               `json:"dstCidrV6,omitempty"` /*  目的IPV6地址列表  */
	NextHop     []*CdaCdaStaticRouteAddNextHopRequest   `json:"nextHop,omitempty"`   /*  下一跳及优先级列表  */
	NextHopV6   []*CdaCdaStaticRouteAddNextHopV6Request `json:"nextHopV6,omitempty"` /*  下一跳及优先级列表  */
}

type CdaCdaStaticRouteAddNextHopRequest struct {
	RemoteGatewayIp string `json:"remoteGatewayIp"` /*  下一跳，即物理专线的远端互联ip  */
	Priority        int32  `json:"priority"`        /*  优先级  */
	Track           *int32 `json:"track,omitempty"` /*  0为关闭，1为开启  */
	Bfd             *bool  `json:"bfd,omitempty"`   /*  是否开启bfd功能，，ture为开启，false为关闭  */
}

type CdaCdaStaticRouteAddNextHopV6Request struct {
	RemoteGatewayIp string `json:"remoteGatewayIp"` /*  下一跳，即物理专线的远端互联ip  */
	Priority        int32  `json:"priority"`        /*  优先级  */
	Track           *int32 `json:"track,omitempty"` /*  0为关闭，1为开启  */
	Bfd             *bool  `json:"bfd,omitempty"`   /*  是否开启bfd功能，，ture为开启，false为关闭  */
}

type CdaCdaStaticRouteAddResponse struct {
	StatusCode  *int32                                   `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaStaticRouteAddReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaStaticRouteAddErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaStaticRouteAddReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为空  */
	ErrorCode *string `json:"errorCode"` /*  错误代码，成功为空  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为空  */
}

type CdaCdaStaticRouteAddErrorDetailResponse struct{}
