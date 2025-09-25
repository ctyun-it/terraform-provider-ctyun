package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateVpcPeerRouteApi
/* 创建对等链接路由
 */type CtvpcCreateVpcPeerRouteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateVpcPeerRouteApi(client *core.CtyunClient) *CtvpcCreateVpcPeerRouteApi {
	return &CtvpcCreateVpcPeerRouteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/vpcpeer/create-route",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateVpcPeerRouteApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateVpcPeerRouteRequest) (*CtvpcCreateVpcPeerRouteResponse, error) {
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
	var resp CtvpcCreateVpcPeerRouteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateVpcPeerRouteRequest struct {
	IpVersion   int32  `json:"ipVersion"`             /*  4 表示 ipv4，6 表示 ipv6  */
	NextHopID   string `json:"nextHopID,omitempty"`   /*  下一跳设备id  */
	VpcID       string `json:"vpcID,omitempty"`       /*  路由表所在 vpc id  */
	Destination string `json:"destination,omitempty"` /*  目的 cidr  */
	RegionID    string `json:"regionID,omitempty"`    /*  区域id  */
}

type CtvpcCreateVpcPeerRouteResponse struct {
	StatusCode  int32                                     `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateVpcPeerRouteReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcCreateVpcPeerRouteReturnObjResponse struct {
	RouteRule *string `json:"routeRule,omitempty"` /*  路由规则 id  */
}
