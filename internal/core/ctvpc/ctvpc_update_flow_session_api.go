package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdateFlowSessionApi
/* 更新流量会话
 */type CtvpcUpdateFlowSessionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdateFlowSessionApi(client *core.CtyunClient) *CtvpcUpdateFlowSessionApi {
	return &CtvpcUpdateFlowSessionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/flowsession/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdateFlowSessionApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdateFlowSessionRequest) (*CtvpcUpdateFlowSessionResponse, error) {
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
	var resp CtvpcUpdateFlowSessionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdateFlowSessionRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  区域ID  */
	FlowSessionID string `json:"flowSessionID,omitempty"` /*  流量镜像会话 ID  */
	Vni           int32  `json:"vni"`                     /*  vni, 0 - 1677215  */
	Name          string `json:"name,omitempty"`          /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
}

type CtvpcUpdateFlowSessionResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
