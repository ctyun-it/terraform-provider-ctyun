package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGetNameserverApi
/* 获取 nameserver
 */type CtvpcGetNameserverApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGetNameserverApi(client *core.CtyunClient) *CtvpcGetNameserverApi {
	return &CtvpcGetNameserverApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone/nameserver",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGetNameserverApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGetNameserverRequest) (*CtvpcGetNameserverResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcGetNameserverResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGetNameserverRequest struct {
	RegionID string /*  资源池ID  */
}

type CtvpcGetNameserverResponse struct {
	StatusCode  int32                                `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGetNameserverReturnObjResponse `json:"returnObj"`             /*  检查结果  */
	Error       *string                              `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcGetNameserverReturnObjResponse struct {
	DnsServer []*string `json:"dnsServer"` /*  nameserver 数组  */
}
