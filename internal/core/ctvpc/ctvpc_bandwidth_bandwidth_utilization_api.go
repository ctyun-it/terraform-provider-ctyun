package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcBandwidthBandwidthUtilizationApi
/* 带宽利用率
 */type CtvpcBandwidthBandwidthUtilizationApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcBandwidthBandwidthUtilizationApi(client *core.CtyunClient) *CtvpcBandwidthBandwidthUtilizationApi {
	return &CtvpcBandwidthBandwidthUtilizationApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/bandwidth/bandwidth-utilization",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcBandwidthBandwidthUtilizationApi) Do(ctx context.Context, credential core.Credential, req *CtvpcBandwidthBandwidthUtilizationRequest) (*CtvpcBandwidthBandwidthUtilizationResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("bandwidthID", req.BandwidthID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcBandwidthBandwidthUtilizationResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcBandwidthBandwidthUtilizationRequest struct {
	RegionID    string /*  资源池 ID  */
	BandwidthID string /*  带宽 ID  */
}

type CtvpcBandwidthBandwidthUtilizationResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
