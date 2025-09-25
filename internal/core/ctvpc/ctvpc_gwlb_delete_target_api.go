package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGwlbDeleteTargetApi
/* 删除target
 */type CtvpcGwlbDeleteTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGwlbDeleteTargetApi(client *core.CtyunClient) *CtvpcGwlbDeleteTargetApi {
	return &CtvpcGwlbDeleteTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/gwlb/delete-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGwlbDeleteTargetApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGwlbDeleteTargetRequest) (*CtvpcGwlbDeleteTargetResponse, error) {
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
	var resp CtvpcGwlbDeleteTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGwlbDeleteTargetRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	TargetID string `json:"targetID,omitempty"` /*  后端服务 ID  */
}

type CtvpcGwlbDeleteTargetResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGwlbDeleteTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcGwlbDeleteTargetReturnObjResponse struct {
	TargetID *string `json:"targetID,omitempty"` /*  后端服务ID  */
}
