package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcEipEgerssBandwidthUtilizationApi
/* 出方向带宽利用率。
 */type CtvpcEipEgerssBandwidthUtilizationApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcEipEgerssBandwidthUtilizationApi(client *core.CtyunClient) *CtvpcEipEgerssBandwidthUtilizationApi {
	return &CtvpcEipEgerssBandwidthUtilizationApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/eip/egress-bandwidth-utilization",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcEipEgerssBandwidthUtilizationApi) Do(ctx context.Context, credential core.Credential, req *CtvpcEipEgerssBandwidthUtilizationRequest) (*CtvpcEipEgerssBandwidthUtilizationResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("eipID", req.EipID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcEipEgerssBandwidthUtilizationResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcEipEgerssBandwidthUtilizationRequest struct {
	RegionID string /*  资源池 ID  */
	EipID    string /*  弹性公网IP的ID  */
}

type CtvpcEipEgerssBandwidthUtilizationResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
