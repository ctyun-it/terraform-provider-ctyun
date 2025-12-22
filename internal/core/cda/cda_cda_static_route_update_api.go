package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaStaticRouteUpdateApi
/* 更新已有的静态路由 */
type CdaCdaStaticRouteUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaStaticRouteUpdateApi(client *core.CtyunClient) *CdaCdaStaticRouteUpdateApi {
	return &CdaCdaStaticRouteUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/static-route/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaStaticRouteUpdateApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaStaticRouteUpdateRequest) (*CdaCdaStaticRouteUpdateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaStaticRouteUpdateRequest
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
	var resp CdaCdaStaticRouteUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaStaticRouteUpdateRequest struct {
	SRID      string                                     `json:"SRID"`                /*  静态路由ID  */
	IpVersion string                                     `json:"ipVersion"`           /*  本参数表示包周期类型。<br>取值范围：<br>IPV4<br>IPV6<br>DUALSTACK  */
	DstCidr   []*string                                  `json:"dstCidr,omitempty"`   /*  目的IPV4地址列表(全量传入)  */
	DstCidrV6 []*string                                  `json:"dstCidrV6,omitempty"` /*  目的IPV6地址列表(全量传入)  */
	NextHop   []*CdaCdaStaticRouteUpdateNextHopRequest   `json:"nextHop,omitempty"`   /*  下一跳及优先级列表(全量传入)  */
	NextHopV6 []*CdaCdaStaticRouteUpdateNextHopV6Request `json:"nextHopV6,omitempty"` /*  下一跳及优先级列表(全量传入)  */
}

type CdaCdaStaticRouteUpdateNextHopRequest struct {
	RemoteGatewayIp string `json:"remoteGatewayIp"` /*  下一跳，即物理专线的远端互联ip  */
	Priority        int32  `json:"priority"`        /*  优先级  */
	Track           *int32 `json:"track,omitempty"` /*  0为关闭，1为开启  */
	Bfd             *bool  `json:"bfd,omitempty"`   /*  是否开启bfd功能，ture为开启，false为关闭  */
}

type CdaCdaStaticRouteUpdateNextHopV6Request struct {
	RemoteGatewayIp string `json:"remoteGatewayIp"` /*  下一跳，即物理专线的远端互联ip  */
	Priority        int32  `json:"priority"`        /*  优先级  */
	Track           *int32 `json:"track,omitempty"` /*  0为关闭，1为开启  */
	Bfd             *bool  `json:"bfd,omitempty"`   /*  是否开启bfd功能，ture为开启，false为关闭  */
}

type CdaCdaStaticRouteUpdateResponse struct {
	StatusCode  *int32                                      `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaStaticRouteUpdateReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaStaticRouteUpdateErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaStaticRouteUpdateReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为空  */
	ErrorCode *string `json:"errorCode"` /*  错误代码，成功为空  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为空  */
}

type CdaCdaStaticRouteUpdateErrorDetailResponse struct{}
