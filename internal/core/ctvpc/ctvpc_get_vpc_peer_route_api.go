package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGetVpcPeerRouteApi
/* 获取对等链接路由
 */type CtvpcGetVpcPeerRouteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGetVpcPeerRouteApi(client *core.CtyunClient) *CtvpcGetVpcPeerRouteApi {
	return &CtvpcGetVpcPeerRouteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/vpcpeer/query-route",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGetVpcPeerRouteApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGetVpcPeerRouteRequest) (*CtvpcGetVpcPeerRouteResponse, error) {
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
	var resp CtvpcGetVpcPeerRouteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGetVpcPeerRouteRequest struct {
	RouteTableID *string `json:"routeTableID,omitempty"` /*  路由表 id  */
	VpcID        string  `json:"vpcID,omitempty"`        /*  虚拟私有云 id  */
	RegionID     string  `json:"regionID,omitempty"`     /*  区域id  */
	PageSize     int32   `json:"pageSize"`               /*  当前页数据条数  */
	PageNumber   int32   `json:"pageNumber"`             /*  当前页  */
}

type CtvpcGetVpcPeerRouteResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGetVpcPeerRouteReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcGetVpcPeerRouteReturnObjResponse struct {
	Results      []*CtvpcGetVpcPeerRouteReturnObjResultsResponse `json:"results"`      /*  路由规则列表  */
	TotalCount   int32                                           `json:"totalCount"`   /*  总条数  */
	TotalPage    int32                                           `json:"totalPage"`    /*  总页数  */
	CurrentCount int32                                           `json:"currentCount"` /*  总页数  */
}

type CtvpcGetVpcPeerRouteReturnObjResultsResponse struct {
	RouteRuleID  *string `json:"routeRuleID,omitempty"`  /*  路由规则 id  */
	Destination  *string `json:"destination,omitempty"`  /*  目的 cidr  */
	RouteTableID *string `json:"routeTableID,omitempty"` /*  路由表 id  */
	VpcID        *string `json:"vpcID,omitempty"`        /*  虚拟私有云 id  */
	NextHopID    *string `json:"nextHopID,omitempty"`    /*  下一跳设备 id  */
}
