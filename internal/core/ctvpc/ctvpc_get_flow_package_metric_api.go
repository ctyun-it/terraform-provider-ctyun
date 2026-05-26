package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGetFlowPackageMetricApi
/* 本接口已下线，新获取共享流量包监控将于近期上线 */
type CtvpcGetFlowPackageMetricApi struct {
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
	_, err := ctReq.WriteJson(struct {
		*CtvpcGetFlowPackageMetricRequest
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
	var resp CtvpcGetFlowPackageMetricResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtvpcGetFlowPackageMetricRequest struct {
	RegionID  string `json:"regionID"`  /*  资源池 ID  */
	SdpID     string `json:"sdpID"`     /*  记录标识  */
	StartTime string `json:"startTime"` /*  开始时间，YYYY-mmm-dd HH:MM:SS  */
	EndTime   string `json:"endTime"`   /*  开始时间，YYYY-mmm-dd HH:MM:SS  */
}

type CtvpcGetFlowPackageMetricResponse struct {
	StatusCode  *int32                                        `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                       `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                       `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                       `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcGetFlowPackageMetricReturnObjResponse `json:"returnObj"`   /*  返回购买的共享流量包列表  */
}

type CtvpcGetFlowPackageMetricReturnObjResponse struct {
	DeductTime   *int32   `json:"deductTime"`   /*  时间戳  */
	DeductAmount *float32 `json:"deductAmount"` /*  消耗流量，单位 Gbps  */
}
