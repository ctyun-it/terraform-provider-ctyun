package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaVpcCountApi
/* 查询用户专线网关下的VPC数量 */
type CdaCdaVpcCountApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaVpcCountApi(client *core.CtyunClient) *CdaCdaVpcCountApi {
	return &CdaCdaVpcCountApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/vpc/count",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaVpcCountApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaVpcCountRequest) (*CdaCdaVpcCountResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaVpcCountRequest
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
	var resp CdaCdaVpcCountResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaVpcCountRequest struct {
	GatewayName string `json:"gatewayName"` /*  专线网关名字  */
	Account     string `json:"account"`     /*  天翼云客户邮箱  */
}

type CdaCdaVpcCountResponse struct {
	StatusCode  *int32                             `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                            `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                            `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaVpcCountReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                            `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaVpcCountErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                            `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaVpcCountReturnObjResponse struct {
	Count *int32 `json:"count"` /*  VPC数量  */
}

type CdaCdaVpcCountErrorDetailResponse struct{}
