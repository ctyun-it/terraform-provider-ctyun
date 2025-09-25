package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteEndpointServiceWhitelistApi
/* 删除终端节点服务白名单
 */type CtvpcDeleteEndpointServiceWhitelistApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteEndpointServiceWhitelistApi(client *core.CtyunClient) *CtvpcDeleteEndpointServiceWhitelistApi {
	return &CtvpcDeleteEndpointServiceWhitelistApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/delete-endpoint-service-whitelist",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteEndpointServiceWhitelistApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteEndpointServiceWhitelistRequest) (*CtvpcDeleteEndpointServiceWhitelistResponse, error) {
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
	var resp CtvpcDeleteEndpointServiceWhitelistResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteEndpointServiceWhitelistRequest struct {
	ClientToken       string  `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string  `json:"regionID,omitempty"`          /*  资源池ID  */
	EndpointServiceID string  `json:"endpointServiceID,omitempty"` /*  终端节点服务ID  */
	Email             *string `json:"email,omitempty"`             /*  账户邮箱，邮箱和账户至少填一个  */
	BssAccountID      *string `json:"bssAccountID,omitempty"`      /*  账户  */
}

type CtvpcDeleteEndpointServiceWhitelistResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
