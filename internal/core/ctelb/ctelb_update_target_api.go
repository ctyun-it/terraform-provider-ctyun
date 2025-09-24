package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbUpdateTargetApi
/* 更新后端服务
 */type CtelbUpdateTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbUpdateTargetApi(client *core.CtyunClient) *CtelbUpdateTargetApi {
	return &CtelbUpdateTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/update-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbUpdateTargetApi) Do(ctx context.Context, credential core.Credential, req *CtelbUpdateTargetRequest) (*CtelbUpdateTargetResponse, error) {
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
	var resp CtelbUpdateTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbUpdateTargetRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域ID  */
	ID           string `json:"ID,omitempty"`           /*  后端服务ID, 该字段后续废弃  */
	TargetID     string `json:"targetID,omitempty"`     /*  后端服务ID, 推荐使用该字段, 当同时使用 ID 和 targetID 时，优先使用 targetID  */
	ProtocolPort int32  `json:"protocolPort,omitempty"` /*  协议端口。取值范围：1-65535  */
	Weight       int32  `json:"weight,omitempty"`       /*  权重。取值范围：1-256，默认为100  */
}

type CtelbUpdateTargetResponse struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbUpdateTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbUpdateTargetReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  后端服务组ID  */
}
