package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteIPv6GatewayApi
/* 调用此接口可删除 IPv6 网关。
 */type CtvpcDeleteIPv6GatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteIPv6GatewayApi(client *core.CtyunClient) *CtvpcDeleteIPv6GatewayApi {
	return &CtvpcDeleteIPv6GatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/delete-ipv6-gateway",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteIPv6GatewayApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteIPv6GatewayRequest) (*CtvpcDeleteIPv6GatewayResponse, error) {
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
	var resp CtvpcDeleteIPv6GatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteIPv6GatewayRequest struct {
	ClientToken   string  `json:"clientToken,omitempty"`   /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	Ipv6GatewayID string  `json:"ipv6GatewayID,omitempty"` /*  ipv6ID值  */
	ProjectID     *string `json:"projectID,omitempty"`     /*  企业项目 ID，默认为0  */
	RegionID      string  `json:"regionID,omitempty"`      /*  资源池ID  */
}

type CtvpcDeleteIPv6GatewayResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
