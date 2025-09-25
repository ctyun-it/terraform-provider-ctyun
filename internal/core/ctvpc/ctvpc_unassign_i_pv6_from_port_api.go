package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUnassignIPv6FromPortApi
/* 单个网卡解绑多个 IPv6 地址
 */type CtvpcUnassignIPv6FromPortApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUnassignIPv6FromPortApi(client *core.CtyunClient) *CtvpcUnassignIPv6FromPortApi {
	return &CtvpcUnassignIPv6FromPortApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ports/unassign-ipv6",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUnassignIPv6FromPortApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUnassignIPv6FromPortRequest) (*CtvpcUnassignIPv6FromPortResponse, error) {
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
	var resp CtvpcUnassignIPv6FromPortResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUnassignIPv6FromPortRequest struct {
	ClientToken        string   `json:"clientToken,omitempty"`        /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID           string   `json:"regionID,omitempty"`           /*  资源池ID  */
	NetworkInterfaceID string   `json:"networkInterfaceID,omitempty"` /*  网卡ID  */
	Ipv6Addresses      []string `json:"ipv6Addresses"`                /*  解绑的IPv6地址列表, 最多支持 1 个  */
}

type CtvpcUnassignIPv6FromPortResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
