package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcSetZoneProxyPatternApi
/* 设置 proxy pattern
 */type CtvpcSetZoneProxyPatternApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcSetZoneProxyPatternApi(client *core.CtyunClient) *CtvpcSetZoneProxyPatternApi {
	return &CtvpcSetZoneProxyPatternApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone/update-proxy-pattern",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcSetZoneProxyPatternApi) Do(ctx context.Context, credential core.Credential, req *CtvpcSetZoneProxyPatternRequest) (*CtvpcSetZoneProxyPatternResponse, error) {
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
	var resp CtvpcSetZoneProxyPatternResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcSetZoneProxyPatternRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID  */
	ZoneID       string `json:"zoneID,omitempty"`       /*  zoneID  */
	ProxyPattern string `json:"proxyPattern,omitempty"` /*  zone：当前可用区不进行递归解析。 record：不完全劫持，进行递归解析代理  */
}

type CtvpcSetZoneProxyPatternResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
