package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCheckPtrApi
/* 检查 ptr 记录合法性
 */type CtvpcCheckPtrApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCheckPtrApi(client *core.CtyunClient) *CtvpcCheckPtrApi {
	return &CtvpcCheckPtrApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone-record/check-ptr",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCheckPtrApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCheckPtrRequest) (*CtvpcCheckPtrResponse, error) {
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
	var resp CtvpcCheckPtrResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCheckPtrRequest struct {
	PtrRecords []string `json:"ptrRecords"` /*  待检查的 ptr 数组，数组长度最大支持 10  */
}

type CtvpcCheckPtrResponse struct {
	StatusCode  int32                           `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                         `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                         `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                         `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCheckPtrReturnObjResponse `json:"returnObj"`             /*  检查结果  */
	Error       *string                         `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCheckPtrReturnObjResponse struct {
	Results []*CtvpcCheckPtrReturnObjResultsResponse `json:"results"` /*  检查结果  */
}

type CtvpcCheckPtrReturnObjResultsResponse struct {
	Ptr     *string `json:"ptr,omitempty"`     /*  被检查的 cname  */
	Valid   *bool   `json:"valid"`             /*  是否合法  */
	Message *string `json:"message,omitempty"` /*  提示信息  */
}
