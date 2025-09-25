package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdateIPv6StatusForVpcApi
/* 修改专有网络VPC的 IPv6 状态：开启、关闭。关闭VPC的IPv6开关前，需要关闭所有子网的IPv6开关。
 */type CtvpcUpdateIPv6StatusForVpcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdateIPv6StatusForVpcApi(client *core.CtyunClient) *CtvpcUpdateIPv6StatusForVpcApi {
	return &CtvpcUpdateIPv6StatusForVpcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/update-ipv6-status",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdateIPv6StatusForVpcApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdateIPv6StatusForVpcRequest) (*CtvpcUpdateIPv6StatusForVpcResponse, error) {
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
	var resp CtvpcUpdateIPv6StatusForVpcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdateIPv6StatusForVpcRequest struct {
	ClientToken string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	ProjectID   *string `json:"projectID,omitempty"`   /*  企业项目 ID，默认为"0"  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池 ID  */
	VpcID       string  `json:"vpcID,omitempty"`       /*  VPC 的  */
	EnableIpv6  bool    `json:"enableIpv6"`            /*  是否开启 IPv6 网段。取值：false（默认值）:不开启，true: 开启  */
}

type CtvpcUpdateIPv6StatusForVpcResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
