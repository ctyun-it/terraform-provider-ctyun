package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaHealthCheckQueryApi
/* 专线网关查询健康检查设置项 */
type CdaCdaHealthCheckQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaHealthCheckQueryApi(client *core.CtyunClient) *CdaCdaHealthCheckQueryApi {
	return &CdaCdaHealthCheckQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/health-check/get",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaHealthCheckQueryApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaHealthCheckQueryRequest) (*CdaCdaHealthCheckQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaHealthCheckQueryRequest
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
	var resp CdaCdaHealthCheckQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaHealthCheckQueryRequest struct {
	GatewayName string `json:"gatewayName"` /*  专线网关名字  */
	VpcID       string `json:"vpcID"`       /*  VPC ID  */
	VpcName     string `json:"vpcName"`     /*  VPC名字  */
}

type CdaCdaHealthCheckQueryResponse struct {
	StatusCode  *int32                                     `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                    `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                    `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ErrorCode   *string                                    `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaHealthCheckQueryErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Result      *string                                    `json:"result"`      /*  1成功， 0失败  */
	Data        *string                                    `json:"data"`        /*  成功为null  */
	ErrorMsg    *string                                    `json:"errorMsg"`    /*  成功为null  */
	TraceId     *string                                    `json:"traceId"`     /*  日志跟踪ID  */
}

type CdaCdaHealthCheckQueryErrorDetailResponse struct{}
