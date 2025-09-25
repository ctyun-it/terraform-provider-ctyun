package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcModifyEndpointServiceApi
/* 修改终端节点服务属性
 */type CtvpcModifyEndpointServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcModifyEndpointServiceApi(client *core.CtyunClient) *CtvpcModifyEndpointServiceApi {
	return &CtvpcModifyEndpointServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/modify-endpoint-service",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcModifyEndpointServiceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcModifyEndpointServiceRequest) (*CtvpcModifyEndpointServiceResponse, error) {
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
	var resp CtvpcModifyEndpointServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcModifyEndpointServiceRequest struct {
	RegionID          string  `json:"regionID,omitempty"`          /*  资源池ID  */
	EndpointServiceID string  `json:"endpointServiceID,omitempty"` /*  终端节点服务ID  */
	AutoConnection    *bool   `json:"autoConnection"`              /*  是否自动连接，true/false  */
	Name              *string `json:"name,omitempty"`              /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	DnsName           *string `json:"dnsName,omitempty"`           /*  dns 名字，仅支持有权限的用户修改  */
	OaType            *string `json:"oaType,omitempty"`            /*  oa 类型，支持: tcp_option / proxy_protocol / close  */
}

type CtvpcModifyEndpointServiceResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
