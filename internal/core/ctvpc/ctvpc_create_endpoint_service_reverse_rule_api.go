package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateEndpointServiceReverseRuleApi
/* 创建终端节点服务中转规则(反向访问规则)
 */type CtvpcCreateEndpointServiceReverseRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateEndpointServiceReverseRuleApi(client *core.CtyunClient) *CtvpcCreateEndpointServiceReverseRuleApi {
	return &CtvpcCreateEndpointServiceReverseRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/create-endpoint-service-reverse-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateEndpointServiceReverseRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateEndpointServiceReverseRuleRequest) (*CtvpcCreateEndpointServiceReverseRuleResponse, error) {
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
	var resp CtvpcCreateEndpointServiceReverseRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateEndpointServiceReverseRuleRequest struct {
	ClientToken       string `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string `json:"regionID,omitempty"`          /*  资源池ID  */
	EndpointServiceID string `json:"endpointServiceID,omitempty"` /*  终端节点关联的终端节点服务  */
	EndpointID        string `json:"endpointID,omitempty"`        /*  节点id  */
	TransitIPAddress  string `json:"transitIPAddress,omitempty"`  /*  中转ip地址  */
	TransitPort       int32  `json:"transitPort"`                 /*  中转端口,1到65535  */
	Protocol          string `json:"protocol,omitempty"`          /*  TCP:TCP协议,UDP:UDP协议  */
	TargetIPAddress   string `json:"targetIPAddress,omitempty"`   /*  目标ip地址  */
	TargetPort        int32  `json:"targetPort"`                  /*  目标端口,1到65535  */
}

type CtvpcCreateEndpointServiceReverseRuleResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *struct {
		ReverseRuleID string `json:"reverseRuleID"`
	} `json:"returnObj"`
}
