package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbElbBindLabelsApi
/* 负载均衡添加标签
 */type CtelbElbBindLabelsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbElbBindLabelsApi(client *core.CtyunClient) *CtelbElbBindLabelsApi {
	return &CtelbElbBindLabelsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/bind-label",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbElbBindLabelsApi) Do(ctx context.Context, credential core.Credential, req *CtelbElbBindLabelsRequest) (*CtelbElbBindLabelsResponse, error) {
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
	var resp CtelbElbBindLabelsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbElbBindLabelsRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  区域ID  */
	ElbID      string `json:"elbID,omitempty"`      /*  负载均衡 ID  */
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签 key  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签 取值  */
}

type CtelbElbBindLabelsResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
