package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteVpcPeerRouteApi
/* 删除对等链接路由
 */type CtvpcDeleteVpcPeerRouteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteVpcPeerRouteApi(client *core.CtyunClient) *CtvpcDeleteVpcPeerRouteApi {
	return &CtvpcDeleteVpcPeerRouteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/vpcpeer/delete-route",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteVpcPeerRouteApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteVpcPeerRouteRequest) (*CtvpcDeleteVpcPeerRouteResponse, error) {
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
	var resp CtvpcDeleteVpcPeerRouteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteVpcPeerRouteRequest struct {
	RouteRuleID string `json:"routeRuleID,omitempty"` /*  路由规则id  */
	RegionID    string `json:"regionID,omitempty"`    /*  区域id  */
}

type CtvpcDeleteVpcPeerRouteResponse struct {
	StatusCode  int32                                     `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDeleteVpcPeerRouteReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcDeleteVpcPeerRouteReturnObjResponse struct {
	RouteRule *string `json:"routeRule,omitempty"` /*  路由规则 id  */
}
