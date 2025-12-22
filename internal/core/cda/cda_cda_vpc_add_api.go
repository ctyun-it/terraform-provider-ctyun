package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaVpcAddApi
/* 创建的云专线(Cloud Dedicated Access)网关。 */
type CdaCdaVpcAddApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaVpcAddApi(client *core.CtyunClient) *CdaCdaVpcAddApi {
	return &CdaCdaVpcAddApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/vpc/add",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaVpcAddApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaVpcAddRequest) (*CdaCdaVpcAddResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaVpcAddRequest
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
	var resp CdaCdaVpcAddResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaVpcAddRequest struct {
	GatewayName      string  `json:"gatewayName"`           /*  专线网关名字  */
	Account          string  `json:"account"`               /*  天翼云客户邮箱  */
	ResourcePool     string  `json:"resourcePool"`          /*  资源池ID  */
	ResourcePoolName string  `json:"resourcePoolName"`      /*  资源池名字  */
	Hostname         string  `json:"hostname"`              /*  交换机hostname  */
	IsSwConfig       *bool   `json:"isSwConfig,omitempty"`  /*  是否下发配置到交换机，默认True  */
	Description      *string `json:"description,omitempty"` /*  描述  */
}

type CdaCdaVpcAddResponse struct {
	StatusCode  *int32                           `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                          `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                          `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaVpcAddReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                          `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaVpcAddErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                          `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaVpcAddReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为null  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为null  */
	ErrorCode *string `json:"errorCode"` /*  成功为null  */
}

type CdaCdaVpcAddErrorDetailResponse struct{}
