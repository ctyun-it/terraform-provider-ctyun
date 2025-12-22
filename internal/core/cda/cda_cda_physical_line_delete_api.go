package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaPhysicalLineDeleteApi
/* 删除已创建的物理专线 */
type CdaCdaPhysicalLineDeleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaPhysicalLineDeleteApi(client *core.CtyunClient) *CdaCdaPhysicalLineDeleteApi {
	return &CdaCdaPhysicalLineDeleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/physical-line/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaPhysicalLineDeleteApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaPhysicalLineDeleteRequest) (*CdaCdaPhysicalLineDeleteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaPhysicalLineDeleteRequest
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
	var resp CdaCdaPhysicalLineDeleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaPhysicalLineDeleteRequest struct {
	LineID string `json:"lineID"` /*  物理专线 ID  */
}

type CdaCdaPhysicalLineDeleteResponse struct {
	StatusCode  *int32                                       `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                      `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                      `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaPhysicalLineDeleteReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                      `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaPhysicalLineDeleteErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                                      `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaPhysicalLineDeleteReturnObjResponse struct {
	Result   *string `json:"result"`   /*  1成功， 0失败  */
	ErrorMsg *string `json:"errorMsg"` /*  成功为空  */
}

type CdaCdaPhysicalLineDeleteErrorDetailResponse struct{}
