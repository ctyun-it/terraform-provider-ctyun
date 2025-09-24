package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbGwlbCreateTargetApi
/* 创建target
 */type CtelbGwlbCreateTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbGwlbCreateTargetApi(client *core.CtyunClient) *CtelbGwlbCreateTargetApi {
	return &CtelbGwlbCreateTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/gwlb/create-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbGwlbCreateTargetApi) Do(ctx context.Context, credential core.Credential, req *CtelbGwlbCreateTargetRequest) (*CtelbGwlbCreateTargetResponse, error) {
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
	var resp CtelbGwlbCreateTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbGwlbCreateTargetRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  资源池 ID  */
	TargetGroupID string `json:"targetGroupID,omitempty"` /*  后端服务组 ID  */
	InstanceID    string `json:"instanceID,omitempty"`    /*  实例 ID  */
	InstanceType  string `json:"instanceType,omitempty"`  /*  支持 VM / BM  */
	Weight        int32  `json:"weight,omitempty"`        /*  权重，仅支持填写 100  */
}

type CtelbGwlbCreateTargetResponse struct {
	StatusCode  int32                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbGwlbCreateTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbGwlbCreateTargetReturnObjResponse struct {
	TargetID string `json:"targetID,omitempty"` /*  后端服务ID  */
}
