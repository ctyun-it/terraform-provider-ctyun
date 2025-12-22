package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaCreateAccountAuthApi
/* 云专线支持跨账号VPC互通，需要先创建跨账号的VPC授权。 */
type CdaCdaCreateAccountAuthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaCreateAccountAuthApi(client *core.CtyunClient) *CdaCdaCreateAccountAuthApi {
	return &CdaCdaCreateAccountAuthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/accountauth/add",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaCreateAccountAuthApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaCreateAccountAuthRequest) (*CdaCdaCreateAccountAuthResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaCreateAccountAuthRequest
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
	var resp CdaCdaCreateAccountAuthResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaCreateAccountAuthRequest struct {
	RegionID      string  `json:"regionID"`              /*  资源池ID  */
	Account       string  `json:"account"`               /*  自己账号邮箱  */
	VpcID         string  `json:"vpcID"`                 /*  授权VPC ID  */
	VpcName       string  `json:"vpcName"`               /*  授权VPC Name  */
	GatewayName   string  `json:"gatewayName"`           /*  授权VPC给对方指定的专线网关  */
	AuthAccountId string  `json:"authAccountId"`         /*  授权对方账号ID  */
	AuthAccount   string  `json:"authAccount"`           /*  授权对方账号邮箱  */
	Description   *string `json:"description,omitempty"` /*  String  */
}

type CdaCdaCreateAccountAuthResponse struct {
	StatusCode  *int32                                      `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaCreateAccountAuthReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaCreateAccountAuthErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
}

type CdaCdaCreateAccountAuthReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为空  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为空  */
	ErrorCode *string `json:"errorCode"` /*  成功为空  */
}

type CdaCdaCreateAccountAuthErrorDetailResponse struct{}
