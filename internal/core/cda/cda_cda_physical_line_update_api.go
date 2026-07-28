package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaPhysicalLineUpdateApi
/* 修改已创建的物理专线 */
type CdaCdaPhysicalLineUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaPhysicalLineUpdateApi(client *core.CtyunClient) *CdaCdaPhysicalLineUpdateApi {
	return &CdaCdaPhysicalLineUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/physical-line/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaPhysicalLineUpdateApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaPhysicalLineUpdateRequest) (*CdaCdaPhysicalLineUpdateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaPhysicalLineUpdateRequest
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
	var resp CdaCdaPhysicalLineUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaPhysicalLineUpdateRequest struct {
	LineID      string  `json:"lineID"`                /*  物理专线 ID  */
	Bandwidth   *int32  `json:"bandwidth,omitempty"`   /*  带宽(M)  */
	LineName    *string `json:"lineName,omitempty"`    /*  物理专线名字  */
	Location    *string `json:"location,omitempty"`    /*  接入位置  */
	LineCode    *string `json:"lineCode,omitempty"`    /*  电路代号  */
	Description *string `json:"description,omitempty"` /*  描述  */
}

type CdaCdaPhysicalLineUpdateResponse struct {
	StatusCode  *int32                                       `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                      `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                      `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaPhysicalLineUpdateReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                      `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaPhysicalLineUpdateErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                                      `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaPhysicalLineUpdateReturnObjResponse struct {
	Result   *string `json:"result"`   /*  1成功， 0失败  */
	ErrorMsg *string `json:"errorMsg"` /*  成功为null  */
}

type CdaCdaPhysicalLineUpdateErrorDetailResponse struct{}
