package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteMirrorFlowFilterApi
/* 删除过滤条件
 */type CtvpcDeleteMirrorFlowFilterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteMirrorFlowFilterApi(client *core.CtyunClient) *CtvpcDeleteMirrorFlowFilterApi {
	return &CtvpcDeleteMirrorFlowFilterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/mirrorflow/delete-filter",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteMirrorFlowFilterApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteMirrorFlowFilterRequest) (*CtvpcDeleteMirrorFlowFilterResponse, error) {
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
	var resp CtvpcDeleteMirrorFlowFilterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteMirrorFlowFilterRequest struct {
	RegionID       string `json:"regionID,omitempty"`       /*  区域ID  */
	MirrorFilterID string `json:"mirrorFilterID,omitempty"` /*  过滤条件 ID  */
}

type CtvpcDeleteMirrorFlowFilterResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
