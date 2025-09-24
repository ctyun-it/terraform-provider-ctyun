package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcModifyEndpointIpVersionApi
/* 转换终端节点IP版本
 */type CtvpcModifyEndpointIpVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcModifyEndpointIpVersionApi(client *core.CtyunClient) *CtvpcModifyEndpointIpVersionApi {
	return &CtvpcModifyEndpointIpVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/modify-endpoint-ip-version",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcModifyEndpointIpVersionApi) Do(ctx context.Context, credential core.Credential, req *CtvpcModifyEndpointIpVersionRequest) (*CtvpcModifyEndpointIpVersionResponse, error) {
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
	var resp CtvpcModifyEndpointIpVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcModifyEndpointIpVersionRequest struct {
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池ID  */
	EndpointID  string  `json:"endpointID,omitempty"`  /*  终端节点id  */
	IpVersion   int32   `json:"ipVersion"`             /*  0:ipv4, 1:ipv6（暂不支持）, 2:双栈  */
	IpAddress   *string `json:"ipAddress,omitempty"`   /*  终端节点ipv4地址  */
	Ipv6Address *string `json:"ipv6Address,omitempty"` /*  终端节点ipv6地址  */
}

type CtvpcModifyEndpointIpVersionResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
