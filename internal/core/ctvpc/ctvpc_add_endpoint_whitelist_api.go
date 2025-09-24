package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcAddEndpointWhitelistApi
/* 添加终端节点白名单
 */type CtvpcAddEndpointWhitelistApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcAddEndpointWhitelistApi(client *core.CtyunClient) *CtvpcAddEndpointWhitelistApi {
	return &CtvpcAddEndpointWhitelistApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/add-endpoint-whitelist",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcAddEndpointWhitelistApi) Do(ctx context.Context, credential core.Credential, req *CtvpcAddEndpointWhitelistRequest) (*CtvpcAddEndpointWhitelistResponse, error) {
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
	var resp CtvpcAddEndpointWhitelistResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcAddEndpointWhitelistRequest struct {
	ClientToken string   `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string   `json:"regionID,omitempty"`    /*  资源池ID  */
	EndpointID  string   `json:"endpointID,omitempty"`  /*  终端节点ID  */
	Whitelist   []string `json:"whitelist"`             /*  白名单列表，最多支持同时添加 20 个 ip  */
}

type CtvpcAddEndpointWhitelistResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
