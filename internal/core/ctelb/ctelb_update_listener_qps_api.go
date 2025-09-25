package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbUpdateListenerQpsApi
/* 设置监听器 QPS
 */type CtelbUpdateListenerQpsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbUpdateListenerQpsApi(client *core.CtyunClient) *CtelbUpdateListenerQpsApi {
	return &CtelbUpdateListenerQpsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/update-listener-qps",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbUpdateListenerQpsApi) Do(ctx context.Context, credential core.Credential, req *CtelbUpdateListenerQpsRequest) (*CtelbUpdateListenerQpsResponse, error) {
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
	var resp CtelbUpdateListenerQpsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbUpdateListenerQpsRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  区域ID  */
	ListenerID  string `json:"listenerID,omitempty"`  /*  监听器ID  */
	ListenerQps int32  `json:"listenerQps,omitempty"` /*  qps 大小  */
}

type CtelbUpdateListenerQpsResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
