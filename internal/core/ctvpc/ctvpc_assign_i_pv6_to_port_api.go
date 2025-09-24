package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcAssignIPv6ToPortApi
/* 单个网卡关联多个IPv6地址
 */type CtvpcAssignIPv6ToPortApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcAssignIPv6ToPortApi(client *core.CtyunClient) *CtvpcAssignIPv6ToPortApi {
	return &CtvpcAssignIPv6ToPortApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ports/assign-ipv6",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcAssignIPv6ToPortApi) Do(ctx context.Context, credential core.Credential, req *CtvpcAssignIPv6ToPortRequest) (*CtvpcAssignIPv6ToPortResponse, error) {
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
	var resp CtvpcAssignIPv6ToPortResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcAssignIPv6ToPortRequest struct {
	ClientToken        string    `json:"clientToken,omitempty"`        /*  客户端存根，用于保证订单幂等性，长度 1 - 64  */
	RegionID           string    `json:"regionID,omitempty"`           /*  资源池ID  */
	NetworkInterfaceID string    `json:"networkInterfaceID,omitempty"` /*  网卡ID  */
	Ipv6AddressesCount int32     `json:"ipv6AddressesCount"`           /*  Ipv6地址数量，新增自动分配地址的IPv6的数量， 与ipv6Addresses 二选一， 最多支持 1 个  */
	Ipv6Addresses      []*string `json:"ipv6Addresses"`                /*  IPv6地址列表，新增指定地址的IPv6列表，与 ipv6AddressesCount 二选一， 最多支持 1 个  */
}

type CtvpcAssignIPv6ToPortResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success，英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功， 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
