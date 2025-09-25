package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGetFlowPackageMetricApi
/* 获取共享流量包监控
 */type CtvpcGetFlowPackageMetricApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGetFlowPackageMetricApi(client *core.CtyunClient) *CtvpcGetFlowPackageMetricApi {
	return &CtvpcGetFlowPackageMetricApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/flow_package/metric",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGetFlowPackageMetricApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGetFlowPackageMetricRequest) (*CtvpcGetFlowPackageMetricResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sdpID", req.SdpID)
	ctReq.AddParam("startTime", req.StartTime)
	ctReq.AddParam("endTime", req.EndTime)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcGetFlowPackageMetricResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGetFlowPackageMetricRequest struct {
	RegionID  string /*  资源池 ID  */
	SdpID     string /*  记录标识  */
	StartTime string /*  开始时间，YYYY-mmm-dd HH:MM:SS  */
	EndTime   string /*  开始时间，YYYY-mmm-dd HH:MM:SS  */
}

type CtvpcGetFlowPackageMetricResponse struct {
	StatusCode  int32                                         `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcGetFlowPackageMetricReturnObjResponse `json:"returnObj"`             /*  返回购买的共享流量包列表  */
}

type CtvpcGetFlowPackageMetricReturnObjResponse struct {
	DeductTime   int32   `json:"deductTime"`   /*  时间戳  */
	DeductAmount float32 `json:"deductAmount"` /*  消耗流量，单位 Gbps  */
}
