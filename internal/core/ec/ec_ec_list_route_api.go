package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcListRouteApi
/* 查询路由 */
type EcEcListRouteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListRouteApi(client *core.CtyunClient) *EcEcListRouteApi {
	return &EcEcListRouteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/route/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListRouteApi) Do(ctx context.Context, credential core.Credential, req *EcEcListRouteRequest) (*EcEcListRouteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("ecID", req.EcID)
	ctReq.AddParam("cgwID", req.CgwID)
	ctReq.AddParam("rtbID", req.RtbID)
	if req.RouteID != nil && *req.RouteID != "" {
		ctReq.AddParam("routeID", *req.RouteID)
	}
	if req.QueryContent != nil && *req.QueryContent != "" {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if req.Status != nil && *req.Status != "" {
		ctReq.AddParam("status", *req.Status)
	}
	if req.RouteType != nil && *req.RouteType != "" {
		ctReq.AddParam("routeType", *req.RouteType)
	}
	if req.NexthopType != nil && *req.NexthopType != "" {
		ctReq.AddParam("nexthopType", *req.NexthopType)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcListRouteResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListRouteRequest struct {
	EcID         string  `json:"ecID"`                   /*  云间高速实例ID  */
	CgwID        string  `json:"cgwID"`                  /*  云网关ID  */
	RtbID        string  `json:"rtbID"`                  /*  路由表ID  */
	RouteID      *string `json:"routeID,omitempty"`      /*  路由ID  */
	QueryContent *string `json:"queryContent,omitempty"` /*  模糊匹配CIDR  */
	Status       *string `json:"status,omitempty"`       /*  运行状态<br/>取值范围:<br/>"1": 正常 <br/>"2":异常  */
	RouteType    *string `json:"routeType,omitempty"`    /*  路由类型<br/>取值范围:<br/>"1":自动学习<br/>"2":自定义（默认）  */
	NexthopType  *string `json:"nexthopType,omitempty"`  /*  下一跳实例类型<br/>取值范围:<br/>"1":vpc<br/>"2":云专线<br/>"3":授权vpc<br/>"4":sdwan<br/>"5":vpn  */
}

type EcEcListRouteResponse struct {
	StatusCode  *int32                          `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                         `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                         `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                         `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                         `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcListRouteReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcListRouteReturnObjResponse struct {
	CurrentCount *int32                                   `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                   `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                   `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcEcListRouteReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcEcListRouteReturnObjResultsResponse struct {
	EcID             *string `json:"ecID"`             /*  云间高速id  */
	CgwID            *string `json:"cgwID"`            /*  云网关id  */
	RtbID            *string `json:"rtbID"`            /*  路由表id  */
	RouteID          *string `json:"routeID"`          /*  路由ID  */
	RouteType        *string `json:"routeType"`        /*  路由类型<br/>取值范围:<br/>"1":自动学习<br/>"2":自定义（默认）  */
	RouteCIDR        *string `json:"routeCIDR"`        /*  子网信息  */
	NexthopType      *string `json:"nexthopType"`      /*  下一跳路由类型<br/>取值范围:<br/>"1":vpc<br/>"2":云专线<br/>"3":SDWAN<br/>"4":VPN<br/>"5":EDS<br/>"20":黑洞路由<br/>"30":跨域连接  */
	NexthopID        *string `json:"nexthopID"`        /*  目的实例ID  */
	IPVersion        *string `json:"IPVersion"`        /*  子网类型<br/>取值范围:<br/>"1":IPv4类型<br/>"2":IPv6类型  */
	RouteDescription *string `json:"routeDescription"` /*  路由描述信息  */
	Status           *string `json:"status"`           /*  运行状态<br/>取值范围<br/>"1":正常<br/>"2":异常  */
	CreateDate       *string `json:"createDate"`       /*  创建时间  */
}
