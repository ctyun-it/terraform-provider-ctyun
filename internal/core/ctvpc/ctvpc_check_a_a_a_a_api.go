package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCheckAAAAApi
/* 检查 aaaa 记录合法性
 */type CtvpcCheckAAAAApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCheckAAAAApi(client *core.CtyunClient) *CtvpcCheckAAAAApi {
	return &CtvpcCheckAAAAApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone-record/check-aaaa",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCheckAAAAApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCheckAAAARequest) (*CtvpcCheckAAAAResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcCheckAAAAResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCheckAAAARequest struct {
	AaaaRecords []string `json:"aaaaRecords"` /*  待检查的 aaaa 数组，数组长度最大支持 10  */
}

type CtvpcCheckAAAAResponse struct {
	StatusCode  int32                            `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCheckAAAAReturnObjResponse `json:"returnObj"`             /*  检查结果  */
	Error       *string                          `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCheckAAAAReturnObjResponse struct {
	Results []*CtvpcCheckAAAAReturnObjResultsResponse `json:"results"` /*  检查结果  */
}

type CtvpcCheckAAAAReturnObjResultsResponse struct {
	Aaaa    *string `json:"aaaa,omitempty"`    /*  被检查的 cname  */
	Valid   *bool   `json:"valid"`             /*  是否合法  */
	Message *string `json:"message,omitempty"` /*  提示信息  */
}
