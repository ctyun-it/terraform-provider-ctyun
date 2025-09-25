package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteDnatEntryApi
/* 删除 dnat 规则
 */type CtvpcDeleteDnatEntryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteDnatEntryApi(client *core.CtyunClient) *CtvpcDeleteDnatEntryApi {
	return &CtvpcDeleteDnatEntryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/delete-dnat-entry",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteDnatEntryApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteDnatEntryRequest) (*CtvpcDeleteDnatEntryResponse, error) {
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
	var resp CtvpcDeleteDnatEntryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteDnatEntryRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  dnat-6a13d8163678  */
	DNatID      string `json:"dNatID,omitempty"`      /*  'eip-xxx'  */
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
}

type CtvpcDeleteDnatEntryResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
