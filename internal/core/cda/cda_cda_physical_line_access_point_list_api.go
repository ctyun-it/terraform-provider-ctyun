package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaPhysicalLineAccessPointListApi
/* 查询物理专线的接入点 */
type CdaCdaPhysicalLineAccessPointListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaPhysicalLineAccessPointListApi(client *core.CtyunClient) *CdaCdaPhysicalLineAccessPointListApi {
	return &CdaCdaPhysicalLineAccessPointListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/physical-line/access-point-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaPhysicalLineAccessPointListApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaPhysicalLineAccessPointListRequest) (*CdaCdaPhysicalLineAccessPointListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaPhysicalLineAccessPointListRequest
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
	var resp CdaCdaPhysicalLineAccessPointListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaPhysicalLineAccessPointListRequest struct {
	LineName string `json:"lineName"` /*  物理专线名字  */
	Account  string `json:"account"`  /*  天翼云客户邮箱  */
}

type CdaCdaPhysicalLineAccessPointListResponse struct {
	StatusCode  *int32                                                `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                               `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                               `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ErrorCode   *string                                               `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaPhysicalLineAccessPointListErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	AccessPoint *string                                               `json:"accessPoint"` /*  接入点  */
	Error       *string                                               `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaPhysicalLineAccessPointListErrorDetailResponse struct{}
