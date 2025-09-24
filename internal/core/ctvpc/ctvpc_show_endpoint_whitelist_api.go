package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowEndpointWhitelistApi
/* 查询终端节点白名单
 */type CtvpcShowEndpointWhitelistApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowEndpointWhitelistApi(client *core.CtyunClient) *CtvpcShowEndpointWhitelistApi {
	return &CtvpcShowEndpointWhitelistApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/show-endpoint-whitelist",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowEndpointWhitelistApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowEndpointWhitelistRequest) (*CtvpcShowEndpointWhitelistResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("endpointID", req.EndpointID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowEndpointWhitelistResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowEndpointWhitelistRequest struct {
	RegionID   string /*  资源池ID  */
	EndpointID string /*  终端节点ID  */
}

type CtvpcShowEndpointWhitelistResponse struct {
	StatusCode  int32     `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*string `json:"returnObj"`             /*  ip 白名单列表  */
	Error       *string   `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
