package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdateSubnetIPv6StatusApi
/* 修改子网subnet的 IPv6 状态：开启、关闭。开启子网IPv6前，需要先开启VPC的IPv6。关闭子网IPv6前，需要删除所有占用IPv6地址的实例。
 */type CtvpcUpdateSubnetIPv6StatusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdateSubnetIPv6StatusApi(client *core.CtyunClient) *CtvpcUpdateSubnetIPv6StatusApi {
	return &CtvpcUpdateSubnetIPv6StatusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/update-subnet-ipv6-status",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdateSubnetIPv6StatusApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdateSubnetIPv6StatusRequest) (*CtvpcUpdateSubnetIPv6StatusResponse, error) {
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
	var resp CtvpcUpdateSubnetIPv6StatusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdateSubnetIPv6StatusRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池 ID  */
	ProjectID   *string `json:"projectID,omitempty"`   /*  企业项目 ID，默认为0  */
	SubnetID    string  `json:"subnetID,omitempty"`    /*  子网 的 ID  */
	EnableIpv6  bool    `json:"enableIpv6"`            /*  是否开启 IPv6 网段。取值：false（默认值）:不开启，true: 开启  */
}

type CtvpcUpdateSubnetIPv6StatusResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
