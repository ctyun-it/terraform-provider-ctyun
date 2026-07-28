package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaStaticsAccountAuthApi
/* 统计账号下已授权的VPC及授权给专线网关数量。 */
type CdaCdaStaticsAccountAuthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaStaticsAccountAuthApi(client *core.CtyunClient) *CdaCdaStaticsAccountAuthApi {
	return &CdaCdaStaticsAccountAuthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/accountauth/statistics",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaStaticsAccountAuthApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaStaticsAccountAuthRequest) (*CdaCdaStaticsAccountAuthResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CdaCdaStaticsAccountAuthResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaStaticsAccountAuthRequest struct{}

type CdaCdaStaticsAccountAuthResponse struct {
	StatusCode  *int32                                       `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                      `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                      `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaStaticsAccountAuthReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                      `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaStaticsAccountAuthErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
}

type CdaCdaStaticsAccountAuthReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为空  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为空  */
	ErrorCode *string `json:"errorCode"` /*  成功为空  */
}

type CdaCdaStaticsAccountAuthErrorDetailResponse struct{}
