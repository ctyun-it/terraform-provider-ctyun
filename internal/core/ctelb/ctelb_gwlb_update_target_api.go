package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbGwlbUpdateTargetApi
/* 更新target
 */type CtelbGwlbUpdateTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbGwlbUpdateTargetApi(client *core.CtyunClient) *CtelbGwlbUpdateTargetApi {
	return &CtelbGwlbUpdateTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/gwlb/update-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbGwlbUpdateTargetApi) Do(ctx context.Context, credential core.Credential, req *CtelbGwlbUpdateTargetRequest) (*CtelbGwlbUpdateTargetResponse, error) {
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
	var resp CtelbGwlbUpdateTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbGwlbUpdateTargetRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	TargetID string `json:"targetID,omitempty"` /*  后端服务 ID  */
	Weight   int32  `json:"weight,omitempty"`   /*  权重，仅支持填 100  */
}

type CtelbGwlbUpdateTargetResponse struct {
	StatusCode  int32                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbGwlbUpdateTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbGwlbUpdateTargetReturnObjResponse struct {
	TargetID string `json:"targetID,omitempty"` /*  后端服务ID  */
}
