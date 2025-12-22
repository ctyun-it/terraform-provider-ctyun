package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaDeleteAccountAuthApi
/* 删除已创建的跨账号授权。 */
type CdaCdaDeleteAccountAuthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaDeleteAccountAuthApi(client *core.CtyunClient) *CdaCdaDeleteAccountAuthApi {
	return &CdaCdaDeleteAccountAuthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/accountauth/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaDeleteAccountAuthApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaDeleteAccountAuthRequest) (*CdaCdaDeleteAccountAuthResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaDeleteAccountAuthRequest
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
	var resp CdaCdaDeleteAccountAuthResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaDeleteAccountAuthRequest struct {
	VpcID         string `json:"vpcID"`         /*  已授权的VPC ID  */
	GatewayName   string `json:"gatewayName"`   /*  已授权VPC给指定的专线网关  */
	AuthAccountId string `json:"authAccountId"` /*  授权对方账号ID  */
}

type CdaCdaDeleteAccountAuthResponse struct {
	StatusCode  *int32                                      `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaDeleteAccountAuthReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaDeleteAccountAuthErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
}

type CdaCdaDeleteAccountAuthReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为空  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为空  */
	ErrorCode *string `json:"errorCode"` /*  成功为空  */
}

type CdaCdaDeleteAccountAuthErrorDetailResponse struct{}
