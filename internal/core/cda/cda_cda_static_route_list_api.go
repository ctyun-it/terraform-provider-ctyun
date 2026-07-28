package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaStaticRouteListApi
/* 查询专线网关下的静态路由 */
type CdaCdaStaticRouteListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaStaticRouteListApi(client *core.CtyunClient) *CdaCdaStaticRouteListApi {
	return &CdaCdaStaticRouteListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/static-route/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaStaticRouteListApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaStaticRouteListRequest) (*CdaCdaStaticRouteListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaStaticRouteListRequest
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
	var resp CdaCdaStaticRouteListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaStaticRouteListRequest struct {
	GatewayName string `json:"gatewayName"` /*  专线网关名字  */
	Account     string `json:"account"`     /*  天翼云客户邮箱  */
}

type CdaCdaStaticRouteListResponse struct {
	StatusCode      *int32                                    `json:"statusCode"`      /*  返回状态码(800为成功，900为失败)  */
	Message         *string                                   `json:"message"`         /*  失败时的错误描述，一般为英文描述  */
	Description     *string                                   `json:"description"`     /*  失败时的错误描述，一般为中文描述  */
	ErrorCode       *string                                   `json:"errorCode"`       /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail     *CdaCdaStaticRouteListErrorDetailResponse `json:"errorDetail"`     /*  错误明细  */
	IpVersion       *string                                   `json:"ipVersion"`       /*  本参数表示包周期类型。<br>取值范围：<br>IPV4<br>IPV6<br>DUALSTACK  */
	GatewayName     *string                                   `json:"gatewayName"`     /*  专线网关名字  */
	DstCidr         []*string                                 `json:"dstCidr"`         /*  目的IPV4地址列表  */
	DstCidrV6       []*string                                 `json:"dstCidrV6"`       /*  目的IPV6地址列表  */
	SRID            *string                                   `json:"SRID"`            /*  静态路由ID  */
	RemoteGatewayIp *string                                   `json:"remoteGatewayIp"` /*  下一跳，即物理专线的远端互联ip  */
	Priority        *int32                                    `json:"priority"`        /*  优先级  */
	Track           *int32                                    `json:"track"`           /*  0为关闭，1为开启  */
	Error           *string                                   `json:"error"`           /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaStaticRouteListErrorDetailResponse struct{}
