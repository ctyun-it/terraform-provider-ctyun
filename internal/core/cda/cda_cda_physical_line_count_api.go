package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaPhysicalLineCountApi
/* 查询用户已创建的物理专线数量 */
type CdaCdaPhysicalLineCountApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaPhysicalLineCountApi(client *core.CtyunClient) *CdaCdaPhysicalLineCountApi {
	return &CdaCdaPhysicalLineCountApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/physical-line/count",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaPhysicalLineCountApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaPhysicalLineCountRequest) (*CdaCdaPhysicalLineCountResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaPhysicalLineCountRequest
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
	var resp CdaCdaPhysicalLineCountResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaPhysicalLineCountRequest struct {
	Account string `json:"account"` /*  天翼云客户邮箱  */
}

type CdaCdaPhysicalLineCountResponse struct {
	StatusCode  *int32                                      `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaPhysicalLineCountErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Count       *int32                                      `json:"count"`       /*  物理专线数量  */
	Error       *string                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaPhysicalLineCountErrorDetailResponse struct{}
