package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbElbUnbindLabelsApi
/* 负载均衡移除标签
 */type CtelbElbUnbindLabelsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbElbUnbindLabelsApi(client *core.CtyunClient) *CtelbElbUnbindLabelsApi {
	return &CtelbElbUnbindLabelsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/unbind-label",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbElbUnbindLabelsApi) Do(ctx context.Context, credential core.Credential, req *CtelbElbUnbindLabelsRequest) (*CtelbElbUnbindLabelsResponse, error) {
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
	var resp CtelbElbUnbindLabelsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbElbUnbindLabelsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  区域ID  */
	ElbID    string `json:"elbID,omitempty"`    /*  负载均衡 ID  */
	LabelID  string `json:"labelID,omitempty"`  /*  标签ID  */
}

type CtelbElbUnbindLabelsResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
