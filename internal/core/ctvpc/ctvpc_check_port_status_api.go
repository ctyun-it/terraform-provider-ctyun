package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCheckPortStatusApi
/* 获取网卡状态接口
 */type CtvpcCheckPortStatusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCheckPortStatusApi(client *core.CtyunClient) *CtvpcCheckPortStatusApi {
	return &CtvpcCheckPortStatusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ports/check-status",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCheckPortStatusApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCheckPortStatusRequest) (*CtvpcCheckPortStatusResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("portID", req.PortID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcCheckPortStatusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCheckPortStatusRequest struct {
	RegionID string /*  区域id  */
	PortID   string /*  port-id  */
}

type CtvpcCheckPortStatusResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCheckPortStatusReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCheckPortStatusReturnObjResponse struct {
	Id     *string `json:"id,omitempty"`     /*  网卡 id  */
	Status *string `json:"status,omitempty"` /*  网卡状态  */
}
