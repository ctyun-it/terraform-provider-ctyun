package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcRemoveIPv6FromIPv6BandwidthApi
/* IPv6 带宽移除 IPv6 地址。
 */type CtvpcRemoveIPv6FromIPv6BandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcRemoveIPv6FromIPv6BandwidthApi(client *core.CtyunClient) *CtvpcRemoveIPv6FromIPv6BandwidthApi {
	return &CtvpcRemoveIPv6FromIPv6BandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ipv6_bandwidth/remove-ipv6",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcRemoveIPv6FromIPv6BandwidthApi) Do(ctx context.Context, credential core.Credential, req *CtvpcRemoveIPv6FromIPv6BandwidthRequest) (*CtvpcRemoveIPv6FromIPv6BandwidthResponse, error) {
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
	var resp CtvpcRemoveIPv6FromIPv6BandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcRemoveIPv6FromIPv6BandwidthRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	BandwidthID string `json:"bandwidthID,omitempty"` /*  IPv6 带宽 ID  */
	Ip          string `json:"ip,omitempty"`          /*  IPv6 地址  */
}

type CtvpcRemoveIPv6FromIPv6BandwidthResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
