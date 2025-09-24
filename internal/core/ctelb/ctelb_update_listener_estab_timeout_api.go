package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbUpdateListenerEstabTimeoutApi
/* 设置监听器 Establish Timeout
 */type CtelbUpdateListenerEstabTimeoutApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbUpdateListenerEstabTimeoutApi(client *core.CtyunClient) *CtelbUpdateListenerEstabTimeoutApi {
	return &CtelbUpdateListenerEstabTimeoutApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/update-listener-estab-timeout",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbUpdateListenerEstabTimeoutApi) Do(ctx context.Context, credential core.Credential, req *CtelbUpdateListenerEstabTimeoutRequest) (*CtelbUpdateListenerEstabTimeoutResponse, error) {
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
	var resp CtelbUpdateListenerEstabTimeoutResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbUpdateListenerEstabTimeoutRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  区域ID  */
	ListenerID       string `json:"listenerID,omitempty"`       /*  监听器ID  */
	EstablishTimeout int32  `json:"establishTimeout,omitempty"` /*  建立连接超时时间，单位秒，取值范围： 1 - 1800  */
}

type CtelbUpdateListenerEstabTimeoutResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
