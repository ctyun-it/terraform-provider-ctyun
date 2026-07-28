package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaBgpRouteDeleteApi
/* 删除已有的BGP动态路由 */
type CdaCdaBgpRouteDeleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaBgpRouteDeleteApi(client *core.CtyunClient) *CdaCdaBgpRouteDeleteApi {
	return &CdaCdaBgpRouteDeleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/bgp-route/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaBgpRouteDeleteApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaBgpRouteDeleteRequest) (*CdaCdaBgpRouteDeleteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaBgpRouteDeleteRequest
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
	var resp CdaCdaBgpRouteDeleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaBgpRouteDeleteRequest struct {
	BGPId string `json:"BGPId"` /*  BGP路由ID  */
}

type CdaCdaBgpRouteDeleteResponse struct {
	StatusCode  *int32                                   `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaBgpRouteDeleteReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaBgpRouteDeleteErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaBgpRouteDeleteReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为空  */
	ErrorCode *string `json:"errorCode"` /*  错误代码，成功为空  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为空  */
}

type CdaCdaBgpRouteDeleteErrorDetailResponse struct{}
