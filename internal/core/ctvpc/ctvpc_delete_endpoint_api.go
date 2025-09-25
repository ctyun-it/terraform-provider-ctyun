package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteEndpointApi
/* 删除终端节点
 */type CtvpcDeleteEndpointApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteEndpointApi(client *core.CtyunClient) *CtvpcDeleteEndpointApi {
	return &CtvpcDeleteEndpointApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/delete-endpoint",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteEndpointApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteEndpointRequest) (*CtvpcDeleteEndpointResponse, error) {
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
	var resp CtvpcDeleteEndpointResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteEndpointRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID  */
	EndpointID  string `json:"endpointID,omitempty"`  /*  终端节点id  */
}

type CtvpcDeleteEndpointResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDeleteEndpointReturnObjResponse `json:"returnObj"`             /*  object  */
}

type CtvpcDeleteEndpointReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  55d531d7bf2d47658897c42ffb918423  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  20221021191602644224  */
	RegionID             *string `json:"regionID,omitempty"`             /*  81f7728662dd11ec810800155d307d5b  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  refuned  */
}
