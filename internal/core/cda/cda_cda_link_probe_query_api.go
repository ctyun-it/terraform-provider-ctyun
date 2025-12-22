package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaLinkProbeQueryApi
/* 展示指定vrf下的所有Ping测历史数据 */
type CdaCdaLinkProbeQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaLinkProbeQueryApi(client *core.CtyunClient) *CdaCdaLinkProbeQueryApi {
	return &CdaCdaLinkProbeQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/link-probe/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaLinkProbeQueryApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaLinkProbeQueryRequest) (*CdaCdaLinkProbeQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaLinkProbeQueryRequest
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
	var resp CdaCdaLinkProbeQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaLinkProbeQueryRequest struct {
	GatewayName string `json:"gatewayName"` /*  专线网关名字  */
}

type CdaCdaLinkProbeQueryResponse struct {
	StatusCode  *int32                                   `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaLinkProbeQueryErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Result      *string                                  `json:"result"`      /*  1成功， 0失败  */
	Data        *string                                  `json:"data"`        /*  成功为空  */
	ErrorMsg    *string                                  `json:"errorMsg"`    /*  成功为空  */
	TraceId     *string                                  `json:"traceId"`     /*  日志跟踪ID  */
}

type CdaCdaLinkProbeQueryErrorDetailResponse struct{}
