package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcEndpointCreateOssPolicyApi
/* 创建oss策略
 */type CtvpcEndpointCreateOssPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcEndpointCreateOssPolicyApi(client *core.CtyunClient) *CtvpcEndpointCreateOssPolicyApi {
	return &CtvpcEndpointCreateOssPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/endpoint/create-oss-policy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcEndpointCreateOssPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtvpcEndpointCreateOssPolicyRequest) (*CtvpcEndpointCreateOssPolicyResponse, error) {
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
	var resp CtvpcEndpointCreateOssPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcEndpointCreateOssPolicyRequest struct {
	EndpointID string `json:"endpointID,omitempty"` /*  终端节点 ID  */
	RegionID   string `json:"regionID,omitempty"`   /*  区域 ID  */
	Policy     string `json:"policy,omitempty"`     /*  JSON 文档形式的存储桶策略  */
}

type CtvpcEndpointCreateOssPolicyResponse struct {
	StatusCode  int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcEndpointCreateOssPolicyReturnObjResponse `json:"returnObj"`             /*  object  */
}

type CtvpcEndpointCreateOssPolicyReturnObjResponse struct {
	EndpointID *string `json:"endpointID,omitempty"` /*  终端节点 ID  */
}
