package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdateIPv6BandwidthAttributeApi
/* 更新 IPv6 带宽信息。
 */type CtvpcUpdateIPv6BandwidthAttributeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdateIPv6BandwidthAttributeApi(client *core.CtyunClient) *CtvpcUpdateIPv6BandwidthAttributeApi {
	return &CtvpcUpdateIPv6BandwidthAttributeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ipv6_bandwidth/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdateIPv6BandwidthAttributeApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdateIPv6BandwidthAttributeRequest) (*CtvpcUpdateIPv6BandwidthAttributeResponse, error) {
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
	var resp CtvpcUpdateIPv6BandwidthAttributeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdateIPv6BandwidthAttributeRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	BandwidthID string `json:"bandwidthID,omitempty"` /*  IPv6 带宽 ID  */
	Name        string `json:"name,omitempty"`        /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
}

type CtvpcUpdateIPv6BandwidthAttributeResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
