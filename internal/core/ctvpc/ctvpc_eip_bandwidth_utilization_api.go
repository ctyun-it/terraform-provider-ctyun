package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcEipBandwidthUtilizationApi
/* eip带宽利用率
 */type CtvpcEipBandwidthUtilizationApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcEipBandwidthUtilizationApi(client *core.CtyunClient) *CtvpcEipBandwidthUtilizationApi {
	return &CtvpcEipBandwidthUtilizationApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/eip/bandwidth-utilization",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcEipBandwidthUtilizationApi) Do(ctx context.Context, credential core.Credential, req *CtvpcEipBandwidthUtilizationRequest) (*CtvpcEipBandwidthUtilizationResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("eipID", req.EipID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcEipBandwidthUtilizationResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcEipBandwidthUtilizationRequest struct {
	RegionID string /*  资源池 ID  */
	EipID    string /*  弹性公网IP的ID  */
}

type CtvpcEipBandwidthUtilizationResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
