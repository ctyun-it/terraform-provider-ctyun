package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaHealthCheckStatusQueryApi
/* 健康检查查询检查结果 */
type CdaCdaHealthCheckStatusQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaHealthCheckStatusQueryApi(client *core.CtyunClient) *CdaCdaHealthCheckStatusQueryApi {
	return &CdaCdaHealthCheckStatusQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/health-check/status/get",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaHealthCheckStatusQueryApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaHealthCheckStatusQueryRequest) (*CdaCdaHealthCheckStatusQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaHealthCheckStatusQueryRequest
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
	var resp CdaCdaHealthCheckStatusQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaHealthCheckStatusQueryRequest struct {
	RegionID     string `json:"regionID"`     /*  资源池ID  */
	ResourcePool string `json:"resourcePool"` /*  资源池ID  */
	GatewayName  string `json:"gatewayName"`  /*  专线网关名字  */
	VpcID        string `json:"vpcID"`        /*  VPC ID  */
}

type CdaCdaHealthCheckStatusQueryResponse struct {
	StatusCode  *int32                                           `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                          `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                          `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ErrorCode   *string                                          `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaHealthCheckStatusQueryErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Result      *string                                          `json:"result"`      /*  1成功， 0失败  */
	ErrorMsg    *string                                          `json:"errorMsg"`    /*  成功为空  */
	Status      *string                                          `json:"status"`      /*  状态中文描述  */
	VpcID       *string                                          `json:"vpcID"`       /*  VPC ID  */
	CtUserID    *string                                          `json:"ctUserID"`    /*  天翼云账号ID  */
	GatewayName *string                                          `json:"gatewayName"` /*  专线网关名字  */
}

type CdaCdaHealthCheckStatusQueryErrorDetailResponse struct{}
