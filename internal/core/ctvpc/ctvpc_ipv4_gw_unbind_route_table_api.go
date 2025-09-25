package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcIpv4GwUnbindRouteTableApi
/* IPv4网关解绑网关路由表
 */type CtvpcIpv4GwUnbindRouteTableApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcIpv4GwUnbindRouteTableApi(client *core.CtyunClient) *CtvpcIpv4GwUnbindRouteTableApi {
	return &CtvpcIpv4GwUnbindRouteTableApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/ipv4-gw/remove-route-table-binding",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcIpv4GwUnbindRouteTableApi) Do(ctx context.Context, credential core.Credential, req *CtvpcIpv4GwUnbindRouteTableRequest) (*CtvpcIpv4GwUnbindRouteTableResponse, error) {
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
	var resp CtvpcIpv4GwUnbindRouteTableResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcIpv4GwUnbindRouteTableRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string `json:"regionID,omitempty"`    /*  区域ID  */
	Ipv4GwID    string `json:"ipv4GwID,omitempty"`    /*  IPv4网关的ID  */
}

type CtvpcIpv4GwUnbindRouteTableResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
