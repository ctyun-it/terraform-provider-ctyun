package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateEndpointServiceTransitIPApi
/* 创建节点服务中转ip
 */type CtvpcCreateEndpointServiceTransitIPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateEndpointServiceTransitIPApi(client *core.CtyunClient) *CtvpcCreateEndpointServiceTransitIPApi {
	return &CtvpcCreateEndpointServiceTransitIPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/create-endpoint-service-transit-ip",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateEndpointServiceTransitIPApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateEndpointServiceTransitIPRequest) (*CtvpcCreateEndpointServiceTransitIPResponse, error) {
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
	var resp CtvpcCreateEndpointServiceTransitIPResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateEndpointServiceTransitIPRequest struct {
	ClientToken       string  `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string  `json:"regionID,omitempty"`          /*  资源池ID  */
	EndpointServiceID string  `json:"endpointServiceID,omitempty"` /*  节点服务id  */
	SubnetID          string  `json:"subnetID,omitempty"`          /*  子网id,子网必须和endpointService属于同一个vpc  */
	TransitIP         *string `json:"transitIP,omitempty"`         /*  中转地址  */
}

type CtvpcCreateEndpointServiceTransitIPResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *struct {
		TransitIP string `json:"transitIP"`
	} `json:"returnObj"`
}
