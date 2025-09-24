package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateIPv6GatewayApi
/* 创建 IPv6 网关
 */type CtvpcCreateIPv6GatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateIPv6GatewayApi(client *core.CtyunClient) *CtvpcCreateIPv6GatewayApi {
	return &CtvpcCreateIPv6GatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/create-ipv6-gateway",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateIPv6GatewayApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateIPv6GatewayRequest) (*CtvpcCreateIPv6GatewayResponse, error) {
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
	var resp CtvpcCreateIPv6GatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateIPv6GatewayRequest struct {
	ClientToken string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	VpcID       string  `json:"vpcID,omitempty"`       /*  VPC 名称  */
	ProjectID   *string `json:"projectID,omitempty"`   /*  企业项目 ID，默认为0  */
	Description *string `json:"description,omitempty"` /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{}  */
	RegionID    string  `json:"regionID,omitempty"`    /*  xj8g-894g-09oi-po09-12ol-6e6a  */
}

type CtvpcCreateIPv6GatewayResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
