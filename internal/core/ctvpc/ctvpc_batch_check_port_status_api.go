package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcBatchCheckPortStatusApi
/* 网卡状态批量查询接口
 */type CtvpcBatchCheckPortStatusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcBatchCheckPortStatusApi(client *core.CtyunClient) *CtvpcBatchCheckPortStatusApi {
	return &CtvpcBatchCheckPortStatusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ports/check-status-batch",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcBatchCheckPortStatusApi) Do(ctx context.Context, credential core.Credential, req *CtvpcBatchCheckPortStatusRequest) (*CtvpcBatchCheckPortStatusResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("portIDs", req.PortIDs)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcBatchCheckPortStatusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcBatchCheckPortStatusRequest struct {
	RegionID string /*  区域id  */
	PortIDs  string /*  多个网卡用 , 拼接起来, port-id,port-id, 最多支持同时检查 10 个网卡  */
}

type CtvpcBatchCheckPortStatusResponse struct {
	StatusCode  int32                                         `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcBatchCheckPortStatusReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                       `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcBatchCheckPortStatusReturnObjResponse struct {
	Id     *string `json:"id,omitempty"`     /*  网卡 id  */
	Status *string `json:"status,omitempty"` /*  网卡状态  */
}
