package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcEndpointGetOssPolicyApi
/* 获取oss策略
 */type CtvpcEndpointGetOssPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcEndpointGetOssPolicyApi(client *core.CtyunClient) *CtvpcEndpointGetOssPolicyApi {
	return &CtvpcEndpointGetOssPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/endpoint/get-oss-policy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcEndpointGetOssPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtvpcEndpointGetOssPolicyRequest) (*CtvpcEndpointGetOssPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("endpointID", req.EndpointID)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcEndpointGetOssPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcEndpointGetOssPolicyRequest struct {
	EndpointID string /*  终端节点 ID  */
	RegionID   string /*  区域 ID  */
}

type CtvpcEndpointGetOssPolicyResponse struct {
	StatusCode  int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcEndpointGetOssPolicyReturnObjResponse `json:"returnObj"`             /*  业务数据  */
}

type CtvpcEndpointGetOssPolicyReturnObjResponse struct {
	Policy *string `json:"policy,omitempty"` /*  JSON 文档形式的Vpce权限策略。  */
}
