package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcV4PortsAsyncDetachBmApi
/* 网卡解绑物理机
 */type CtvpcV4PortsAsyncDetachBmApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcV4PortsAsyncDetachBmApi(client *core.CtyunClient) *CtvpcV4PortsAsyncDetachBmApi {
	return &CtvpcV4PortsAsyncDetachBmApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ports/async-detach-bm",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcV4PortsAsyncDetachBmApi) Do(ctx context.Context, credential core.Credential, req *CtvpcV4PortsAsyncDetachBmRequest) (*CtvpcV4PortsAsyncDetachBmResponse, error) {
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
	var resp CtvpcV4PortsAsyncDetachBmResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcV4PortsAsyncDetachBmRequest struct {
	RegionID           string `json:"regionID,omitempty"`           /*  资源池ID  */
	AzName             string `json:"azName,omitempty"`             /*  可用区  */
	NetworkInterfaceID string `json:"networkInterfaceID,omitempty"` /*  网卡ID  */
	InstanceID         string `json:"instanceID,omitempty"`         /*  解绑的物理机ID  */
}

type CtvpcV4PortsAsyncDetachBmResponse struct {
	StatusCode  int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcV4PortsAsyncDetachBmReturnObjResponse `json:"returnObj"`             /*  见表returnObj  */
}

type CtvpcV4PortsAsyncDetachBmReturnObjResponse struct {
	Status *string `json:"status,omitempty"` /*  状态。in_progress表示在异步处理中，done成功  */
}
